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

// OpenAIProvider calls any OpenAI-compatible /v1/chat/completions endpoint.
type OpenAIProvider struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
	health  ProviderHealth
	mu      sync.Mutex
}

// NewOpenAIProvider creates a new OpenAI-compatible provider.
// baseURL should include the version path, e.g. "https://api.openai.com/v1".
// If model is empty it defaults to "gpt-4o-mini".
func NewOpenAIProvider(apiKey, baseURL, model string) *OpenAIProvider {
	if model == "" {
		model = "gpt-4o-mini"
	}
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	p := &OpenAIProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		client:  &http.Client{Timeout: 60 * time.Second},
	}
	p.health.Name = "openai-" + baseURL
	if len(p.health.Name) > 30 {
		p.health.Name = p.health.Name[:30]
	}
	p.health.UpdatedAt = time.Now()
	if apiKey != "" {
		p.health.Status = "healthy"
	} else {
		p.health.Status = "unhealthy"
		p.health.LastError = "no API key"
	}
	return p
}

// Complete sends a chat prompt and returns the assistant content.
func (o *OpenAIProvider) Complete(prompt string) (string, error) {
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
		return "", fmt.Errorf("openai api returned status %d", resp.StatusCode)
	}

	var r openAIChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		o.updateHealthWithError(err)
		return "", err
	}
	if len(r.Choices) == 0 {
		o.updateHealthWithError(fmt.Errorf("no choices in response"))
		return "", fmt.Errorf("no choices in response")
	}

	o.mu.Lock()
	o.health.Status = "healthy"
	o.health.LastError = ""
	o.health.Latency = time.Since(start)
	o.health.UpdatedAt = time.Now()
	o.mu.Unlock()

	return r.Choices[0].Message.Content, nil
}

// Stream implements streaming for the OpenAI-compatible provider.
// It sends the request with stream=true and parses SSE data events,
// forwarding content chunks to chunkCh and signaling completion on doneCh.
func (o *OpenAIProvider) Stream(prompt string, chunkCh chan<- string, doneCh chan<- error) {
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
    // Stream content chunks as they arrive
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 {
            continue
        }
        // Expect lines like: data: {"choices":[{"delta":{"content":"..."}}]}
        if strings.HasPrefix(line, "data:") {
            payloadLine := strings.TrimSpace(line[len("data:") :])
            if payloadLine == "[DONE]" {
                break
            }
            var chunk openAIStreamChunk
            if err := json.Unmarshal([]byte(payloadLine), &chunk); err != nil {
                // skip malformed lines
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

    // Update health on success
    o.mu.Lock()
    o.health.Status = "healthy"
    o.health.LastError = ""
    o.health.Latency = time.Since(start)
    o.health.UpdatedAt = time.Now()
    o.mu.Unlock()

    doneCh <- nil
}

// internal streaming struct for OpenAI-like responses
type openAIStreamChunk struct {
    Choices []struct {
        Delta struct {
            Content string `json:"content"`
        } `json:"delta"`
    } `json:"choices"`
}

// Health returns the provider health.
func (o *OpenAIProvider) Health() ProviderHealth {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.health
}

func (o *OpenAIProvider) updateHealthWithError(err error) {
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

type openAIChatResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
