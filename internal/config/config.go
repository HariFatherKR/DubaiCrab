package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds all application configuration
type Config struct {
	// General settings
	AppName string `json:"appName"`
	Version string `json:"version"`

	// Ollama settings
	OllamaURL   string `json:"ollamaUrl"`
	OllamaModel string `json:"ollamaModel"`

	// Kakao settings
	KakaoEnabled     bool     `json:"kakaoEnabled"`
	KakaoPort        int      `json:"kakaoPort"`
	KakaoWebhookPath string   `json:"kakaoWebhookPath"`
	KakaoDMPolicy    string   `json:"kakaoDmPolicy"`
	KakaoAllowFrom   []string `json:"kakaoAllowFrom"`
	KakaoSystemPrompt string  `json:"kakaoSystemPrompt"`
	KakaoModel       string   `json:"kakaoModel"`

	// Relay settings
	RelayURL   string `json:"relayUrl"`
	RelayToken string `json:"relayToken"`

	// Auth settings
	AuthProvider string `json:"authProvider"`
	AccessToken  string `json:"accessToken"`
}

// DefaultConfig returns a new Config with default values
func DefaultConfig() *Config {
	return &Config{
		AppName:          "Dubai Crab",
		Version:          "0.2.0",
		OllamaURL:        "http://localhost:11434",
		OllamaModel:      "qwen2.5:0.5b",
		KakaoEnabled:     true,
		KakaoPort:        3847,
		KakaoWebhookPath: "/kakao/webhook",
		KakaoDMPolicy:    "open",
		KakaoAllowFrom:   []string{},
		KakaoSystemPrompt: "당신은 Dubai Crab, 한국 사무직을 위한 친절한 AI 비서입니다. 간결하고 도움이 되는 답변을 제공하세요.",
		KakaoModel:       "qwen2.5:0.5b",
		RelayURL:         "wss://relay.dubaicrab.io",
	}
}

// Load loads configuration from the default config file
func Load() (*Config, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}

	config := DefaultConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// Save saves configuration to the default config file
func (c *Config) Save() error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// ConfigPath returns the path to the config file
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "dubai-crab", "config.json"), nil
}
