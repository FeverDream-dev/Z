package llm

import (
	"fmt"
	"strings"
)

// ProviderInfo describes a supported LLM provider and its available models.
type ProviderInfo struct {
	Name        string   `json:"name"`
	BaseURL     string   `json:"base_url"`
	DefaultModel string  `json:"default_model"`
	Models      []ModelInfo `json:"models"`
}

// ModelInfo describes a single model from a provider.
type ModelInfo struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ContextWindow int   `json:"context_window"`
	Capabilities []string `json:"capabilities"`
	Tier         string `json:"tier"` // cheap, strong, coding, vision
}

// BuiltInProviders is the static catalog of 50+ providers and their models.
var BuiltInProviders = []ProviderInfo{
	{
		Name: "openai", BaseURL: "https://api.openai.com/v1",
		DefaultModel: "gpt-4o-mini",
		Models: []ModelInfo{
			{ID: "gpt-4o", Name: "GPT-4o", ContextWindow: 128000, Capabilities: []string{"vision", "tools", "json"}, Tier: "strong"},
			{ID: "gpt-4o-mini", Name: "GPT-4o Mini", ContextWindow: 128000, Capabilities: []string{"vision", "tools", "json"}, Tier: "cheap"},
			{ID: "gpt-4-turbo", Name: "GPT-4 Turbo", ContextWindow: 128000, Capabilities: []string{"tools", "json"}, Tier: "strong"},
			{ID: "gpt-4", Name: "GPT-4", ContextWindow: 8192, Capabilities: []string{"json"}, Tier: "strong"},
			{ID: "gpt-3.5-turbo", Name: "GPT-3.5 Turbo", ContextWindow: 16385, Capabilities: []string{"json"}, Tier: "cheap"},
			{ID: "o1-preview", Name: "o1 Preview", ContextWindow: 128000, Capabilities: []string{"reasoning"}, Tier: "strong"},
			{ID: "o1-mini", Name: "o1 Mini", ContextWindow: 128000, Capabilities: []string{"reasoning"}, Tier: "coding"},
			{ID: "text-embedding-3-small", Name: "Embedding Small", ContextWindow: 8191, Capabilities: []string{"embedding"}, Tier: "cheap"},
			{ID: "text-embedding-3-large", Name: "Embedding Large", ContextWindow: 8191, Capabilities: []string{"embedding"}, Tier: "cheap"},
			{ID: "whisper-1", Name: "Whisper", ContextWindow: 0, Capabilities: []string{"audio"}, Tier: "cheap"},
			{ID: "tts-1", Name: "TTS", ContextWindow: 0, Capabilities: []string{"audio"}, Tier: "cheap"},
			{ID: "dall-e-3", Name: "DALL-E 3", ContextWindow: 0, Capabilities: []string{"image"}, Tier: "strong"},
		},
	},
	{
		Name: "anthropic", BaseURL: "https://api.anthropic.com/v1",
		DefaultModel: "claude-3-5-sonnet-20241022",
		Models: []ModelInfo{
			{ID: "claude-3-5-sonnet-20241022", Name: "Claude 3.5 Sonnet", ContextWindow: 200000, Capabilities: []string{"vision", "tools", "json"}, Tier: "strong"},
			{ID: "claude-3-5-haiku-20241022", Name: "Claude 3.5 Haiku", ContextWindow: 200000, Capabilities: []string{"vision", "tools"}, Tier: "cheap"},
			{ID: "claude-3-opus-20240229", Name: "Claude 3 Opus", ContextWindow: 200000, Capabilities: []string{"vision", "tools"}, Tier: "strong"},
			{ID: "claude-3-sonnet-20240229", Name: "Claude 3 Sonnet", ContextWindow: 200000, Capabilities: []string{"vision", "tools"}, Tier: "strong"},
			{ID: "claude-3-haiku-20240307", Name: "Claude 3 Haiku", ContextWindow: 200000, Capabilities: []string{"vision"}, Tier: "cheap"},
		},
	},
	{
		Name: "google", BaseURL: "https://generativelanguage.googleapis.com/v1beta",
		DefaultModel: "gemini-1.5-flash",
		Models: []ModelInfo{
			{ID: "gemini-1.5-pro", Name: "Gemini 1.5 Pro", ContextWindow: 2000000, Capabilities: []string{"vision", "tools", "json"}, Tier: "strong"},
			{ID: "gemini-1.5-flash", Name: "Gemini 1.5 Flash", ContextWindow: 1000000, Capabilities: []string{"vision", "tools"}, Tier: "cheap"},
			{ID: "gemini-1.0-pro", Name: "Gemini 1.0 Pro", ContextWindow: 32000, Capabilities: []string{"vision"}, Tier: "strong"},
			{ID: "gemini-1.0-ultra", Name: "Gemini 1.0 Ultra", ContextWindow: 32000, Capabilities: []string{"vision"}, Tier: "strong"},
			{ID: "gemini-pro-vision", Name: "Gemini Pro Vision", ContextWindow: 16384, Capabilities: []string{"vision"}, Tier: "vision"},
		},
	},
	{
		Name: "cohere", BaseURL: "https://api.cohere.com/v1",
		DefaultModel: "command-r",
		Models: []ModelInfo{
			{ID: "command-r-plus", Name: "Command R+", ContextWindow: 128000, Capabilities: []string{"tools", "json"}, Tier: "strong"},
			{ID: "command-r", Name: "Command R", ContextWindow: 128000, Capabilities: []string{"tools", "json"}, Tier: "cheap"},
			{ID: "command", Name: "Command", ContextWindow: 4096, Capabilities: []string{}, Tier: "cheap"},
			{ID: "command-light", Name: "Command Light", ContextWindow: 4096, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "mistral", BaseURL: "https://api.mistral.ai/v1",
		DefaultModel: "mistral-small-latest",
		Models: []ModelInfo{
			{ID: "mistral-large-latest", Name: "Mistral Large", ContextWindow: 128000, Capabilities: []string{"tools", "json"}, Tier: "strong"},
			{ID: "mistral-medium-latest", Name: "Mistral Medium", ContextWindow: 32000, Capabilities: []string{"json"}, Tier: "strong"},
			{ID: "mistral-small-latest", Name: "Mistral Small", ContextWindow: 32000, Capabilities: []string{"json"}, Tier: "cheap"},
			{ID: "codestral-latest", Name: "Codestral", ContextWindow: 32000, Capabilities: []string{"coding"}, Tier: "coding"},
			{ID: "mixtral-8x22b", Name: "Mixtral 8x22B", ContextWindow: 64000, Capabilities: []string{}, Tier: "strong"},
			{ID: "mixtral-8x7b", Name: "Mixtral 8x7B", ContextWindow: 32000, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "groq", BaseURL: "https://api.groq.com/openai/v1",
		DefaultModel: "llama3-8b-8192",
		Models: []ModelInfo{
			{ID: "llama-3.1-70b-versatile", Name: "Llama 3.1 70B", ContextWindow: 128000, Capabilities: []string{"tools"}, Tier: "strong"},
			{ID: "llama-3.1-8b-instant", Name: "Llama 3.1 8B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "mixtral-8x7b-32768", Name: "Mixtral 8x7B", ContextWindow: 32768, Capabilities: []string{}, Tier: "cheap"},
			{ID: "gemma2-9b-it", Name: "Gemma 2 9B", ContextWindow: 8192, Capabilities: []string{}, Tier: "cheap"},
			{ID: "gemma-7b-it", Name: "Gemma 7B", ContextWindow: 8192, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "together", BaseURL: "https://api.together.xyz/v1",
		DefaultModel: "meta-llama/Llama-3.2-3B-Instruct-Turbo",
		Models: []ModelInfo{
			{ID: "meta-llama/Llama-3.3-70B-Instruct-Turbo", Name: "Llama 3.3 70B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "meta-llama/Llama-3.2-3B-Instruct-Turbo", Name: "Llama 3.2 3B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "mistralai/Mixtral-8x7B-Instruct-v0.1", Name: "Mixtral 8x7B", ContextWindow: 32768, Capabilities: []string{}, Tier: "cheap"},
			{ID: "Qwen/Qwen2.5-72B-Instruct-Turbo", Name: "Qwen2.5 72B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
		},
	},
	{
		Name: "fireworks", BaseURL: "https://api.fireworks.ai/inference/v1",
		DefaultModel: "accounts/fireworks/models/llama-v3p2-3b-instruct",
		Models: []ModelInfo{
			{ID: "accounts/fireworks/models/llama-v3p1-405b-instruct", Name: "Llama 3.1 405B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "accounts/fireworks/models/llama-v3p2-3b-instruct", Name: "Llama 3.2 3B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "accounts/fireworks/models/mixtral-8x22b-instruct", Name: "Mixtral 8x22B", ContextWindow: 64000, Capabilities: []string{}, Tier: "strong"},
			{ID: "accounts/fireworks/models/qwen2p5-72b-instruct", Name: "Qwen2.5 72B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
		},
	},
	{
		Name: "perplexity", BaseURL: "https://api.perplexity.ai",
		DefaultModel: "llama-3.1-sonar-small-128k-online",
		Models: []ModelInfo{
			{ID: "llama-3.1-sonar-huge-128k-online", Name: "Sonar Huge", ContextWindow: 128000, Capabilities: []string{"search"}, Tier: "strong"},
			{ID: "llama-3.1-sonar-large-128k-online", Name: "Sonar Large", ContextWindow: 128000, Capabilities: []string{"search"}, Tier: "strong"},
			{ID: "llama-3.1-sonar-small-128k-online", Name: "Sonar Small", ContextWindow: 128000, Capabilities: []string{"search"}, Tier: "cheap"},
			{ID: "llama-3.1-sonar-pro-128k-online", Name: "Sonar Pro", ContextWindow: 128000, Capabilities: []string{"search"}, Tier: "strong"},
		},
	},
	{
		Name: "deepseek", BaseURL: "https://api.deepseek.com/v1",
		DefaultModel: "deepseek-chat",
		Models: []ModelInfo{
			{ID: "deepseek-chat", Name: "DeepSeek Chat", ContextWindow: 64000, Capabilities: []string{"json"}, Tier: "cheap"},
			{ID: "deepseek-coder", Name: "DeepSeek Coder", ContextWindow: 64000, Capabilities: []string{"coding"}, Tier: "coding"},
			{ID: "deepseek-reasoner", Name: "DeepSeek Reasoner", ContextWindow: 64000, Capabilities: []string{"reasoning"}, Tier: "strong"},
		},
	},
	{
		Name: "xai", BaseURL: "https://api.x.ai/v1",
		DefaultModel: "grok-beta",
		Models: []ModelInfo{
			{ID: "grok-beta", Name: "Grok Beta", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "grok-vision-beta", Name: "Grok Vision", ContextWindow: 128000, Capabilities: []string{"vision"}, Tier: "vision"},
		},
	},
	{
		Name: "azure-openai", BaseURL: "https://{resource}.openai.azure.com/openai/deployments/{deployment}",
		DefaultModel: "gpt-4o",
		Models: []ModelInfo{
			{ID: "gpt-4o", Name: "GPT-4o", ContextWindow: 128000, Capabilities: []string{"vision", "tools"}, Tier: "strong"},
			{ID: "gpt-4o-mini", Name: "GPT-4o Mini", ContextWindow: 128000, Capabilities: []string{"vision", "tools"}, Tier: "cheap"},
			{ID: "gpt-4-turbo", Name: "GPT-4 Turbo", ContextWindow: 128000, Capabilities: []string{"tools"}, Tier: "strong"},
			{ID: "gpt-35-turbo", Name: "GPT-3.5 Turbo", ContextWindow: 16385, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "ollama-local", BaseURL: "http://localhost:11434/v1",
		DefaultModel: "llama3.2",
		Models: []ModelInfo{
			{ID: "llama3.2", Name: "Llama 3.2", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "llama3.1", Name: "Llama 3.1", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "mistral", Name: "Mistral", ContextWindow: 32000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "codellama", Name: "CodeLlama", ContextWindow: 16000, Capabilities: []string{"coding"}, Tier: "coding"},
			{ID: "phi4", Name: "Phi-4", ContextWindow: 16000, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "openrouter", BaseURL: "https://openrouter.ai/api/v1",
		DefaultModel: "openai/gpt-4o-mini",
		Models: []ModelInfo{
			{ID: "openai/gpt-4o", Name: "GPT-4o", ContextWindow: 128000, Capabilities: []string{"vision", "tools"}, Tier: "strong"},
			{ID: "anthropic/claude-3.5-sonnet", Name: "Claude 3.5 Sonnet", ContextWindow: 200000, Capabilities: []string{"vision", "tools"}, Tier: "strong"},
			{ID: "google/gemini-1.5-pro", Name: "Gemini 1.5 Pro", ContextWindow: 2000000, Capabilities: []string{"vision", "tools"}, Tier: "strong"},
			{ID: "meta-llama/llama-3.1-70b-instruct", Name: "Llama 3.1 70B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "nousresearch/hermes-3-llama-3.1-405b", Name: "Hermes 3 405B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
		},
	},
	{
		Name: "ai21", BaseURL: "https://api.ai21.com/studio/v1",
		DefaultModel: "jamba-1.5-mini",
		Models: []ModelInfo{
			{ID: "jamba-1.5-large", Name: "Jamba 1.5 Large", ContextWindow: 256000, Capabilities: []string{"tools"}, Tier: "strong"},
			{ID: "jamba-1.5-mini", Name: "Jamba 1.5 Mini", ContextWindow: 256000, Capabilities: []string{"tools"}, Tier: "cheap"},
			{ID: "jamba-instruct", Name: "Jamba Instruct", ContextWindow: 256000, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "replicate", BaseURL: "https://api.replicate.com/v1",
		DefaultModel: "meta/meta-llama-3-8b-instruct",
		Models: []ModelInfo{
			{ID: "meta/meta-llama-3-70b-instruct", Name: "Llama 3 70B", ContextWindow: 8000, Capabilities: []string{}, Tier: "strong"},
			{ID: "meta/meta-llama-3-8b-instruct", Name: "Llama 3 8B", ContextWindow: 8000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "mistralai/mixtral-8x7b-instruct-v0.1", Name: "Mixtral 8x7B", ContextWindow: 32000, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "novita", BaseURL: "https://api.novita.ai/v3/openai",
		DefaultModel: "meta-llama/llama-3.1-8b-instruct",
		Models: []ModelInfo{
			{ID: "meta-llama/llama-3.1-70b-instruct", Name: "Llama 3.1 70B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "meta-llama/llama-3.1-8b-instruct", Name: "Llama 3.1 8B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "mistralai/mixtral-8x7b-instruct", Name: "Mixtral 8x7B", ContextWindow: 32000, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "hyperbolic", BaseURL: "https://api.hyperbolic.xyz/v1",
		DefaultModel: "meta-llama/Meta-Llama-3.1-8B-Instruct",
		Models: []ModelInfo{
			{ID: "meta-llama/Meta-Llama-3.1-70B-Instruct", Name: "Llama 3.1 70B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "meta-llama/Meta-Llama-3.1-8B-Instruct", Name: "Llama 3.1 8B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "Qwen/Qwen2.5-72B-Instruct", Name: "Qwen2.5 72B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
		},
	},
	{
		Name: "siliconflow", BaseURL: "https://api.siliconflow.cn/v1",
		DefaultModel: "meta-llama/Meta-Llama-3.1-8B-Instruct",
		Models: []ModelInfo{
			{ID: "meta-llama/Meta-Llama-3.1-70B-Instruct", Name: "Llama 3.1 70B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "meta-llama/Meta-Llama-3.1-8B-Instruct", Name: "Llama 3.1 8B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "deepseek-ai/DeepSeek-V2.5", Name: "DeepSeek V2.5", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
		},
	},
	{
		Name: "deepinfra", BaseURL: "https://api.deepinfra.com/v1/openai",
		DefaultModel: "meta-llama/Meta-Llama-3.1-8B-Instruct",
		Models: []ModelInfo{
			{ID: "meta-llama/Meta-Llama-3.1-70B-Instruct", Name: "Llama 3.1 70B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "meta-llama/Meta-Llama-3.1-8B-Instruct", Name: "Llama 3.1 8B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "Qwen/Qwen2.5-72B-Instruct", Name: "Qwen2.5 72B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
		},
	},
	{
		Name: "nvidia", BaseURL: "https://integrate.api.nvidia.com/v1",
		DefaultModel: "meta/llama-3.1-8b-instruct",
		Models: []ModelInfo{
			{ID: "meta/llama-3.1-70b-instruct", Name: "Llama 3.1 70B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "meta/llama-3.1-8b-instruct", Name: "Llama 3.1 8B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "nvidia/nemotron-4-340b-instruct", Name: "Nemotron 4 340B", ContextWindow: 4000, Capabilities: []string{}, Tier: "strong"},
		},
	},
	{
		Name: "sambanova", BaseURL: "https://api.sambanova.ai/v1",
		DefaultModel: "Meta-Llama-3.1-8B-Instruct",
		Models: []ModelInfo{
			{ID: "Meta-Llama-3.1-70B-Instruct", Name: "Llama 3.1 70B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "Meta-Llama-3.1-8B-Instruct", Name: "Llama 3.1 8B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "Meta-Llama-3.2-1B-Instruct", Name: "Llama 3.2 1B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "lambda", BaseURL: "https://api.lambdalabs.com/v1",
		DefaultModel: "hermes-3-llama-3.1-405b-fp8",
		Models: []ModelInfo{
			{ID: "hermes-3-llama-3.1-405b-fp8", Name: "Hermes 3 405B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "llama3.1-70b-instruct-fp8", Name: "Llama 3.1 70B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "llama3.1-8b-instruct-fp8", Name: "Llama 3.1 8B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "friendliai", BaseURL: "https://inference.friendli.ai/v1",
		DefaultModel: "meta-llama-3.1-8b-instruct",
		Models: []ModelInfo{
			{ID: "meta-llama-3.1-70b-instruct", Name: "Llama 3.1 70B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "meta-llama-3.1-8b-instruct", Name: "Llama 3.1 8B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "chutes", BaseURL: "https://chutes-api.ai/v1",
		DefaultModel: "llama-3.1-8b",
		Models: []ModelInfo{
			{ID: "llama-3.1-8b", Name: "Llama 3.1 8B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "llama-3.1-70b", Name: "Llama 3.1 70B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
		},
	},
	{
		Name: "cloudflare", BaseURL: "https://api.cloudflare.com/client/v4/accounts/{account_id}/ai/run",
		DefaultModel: "@cf/meta/llama-3.1-8b-instruct",
		Models: []ModelInfo{
			{ID: "@cf/meta/llama-3.1-8b-instruct", Name: "Llama 3.1 8B", ContextWindow: 32000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "@cf/meta/llama-3.1-70b-instruct", Name: "Llama 3.1 70B", ContextWindow: 32000, Capabilities: []string{}, Tier: "strong"},
		},
	},
	{
		Name: "octoai", BaseURL: "https://text.octoai.run/v1",
		DefaultModel: "meta-llama-3.1-8b-instruct",
		Models: []ModelInfo{
			{ID: "meta-llama-3.1-70b-instruct", Name: "Llama 3.1 70B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "meta-llama-3.1-8b-instruct", Name: "Llama 3.1 8B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "predibase", BaseURL: "https://serving.predibase.com/v1",
		DefaultModel: "llama-3.1-8b",
		Models: []ModelInfo{
			{ID: "llama-3.1-8b", Name: "Llama 3.1 8B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "llama-3.1-70b", Name: "Llama 3.1 70B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
		},
	},
	{
		Name: "poe", BaseURL: "https://api.poe.com/bot",
		DefaultModel: "ChatGPT",
		Models: []ModelInfo{
			{ID: "ChatGPT", Name: "ChatGPT", ContextWindow: 32000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "Claude-3.5-Sonnet", Name: "Claude 3.5 Sonnet", ContextWindow: 200000, Capabilities: []string{}, Tier: "strong"},
			{ID: "GPT-4o", Name: "GPT-4o", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
		},
	},
	{
		Name: "moonshot", BaseURL: "https://api.moonshot.cn/v1",
		DefaultModel: "moonshot-v1-8k",
		Models: []ModelInfo{
			{ID: "moonshot-v1-8k", Name: "Moonshot 8K", ContextWindow: 8192, Capabilities: []string{}, Tier: "cheap"},
			{ID: "moonshot-v1-32k", Name: "Moonshot 32K", ContextWindow: 32000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "moonshot-v1-128k", Name: "Moonshot 128K", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
		},
	},
	{
		Name: "01-ai", BaseURL: "https://api.01.ai/v1",
		DefaultModel: "yi-large",
		Models: []ModelInfo{
			{ID: "yi-large", Name: "Yi Large", ContextWindow: 32000, Capabilities: []string{}, Tier: "strong"},
			{ID: "yi-medium", Name: "Yi Medium", ContextWindow: 16000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "yi-vision", Name: "Yi Vision", ContextWindow: 16000, Capabilities: []string{"vision"}, Tier: "vision"},
		},
	},
	{
		Name: "baidu", BaseURL: "https://qianfan.baidubce.com/v2",
		DefaultModel: "ernie-4.0-turbo-8k",
		Models: []ModelInfo{
			{ID: "ernie-4.0-turbo-8k", Name: "Ernie 4.0 Turbo", ContextWindow: 8192, Capabilities: []string{}, Tier: "strong"},
			{ID: "ernie-speed-128k", Name: "Ernie Speed", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "ernie-lite-8k", Name: "Ernie Lite", ContextWindow: 8192, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "alibaba", BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
		DefaultModel: "qwen-turbo",
		Models: []ModelInfo{
			{ID: "qwen-max", Name: "Qwen Max", ContextWindow: 32000, Capabilities: []string{}, Tier: "strong"},
			{ID: "qwen-plus", Name: "Qwen Plus", ContextWindow: 32000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "qwen-turbo", Name: "Qwen Turbo", ContextWindow: 32000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "qwen-coder-plus", Name: "Qwen Coder", ContextWindow: 32000, Capabilities: []string{"coding"}, Tier: "coding"},
		},
	},
	{
		Name: "tencent", BaseURL: "https://hunyuan.tencentcloudapi.com",
		DefaultModel: "hunyuan-lite",
		Models: []ModelInfo{
			{ID: "hunyuan-pro", Name: "Hunyuan Pro", ContextWindow: 32000, Capabilities: []string{}, Tier: "strong"},
			{ID: "hunyuan-standard", Name: "Hunyuan Standard", ContextWindow: 32000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "hunyuan-lite", Name: "Hunyuan Lite", ContextWindow: 32000, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "minimax", BaseURL: "https://api.minimax.chat/v1",
		DefaultModel: "abab6.5s-chat",
		Models: []ModelInfo{
			{ID: "abab6.5-chat", Name: "ABaB 6.5", ContextWindow: 32000, Capabilities: []string{}, Tier: "strong"},
			{ID: "abab6.5s-chat", Name: "ABaB 6.5s", ContextWindow: 32000, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "stepfun", BaseURL: "https://api.stepfun.com/v1",
		DefaultModel: "step-1-8k",
		Models: []ModelInfo{
			{ID: "step-2-16k", Name: "Step 2", ContextWindow: 16000, Capabilities: []string{}, Tier: "strong"},
			{ID: "step-1-128k", Name: "Step 1", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "step-1-8k", Name: "Step 1 8K", ContextWindow: 8192, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "zhipu", BaseURL: "https://open.bigmodel.cn/api/paas/v4",
		DefaultModel: "glm-4-flash",
		Models: []ModelInfo{
			{ID: "glm-4-plus", Name: "GLM-4 Plus", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "glm-4", Name: "GLM-4", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "glm-4-flash", Name: "GLM-4 Flash", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "glm-4v", Name: "GLM-4V", ContextWindow: 2000, Capabilities: []string{"vision"}, Tier: "vision"},
		},
	},
	{
		Name: "ai360", BaseURL: "https://api.360.cn/v1",
		DefaultModel: "360gpt-turbo",
		Models: []ModelInfo{
			{ID: "360gpt-pro", Name: "360 GPT Pro", ContextWindow: 32000, Capabilities: []string{}, Tier: "strong"},
			{ID: "360gpt-turbo", Name: "360 GPT Turbo", ContextWindow: 32000, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "lingyiwanwu", BaseURL: "https://api.lingyiwanwu.com/v1",
		DefaultModel: "yi-34b-chat-0205",
		Models: []ModelInfo{
			{ID: "yi-34b-chat-0205", Name: "Yi 34B", ContextWindow: 32000, Capabilities: []string{}, Tier: "strong"},
			{ID: "yi-6b-chat", Name: "Yi 6B", ContextWindow: 32000, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "jina", BaseURL: "https://api.jina.ai/v1",
		DefaultModel: "jina-embeddings-v2-base-en",
		Models: []ModelInfo{
			{ID: "jina-embeddings-v2-base-en", Name: "Jina Embeddings", ContextWindow: 8192, Capabilities: []string{"embedding"}, Tier: "cheap"},
			{ID: "jina-embeddings-v3", Name: "Jina Embeddings V3", ContextWindow: 8192, Capabilities: []string{"embedding"}, Tier: "cheap"},
		},
	},
	{
		Name: "voyage", BaseURL: "https://api.voyageai.com/v1",
		DefaultModel: "voyage-3-lite",
		Models: []ModelInfo{
			{ID: "voyage-3", Name: "Voyage 3", ContextWindow: 32000, Capabilities: []string{"embedding"}, Tier: "cheap"},
			{ID: "voyage-3-lite", Name: "Voyage 3 Lite", ContextWindow: 32000, Capabilities: []string{"embedding"}, Tier: "cheap"},
			{ID: "voyage-code-3", Name: "Voyage Code 3", ContextWindow: 16000, Capabilities: []string{"embedding", "coding"}, Tier: "cheap"},
		},
	},
	{
		Name: "qdrant", BaseURL: "https://qdrant.io/api",
		DefaultModel: "dense",
		Models: []ModelInfo{
			{ID: "dense", Name: "Dense Embeddings", ContextWindow: 0, Capabilities: []string{"embedding"}, Tier: "cheap"},
		},
	},
	{
		Name: "anyscale", BaseURL: "https://api.endpoints.anyscale.com/v1",
		DefaultModel: "meta-llama/Llama-2-7b-chat-hf",
		Models: []ModelInfo{
			{ID: "meta-llama/Llama-2-7b-chat-hf", Name: "Llama 2 7B", ContextWindow: 4096, Capabilities: []string{}, Tier: "cheap"},
			{ID: "meta-llama/Llama-2-13b-chat-hf", Name: "Llama 2 13B", ContextWindow: 4096, Capabilities: []string{}, Tier: "cheap"},
			{ID: "meta-llama/Llama-2-70b-chat-hf", Name: "Llama 2 70B", ContextWindow: 4096, Capabilities: []string{}, Tier: "strong"},
		},
	},
	{
		Name: "monsterapi", BaseURL: "https://api.monsterapi.ai/v1",
		DefaultModel: "llama3-8b",
		Models: []ModelInfo{
			{ID: "llama3-8b", Name: "Llama 3 8B", ContextWindow: 8000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "llama3-70b", Name: "Llama 3 70B", ContextWindow: 8000, Capabilities: []string{}, Tier: "strong"},
		},
	},
	{
		Name: "inference.net", BaseURL: "https://api.inference.net/v1",
		DefaultModel: "llama3.1-8b",
		Models: []ModelInfo{
			{ID: "llama3.1-8b", Name: "Llama 3.1 8B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "llama3.1-70b", Name: "Llama 3.1 70B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
		},
	},
	{
		Name: "ppio", BaseURL: "https://api.ppio.ai/v1",
		DefaultModel: "meta-llama/llama-3.1-8b-instruct",
		Models: []ModelInfo{
			{ID: "meta-llama/llama-3.1-8b-instruct", Name: "Llama 3.1 8B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "meta-llama/llama-3.1-70b-instruct", Name: "Llama 3.1 70B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "deepseek-ai/deepseek-llm-67b-chat", Name: "DeepSeek 67B", ContextWindow: 32000, Capabilities: []string{}, Tier: "strong"},
		},
	},
	{
		Name: "nebius", BaseURL: "https://api.studio.nebius.ai/v1",
		DefaultModel: "meta-llama/Meta-Llama-3.1-8B-Instruct",
		Models: []ModelInfo{
			{ID: "meta-llama/Meta-Llama-3.1-70B-Instruct", Name: "Llama 3.1 70B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "meta-llama/Meta-Llama-3.1-8B-Instruct", Name: "Llama 3.1 8B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "cerebras", BaseURL: "https://api.cerebras.ai/v1",
		DefaultModel: "llama3.1-8b",
		Models: []ModelInfo{
			{ID: "llama3.1-8b", Name: "Llama 3.1 8B", ContextWindow: 8192, Capabilities: []string{}, Tier: "cheap"},
			{ID: "llama3.1-70b", Name: "Llama 3.1 70B", ContextWindow: 8192, Capabilities: []string{}, Tier: "strong"},
		},
	},
	{
		Name: "klusterai", BaseURL: "https://api.kluster.ai/v1",
		DefaultModel: "klusterai/Meta-Llama-3.1-8B-Instruct-Turbo",
		Models: []ModelInfo{
			{ID: "klusterai/Meta-Llama-3.1-8B-Instruct-Turbo", Name: "Llama 3.1 8B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "klusterai/Meta-Llama-3.1-70B-Instruct-Turbo", Name: "Llama 3.1 70B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
		},
	},
	{
		Name: "huggingface", BaseURL: "https://api-inference.huggingface.co/v1",
		DefaultModel: "meta-llama/Meta-Llama-3.1-8B-Instruct",
		Models: []ModelInfo{
			{ID: "meta-llama/Meta-Llama-3.1-70B-Instruct", Name: "Llama 3.1 70B", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "meta-llama/Meta-Llama-3.1-8B-Instruct", Name: "Llama 3.1 8B", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "mistralai/Mixtral-8x7B-Instruct-v0.1", Name: "Mixtral 8x7B", ContextWindow: 32000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "microsoft/Phi-3-mini-4k-instruct", Name: "Phi-3 Mini", ContextWindow: 4000, Capabilities: []string{}, Tier: "cheap"},
		},
	},
	{
		Name: "zai-coding", BaseURL: "https://api.z.ai/api/coding/paas/v4",
		DefaultModel: "GLM-5.1",
		Models: []ModelInfo{
			{ID: "GLM-5.1", Name: "GLM 5.1", ContextWindow: 128000, Capabilities: []string{"coding", "tools"}, Tier: "coding"},
			{ID: "GLM-5", Name: "GLM 5", ContextWindow: 128000, Capabilities: []string{"coding", "tools"}, Tier: "coding"},
			{ID: "GLM-4.7", Name: "GLM 4.7", ContextWindow: 128000, Capabilities: []string{"coding", "thinking"}, Tier: "coding"},
			{ID: "GLM-4.5-air", Name: "GLM 4.5 Air", ContextWindow: 64000, Capabilities: []string{"coding"}, Tier: "coding"},
		},
	},
	{
		Name: "opencode", BaseURL: "https://opencode.ai/zen/v1",
		DefaultModel: "gpt-5.1",
		Models: []ModelInfo{
			{ID: "gpt-5.5", Name: "GPT 5.5", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "gpt-5.5-pro", Name: "GPT 5.5 Pro", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "gpt-5.4", Name: "GPT 5.4", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "gpt-5.4-pro", Name: "GPT 5.4 Pro", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "gpt-5.4-mini", Name: "GPT 5.4 Mini", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "gpt-5.4-nano", Name: "GPT 5.4 Nano", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "gpt-5.3-codex", Name: "GPT 5.3 Codex", ContextWindow: 128000, Capabilities: []string{"coding"}, Tier: "coding"},
			{ID: "gpt-5.3-codex-spark", Name: "GPT 5.3 Codex Spark", ContextWindow: 128000, Capabilities: []string{"coding"}, Tier: "coding"},
			{ID: "gpt-5.2", Name: "GPT 5.2", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "gpt-5.2-codex", Name: "GPT 5.2 Codex", ContextWindow: 128000, Capabilities: []string{"coding"}, Tier: "coding"},
			{ID: "gpt-5.1", Name: "GPT 5.1", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "gpt-5.1-codex", Name: "GPT 5.1 Codex", ContextWindow: 128000, Capabilities: []string{"coding"}, Tier: "coding"},
			{ID: "gpt-5.1-codex-max", Name: "GPT 5.1 Codex Max", ContextWindow: 128000, Capabilities: []string{"coding"}, Tier: "coding"},
			{ID: "gpt-5.1-codex-mini", Name: "GPT 5.1 Codex Mini", ContextWindow: 128000, Capabilities: []string{"coding"}, Tier: "coding"},
			{ID: "gpt-5", Name: "GPT 5", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "gpt-5-codex", Name: "GPT 5 Codex", ContextWindow: 128000, Capabilities: []string{"coding"}, Tier: "coding"},
			{ID: "gpt-5-nano", Name: "GPT 5 Nano", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "claude-opus-4-7", Name: "Claude Opus 4.7", ContextWindow: 200000, Capabilities: []string{}, Tier: "strong"},
			{ID: "claude-opus-4-6", Name: "Claude Opus 4.6", ContextWindow: 200000, Capabilities: []string{}, Tier: "strong"},
			{ID: "claude-sonnet-4-5", Name: "Claude Sonnet 4.5", ContextWindow: 200000, Capabilities: []string{}, Tier: "strong"},
			{ID: "claude-haiku-4-5", Name: "Claude Haiku 4.5", ContextWindow: 200000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "gemini-3.1-pro", Name: "Gemini 3.1 Pro", ContextWindow: 2000000, Capabilities: []string{"vision"}, Tier: "strong"},
			{ID: "gemini-3-flash", Name: "Gemini 3 Flash", ContextWindow: 1000000, Capabilities: []string{"vision"}, Tier: "cheap"},
			{ID: "glm-5.1", Name: "GLM 5.1", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "glm-5", Name: "GLM 5", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "glm-4.7", Name: "GLM 4.7", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "glm-4.7-free", Name: "GLM 4.7 Free", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "glm-4.6", Name: "GLM 4.6", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "kimi-k2.5", Name: "Kimi K2.5", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "kimi-k2.6", Name: "Kimi K2.6", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "kimi-k2-thinking", Name: "Kimi K2 Thinking", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "kimi-k2", Name: "Kimi K2", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "kimi-k2.5-free", Name: "Kimi K2.5 Free", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "qwen3-coder", Name: "Qwen3 Coder 480B", ContextWindow: 128000, Capabilities: []string{"coding"}, Tier: "coding"},
			{ID: "big-pickle", Name: "Big Pickle", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "ling-2.6-flash", Name: "Ling 2.6 Flash", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "hy3-preview-free", Name: "Hy3 Preview Free", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "nemotron-3-super-free", Name: "Nemotron 3 Super Free", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
			{ID: "minimax-m2.1", Name: "MiniMax M2.1", ContextWindow: 128000, Capabilities: []string{}, Tier: "strong"},
			{ID: "minimax-m2.1-free", Name: "MiniMax M2.1 Free", ContextWindow: 128000, Capabilities: []string{}, Tier: "cheap"},
		},
	},
}

// TotalProviderCount returns the number of providers in the catalog.
func TotalProviderCount() int {
	return len(BuiltInProviders)
}

// TotalModelCount returns the total number of models across all providers.
func TotalModelCount() int {
	count := 0
	for _, p := range BuiltInProviders {
		count += len(p.Models)
	}
	return count
}

// FindProvider looks up a provider by name.
func FindProvider(name string) (*ProviderInfo, error) {
	for _, p := range BuiltInProviders {
		if strings.EqualFold(p.Name, name) {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("provider %q not found in catalog", name)
}

// FindModel looks up a model across all providers.
func FindModel(modelID string) (*ProviderInfo, *ModelInfo, error) {
	for _, p := range BuiltInProviders {
		for i := range p.Models {
			if strings.EqualFold(p.Models[i].ID, modelID) {
				return &p, &p.Models[i], nil
			}
		}
	}
	return nil, nil, fmt.Errorf("model %q not found in catalog", modelID)
}

// ProvidersByTier returns all providers that have at least one model in the given tier.
func ProvidersByTier(tier string) []ProviderInfo {
	var out []ProviderInfo
	for _, p := range BuiltInProviders {
		for _, m := range p.Models {
			if strings.EqualFold(m.Tier, tier) {
				out = append(out, p)
				break
			}
		}
	}
	return out
}
