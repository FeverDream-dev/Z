package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/FeverDream-dev/zsistant/internal/activity"
	"github.com/FeverDream-dev/zsistant/internal/approvals"
	"github.com/FeverDream-dev/zsistant/internal/assistant"
	"github.com/FeverDream-dev/zsistant/internal/browser"
	cfgpkg "github.com/FeverDream-dev/zsistant/internal/config"
	"github.com/FeverDream-dev/zsistant/internal/devmode"
	"github.com/FeverDream-dev/zsistant/internal/jobs"
	"github.com/FeverDream-dev/zsistant/internal/knowledge"
	"github.com/FeverDream-dev/zsistant/internal/llm"
	"github.com/FeverDream-dev/zsistant/internal/memory"
	rt "github.com/FeverDream-dev/zsistant/internal/runtime"
	"github.com/FeverDream-dev/zsistant/internal/tools"
)

// Server is a fully integrated ASSISTANT-FIRST HTTP server.
type Server struct {
	addr           string
	dataPath       string
	uiPath         string
	assistantReg   *assistant.Registry
	memoryStore    *memory.Store
	knowledgeStore *knowledge.Store
	activityLog    string
	config         *cfgpkg.Config
	mux            *http.ServeMux
	brok           *tools.Broker
	diag           devmode.Diagnostics
	startTime      time.Time
	engine         *rt.Engine
	approvalStore  *approvals.Store
}

// New creates a new Server. basePath is the data directory.
// The UI is served from the project root (next to the go.mod).
func New(addr, basePath string) *Server {
	uiPath := filepath.Join(basePath, "ui")
	if _, err := os.Stat(uiPath); err != nil {
		// Fallback: locate ui next to executable or project root
		_, thisFile, _, _ := runtime.Caller(0)
		projectRoot := filepath.Join(thisFile, "..", "..", "..")
		uiPath = filepath.Join(projectRoot, "ui")
		if _, err2 := os.Stat(uiPath); err2 != nil {
			uiPath = "/mnt/projects-ssd/Zsisstant/ui"
		}
	}

	cfgPath, _ := cfgpkg.DefaultPath()
	cfg, _ := cfgpkg.Load(cfgPath)
	if cfg == nil {
		cfg = &cfgpkg.Config{}
	}
	if cfg.ProviderKeys == nil {
		cfg.ProviderKeys = map[string]string{}
	}
	// Also load from env vars
	loadEnvKeys(cfg.ProviderKeys)

	reg := assistant.NewRegistry(basePath)
	memStore := memory.NewStore(basePath)
	knowStore := knowledge.NewStore(basePath)

	// Create tools broker and register built-in tools
	broker := tools.NewBroker(basePath)
	// Register noop built-in tools for demonstration
	broker.Register(&noopTool{name: "search", desc: "Search the web for information"})
	broker.Register(&noopTool{name: "file_read", desc: "Read contents of a file"})
	broker.Register(&noopTool{name: "file_write", desc: "Write content to a file"})
	broker.Register(&noopTool{name: "execute_command", desc: "Execute a shell command (requires approval)"})
	broker.Register(&noopTool{name: "browser_navigate", desc: "Navigate browser to a URL"})
	broker.Register(&noopTool{name: "screenshot", desc: "Take a screenshot of the current page"})

	approvalStore := approvals.NewStore(basePath)
	factory := llm.NewFactory(cfg.ProviderKeys)
	engine := rt.NewEngine(basePath, reg, approvalStore, factory, broker)
	engine.Start()

	s := &Server{
		addr:          addr,
		dataPath:      basePath,
		uiPath:        uiPath,
		assistantReg:  reg,
		memoryStore:   memStore,
		knowledgeStore: knowStore,
		activityLog:   filepath.Join(basePath, "activity.jsonl"),
		config:        cfg,
		mux:           http.NewServeMux(),
		brok:          broker,
		startTime:     time.Now(),
		engine:        engine,
		approvalStore: approvalStore,
	}
	s.diag = devmode.Diagnostics{
		WorkspaceRoot: cfgpkg.ExpandPath(cfg.DataPath),
		ConfigPath:    cfgPath,
		GoVersion:     runtime.Version(),
	}
	s.updateDiagnostics()
	s.registerRoutes()
	return s
}

