package bus

import (
	"encoding/json"
	"fmt"
	"time"
)

// RequestType categorizes inter-agent requests.
type RequestType string

const (
	ReqSummary   RequestType = "summary"
	ReqFileRead  RequestType = "file_read"
	ReqFileWrite RequestType = "file_write"
	ReqExec      RequestType = "exec"
)

// Envelope is a typed request between agents.
type Envelope struct {
	ID        string      `json:"id"`
	From      string      `json:"from"`
	To        string      `json:"to"`
	Type      RequestType `json:"type"`
	Payload   string      `json:"payload"`
	CreatedAt time.Time   `json:"created_at"`
}

// Response is the result of an envelope request.
type Response struct {
	RequestID string    `json:"request_id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Allowed   bool      `json:"allowed"`
	Result    string    `json:"result,omitempty"`
	Error     string    `json:"error,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// NewEnvelope creates a new request envelope.
func NewEnvelope(from, to string, reqType RequestType, payload string) *Envelope {
	return &Envelope{
		ID:        fmt.Sprintf("req-%d", time.Now().UnixNano()),
		From:      from,
		To:        to,
		Type:      reqType,
		Payload:   payload,
		CreatedAt: time.Now(),
	}
}

// NewResponse creates a response to an envelope.
func NewResponse(req *Envelope, allowed bool, result, err string) *Response {
	return &Response{
		RequestID: req.ID,
		From:      req.To,
		To:        req.From,
		Allowed:   allowed,
		Result:    result,
		Error:     err,
		CreatedAt: time.Now(),
	}
}

// String returns a JSON representation.
func (e *Envelope) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func (r *Response) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}
