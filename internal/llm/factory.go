package llm

// ProviderFactory creates real implementations from registry catalog entries.
// This bridges the static catalog (registry.go) to the active provider implementations.

type Factory struct {
	// API keys indexed by lowercased provider name
	Keys map[string]string
}

// NewFactory creates a provider factory with the given API keys.
func NewFactory(keys map[string]string) *Factory {
	return &Factory{Keys: keys}
}

// Create instantiates a Provider for the given model ID by looking up the
// provider in the catalog, resolving the API key, and constructing the
// correct implementation.
func (f *Factory) Create(modelID string) (Provider, *ProviderInfo, *ModelInfo, error) {
	pi, mi, err := FindModel(modelID)
	if err != nil {
		return nil, nil, nil, err
	}
	key := f.Keys[pi.Name]
	if key == "" {
		key = f.Keys[pi.Name] // fallback exact
	}
	prov := f.buildProvider(pi, mi, key)
	return prov, pi, mi, nil
}

// CreateByProvider creates a provider by provider name instead of model ID.
func (f *Factory) CreateByProvider(providerName string) (Provider, *ProviderInfo, error) {
	pi, err := FindProvider(providerName)
	if err != nil {
		return nil, nil, err
	}
	key := f.Keys[providerName]
	if key == "" {
		key = f.Keys[lowerName(providerName)]
	}
	var mi *ModelInfo
	if len(pi.Models) > 0 {
		mi = &pi.Models[0]
	}
	prov := f.buildProvider(pi, mi, key)
	return prov, pi, nil
}

func (f *Factory) buildProvider(pi *ProviderInfo, mi *ModelInfo, key string) Provider {
	model := ""
	if mi != nil {
		model = mi.ID
	}
	if model == "" {
		model = pi.DefaultModel
	}
	switch pi.Name {
	case "ollama-local":
		return NewOllamaProvider(key, model)
	case "vllm-local", "tgi-local", "tabbyml-local", "lmstudio-local", "jan-local":
		return NewOpenAIProvider(key, pi.BaseURL, model)
	default:
		return NewOpenAIProvider(key, pi.BaseURL, model)
	}
}

func lowerName(s string) string {
	out := []byte(s)
	for i, c := range out {
		if 'A' <= c && c <= 'Z' {
			out[i] = c + ('a' - 'A')
		}
	}
	return string(out)
}

// AllProviderHealth returns the health status of all registered providers
// based on whether their API key is configured.
func AllProviderHealth(keys map[string]string) []ProviderHealth {
	var out []ProviderHealth
	for _, pi := range BuiltInProviders {
		key := keys[pi.Name]
		if key == "" {
			key = keys[lowerName(pi.Name)]
		}
		status := "unconfigured"
		if pi.Name == "ollama-local" {
			status = "available"
		} else if key != "" {
			status = "configured"
		}
		out = append(out, ProviderHealth{
			Name:   pi.Name,
			Status: status,
		})
	}
	return out
}