func (s *Server) updateDiagnostics() {
	list, _ := s.assistantReg.List()
	s.diag.ActiveAssistants = len(list)
	s.diag.RegisteredTools = len(s.brok.ListTools())
	// Count connected channels
	connected := 0
	for _, a := range list {
		for _, c := range a.Channels {
			if c.Status == "connected" {
				connected++
			}
		}
	}
	s.diag.ConnectedChannels = connected
	s.diag.UptimeSeconds = int64(time.Since(s.startTime).Seconds())
}

func (s *Server) registerRoutes() {
	// Root SPA redirect
	s.mux.HandleFunc("/", s.handleIndex)
	s.mux.HandleFunc("/health", s.handleHealth)

	// Static files
	s.mux.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir(s.uiPath))))

	// API Routes
	s.mux.HandleFunc("/api/assistants", s.handleAssistantsRoot)
	s.mux.HandleFunc("/api/assistants/", s.handleAssistantsScoped)
	s.mux.HandleFunc("/api/providers", s.handleProviders)
	s.mux.HandleFunc("/api/models", s.handleModels)
	s.mux.HandleFunc("/api/settings", s.handleSettings)
	s.mux.HandleFunc("/api/jobs", s.handleGlobalJobs)
	s.mux.HandleFunc("/api/tools", s.handleTools)
	s.mux.HandleFunc("/api/activity", s.handleActivity)
	s.mux.HandleFunc("/api/runtime/status", s.handleRuntimeStatus)
	s.mux.HandleFunc("/api/approvals", s.handleApprovals)
	s.mux.HandleFunc("/api/approvals/", s.handleApprovalsScoped)
	s.mux.HandleFunc("/api/runtime/events", s.handleRuntimeEvents)
	s.mux.HandleFunc("/api/runtime/activity", s.handleRuntimeActivity)
}

