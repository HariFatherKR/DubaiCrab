package kakao

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"DubaiCrab-Go/internal/ollama"
)

// Config holds Kakao webhook server configuration
type Config struct {
	Enabled      bool     `json:"enabled"`
	Port         int      `json:"port"`
	WebhookPath  string   `json:"webhookPath"`
	DMPolicy     string   `json:"dmPolicy"` // "open", "allowlist", "disabled"
	AllowFrom    []string `json:"allowFrom"`
	SystemPrompt string   `json:"systemPrompt"`
	Model        string   `json:"model"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Enabled:      true,
		Port:         3847,
		WebhookPath:  "/kakao/webhook",
		DMPolicy:     "open",
		AllowFrom:    []string{},
		SystemPrompt: "당신은 Dubai Crab, 한국 사무직을 위한 친절한 AI 비서입니다. 간결하고 도움이 되는 답변을 제공하세요.",
		Model:        "qwen2.5:3b",
	}
}

// Status represents server status
type Status struct {
	Running     bool   `json:"running"`
	Enabled     bool   `json:"enabled"`
	Port        int    `json:"port"`
	WebhookPath string `json:"webhookPath"`
}

// KakaoBot represents bot info
type KakaoBot struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// KakaoUser represents user info
type KakaoUser struct {
	ID         string                 `json:"id,omitempty"`
	Type       string                 `json:"type,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// KakaoBlock represents block info
type KakaoBlock struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// KakaoUserRequest represents user request
type KakaoUserRequest struct {
	Utterance   string      `json:"utterance,omitempty"`
	User        *KakaoUser  `json:"user,omitempty"`
	Params      interface{} `json:"params,omitempty"`
	CallbackURL string      `json:"callbackUrl,omitempty"`
	Lang        string      `json:"lang,omitempty"`
	Timezone    string      `json:"timezone,omitempty"`
	Block       *KakaoBlock `json:"block,omitempty"`
}

// KakaoRequest represents incoming webhook request
type KakaoRequest struct {
	Bot         *KakaoBot         `json:"bot,omitempty"`
	UserRequest *KakaoUserRequest `json:"userRequest,omitempty"`
}

// SimpleText represents simple text output
type SimpleText struct {
	Text string `json:"text"`
}

// KakaoOutput represents output wrapper
type KakaoOutput struct {
	SimpleText SimpleText `json:"simpleText"`
}

// KakaoTemplate represents response template
type KakaoTemplate struct {
	Outputs []KakaoOutput `json:"outputs"`
}

// KakaoData represents callback data
type KakaoData struct {
	Text string `json:"text,omitempty"`
}

// KakaoResponse represents webhook response
type KakaoResponse struct {
	Version     string         `json:"version"`
	UseCallback *bool          `json:"useCallback,omitempty"`
	Template    KakaoTemplate  `json:"template"`
	Data        *KakaoData     `json:"data,omitempty"`
}

// NewSimpleTextResponse creates a simple text response
func NewSimpleTextResponse(text string) *KakaoResponse {
	return &KakaoResponse{
		Version: "2.0",
		Template: KakaoTemplate{
			Outputs: []KakaoOutput{
				{SimpleText: SimpleText{Text: text}},
			},
		},
	}
}

// NewCallbackInitialResponse creates initial callback response
func NewCallbackInitialResponse() *KakaoResponse {
	useCallback := true
	return &KakaoResponse{
		Version:     "2.0",
		UseCallback: &useCallback,
		Template:    KakaoTemplate{Outputs: []KakaoOutput{}},
		Data:        &KakaoData{Text: "처리중입니다..."},
	}
}

// Server handles Kakao webhook requests
type Server struct {
	ollama    *ollama.Manager
	config    *Config
	server    *http.Server
	running   bool
	mu        sync.RWMutex
	cancelFn  context.CancelFunc
}

// NewServer creates a new Kakao webhook server
func NewServer(ollamaManager *ollama.Manager) *Server {
	return &Server{
		ollama: ollamaManager,
		config: DefaultConfig(),
	}
}

// UpdateConfig updates the server configuration
func (s *Server) UpdateConfig(config *Config) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config = config
}

// GetConfig returns current configuration
func (s *Server) GetConfig() *Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config
}

