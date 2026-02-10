package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	defaultBaseURL       = "http://localhost:11434"
	startupTimeout       = 30 * time.Second
	healthCheckTimeout   = 2 * time.Second
	chatTimeout          = 120 * time.Second
	listModelsTimeout    = 10 * time.Second
)

// Manager handles Ollama process lifecycle and API calls
type Manager struct {
	baseURL string
	process *exec.Cmd
	mu      sync.Mutex
	client  *http.Client
}

// NewManager creates a new Ollama manager
func NewManager() *Manager {
	return &Manager{
		baseURL: defaultBaseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetBaseURL allows overriding the default Ollama URL
func (m *Manager) SetBaseURL(url string) {
	m.baseURL = url
}

// Start starts the Ollama server
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if already running
	if m.IsRunning() {
		return nil
	}

	// Find Ollama binary
	ollamaPath := m.findOllamaPath()
	if ollamaPath == "" {
		return fmt.Errorf("ollama not found. Please install Ollama from https://ollama.ai")
	}

	// Start Ollama serve
	m.process = exec.CommandContext(ctx, ollamaPath, "serve")
	m.process.Stdout = io.Discard
	m.process.Stderr = io.Discard

	if err := m.process.Start(); err != nil {
		return fmt.Errorf("failed to start Ollama: %w", err)
	}

	// Wait for server to be ready
	return m.waitForReady(ctx, startupTimeout)
}

// Stop stops the Ollama server
func (m *Manager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.process == nil {
		return
	}

	// Try graceful shutdown first
	if runtime.GOOS != "windows" {
		m.process.Process.Signal(os.Interrupt)
	} else {
		m.process.Process.Kill()
	}

	// Wait with timeout
	done := make(chan error, 1)
	go func() {
		done <- m.process.Wait()
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		m.process.Process.Kill()
	}

	m.process = nil
}

// IsRunning checks if Ollama server is running
func (m *Manager) IsRunning() bool {
	ctx, cancel := context.WithTimeout(context.Background(), healthCheckTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", m.baseURL+"/api/tags", nil)
	if err != nil {
		return false
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// ListModels returns available models
func (m *Manager) ListModels() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), listModelsTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", m.baseURL+"/api/tags", nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	models := make([]string, len(result.Models))
	for i, m := range result.Models {
		models[i] = m.Name
	}

	return models, nil
}

// PullModel downloads a model
func (m *Manager) PullModel(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	body, _ := json.Marshal(map[string]string{"name": name})
	req, err := http.NewRequestWithContext(ctx, "POST", m.baseURL+"/api/pull", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to pull model: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to pull model: HTTP %d", resp.StatusCode)
	}

	return nil
}

// ChatMessage represents a chat message
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents a chat request
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

// ChatResponse represents a chat response
type ChatResponse struct {
	Message ChatMessage `json:"message"`
}

// Chat sends a chat message and returns the response
func (m *Manager) Chat(ctx context.Context, model, message string, systemPrompt *string) (string, error) {
	messages := []ChatMessage{}

	if systemPrompt != nil && *systemPrompt != "" {
		messages = append(messages, ChatMessage{
			Role:    "system",
			Content: *systemPrompt,
		})
	}

	messages = append(messages, ChatMessage{
		Role:    "user",
		Content: message,
	})

	request := ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(ctx, chatTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", m.baseURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("chat request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("chat failed: HTTP %d", resp.StatusCode)
	}

	var result ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Message.Content, nil
}

// findOllamaPath finds the Ollama binary path
func (m *Manager) findOllamaPath() string {
	// 1. Check PATH
	if path, err := exec.LookPath("ollama"); err == nil {
		return path
	}

	// 2. Platform-specific default paths
	var paths []string
	switch runtime.GOOS {
	case "darwin":
		paths = []string{
			"/usr/local/bin/ollama",
			"/opt/homebrew/bin/ollama",
		}
	case "linux":
		paths = []string{
			"/usr/local/bin/ollama",
			"/usr/bin/ollama",
		}
	case "windows":
		localApp := os.Getenv("LOCALAPPDATA")
		if localApp != "" {
			paths = append(paths, filepath.Join(localApp, "Programs", "Ollama", "ollama.exe"))
		}
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return ""
}

// waitForReady waits for the Ollama server to be ready
func (m *Manager) waitForReady(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for Ollama to start")
		case <-ticker.C:
			if m.IsRunning() {
				return nil
			}
		}
	}
}

// Generate generates text (for simple prompts without chat)
func (m *Manager) Generate(ctx context.Context, model, prompt string) (string, error) {
	request := map[string]interface{}{
		"model":  model,
		"prompt": prompt,
		"stream": false,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(ctx, chatTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", m.baseURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("generate request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("generate failed: HTTP %d", resp.StatusCode)
	}

	var result struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return strings.TrimSpace(result.Response), nil
}
