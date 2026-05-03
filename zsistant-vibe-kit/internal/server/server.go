package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/FeverDream-dev/zsistant/internal/agents"
	zkconfig "github.com/FeverDream-dev/zsistant/internal/config"
	"github.com/FeverDream-dev/zsistant/internal/jobs"
	"github.com/FeverDream-dev/zsistant/internal/llm"
	"github.com/FeverDream-dev/zsistant/internal/tools"
)

// Server is the Zsistant HTTP server.
type Server struct {
	addr       string
	basePath   string
	assetsPath string
	mux        *http.ServeMux
}

// New creates a new server instance.
func New(addr, basePath string) *Server {
    s := &Server{
        addr:       addr,
        basePath:   basePath,
        assetsPath: filepath.Join("assets", "brand"),
        mux:        http.NewServeMux(),
    }
    s.registerRoutes()
    return s
}

// NewWithAssets creates a server with a custom assets directory (useful for tests).
func NewWithAssets(addr, basePath, assetsPath string) *Server {
	s := &Server{
		addr:       addr,
		basePath:   basePath,
		assetsPath: assetsPath,
		mux:        http.NewServeMux(),
	}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/", s.handleApp())
	s.mux.HandleFunc("/chat", s.handleApp())
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/api/agents", s.handleAgents)
	s.mux.HandleFunc("/api/providers", s.handleProviders)
	s.mux.HandleFunc("/api/models", s.handleModels)
	s.mux.HandleFunc("/api/chat", s.handleChatAPI)
	s.mux.HandleFunc("/api/chat/stream", s.handleChatStream)
	s.mux.HandleFunc("/api/jobs/", s.handleJobsAPI)
	s.mux.HandleFunc("/api/settings", s.handleSettings)
	s.mux.HandleFunc("/api/conversations", s.handleConversations)
	s.mux.HandleFunc("/api/tools", s.handleTools)
	s.mux.HandleFunc("/api/tools/execute", s.handleToolsExecute)
	s.mux.HandleFunc("/api/agents/", s.handleAgentPermissions)
	s.mux.Handle("/assets/brand/", http.StripPrefix("/assets/brand/", http.FileServer(http.Dir(s.assetsPath))))
	s.mux.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("ui"))))
}

// Run starts the HTTP server.
func (s *Server) Run() error {
	fmt.Printf("Zsistant server listening on http://%s\n", s.addr)
	return http.ListenAndServe(s.addr, s.mux)
}

// buildRouter creates an LLM router configured from the user's config file.
// It registers real providers when API keys are present, with mock as final fallback.
func (s *Server) buildRouter() *llm.Router {
	router := llm.NewRouter(60 * time.Second)
	cfgPath, _ := zkconfig.DefaultPath()
	cfg, _ := zkconfig.Load(cfgPath)

	if cfg == nil || cfg.Secrets == nil {
		router.Register(llm.TaskCheap, llm.NewMockProvider())
		router.Register(llm.TaskStrong, llm.NewMockProvider())
		router.Register(llm.TaskCoding, llm.NewMockProvider())
		return router
	}

	ollamaKey := cfg.Secrets["ollama_api_key"]
	if ollamaKey != "" {
		router.Register(llm.TaskCheap, llm.NewOllamaProvider(ollamaKey, ""))
		router.Register(llm.TaskStrong, llm.NewOllamaProvider(ollamaKey, ""))
		router.Register(llm.TaskCoding, llm.NewOllamaProvider(ollamaKey, ""))
	}

	openaiKey := cfg.Secrets["openai_api_key"]
	if openaiKey != "" {
		openaiURL := cfg.Secrets["openai_base_url"]
		if openaiURL == "" {
			openaiURL = "https://api.openai.com/v1"
		}
		openaiModel := cfg.Secrets["openai_model"]
		if openaiModel == "" {
			openaiModel = "gpt-4o-mini"
		}
		router.Register(llm.TaskCheap, llm.NewOpenAIProvider(openaiKey, openaiURL, openaiModel))
		router.Register(llm.TaskStrong, llm.NewOpenAIProvider(openaiKey, openaiURL, openaiModel))
	}

	zaiKey := cfg.Secrets["zai_api_key"]
	if zaiKey != "" {
		zaiModel := cfg.Secrets["zai_model"]
		if zaiModel == "" {
			zaiModel = "GLM-5.1"
		}
		router.Register(llm.TaskCoding, llm.NewZAICodingProvider(zaiKey, zaiModel))
	}

	opencodeKey := cfg.Secrets["opencode_api_key"]
	if opencodeKey != "" {
		opencodeModel := cfg.Secrets["opencode_model"]
		if opencodeModel == "" {
			opencodeModel = "gpt-5.1"
		}
		router.Register(llm.TaskCheap, llm.NewOpenCodeProvider(opencodeKey, opencodeModel))
		router.Register(llm.TaskStrong, llm.NewOpenCodeProvider(opencodeKey, opencodeModel))
	}

	// Always register mock as final fallback
	router.Register(llm.TaskCheap, llm.NewMockProvider())
	router.Register(llm.TaskStrong, llm.NewMockProvider())
	router.Register(llm.TaskCoding, llm.NewMockProvider())
	return router
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleProviders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	router := s.buildRouter()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(router.Health())
}

