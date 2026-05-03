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

// ZAICodingProvider calls the Z.AI GLM Coding Plan endpoint.
// This is a separate endpoint from the general Z.AI API:
//   General: https://api.z.ai/api/paas/v4
//   Coding:  https://api.z.ai/api/coding/paas/v4
type ZAICodingProvider struct {
    apiKey string
    model  string
    client *http.Client
    health ProviderHealth
    mu     sync.Mutex
}

// NewZAICodingProvider creates a Z.AI Coding Plan provider.
// Default model is GLM-5.1. Other options: GLM-5, GLM-4.7, GLM-4.5-air.
func NewZAICodingProvider(apiKey, model string) *ZAICodingProvider {
    if model == "" {
        model = "GLM-5.1"
    }
    p := &ZAICodingProvider{
        apiKey: apiKey,
        model:  model,
        client: &http.Client{Timeout: 120 * time.Second},
    }
    p.health.Name = "zai-coding"
    p.health.UpdatedAt = time.Now()
    if apiKey != "" {
        p.health.Status = "healthy"
    } else {
        p.health.Status = "unhealthy"
        p.health.LastError = "no API key"
    }
    return p
}

// Complete sends a chat prompt to Z.AI Coding Plan and returns the assistant content.
func (z *ZAICodingProvider) Complete(prompt string) (string, error) {
    payload := map[string]interface{}{
        "model": z.model,
        "messages": []map[string]string{
            {"role": "user", "content": prompt},
        },
    }
    body, err := json.Marshal(payload)
    if err != nil {
        z.updateHealthWithError(err)
        return "", err
    }

    req, err := http.NewRequest("POST", "https://api.z.ai/api/coding/paas/v4/chat/completions", bytes.NewBuffer(body))
    if err != nil {
        z.updateHealthWithError(err)
        return "", err
    }
    req.Header.Set("Content-Type", "application/json")
    if z.apiKey != "" {
        req.Header.Set("Authorization", "Bearer "+z.apiKey)
    }

    start := time.Now()
    resp, err := z.client.Do(req)
    if err != nil {
        z.updateHealthWithError(err)
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        z.updateHealthWithError(fmt.Errorf("status %d", resp.StatusCode))
        return "", fmt.Errorf("zai coding api returned status %d", resp.StatusCode)
    }

    var r openAIChatResponse
    if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
        z.updateHealthWithError(err)
        return "", err
    }
    if len(r.Choices) == 0 {
        z.updateHealthWithError(fmt.Errorf("no choices"))
        return "", fmt.Errorf("no choices in zai response")
    }

    z.mu.Lock()
    z.health.Status = "healthy"
    z.health.LastError = ""
    z.health.Latency = time.Since(start)
    z.health.UpdatedAt = time.Now()
    z.mu.Unlock()

    return r.Choices[0].Message.Content, nil
}

// Stream implements streaming for the Z.AI Coding Plan provider.
// It uses the same OpenAI-like streaming format when available.
func (z *ZAICodingProvider) Stream(prompt string, chunkCh chan<- string, doneCh chan<- error) {
    payload := map[string]interface{}{
        "model": z.model,
        "messages": []map[string]string{
            {"role": "user", "content": prompt},
        },
        "stream": true,
    }
    body, err := json.Marshal(payload)
    if err != nil {
        z.updateHealthWithError(err)
        doneCh <- err
        return
    }

    req, err := http.NewRequest("POST", "https://api.z.ai/api/coding/paas/v4/chat/completions", bytes.NewBuffer(body))
    if err != nil {
        z.updateHealthWithError(err)
        doneCh <- err
        return
    }
    req.Header.Set("Content-Type", "application/json")
    if z.apiKey != "" {
        req.Header.Set("Authorization", "Bearer "+z.apiKey)
    }

    start := time.Now()
    resp, err := z.client.Do(req)
    if err != nil {
        z.updateHealthWithError(err)
        doneCh <- err
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        e := fmt.Errorf("status %d", resp.StatusCode)
        z.updateHealthWithError(e)
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
        z.updateHealthWithError(err)
        doneCh <- err
        return
    }

    z.mu.Lock()
    z.health.Status = "healthy"
    z.health.LastError = ""
    z.health.Latency = time.Since(start)
    z.health.UpdatedAt = time.Now()
    z.mu.Unlock()

    doneCh <- nil
}

// openAIStreamChunk type is defined in openai.go for reuse across providers.

// Health returns the provider health.
func (z *ZAICodingProvider) Health() ProviderHealth {
    z.mu.Lock()
    defer z.mu.Unlock()
    return z.health
}

func (z *ZAICodingProvider) updateHealthWithError(err error) {
    z.mu.Lock()
    defer z.mu.Unlock()
    z.health.Status = "unhealthy"
    if err != nil {
        z.health.LastError = err.Error()
    } else {
        z.health.LastError = "unknown error"
    }
    z.health.UpdatedAt = time.Now()
}
