package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

const (
	// OAuth endpoints
	defaultAuthURL     = "https://dubaicrab.io/oauth/authorize"
	defaultTokenURL    = "https://dubaicrab.io/oauth/token"
	defaultUserInfoURL = "https://dubaicrab.io/api/user/me"
	defaultClientID    = "dubaicrab-desktop"
	defaultRedirectURI = "http://localhost:38470/callback"
	callbackPort       = 38470
)

// TokenData represents OAuth token data
type TokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExpiresAt    int64  `json:"expires_at"`
	Scope        string `json:"scope,omitempty"`
}

// UserInfo represents authenticated user info
type UserInfo struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar,omitempty"`
	Provider string `json:"provider"`
}

// AuthStatus represents authentication status
type AuthStatus struct {
	Authenticated bool      `json:"authenticated"`
	User          *UserInfo `json:"user,omitempty"`
	ExpiresAt     int64     `json:"expires_at,omitempty"`
}

// OAuthManager manages OAuth authentication
type OAuthManager struct {
	mu           sync.RWMutex
	token        *TokenData
	user         *UserInfo
	authURL      string
	tokenURL     string
	userInfoURL  string
	clientID     string
	clientSecret string
	redirectURI  string
	server       *http.Server
	authCodeChan chan string
	errorChan    chan error
}

// NewOAuthManager creates a new OAuth manager
func NewOAuthManager() *OAuthManager {
	return &OAuthManager{
		authURL:      defaultAuthURL,
		tokenURL:     defaultTokenURL,
		userInfoURL:  defaultUserInfoURL,
		clientID:     defaultClientID,
		redirectURI:  defaultRedirectURI,
		authCodeChan: make(chan string, 1),
		errorChan:    make(chan error, 1),
	}
}

// Configure configures the OAuth manager
func (m *OAuthManager) Configure(authURL, tokenURL, userInfoURL, clientID, clientSecret string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if authURL != "" {
		m.authURL = authURL
	}
	if tokenURL != "" {
		m.tokenURL = tokenURL
	}
	if userInfoURL != "" {
		m.userInfoURL = userInfoURL
	}
	if clientID != "" {
		m.clientID = clientID
	}
	if clientSecret != "" {
		m.clientSecret = clientSecret
	}
}

// Login initiates the OAuth login flow
func (m *OAuthManager) Login(ctx context.Context) (*AuthStatus, error) {
	// Start callback server
	if err := m.startCallbackServer(); err != nil {
		return nil, fmt.Errorf("콜백 서버 시작 실패: %v", err)
	}
	defer m.stopCallbackServer()

	// Build authorization URL
	params := url.Values{
		"client_id":     {m.clientID},
		"redirect_uri":  {m.redirectURI},
		"response_type": {"code"},
		"scope":         {"read write"},
	}
	authURL := fmt.Sprintf("%s?%s", m.authURL, params.Encode())

	// Open browser
	if err := openBrowser(authURL); err != nil {
		return nil, fmt.Errorf("브라우저 열기 실패: %v", err)
	}

	// Wait for auth code
	select {
	case code := <-m.authCodeChan:
		// Exchange code for token
		token, err := m.exchangeCode(ctx, code)
		if err != nil {
			return nil, err
		}
		m.mu.Lock()
		m.token = token
		m.mu.Unlock()

		// Get user info
		user, err := m.getUserInfo(ctx)
		if err != nil {
			return nil, err
		}

		m.mu.Lock()
		m.user = user
		m.mu.Unlock()

		// Save token
		if err := m.saveToken(); err != nil {
			// Non-fatal
			fmt.Printf("토큰 저장 실패: %v\n", err)
		}

		return &AuthStatus{
			Authenticated: true,
			User:          user,
			ExpiresAt:     token.ExpiresAt,
		}, nil

	case err := <-m.errorChan:
		return nil, err

	case <-ctx.Done():
		return nil, ctx.Err()

	case <-time.After(5 * time.Minute):
		return nil, fmt.Errorf("로그인 시간 초과")
	}
}

// Logout logs out the user
func (m *OAuthManager) Logout() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.token = nil
	m.user = nil

	// Remove saved token
	tokenPath, err := tokenFilePath()
	if err == nil {
		os.Remove(tokenPath)
	}

	return nil
}