// IsRunning returns whether the server is running
func (s *Server) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// GetStatus returns server status
func (s *Server) GetStatus() Status {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return Status{
		Running:     s.running,
		Enabled:     s.config.Enabled,
		Port:        s.config.Port,
		WebhookPath: s.config.WebhookPath,
	}
}

// Start starts the webhook server
func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return nil
	}

	if !s.config.Enabled {
		log.Println("[kakao] Server disabled in config")
		return nil
	}

	mux := http.NewServeMux()
	mux.HandleFunc(s.config.WebhookPath, s.handleWebhook)
	mux.HandleFunc("/health", s.handleHealth)

	addr := fmt.Sprintf(":%d", s.config.Port)
	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFn = cancel

	go func() {
		log.Printf("[kakao] Starting webhook server on %s%s", addr, s.config.WebhookPath)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[kakao] Server error: %v", err)
		}
		s.mu.Lock()
		s.running = false
		s.mu.Unlock()
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)
	s.running = true
	
	go func() {
		<-ctx.Done()
	}()

	log.Printf("[kakao] Webhook server started on http://0.0.0.0%s%s", addr, s.config.WebhookPath)
	return nil
}

// Stop stops the webhook server
func (s *Server) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running || s.server == nil {
		return
	}

	if s.cancelFn != nil {
		s.cancelFn()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("[kakao] Shutdown error: %v", err)
	}

	s.running = false
	log.Println("[kakao] Webhook server stopped")
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req KakaoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[kakao] Failed to decode request: %v", err)
		sendJSON(w, http.StatusBadRequest, NewSimpleTextResponse("잘못된 요청입니다."))
		return
	}

	if req.UserRequest == nil {
		sendJSON(w, http.StatusBadRequest, NewSimpleTextResponse("잘못된 요청입니다."))
		return
	}

	message := strings.TrimSpace(req.UserRequest.Utterance)
	if message == "" {
		sendJSON(w, http.StatusBadRequest, NewSimpleTextResponse("메시지가 비어있습니다."))
		return
	}

	// Extract sender ID
	senderID := "unknown"
	if req.UserRequest.User != nil {
		if req.UserRequest.User.ID != "" {
			senderID = req.UserRequest.User.ID
		} else if req.UserRequest.User.Properties != nil {
			if botUserKey, ok := req.UserRequest.User.Properties["botUserKey"].(string); ok {
				senderID = botUserKey
			}
		}
	}

	log.Printf("[kakao] Message from %s: %s", senderID, message)

	// Handle callback URL (async processing)
	if req.UserRequest.CallbackURL != "" {
		go s.processAsync(req.UserRequest.CallbackURL, senderID, message)
		sendJSON(w, http.StatusOK, NewCallbackInitialResponse())
		return
	}

	// Synchronous processing
	response := s.processMessage(senderID, message)
	sendJSON(w, http.StatusOK, NewSimpleTextResponse(response))
}

func (s *Server) processMessage(senderID, message string) string {
	config := s.GetConfig()

	// Access control
	if config.DMPolicy == "disabled" {
		return "이 봇은 현재 비활성화되어 있습니다."
	}

	if config.DMPolicy == "allowlist" && len(config.AllowFrom) > 0 {
		allowed := false
		for _, id := range config.AllowFrom {
			if id == "*" || strings.EqualFold(id, senderID) {
				allowed = true
				break
			}
		}
		if !allowed {
			return "접근이 허용되지 않았습니다."
		}
	}

	// Call Ollama
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var systemPrompt *string
	if config.SystemPrompt != "" {
		systemPrompt = &config.SystemPrompt
	}

	response, err := s.ollama.Chat(ctx, config.Model, message, systemPrompt)
	if err != nil {
		log.Printf("[kakao] Ollama error: %v", err)
		return fmt.Sprintf("AI 응답 생성 중 오류가 발생했습니다: %v", err)
	}

	// Truncate if too long (Kakao limit: 1000 chars)
	if len(response) > 900 {
		response = response[:900] + "...\n\n(응답이 너무 길어 일부만 표시)"
	}

	return response
}

func (s *Server) processAsync(callbackURL, senderID, message string) {
	response := s.processMessage(senderID, message)

	payload := NewSimpleTextResponse(response)
	body, _ := json.Marshal(payload)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", callbackURL, strings.NewReader(string(body)))
	if err != nil {
		log.Printf("[kakao] Failed to create callback request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[kakao] Callback failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[kakao] Callback returned status %d", resp.StatusCode)
	}
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
