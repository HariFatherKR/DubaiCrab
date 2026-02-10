package main

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"DubaiCrab-Go/internal/agent"
	"DubaiCrab-Go/internal/auth"
	"DubaiCrab-Go/internal/config"
	"DubaiCrab-Go/internal/kakao"
	"DubaiCrab-Go/internal/ollama"
	"DubaiCrab-Go/internal/relay"
	"DubaiCrab-Go/internal/tools"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	config       *config.Config
	ollama       *ollama.Manager
	kakao        *kakao.Server
	agent        *agent.Agent
	toolRegistry *tools.Registry
	relay        *relay.Client
	auth         *auth.OAuthManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		cfg = config.DefaultConfig()
	}
	a.config = cfg

	// Initialize Ollama manager
	a.ollama = ollama.NewManager()

	// Initialize tool registry
	a.toolRegistry = tools.NewRegistry()
	tools.RegisterBuiltinTools(a.toolRegistry)
	tools.RegisterOcrTools(a.toolRegistry)

	// Initialize relay client
	a.relay = relay.NewClient(cfg.RelayURL)

	// Initialize OAuth manager
	a.auth = auth.NewOAuthManager()
	// Try to load saved token
	if err := a.auth.LoadSavedToken(); err != nil {
		log.Printf("No saved auth token: %v", err)
	}

	// Initialize agent
	a.agent = agent.NewAgent(a.ollama, a.toolRegistry)

	// Initialize Kakao server
	a.kakao = kakao.NewServer(a.ollama)
	if cfg.KakaoEnabled {
		a.kakao.UpdateConfig(&kakao.Config{
			Enabled:      cfg.KakaoEnabled,
			Port:         cfg.KakaoPort,
			WebhookPath:  cfg.KakaoWebhookPath,
			DMPolicy:     cfg.KakaoDMPolicy,
			AllowFrom:    cfg.KakaoAllowFrom,
			SystemPrompt: cfg.KakaoSystemPrompt,
			Model:        cfg.KakaoModel,
		})
	}

	// Start services in background
	go func() {
		log.Println("Starting Ollama...")
		if err := a.ollama.Start(ctx); err != nil {
			log.Printf("Ollama start error: %v", err)
			wailsRuntime.EventsEmit(ctx, "ollama:error", err.Error())
		} else {
			log.Println("Ollama started successfully")
			wailsRuntime.EventsEmit(ctx, "ollama:ready", true)

			// List models
			models, err := a.ollama.ListModels()
			if err == nil {
				wailsRuntime.EventsEmit(ctx, "ollama:models", models)
			}

			// Start Kakao server
			if cfg.KakaoEnabled {
				log.Println("Starting Kakao webhook server...")
				if err := a.kakao.Start(); err != nil {
					log.Printf("Kakao server error: %v", err)
				}
			}
		}
	}()
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	log.Println("Shutting down...")

	// Stop Kakao server
	if a.kakao != nil {
		a.kakao.Stop()
	}

	// Stop Ollama
	if a.ollama != nil {
		a.ollama.Stop()
	}

	// Save config
	if a.config != nil {
		a.config.Save()
	}
}

// ============================================
// System Commands
// ============================================

// SystemInfo represents system information
type SystemInfo struct {
	OS       string  `json:"os"`
	Arch     string  `json:"arch"`
	MemoryGB float64 `json:"memoryGb"`
}

// GetSystemInfo returns system information
func (a *App) GetSystemInfo() SystemInfo {
	return SystemInfo{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		MemoryGB: 0, // TODO: actual memory
	}
}

// ============================================
// Ollama Commands
// ============================================

// CheckOllama checks if Ollama is running
func (a *App) CheckOllama() bool {
	return a.ollama.IsRunning()
}

// GetOllamaModels returns available Ollama models
func (a *App) GetOllamaModels() ([]string, error) {
	return a.ollama.ListModels()
}

// PullOllamaModel downloads a model
func (a *App) PullOllamaModel(modelName string) error {
	return a.ollama.PullModel(modelName)
}

// StartOllama starts Ollama server
func (a *App) StartOllama() error {
	return a.ollama.Start(a.ctx)
}

// ============================================
// Chat Commands
// ============================================

// ChatMessage represents a chat message for frontend
type ChatMessage struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// Chat sends a message and returns the response
func (a *App) Chat(sessionID, message string) (string, error) {
	if sessionID == "" {
		sessionID = "default"
	}

	response, err := a.agent.ProcessMessage(a.ctx, sessionID, message)
	if err != nil {
		return "", err
	}

	return response, nil
}

// SimpleChat sends a simple chat without session management
func (a *App) SimpleChat(model, message string) (string, error) {
	if model == "" {
		model = a.config.OllamaModel
	}

	systemPrompt := "당신은 Dubai Crab, 한국 사무직을 위한 친절한 AI 비서입니다."
	return a.ollama.Chat(a.ctx, model, message, &systemPrompt)
}

