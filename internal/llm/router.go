package llm

import (
    "context"
    "fmt"
    "time"
)

// Router routes prompts to providers with timeout and fallback support.
type Router struct {
    providers map[TaskType][]Provider // ordered by priority
    timeout   time.Duration
}

// NewRouter creates a router with the given per-request timeout.
func NewRouter(timeout time.Duration) *Router {
    return &Router{
        providers: make(map[TaskType][]Provider),
        timeout:   timeout,
    }
}

// Register adds a provider for a task type. Providers are tried in registration order.
func (r *Router) Register(taskType TaskType, p Provider) {
    r.providers[taskType] = append(r.providers[taskType], p)
}

// Route sends the prompt to the first available provider for the task type.
// If the primary fails or times out, it tries fallback providers.
// Returns the response, the name of the successful provider, and any error.
func (r *Router) Route(taskType TaskType, prompt string) (response, providerName string, err error) {
    chain, ok := r.providers[taskType]
    if !ok || len(chain) == 0 {
        return "", "", fmt.Errorf("no providers registered for task type %s", taskType)
    }

    for i, p := range chain {
        ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
        resultCh := make(chan struct {
            resp string
            err  error
        }, 1)

        go func(prov Provider) {
            resp, err := prov.Complete(prompt)
            resultCh <- struct {
                resp string
                err  error
            }{resp, err}
        }(p)

        select {
        case <-ctx.Done():
            cancel()
            if i < len(chain)-1 {
                continue // fallback
            }
            return "", "", fmt.Errorf("all providers timed out for task type %s", taskType)
        case res := <-resultCh:
            cancel()
            if res.err == nil {
                return res.resp, p.Health().Name, nil
            }
            if i < len(chain)-1 {
                continue // fallback
            }
            return "", "", fmt.Errorf("all providers failed for task type %s: last error: %w", taskType, res.err)
        }
    }

    return "", "", fmt.Errorf("all providers failed for task type %s", taskType)
}

// ChainFor returns the ordered provider chain for a task type (nil if none).
func (r *Router) ChainFor(taskType TaskType) []Provider {
    return r.providers[taskType]
}

// Health returns the health of all registered providers.
func (r *Router) Health() []ProviderHealth {
    var health []ProviderHealth
    seen := make(map[string]bool)
    for _, chain := range r.providers {
        for _, p := range chain {
            h := p.Health()
            if !seen[h.Name] {
                seen[h.Name] = true
                health = append(health, h)
            }
        }
    }
    return health
}
