package llm

import (
    "bytes"
    "encoding/json"
    "net/http"
    "sync"
    "time"
    "fmt"
)

// OllamaProvider calls the Ollama Cloud API for chat completions.
type OllamaProvider struct {
    apiKey  string
    baseURL string
    model   string
    client  *http.Client
    health  ProviderHealth
    mu      sync.Mutex
}

// NewOllamaProvider creates a new Ollama Cloud provider.
// If model is empty, defaults to gemma3:4b. If apiKey is empty, health will be unhealthy.
func NewOllamaProvider(apiKey, model string) *OllamaProvider {
    if model == "" {
        model = "gemma3:4b"
    }
    p := &OllamaProvider{
        apiKey:  apiKey,
        baseURL: "https://ollama.com",
        model:   model,
        client:  &http.Client{Timeout: 60 * time.Second},
    }

    p.health.Name = "ollama"
    p.health.Latency = 0
    p.health.UpdatedAt = time.Now()
    if apiKey != "" {
        p.health.Status = "healthy"
        p.health.LastError = ""
    } else {
        p.health.Status = "unhealthy"
        p.health.LastError = "no API key"
    }

    return p
}

// Complete sends a chat prompt to Ollama Cloud and returns the assistant content.
func (o *OllamaProvider) Complete(prompt string) (string, error) {
    // Prepare request payload
    payload := map[string]interface{}{
        "model": o.model,
        "messages": []map[string]string{
            {"role": "user", "content": prompt},
        },
        "stream": false,
    }
    body, err := json.Marshal(payload)
    if err != nil {
        o.updateHealthWithError(err)
        return "", err
    }

    req, err := http.NewRequest("POST", o.baseURL+"/api/chat", bytes.NewBuffer(body))
    if err != nil {
        o.updateHealthWithError(err)
        return "", err
    }
    req.Header.Set("Content-Type", "application/json")
    if o.apiKey != "" {
        req.Header.Set("Authorization", "Bearer "+o.apiKey)
    }

    start := time.Now()
    resp, err := o.client.Do(req)
    if err != nil {
        o.updateHealthWithError(err)
        return "", err
    }
    defer resp.Body.Close()

    // Expect 200 OK
    if resp.StatusCode != http.StatusOK {
        o.updateHealthWithError(fmt.Errorf("unexpected status: %d", resp.StatusCode))
        return "", fmt.Errorf("ollama api returned status %d", resp.StatusCode)
    }

    var r ollamaResponse
    if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
        o.updateHealthWithError(err)
        return "", err
    }

    // Update health to healthy on success with latency
    o.mu.Lock()
    o.health.Status = "healthy"
    o.health.LastError = ""
    o.health.Latency = time.Since(start)
    o.health.UpdatedAt = time.Now()
    o.mu.Unlock()

    return r.Message.Content, nil
}

// Health returns the provider health information.
func (o *OllamaProvider) Health() ProviderHealth {
    o.mu.Lock()
    defer o.mu.Unlock()
    return o.health
}

// updateHealthWithError marks the health as unhealthy with the last error.
func (o *OllamaProvider) updateHealthWithError(err error) {
    o.mu.Lock()
    defer o.mu.Unlock()
    o.health.Status = "unhealthy"
    if err != nil {
        o.health.LastError = err.Error()
    } else {
        o.health.LastError = "unknown error"
    }
    o.health.UpdatedAt = time.Now()
}

// Internal response shapes
type ollamaResponse struct {
    Model   string       `json:"model"`
    Message ollamaMessage `json:"message"`
    Done    bool         `json:"done"`
    Error   string       `json:"error,omitempty"`
}

type ollamaMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}