// GetChatHistory returns chat history for a session
func (a *App) GetChatHistory(sessionID string) []ChatMessage {
	if sessionID == "" {
		sessionID = "default"
	}

	messages := a.agent.GetSessionHistory(sessionID)
	result := make([]ChatMessage, len(messages))
	for i, msg := range messages {
		result[i] = ChatMessage{
			Role:      msg.Role,
			Content:   msg.Content,
			Timestamp: msg.Timestamp.UnixMilli(),
		}
	}
	return result
}

// ClearChatHistory clears chat history for a session
func (a *App) ClearChatHistory(sessionID string) {
	if sessionID == "" {
		sessionID = "default"
	}
	a.agent.ClearSession(sessionID)
}

// ============================================
// Kakao Commands
// ============================================

// KakaoStatus represents Kakao server status
type KakaoStatus struct {
	Running     bool   `json:"running"`
	Enabled     bool   `json:"enabled"`
	Port        int    `json:"port"`
	WebhookPath string `json:"webhookPath"`
}

// StartKakaoServer starts the Kakao webhook server
func (a *App) StartKakaoServer() error {
	return a.kakao.Start()
}

// StopKakaoServer stops the Kakao webhook server
func (a *App) StopKakaoServer() error {
	a.kakao.Stop()
	return nil
}

// GetKakaoStatus returns Kakao server status
func (a *App) GetKakaoStatus() KakaoStatus {
	status := a.kakao.GetStatus()
	return KakaoStatus{
		Running:     status.Running,
		Enabled:     status.Enabled,
		Port:        status.Port,
		WebhookPath: status.WebhookPath,
	}
}

// KakaoConfig represents Kakao configuration for frontend
type KakaoConfigJS struct {
	Enabled      bool     `json:"enabled"`
	Port         int      `json:"port"`
	WebhookPath  string   `json:"webhookPath"`
	DMPolicy     string   `json:"dmPolicy"`
	AllowFrom    []string `json:"allowFrom"`
	SystemPrompt string   `json:"systemPrompt"`
	Model        string   `json:"model"`
}

// GetKakaoConfig returns Kakao configuration
func (a *App) GetKakaoConfig() KakaoConfigJS {
	cfg := a.kakao.GetConfig()
	return KakaoConfigJS{
		Enabled:      cfg.Enabled,
		Port:         cfg.Port,
		WebhookPath:  cfg.WebhookPath,
		DMPolicy:     cfg.DMPolicy,
		AllowFrom:    cfg.AllowFrom,
		SystemPrompt: cfg.SystemPrompt,
		Model:        cfg.Model,
	}
}

// UpdateKakaoConfig updates Kakao configuration
func (a *App) UpdateKakaoConfig(cfg KakaoConfigJS) error {
	a.kakao.UpdateConfig(&kakao.Config{
		Enabled:      cfg.Enabled,
		Port:         cfg.Port,
		WebhookPath:  cfg.WebhookPath,
		DMPolicy:     cfg.DMPolicy,
		AllowFrom:    cfg.AllowFrom,
		SystemPrompt: cfg.SystemPrompt,
		Model:        cfg.Model,
	})

	// Update app config
	a.config.KakaoEnabled = cfg.Enabled
	a.config.KakaoPort = cfg.Port
	a.config.KakaoWebhookPath = cfg.WebhookPath
	a.config.KakaoDMPolicy = cfg.DMPolicy
	a.config.KakaoAllowFrom = cfg.AllowFrom
	a.config.KakaoSystemPrompt = cfg.SystemPrompt
	a.config.KakaoModel = cfg.Model

	return a.config.Save()
}

// ============================================
// File Commands
// ============================================

// HWPParseResult represents HWP parse result
type HWPParseResult struct {
	Success bool    `json:"success"`
	Text    *string `json:"text,omitempty"`
	Error   *string `json:"error,omitempty"`
}

// ParseHWP parses an HWP file
func (a *App) ParseHWP(path string, includeTables bool) HWPParseResult {
	tool := &tools.HWPParserTool{}
	result, err := tool.Execute(a.ctx, map[string]interface{}{
		"path":           path,
		"include_tables": includeTables,
	})

	if err != nil {
		errStr := err.Error()
		return HWPParseResult{
			Success: false,
			Error:   &errStr,
		}
	}

	return HWPParseResult{
		Success: true,
		Text:    &result,
	}
}

// ============================================
// Utility Commands
// ============================================

// OpenURL opens a URL in the default browser
func (a *App) OpenURL(url string) error {
	tool := &tools.OpenURLTool{}
	_, err := tool.Execute(a.ctx, map[string]interface{}{"url": url})
	return err
}

// CopyToClipboard copies text to clipboard
func (a *App) CopyToClipboard(text string) error {
	tool := &tools.ClipboardTool{}
	_, err := tool.Execute(a.ctx, map[string]interface{}{
		"action": "copy",
		"text":   text,
	})
	return err
}

// GetToolList returns available tools
func (a *App) GetToolList() []string {
	return a.toolRegistry.List()
}

// ============================================
// Config Commands
// ============================================

// GetConfig returns the current configuration
func (a *App) GetConfig() *config.Config {
	return a.config
}

// SaveConfig saves the configuration
func (a *App) SaveConfig(cfg *config.Config) error {
	a.config = cfg
	return cfg.Save()
}