// GetAuthStatus returns the current authentication status
func (m *OAuthManager) GetAuthStatus() AuthStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.token == nil {
		return AuthStatus{Authenticated: false}
	}

	// Check if token is expired
	if m.token.ExpiresAt > 0 && time.Now().Unix() > m.token.ExpiresAt {
		return AuthStatus{Authenticated: false}
	}

	return AuthStatus{
		Authenticated: true,
		User:          m.user,
		ExpiresAt:     m.token.ExpiresAt,
	}
}

// GetAccessToken returns the current access token
func (m *OAuthManager) GetAccessToken() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.token == nil {
		return ""
	}
	return m.token.AccessToken
}

// LoadSavedToken loads token from file
func (m *OAuthManager) LoadSavedToken() error {
	tokenPath, err := tokenFilePath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return err
	}

	var token TokenData
	if err := json.Unmarshal(data, &token); err != nil {
		return err
	}

	// Check if expired
	if token.ExpiresAt > 0 && time.Now().Unix() > token.ExpiresAt {
		os.Remove(tokenPath)
		return fmt.Errorf("토큰이 만료되었습니다")
	}

	m.mu.Lock()
	m.token = &token
	m.mu.Unlock()

	// Try to get user info
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if user, err := m.getUserInfo(ctx); err == nil {
		m.mu.Lock()
		m.user = user
		m.mu.Unlock()
	}

	return nil
}

// saveToken saves token to file
func (m *OAuthManager) saveToken() error {
	m.mu.RLock()
	token := m.token
	m.mu.RUnlock()

	if token == nil {
		return nil
	}

	tokenPath, err := tokenFilePath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(tokenPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	return os.WriteFile(tokenPath, data, 0600)
}

// startCallbackServer starts the OAuth callback server
func (m *OAuthManager) startCallbackServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		errMsg := r.URL.Query().Get("error")

		if errMsg != "" {
			m.errorChan <- fmt.Errorf("OAuth 오류: %s", errMsg)
			w.Write([]byte("<html><body><h1>로그인 실패</h1><p>창을 닫아주세요.</p></body></html>"))
			return
		}

		if code == "" {
			m.errorChan <- fmt.Errorf("인증 코드를 받지 못했습니다")
			w.Write([]byte("<html><body><h1>로그인 실패</h1><p>창을 닫아주세요.</p></body></html>"))
			return
		}

		m.authCodeChan <- code
		w.Write([]byte("<html><body><h1>로그인 성공!</h1><p>창을 닫아주세요.</p><script>window.close()</script></body></html>"))
	})

	m.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", callbackPort),
		Handler: mux,
	}

	go m.server.ListenAndServe()
	return nil
}

// stopCallbackServer stops the callback server
func (m *OAuthManager) stopCallbackServer() {
	if m.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		m.server.Shutdown(ctx)
		m.server = nil
	}
}

// exchangeCode exchanges auth code for token
func (m *OAuthManager) exchangeCode(ctx context.Context, code string) (*TokenData, error) {
	data := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"client_id":    {m.clientID},
		"redirect_uri": {m.redirectURI},
	}

	m.mu.RLock()
	if m.clientSecret != "" {
		data.Set("client_secret", m.clientSecret)
	}
	tokenURL := m.tokenURL
	m.mu.RUnlock()

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("토큰 요청 실패: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("응답 읽기 실패: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("토큰 요청 실패: %s", string(body))
	}

	var token TokenData
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("토큰 파싱 실패: %v", err)
	}

	// Calculate expiration time
	if token.ExpiresIn > 0 {
		token.ExpiresAt = time.Now().Unix() + int64(token.ExpiresIn)
	}

	return &token, nil
}

// getUserInfo gets user info from API
func (m *OAuthManager) getUserInfo(ctx context.Context) (*UserInfo, error) {
	m.mu.RLock()
	token := m.token
	userInfoURL := m.userInfoURL
	m.mu.RUnlock()

	if token == nil {
		return nil, fmt.Errorf("로그인되지 않았습니다")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("사용자 정보 요청 실패: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("응답 읽기 실패: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("사용자 정보 요청 실패: %s", string(body))
	}

	var user UserInfo
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("사용자 정보 파싱 실패: %v", err)
	}

	return &user, nil
}

// tokenFilePath returns the path to the token file
func tokenFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "dubai-crab", "oauth_token.json"), nil
}

// openBrowser opens a URL in the default browser
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("지원되지 않는 운영체제입니다")
	}

	return cmd.Start()
}