// Run starts the HTTP server.
func (s *Server) Run() error {
	fmt.Printf("Zsistant v0.1.0 (ASSISTANT-FIRST) listening on http://%s\n", s.addr)
	fmt.Printf("Data path: %s\n", s.dataPath)
	fmt.Printf("UI path:   %s\n", s.uiPath)
	return http.ListenAndServe(s.addr, cors(s.mux))
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func loadEnvKeys(keys map[string]string) {
	mappings := map[string]string{
		"OPENAI_API_KEY":         "openai",
		"ANTHROPIC_API_KEY":        "anthropic",
		"GOOGLE_API_KEY":           "google",
		"COHERE_API_KEY":           "cohere",
		"MISTRAL_API_KEY":          "mistral",
		"GROQ_API_KEY":             "groq",
		"TOGETHER_API_KEY":         "together",
		"FIREWORKS_API_KEY":        "fireworks",
		"PERPLEXITY_API_KEY":       "perplexity",
		"DEEPSEEK_API_KEY":         "deepseek",
		"XAI_API_KEY":              "xai",
		"AZURE_OPENAI_KEY":         "azure-openai",
		"OPENROUTER_API_KEY":       "openrouter",
		"AI21_API_KEY":             "ai21",
		"REPLICATE_API_TOKEN":      "replicate",
		"NOVITA_API_KEY":           "novita",
		"HYPERBOLIC_API_KEY":       "hyperbolic",
		"SILICONFLOW_API_KEY":      "siliconflow",
		"DEEPINFRA_API_KEY":        "deepinfra",
		"NVIDIA_API_KEY":           "nvidia",
		"SAMBANOVA_API_KEY":        "sambanova",
		"LAMBDA_API_KEY":           "lambda",
		"FRIENDLI_API_KEY":         "friendliai",
		"CHUTES_API_KEY":           "chutes",
		"CLOUDFLARE_API_KEY":       "cloudflare",
		"OCTOAI_API_KEY":           "octoai",
		"PREDIBASE_API_TOKEN":      "predibase",
		"POE_API_KEY":              "poe",
		"MOONSHOT_API_KEY":         "moonshot",
		"YI_API_KEY":               "01-ai",
		"BAIDU_API_KEY":            "baidu",
		"DASHSCOPE_API_KEY":        "alibaba",
		"TENCENT_API_KEY":              "tencent",
		"AWS_BEDROCK_KEY":              "aws-bedrock",
		"GOOGLE_VERTEX_KEY":            "google-vertex",
		"VLLM_KEY":                     "vllm-local",
		"TGI_KEY":                      "tgi-local",
		"TABBYML_KEY":                  "tabbyml-local",
		"LMSTUDIO_KEY":                 "lmstudio-local",
		"JAN_KEY":                      "jan-local",
		"HUGGINGFACE_TOKEN":            "huggingface",
		"BASETEN_KEY":                  "baseten",
		"ANYSCALE_KEY":                 "anyscale",
		"ZHIPU_KEY":                    "zhipu",
		"MINIMAX_KEY":                  "minimax",
		"STEPFUN_KEY":                  "stepfun",
		"SENSETIME_KEY":                "sensetime",
		"BAICHUAN_KEY":                 "baichuan",
		"BYTEDANCE_KEY":                "bytedance",
		"YOU_KEY":                      "you",
		"PHIND_KEY":                    "phind",
		"CODEIUM_KEY":                  "codeium",
	}
	for env, provider := range mappings {
		if v := os.Getenv(env); v != "" {
			keys[provider] = v
		}
	}
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/ui/", http.StatusFound)
		return
	}
	http.NotFound(w, r)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	s.updateDiagnostics()
	resp := map[string]interface{}{
		"status":      "ok",
		"version":     "v0.1.0",
		"uptime_sec":  s.diag.UptimeSeconds,
		"assistants":  s.diag.ActiveAssistants,
		"tools":       s.diag.RegisteredTools,
		"channels":    s.diag.ConnectedChannels,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleProviders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	type provWithConfig struct {
		llm.ProviderInfo
		Configured bool `json:"configured"`
	}
	var out []provWithConfig
	for _, pi := range llm.BuiltInProviders {
		configured := pi.Name == "ollama-local"
		if !configured {
			configured = s.config.ProviderKeys[pi.Name] != ""
		}
		out = append(out, provWithConfig{ProviderInfo: pi, Configured: configured})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (s *Server) handleModels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	type flatModel struct {
		llm.ModelInfo
		Provider string `json:"provider"`
	}
	var out []flatModel
	for _, pi := range llm.BuiltInProviders {
		for _, m := range pi.Models {
			// Set provider name on model
			mCopy := m
			out = append(out, flatModel{mCopy, pi.Name})
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (s *Server) handleSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		redacted := map[string]interface{}{
			"theme":          s.config.Theme,
			"dev_mode":       s.config.DevMode,
			"data_path":      s.config.DataPath,
			"default_model":  s.config.DefaultModel,
			"provider_keys":  redactMap(s.config.ProviderKeys),
			"providers":      s.config.Providers,
		}
		json.NewEncoder(w).Encode(redacted)
	case http.MethodPost:
		var body struct {
			DevMode       bool              `json:"dev_mode"`
			Theme         string            `json:"theme"`
			ProviderKeys  map[string]string `json:"provider_keys"`
			DefaultModel  string            `json:"default_model"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
			return
		}
		s.config.DevMode = body.DevMode
		if body.Theme != "" {
			s.config.Theme = body.Theme
		}
		if body.DefaultModel != "" {
			s.config.DefaultModel = body.DefaultModel
		}
		for k, v := range body.ProviderKeys {
			if v != "" && v != "[redacted]" && v != "__redacted" {
				s.config.ProviderKeys[k] = v
			}
		}
		cfgPath, _ := cfgpkg.DefaultPath()
		_ = cfgpkg.Save(cfgPath, s.config)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func redactMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	out := make(map[string]string, len(m))
	for k := range m {
		out[k] = "[configured]"
	}
	return out
}

func (s *Server) handleAssistantsRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/assistants" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		list, err := s.assistantReg.List()
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var a assistant.Assistant
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
			return
		}
		if a.ID == "" || a.Name == "" {
			http.Error(w, `{"error":"id and name required"}`, http.StatusBadRequest)
			return
		}
		created, err := s.assistantReg.Create(a.ID, a.Name)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
			return
		}
		// Apply extra fields
		if a.Description != "" {
			created.Description = a.Description
		}
		if a.Purpose != "" {
			created.Purpose = a.Purpose
		}
		if a.Persona.Tone != "" {
			created.Persona = a.Persona
		}
		if a.MemoryPolicy.Scope != "" {
			created.MemoryPolicy = a.MemoryPolicy
		}
		if a.DefaultModel != "" {
			created.DefaultModel = a.DefaultModel
			// Infer provider
			if pi, _, err := llm.FindModel(a.DefaultModel); err == nil {
				created.ProviderName = pi.Name
			}
		}
		if len(a.Channels) > 0 {
			created.Channels = a.Channels
		} else {
			created.Channels = []assistant.ChannelConfig{
				{ChannelType: "web_ui", Status: "connected", CreatedAt: time.Now()},
			}
		}
		created.JobsEnabled = a.JobsEnabled
		created.Status = assistant.StatusActive
		err = s.assistantReg.Update(created)
		if err != nil {
			// best effort
		}
		// Log activity
		activity.Log(s.dataPath, activity.ActivityEvent{
			AssistantID: a.ID,
			EventType:   "assistant.created",
			Message:     fmt.Sprintf("Assistant %s created", a.Name),
			Severity:    "info",
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(created)
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleAssistantsScoped(w http.ResponseWriter, r *http.Request) {
	trimmed := strings.TrimPrefix(r.URL.Path, "/api/assistants/")
	parts := strings.Split(trimmed, "/")
	if len(parts) == 0 {
		http.NotFound(w, r)
		return
	}
	id := parts[0]

	switch len(parts) {
	case 1:
		s.handleAssistantDetail(w, r, id)
	case 2:
		switch parts[1] {
		case "jobs":
			s.handleJobsForAssistant(w, r, id)
		case "memory":
			s.handleMemoryForAssistant(w, r, id)
		case "knowledge":
			s.handleKnowledgeForAssistant(w, r, id)
		case "channels":
			s.handleChannelsForAssistant(w, r, id)
		case "logs":
			s.handleLogsForAssistant(w, r, id)
		case "chat":
			s.handleChatForAssistant(w, r, id, false)
		case "browser":
			s.handleBrowserForAssistant(w, r, id)
		case "run":
			s.handleRunAssistant(w, r, id)
		case "pause":
			s.handlePauseAssistant(w, r, id)
		case "resume":
			s.handleResumeAssistant(w, r, id)
		case "state":
			s.handleAssistantState(w, r, id)
		default:
			http.NotFound(w, r)
		}
	case 3:
		if parts[1] == "chat" && parts[2] == "stream" {
			s.handleChatForAssistant(w, r, id, true)
			return
		}
		http.NotFound(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) handleAssistantDetail(w http.ResponseWriter, r *http.Request, id string) {
	switch r.Method {
	case http.MethodGet:
		a, err := s.assistantReg.Get(id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(a)
	case http.MethodPut:
		var a assistant.Assistant
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
			return
		}
		a.ID = id
		if err := s.assistantReg.Update(&a); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
			return
		}
		activity.Log(s.dataPath, activity.ActivityEvent{
			AssistantID: id,
			EventType:   "assistant.updated",
			Message:     fmt.Sprintf("Assistant %s updated", a.Name),
			Severity:    "info",
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(a)
	case http.MethodDelete:
		if err := s.assistantReg.Delete(id); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
			return
		}
		activity.Log(s.dataPath, activity.ActivityEvent{
			AssistantID: id,
			EventType:   "assistant.deleted",
			Message:     fmt.Sprintf("Assistant %s deleted", id),
			Severity:    "info",
		})
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleJobsForAssistant(w http.ResponseWriter, r *http.Request, id string) {
	q := jobs.NewQueue(filepath.Join(s.dataPath, "assistants", id))
	switch r.Method {
	case http.MethodGet:
		list, err := q.List()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("[]"))
			return
		}
		// Dereference
		out := []jobs.Job{}
		for _, j := range list {
			out = append(out, *j)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(out)
	case http.MethodPost:
		var body struct {
			Name     string   `json:"name"`
			Purpose  string   `json:"purpose"`
			Schedule string   `json:"schedule,omitempty"`
			Type     string   `json:"type"` // "manual" or "cron"
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
			return
		}
		var j *jobs.Job
		if body.Type == "manual" {
			j = jobs.NewJob(id, body.Name, body.Purpose)
		} else {
			j = jobs.NewScheduledJob(id, body.Name, body.Purpose, body.Schedule)
		}
		err := q.Enqueue(j)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
			return
		}
		activity.Log(s.dataPath, activity.ActivityEvent{
			AssistantID: id,
			EventType:   "job.scheduled",
			Message:     fmt.Sprintf("Job %s scheduled for assistant %s", j.Name, id),
			Severity:    "info",
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(j)
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleMemoryForAssistant(w http.ResponseWriter, r *http.Request, id string) {
	switch r.Method {
	case http.MethodGet:
		entries, err := s.memoryStore.List(id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entries)
	case http.MethodPost:
		var body struct {
			Category string `json:"category"`
			Content  string `json:"content"`
			Source   string `json:"source,omitempty"`
			Approved bool   `json:"approved"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
			return
		}
		m, err := s.memoryStore.Add(id, body.Category, body.Content, body.Source, body.Approved)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(m)
	case http.MethodDelete:
		// /api/assistants/{id}/memory/{memid}
		// This path format uses extra segments — but the router only passed us 2 parts.
		// We'll check the raw path for extra segments.
		extra := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("/api/assistants/%s/memory", id))
		extra = strings.TrimPrefix(extra, "/")
		if extra != "" {
			if err := s.memoryStore.Delete(extra); err != nil {
				http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, `{"error":"memory id required"}`, http.StatusBadRequest)
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleKnowledgeForAssistant(w http.ResponseWriter, r *http.Request, id string) {
	switch r.Method {
	case http.MethodGet:
		entries, err := s.knowledgeStore.List(id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entries)
	case http.MethodPut:
		var body struct {
			AttachIDs []string `json:"attach_ids"`   // IDs to attach
			DetachIDs []string `json:"detach_ids"`   // IDs to detach
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
			return
		}
		for _, kid := range body.AttachIDs {
			s.knowledgeStore.AttachToAssistant(kid, id)
		}
		for _, kid := range body.DetachIDs {
			s.knowledgeStore.DetachFromAssistant(kid, id)
		}
		entries, _ := s.knowledgeStore.List(id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entries)
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleChannelsForAssistant(w http.ResponseWriter, r *http.Request, id string) {
	switch r.Method {
	case http.MethodGet:
		a, err := s.assistantReg.Get(id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(a.Channels)
	case http.MethodPut:
		var payload struct {
			Channels []assistant.ChannelConfig `json:"channels"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
			return
		}
		a, err := s.assistantReg.Get(id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusNotFound)
			return
		}
		a.Channels = payload.Channels
		for i := range a.Channels {
			if a.Channels[i].CreatedAt.IsZero() {
				a.Channels[i].CreatedAt = time.Now()
			}
		}
		_ = s.assistantReg.Update(a)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(a.Channels)
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleLogsForAssistant(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	limit := 50
	entries, err := activity.ReadAssistant(s.dataPath, id, limit)
	if err != nil {
		entries = []activity.ActivityEvent{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func (s *Server) handleBrowserForAssistant(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	cfg := browser.DefaultConfig()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
}

func (s *Server) handleRunAssistant(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	state := rt.LoadAssistantState(s.dataPath, id)
	state.Status = "running"
	rt.SaveAssistantState(s.dataPath, id, state)
	go func() {
		a, _ := s.assistantReg.Get(id)
		if a != nil {
			s.engine.RunNow(*a)
		}
	}()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "running", "assistant_id": id})
}

func (s *Server) handlePauseAssistant(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	state := rt.LoadAssistantState(s.dataPath, id)
	state.Status = "paused"
	state.Enabled = false
	rt.SaveAssistantState(s.dataPath, id, state)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "paused", "assistant_id": id})
}

func (s *Server) handleResumeAssistant(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	state := rt.LoadAssistantState(s.dataPath, id)
	state.Status = "idle"
	state.Enabled = true
	rt.SaveAssistantState(s.dataPath, id, state)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "resumed", "assistant_id": id})
}

func (s *Server) handleAssistantState(w http.ResponseWriter, r *http.Request, id string) {
	switch r.Method {
	case http.MethodGet:
		state := rt.LoadAssistantState(s.dataPath, id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(state)
	case http.MethodPut:
		var body rt.RuntimeState
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
			return
		}
		if err := rt.SaveAssistantState(s.dataPath, id, body); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleTools(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	infos := s.brok.ListTools()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(infos)
}

func (s *Server) handleGlobalJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var all []jobs.Job
	// No efficient global scan without knowing assistant IDs — return empty for now.
	// In production, we'd maintain a global job index.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(all)
}

func (s *Server) handleActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	entries, err := activity.ReadGlobal(s.dataPath, 100)
	if err != nil {
		entries = []activity.ActivityEvent{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func (s *Server) handleRuntimeStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	resp := map[string]interface{}{
		"running":  s.engine.IsRunning(),
		"uptime":   time.Since(s.startTime).Seconds(),
		"tick_sec": 30,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleApprovals(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list, err := s.approvalStore.List("", "", 100)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var req approvals.Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
			return
		}
		if err := s.approvalStore.Create(&req); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(req)
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleApprovalsScoped(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/approvals/")
	if id == "" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodPost:
		var body struct {
			Decision string `json:"decision"`
			By       string `json:"by"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
			return
		}
		if body.Decision != approvals.StatusApproved && body.Decision != approvals.StatusDenied {
			http.Error(w, `{"error":"decision must be approved or denied"}`, http.StatusBadRequest)
			return
		}
		req, err := s.approvalStore.Resolve(id, body.By, body.Decision)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(req)
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleRuntimeEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	activityPath := filepath.Join(s.dataPath, "activity.jsonl")
	if _, err := os.Stat(activityPath); os.IsNotExist(err) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}
	data, err := os.ReadFile(activityPath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	var out []map[string]interface{}
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		var ev map[string]interface{}
		if json.Unmarshal([]byte(l), &ev) == nil {
			out = append(out, ev)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (s *Server) handleRuntimeActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	limit := 100
	entries, err := activity.ReadGlobal(s.dataPath, limit)
	if err != nil {
		entries = []activity.ActivityEvent{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

// handleChatForAssistant routes chat through the LLM layer.
func (s *Server) handleChatForAssistant(w http.ResponseWriter, r *http.Request, id string, streaming bool) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var payload struct{ Message string `json:"message"` }
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
		return
	}
	if payload.Message == "" {
		http.Error(w, `{"error":"message required"}`, http.StatusBadRequest)
		return
	}

	// Get or create assistant
	a, err := s.assistantReg.Get(id)
	if err != nil {
		a, _ = s.assistantReg.Create(id, id)
		a.Status = assistant.StatusActive
		a.Channels = []assistant.ChannelConfig{{ChannelType: "web_ui", Status: "connected", CreatedAt: time.Now()}}
		s.assistantReg.Update(a)
	}

	now := time.Now()
	a.LastActivityAt = &now
	s.assistantReg.Update(a)

	// Record job
	jobID := fmt.Sprintf("chat-%d", time.Now().UnixNano())
	job := jobs.NewJob(id, "chat", payload.Message)
	job.ID = jobID
	job.Status = jobs.StatusRunning
	q := jobs.NewQueue(filepath.Join(s.dataPath, "assistants", id))
	_ = q.Enqueue(job)

	// Determine provider and model
	modelID := a.DefaultModel
	if modelID == "" {
		modelID = "gpt-4o-mini"
	}
	providerName := a.ProviderName

	// Resolve API key
	factory := llm.NewFactory(s.config.ProviderKeys)
	var prov llm.Provider
	var pInfo *llm.ProviderInfo
	var mInfo *llm.ModelInfo
	var provErr error

	if providerName != "" {
		prov, pInfo, provErr = factory.CreateByProvider(providerName)
		if provErr == nil && mInfo != nil {
			modelID = mInfo.ID
		}
	}
	if prov == nil {
		prov, pInfo, mInfo, provErr = factory.Create(modelID)
	}

	if provErr != nil || prov == nil {
		job.Status = jobs.StatusFailed
		job.LastError = "Provider not configured"
		job.Result = "Provider not configured"
		_ = q.Update(job)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   fmt.Sprintf("Provider not configured for model %q. Set the API key in Settings.", modelID),
			"model":   modelID,
			"message": payload.Message,
		})
		return
	}

	key := ""
	if pInfo != nil {
		key = s.config.ProviderKeys[pInfo.Name]
	}
	if key == "" {
		job.Status = jobs.StatusFailed
		job.LastError = "API key missing"
		_ = q.Update(job)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("API key not configured for provider %q. Add it in Settings or via environment variable.", pInfo.Name),
			"model": modelID,
		})
		return
	}

	if streaming {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		flusher, ok := w.(http.Flusher)
		if !ok {
			job.Status = jobs.StatusFailed
			job.LastError = "streaming not supported"
			_ = q.Update(job)
			return
		}

		chunkCh := make(chan string, 64)
		doneCh := make(chan error, 1)

		// Type assert to streaming interface
		start := time.Now()
		switch sp := prov.(type) {
		case interface{ Stream(string, chan<- string, chan<- error) }:
			go sp.Stream(payload.Message, chunkCh, doneCh)
		default:
			// Fallback: Complete and emit chunks
			go func() {
				resp, err := prov.Complete(payload.Message)
				if err != nil {
					doneCh <- err
					return
				}
				parts := splitResponse(resp, 8)
				for _, part := range parts {
					chunkCh <- part
				}
				doneCh <- nil
			}()
		}

		var fullResp strings.Builder
		for {
			select {
			case chunk := <-chunkCh:
				fmt.Fprintf(w, "data: %s\n\n", jsonEscape(chunk))
				flusher.Flush()
				fullResp.WriteString(chunk)
			case err := <-doneCh:
				if err != nil {
					fmt.Fprintf(w, "event: error\ndata: %s\n\n", jsonEscape(err.Error()))
					flusher.Flush()
					job.Status = jobs.StatusFailed
					job.LastError = err.Error()
				} else {
					fmt.Fprintf(w, "data: [DONE]\n\n")
					flusher.Flush()
					job.Status = jobs.StatusCompleted
					job.Result = fullResp.String()
				}
				_ = q.Update(job)
				// Log activity
				activity.Log(s.dataPath, activity.ActivityEvent{
					AssistantID: id,
					EventType:   "message.sent",
					Message:     fmt.Sprintf("Chat message from %s", r.RemoteAddr),
					Severity:    "info",
				})
				activity.Log(s.dataPath, activity.ActivityEvent{
					AssistantID: id,
					EventType:   "job.completed",
					Message:     fmt.Sprintf("Chat job completed in %v", time.Since(start)),
					Severity:    "info",
				})
				return
			}
		}
	}

	// Non-streaming
	start := time.Now()
	resp, err := prov.Complete(payload.Message)
	if err != nil {
		job.Status = jobs.StatusFailed
		job.LastError = err.Error()
		_ = q.Update(job)
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
		return
	}

	job.Status = jobs.StatusCompleted
	job.Result = resp
	_ = q.Update(job)

	activity.Log(s.dataPath, activity.ActivityEvent{
		AssistantID: id,
		EventType:   "job.completed",
		Message:     fmt.Sprintf("Chat completed in %v", time.Since(start)),
		Severity:    "info",
	})
	activity.Log(s.dataPath, activity.ActivityEvent{
		AssistantID: id,
		EventType:   "message.sent",
		Message:     payload.Message,
		Severity:    "info",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"job_id":    jobID,
		"response":  resp,
		"status":    job.Status,
		"model":     modelID,
		"provider":  pInfo.Name,
		"duration":  time.Since(start).Milliseconds(),
	})
}

func splitResponse(s string, n int) []string {
	if n <= 1 || len(s) == 0 {
		return []string{s}
	}
	chunk := len(s) / n
	if chunk == 0 {
		chunk = len(s)
	}
	var out []string
	for i := 0; i < len(s); i += chunk {
		end := i + chunk
		if end > len(s) {
			end = len(s)
		}
		out = append(out, s[i:end])
	}
	return out
}

func jsonEscape(s string) string {
	b, _ := json.Marshal(s)
	return string(b[1 : len(b)-1])
}

// --- no-op tool for built-in demos ---

type noopTool struct {
	name string
	desc string
}

func (t *noopTool) Name() string        { return t.name }
func (t *noopTool) Description() string   { return t.desc }
func (t *noopTool) Parameters() []string { return nil }
func (t *noopTool) Execute(call tools.ToolCall) tools.ToolResult {
	return tools.ToolResult{
		CallID:  call.ID,
		ToolName: t.name,
		Success: true,
		Output:  fmt.Sprintf("[demo] %s would execute here", t.name),
	}
}