// SetOllamaModel sets the default Ollama model
func (a *App) SetOllamaModel(model string) error {
	a.config.OllamaModel = model

	// Also update agent
	a.agent.Configure(agent.Config{
		Model: model,
	})

	return a.config.Save()
}

// ============================================
// Development / Debug Commands
// ============================================

// Greet returns a greeting (for testing bindings)
func (a *App) Greet(name string) string {
	return fmt.Sprintf("안녕하세요 %s님! Dubai Crab이 도와드리겠습니다! 🦀", name)
}

// ============================================
// HWP to PDF Commands
// ============================================

// HWPConvertResult represents HWP conversion result
type HWPConvertResult struct {
	Success bool    `json:"success"`
	Message *string `json:"message,omitempty"`
	Error   *string `json:"error,omitempty"`
}

// ConvertHWPToPDF converts HWP file to PDF
func (a *App) ConvertHWPToPDF(inputPath, outputPath string) HWPConvertResult {
	tool := &tools.HWPToPDFTool{}
	result, err := tool.Execute(a.ctx, map[string]interface{}{
		"input_path":  inputPath,
		"output_path": outputPath,
	})

	if err != nil {
		errStr := err.Error()
		return HWPConvertResult{
			Success: false,
			Error:   &errStr,
		}
	}

	return HWPConvertResult{
		Success: true,
		Message: &result,
	}
}

// ============================================
// OCR Commands
// ============================================

// OcrResult represents OCR result
type OcrResult struct {
	Success   bool     `json:"success"`
	Text      *string  `json:"text,omitempty"`
	Error     *string  `json:"error,omitempty"`
	LineCount *int     `json:"line_count,omitempty"`
}

// OcrFromFile performs OCR on an image file
func (a *App) OcrFromFile(imagePath string) OcrResult {
	result := tools.OcrFromFile(imagePath, "kor+eng")
	return OcrResult{
		Success:   result.Success,
		Text:      result.Text,
		Error:     result.Error,
		LineCount: result.LineCount,
	}
}

// OcrFromBase64 performs OCR on a base64-encoded image
func (a *App) OcrFromBase64(base64Data string) OcrResult {
	result := tools.OcrFromBase64(base64Data, "kor+eng")
	return OcrResult{
		Success:   result.Success,
		Text:      result.Text,
		Error:     result.Error,
		LineCount: result.LineCount,
	}
}

// ============================================
// Relay Commands
// ============================================

// RelayStatus represents relay connection status
type RelayStatus struct {
	Connected bool   `json:"connected"`
	Code      string `json:"code"`
	URL       string `json:"url"`
}

// ConnectRelay connects to the relay server
func (a *App) ConnectRelay(code string) (RelayStatus, error) {
	// Set message handler
	a.relay.SetHandler(func(msg relay.RelayMessage) (string, error) {
		// Process message through Ollama
		response, err := a.agent.ProcessMessage(a.ctx, msg.SessionID, msg.Content)
		if err != nil {
			return "", err
		}
		wailsRuntime.EventsEmit(a.ctx, "relay:message", map[string]interface{}{
			"from":    msg.From,
			"content": msg.Content,
			"response": response,
		})
		return response, nil
	})

	if err := a.relay.Connect(code); err != nil {
		return RelayStatus{}, err
	}

	if err := a.relay.Start(); err != nil {
		return RelayStatus{}, err
	}

	status := a.relay.GetStatus()
	return RelayStatus{
		Connected: status.Connected,
		Code:      status.Code,
		URL:       status.URL,
	}, nil
}

// DisconnectRelay disconnects from the relay server
func (a *App) DisconnectRelay() {
	a.relay.Disconnect()
}

// GetRelayStatus returns the relay connection status
func (a *App) GetRelayStatus() RelayStatus {
	status := a.relay.GetStatus()
	return RelayStatus{
		Connected: status.Connected,
		Code:      status.Code,
		URL:       status.URL,
	}
}

// GenerateRelayCode generates a new relay connection code
func (a *App) GenerateRelayCode() string {
	return relay.GenerateCode()
}

// ============================================
// OAuth Commands
// ============================================

// AuthStatus represents authentication status
type AuthStatus struct {
	Authenticated bool       `json:"authenticated"`
	User          *auth.UserInfo `json:"user,omitempty"`
	ExpiresAt     int64      `json:"expires_at,omitempty"`
}

// Login initiates OAuth login
func (a *App) Login() (AuthStatus, error) {
	status, err := a.auth.Login(a.ctx)
	if err != nil {
		return AuthStatus{}, err
	}

	return AuthStatus{
		Authenticated: status.Authenticated,
		User:          status.User,
		ExpiresAt:     status.ExpiresAt,
	}, nil
}

// Logout logs out the user
func (a *App) Logout() error {
	return a.auth.Logout()
}

// GetAuthStatus returns the current auth status
func (a *App) GetAuthStatus() AuthStatus {
	status := a.auth.GetAuthStatus()
	return AuthStatus{
		Authenticated: status.Authenticated,
		User:          status.User,
		ExpiresAt:     status.ExpiresAt,
	}
}
