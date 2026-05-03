package llm

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// Streamer is an optional interface for providers that support streaming.
type Streamer interface {
    Provider
    // Stream sends prompt chunks to chunkCh and signals completion (or error) on doneCh.
    Stream(prompt string, chunkCh chan<- string, doneCh chan<- error)
}

// ollamaStreamResponse is the per-line JSON shape from Ollama when streaming.
type ollamaStreamResponse struct {
    Model   string        `json:"model"`
    Message ollamaMessage `json:"message"`
    Done    bool          `json:"done"`
    Error   string        `json:"error,omitempty"`
}

// Stream implements Streamer for OllamaProvider.
func (o *OllamaProvider) Stream(prompt string, chunkCh chan<- string, doneCh chan<- error) {
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

    req, err := http.NewRequest("POST", o.baseURL+"/api/chat", bytes.NewBuffer(body))
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
    var fullContent string
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 {
            continue
        }
        var delta ollamaStreamResponse
        if err := json.Unmarshal([]byte(line), &delta); err != nil {
            continue
        }
        if delta.Error != "" {
            doneCh <- fmt.Errorf("ollama stream error: %s", delta.Error)
            return
        }
        if delta.Done {
            break
        }
        if delta.Message.Content != "" {
            fullContent += delta.Message.Content
            chunkCh <- delta.Message.Content
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

    _ = fullContent
    doneCh <- nil
}

// openAIStreamChunk type is defined in openai.go for reuse across providers.