func (s *Server) handleAgents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		reg := agents.NewRegistry(s.basePath)
		list, err := reg.List()
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var body struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Role string `json:"role"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
			return
		}
		if body.ID == "" || body.Name == "" || body.Role == "" {
			http.Error(w, `{"error":"id, name, and role are required"}`, http.StatusBadRequest)
			return
		}
		reg := agents.NewRegistry(s.basePath)
		a, err := reg.Create(body.ID, body.Name, body.Role)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(a)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

type chatRequest struct {
	AgentID string `json:"agent_id"`
	Message string `json:"message"`
	Model   string `json:"model,omitempty"`
}

type chatResponse struct {
	JobID    string `json:"job_id"`
	Response string `json:"response"`
	Status   string `json:"status"`
}

func (s *Server) handleChatAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
		return
	}
	if req.Message == "" {
		http.Error(w, `{"error":"message is required"}`, http.StatusBadRequest)
		return
	}

	if req.AgentID == "" {
		req.AgentID = "default"
	}

	reg := agents.NewRegistry(s.basePath)
	a, err := reg.Get(req.AgentID)
	if err != nil {
		a, err = reg.Create(req.AgentID, req.AgentID, "web chat agent")
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, "failed to create agent"), http.StatusInternalServerError)
			return
		}
	}

	job := jobs.NewJob(req.AgentID, "web", req.Message)

	// Persist job to queue
	q := jobs.NewQueue(a.WorkspaceRoot)
	_ = q.Enqueue(job)

	// Log job created event
	evt := jobs.JobEvent{
		JobID:     job.ID,
		AgentID:   req.AgentID,
		EventType: "job.created",
		Message:   "Job created from web chat",
		CreatedAt: job.CreatedAt,
	}
	_ = jobs.AppendEvent(a.WorkspaceRoot, evt)

	// Use LLM router with fallback support
	router := s.buildRouter()

	response, providerName, err := router.Route(llm.TaskCheap, req.Message)
	if err != nil {
		job.Status = "failed"
		job.LastError = err.Error()
		job.UpdatedAt = time.Now()
		_ = q.Update(job)

		// Check for loop
		ld := jobs.NewLoopDetector(3)
		isLoop, reason := ld.Check([]jobs.JobEvent{evt, {EventType: "job.failed", Message: err.Error()}})
		if isLoop {
			_, _ = q.Pause(job.ID, reason)
			// Log loop detection event
			loopEvt := jobs.JobEvent{
				JobID:     job.ID,
				AgentID:   req.AgentID,
				EventType: "job.paused",
				Message:   reason,
				CreatedAt: time.Now(),
			}
			_ = jobs.AppendEvent(a.WorkspaceRoot, loopEvt)
		}

		// Log error event
		evt = jobs.JobEvent{
			JobID:     job.ID,
			AgentID:   req.AgentID,
			EventType: "job.failed",
			Message:   err.Error(),
			CreatedAt: time.Now(),
		}
		_ = jobs.AppendEvent(a.WorkspaceRoot, evt)
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
		return
	}

	job.Status = "completed"
	job.Result = response
	now := time.Now()
	job.UpdatedAt = now
	job.CompletedAt = &now
	_ = q.Update(job)

	// Log completion event with provider info
	evt = jobs.JobEvent{
		JobID:     job.ID,
		AgentID:   req.AgentID,
		EventType: "job.completed",
		Message:   response,
		Metadata:  map[string]string{"provider": providerName},
		CreatedAt: now,
	}
	_ = jobs.AppendEvent(a.WorkspaceRoot, evt)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatResponse{
		JobID:    job.ID,
		Response: response,
		Status:   job.Status,
	})
}

func (s *Server) handleChatStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
		return
	}
	if req.Message == "" {
		http.Error(w, `{"error":"message is required"}`, http.StatusBadRequest)
		return
	}
	if req.AgentID == "" {
		req.AgentID = "default"
	}

	reg := agents.NewRegistry(s.basePath)
	a, err := reg.Get(req.AgentID)
	if err != nil {
		a, err = reg.Create(req.AgentID, req.AgentID, "web chat agent")
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, "failed to create agent"), http.StatusInternalServerError)
			return
		}
	}

	job := jobs.NewJob(req.AgentID, "web-stream", req.Message)
	q := jobs.NewQueue(a.WorkspaceRoot)
	_ = q.Enqueue(job)

	// Setup SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, `{"error":"streaming not supported"}`, http.StatusInternalServerError)
		return
	}

	router := s.buildRouter()
	chain := router.ChainFor(llm.TaskCheap)
	if len(chain) == 0 {
		fmt.Fprintf(w, "event: error\ndata: %s\n\n", `{"error":"no providers available"}`)
		flusher.Flush()
		return
	}

	for _, p := range chain {
		streamer, canStream := p.(llm.Streamer)
		if !canStream {
			continue
		}
		chunkCh := make(chan string, 32)
		doneCh := make(chan error, 1)
		go streamer.Stream(req.Message, chunkCh, doneCh)

		var fullContent string
		var streamErr error
	StreamLoop:
		for {
			select {
			case chunk := <- chunkCh:
				if chunk == "" {
					continue
				}
				fullContent += chunk
				escaped, _ := json.Marshal(chunk)
				fmt.Fprintf(w, "event: chunk\ndata: %s\n\n", escaped)
				flusher.Flush()
			case streamErr = <- doneCh:
				break StreamLoop
			}
		}

		if streamErr == nil {
			job.Status = "completed"
			job.Result = fullContent
			now := time.Now()
			job.UpdatedAt = now
			job.CompletedAt = &now
			_ = q.Update(job)
			h := p.Health()
			fmt.Fprintf(w, "event: done\ndata: %s\n\n", fmt.Sprintf(`{"job_id":"%s","provider":"%s","status":"completed"}`, job.ID, h.Name))
			flusher.Flush()
			return
		}
		// Fallback: try next provider
		fmt.Fprintf(w, "event: provider_error\ndata: %s\n\n", fmt.Sprintf(`{"provider":"%s","error":%q}`, p.Health().Name, streamErr.Error()))
		flusher.Flush()
	}

	// All providers failed
	job.Status = "failed"
	job.LastError = "all streaming providers failed"
	job.UpdatedAt = time.Now()
	_ = q.Update(job)
	fmt.Fprintf(w, "event: error\ndata: %s\n\n", `{"error":"all providers failed"}`)
	flusher.Flush()
}

func (s *Server) handleJobsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Extract agent_id from path: /api/jobs/<agent_id>
	agentID := filepath.Base(r.URL.Path)
	if agentID == "" || agentID == "jobs" {
		http.Error(w, `{"error":"agent_id is required"}`, http.StatusBadRequest)
		return
	}

	reg := agents.NewRegistry(s.basePath)
	a, err := reg.Get(agentID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, "agent not found"), http.StatusNotFound)
		return
	}

	// Read audit.jsonl for this agent
	var events []jobs.JobEvent
	auditPath := filepath.Join(a.WorkspaceRoot, "audit.jsonl")
	f, err := os.Open(auditPath)
	if err != nil {
		if os.IsNotExist(err) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]jobs.JobEvent{})
			return
		}
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var evt jobs.JobEvent
		if err := json.Unmarshal(scanner.Bytes(), &evt); err != nil {
			continue
		}
		events = append(events, evt)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (s *Server) handleApp() http.HandlerFunc {
	var uiHTML []byte
	var err error
	for _, p := range []string{"ui/index.html", "../../ui/index.html"} {
		uiHTML, err = os.ReadFile(p)
		if err == nil {
			break
		}
	}
	if err != nil {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(`<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title>Zsistant</title></head><body><div id="conversationList"></div><div id="messages"></div></body></html>`))
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(uiHTML)
	}
}

func (s *Server) handleModels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	grouped := r.URL.Query().Get("grouped") == "true"

	if grouped {
		type providerGroup struct {
			Name    string           `json:"name"`
			Models  []llm.ModelInfo  `json:"models"`
		}
		var groups []providerGroup
		for _, p := range llm.BuiltInProviders {
			groups = append(groups, providerGroup{Name: p.Name, Models: p.Models})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(groups)
		return
	}

	type modelInfo struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Provider string `json:"provider"`
	}
	var out []modelInfo
	for _, p := range llm.BuiltInProviders {
		for _, m := range p.Models {
			out = append(out, modelInfo{ID: m.ID, Name: m.Name, Provider: p.Name})
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (s *Server) handleSettings(w http.ResponseWriter, r *http.Request) {
	cfgPath, _ := zkconfig.DefaultPath()
	cfg, _ := zkconfig.Load(cfgPath)
	if cfg == nil {
		cfg = zkconfig.NewDefault()
	}
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cfg.Secrets)
	case http.MethodPost:
		var body struct {
			Secrets map[string]string `json:"secrets"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
			return
		}
		if cfg.Secrets == nil {
			cfg.Secrets = make(map[string]string)
		}
		for k, v := range body.Secrets {
			cfg.Secrets[k] = v
		}
		if err := zkconfig.Save(cfgPath, cfg); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

type conversation struct {
	ID        string                 `json:"id"`
	AgentID   string                 `json:"agent_id"`
	Title     string                 `json:"title"`
	UpdatedAt string                 `json:"updated_at"`
	Messages  []map[string]interface{} `json:"messages"`
}

func conversationsPath() string {
	cfgPath, _ := zkconfig.DefaultPath()
	return filepath.Join(filepath.Dir(cfgPath), "conversations.jsonl")
}

func (s *Server) handleConversations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		path := conversationsPath()
		f, err := os.Open(path)
		if err != nil {
			if os.IsNotExist(err) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode([]conversation{})
				return
			}
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		var list []conversation
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			var c conversation
			if err := json.Unmarshal(scanner.Bytes(), &c); err == nil {
				list = append(list, c)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var c conversation
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
			return
		}
		b, _ := json.Marshal(c)
		path := conversationsPath()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
			return
		}
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		f.Write(b)
		f.Write([]byte("\n"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}
