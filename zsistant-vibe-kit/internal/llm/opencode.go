package llm

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "sync"
    "time"
)

// OpenCodeProvider calls the OpenCode Zen API endpoint.
// Endpoint: https://opencode.ai/zen/v1/chat/completions (OpenAI-compatible)
// Also supports /responses and /messages paths for specific models.
type OpenCodeProvider struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
	health  ProviderHealth
	mu      sync.Mutex
}

// NewOpenCodeProvider creates an OpenCode Zen provider.
// Default model is gpt-5.1. Endpoint defaults to https://opencode.ai/zen/v1.
func NewOpenCodeProvider(apiKey, model string) *OpenCodeProvider {
	if model == "" {
		model = "gpt-5.1"
	}
	p := &OpenCodeProvider{
		apiKey:  apiKey,
		baseURL: "https://opencode.ai/zen/v1",
		model:   model,
		client:  &http.Client{Timeout: 120 * time.Second},
	}
	p.health.Name = "opencode"
	p.health.UpdatedAt = time.Now()
	if apiKey != "" {
		p.health.Status = "healthy"
	} else {
		p.health.Status = "unhealthy"
		p.health.LastError = "no API key"
	}
	return p
}

// Complete sends a chat prompt to OpenCode Zen.
func (o *OpenCodeProvider) Complete(prompt string) (string, error) {
	payload := map[string]interface{}{
		"model": o.model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		o.updateHealthWithError(err)
		return "", err
	}

	req, err := http.NewRequest("POST", o.baseURL+"/chat/completions", bytes.NewBuffer(body))
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

	if resp.StatusCode != http.StatusOK {
		o.updateHealthWithError(fmt.Errorf("status %d", resp.StatusCode))
		return "", fmt.Errorf("opencode api returned status %d", resp.StatusCode)
	}

	var r openAIChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		o.updateHealthWithError(err)
		return "", err
	}
	if len(r.Choices) == 0 {
		o.updateHealthWithError(fmt.Errorf("no choices"))
		return "", fmt.Errorf("no choices in opencode response")
	}

	o.mu.Lock()
	o.health.Status = "healthy"
	o.health.LastError = ""
	o.health.Latency = time.Since(start)
	o.health.UpdatedAt = time.Now()
	o.mu.Unlock()

	return r.Choices[0].Message.Content, nil
}

// Stream implements streaming for the OpenCode Zen provider.
func (o *OpenCodeProvider) Stream(prompt string, chunkCh chan<- string, doneCh chan<- error) {
    payload := map[string]interface{}{
        "model": o.model,
        "messages": []map[string]string{
            {"role": "user", "content": prompt},
        },
        "stream": true,
    }
    body, err := json.Marshal(payload)
    if err != nil {
        o.updateHealthWithError(err)
        doneCh <- err
        return
    }

    req, err := http.NewRequest("POST", o.baseURL+"/chat/completions", bytes.NewBuffer(body))
    if err != nil {
        o.updateHealthWithError(err)
        doneCh <- err
        return
    }
    req.Header.Set("Content-Type", "application/json")
    if o.apiKey != "" {
        req.Header.Set("Authorization", "Bearer "+o.apiKey)
    }

    start := time.Now()
    resp, err := o.client.Do(req)
    if err != nil {
        o.updateHealthWithError(err)
        doneCh <- err
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        e := fmt.Errorf("status %d", resp.StatusCode)
        o.updateHealthWithError(e)
        doneCh <- e
        return
    }

    scanner := bufio.NewScanner(resp.Body)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 {
            continue
        }
        if strings.HasPrefix(line, "data:") {
            payloadLine := strings.TrimSpace(line[len("data:"):])
            if payloadLine == "[DONE]" {
                break
            }
            var chunk openAIStreamChunk
            if err := json.Unmarshal([]byte(payloadLine), &chunk); err != nil {
                continue
            }
            if len(chunk.Choices) > 0 {
                content := chunk.Choices[0].Delta.Content
                if content != "" {
                    chunkCh <- content
                }
            }
        }
    }
    if err := scanner.Err(); err != nil {
        o.updateHealthWithError(err)
        doneCh <- err
        return
    }

    o.mu.Lock()
    o.health.Status = "healthy"
    o.health.LastError = ""
    o.health.Latency = time.Since(start)
    o.health.UpdatedAt = time.Now()
    o.mu.Unlock()

    doneCh <- nil
}

// openAIStreamChunk type is defined in openai.go for reuse across providers.

// Health returns the provider health.
func (o *OpenCodeProvider) Health() ProviderHealth {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.health
}

func (o *OpenCodeProvider) updateHealthWithError(err error) {
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
