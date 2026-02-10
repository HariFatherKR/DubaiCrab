package relay

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// MessageHandler handles incoming messages from relay
type MessageHandler func(msg RelayMessage) (string, error)

// Client represents a WebSocket client for the relay server
type Client struct {
	url            string
	code           string
	conn           *websocket.Conn
	mu             sync.RWMutex
	connected      bool
	reconnecting   bool
	handler        MessageHandler
	ctx            context.Context
	cancel         context.CancelFunc
	reconnectDelay time.Duration
	maxReconnect   int
	reconnectCount int
}

// RelayMessage represents a message from the relay server
type RelayMessage struct {
	Type      string `json:"type"`
	From      string `json:"from"`
	Content   string `json:"content"`
	SessionID string `json:"session_id,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// RelayResponse represents a response to send to the relay server
type RelayResponse struct {
	Type      string `json:"type"`
	To        string `json:"to"`
	Content   string `json:"content"`
	SessionID string `json:"session_id,omitempty"`
	Code      string `json:"code"`
}

// NewClient creates a new relay client
func NewClient(url string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		url:            url,
		reconnectDelay: 5 * time.Second,
		maxReconnect:   -1, // infinite
		ctx:            ctx,
		cancel:         cancel,
	}
}

// GenerateCode generates a 6-digit connection code
func GenerateCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}

// Connect connects to the relay server with a connection code
func (c *Client) Connect(code string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if code == "" {
		code = GenerateCode()
	}
	c.code = code

	conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return fmt.Errorf("연결 실패: %v", err)
	}

	c.conn = conn
	c.connected = true
	c.reconnectCount = 0

	// Send registration message
	regMsg := map[string]string{
		"type": "register",
		"code": c.code,
	}
	if err := conn.WriteJSON(regMsg); err != nil {
		conn.Close()
		c.connected = false
		return fmt.Errorf("등록 메시지 전송 실패: %v", err)
	}

	log.Printf("릴레이 서버 연결됨 (코드: %s)", c.code)
	return nil
}

// Disconnect disconnects from the relay server
func (c *Client) Disconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cancel()

	if c.conn != nil {
		c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		c.conn.Close()
		c.conn = nil
	}
	c.connected = false
	log.Println("릴레이 서버 연결 해제됨")
}

// SetHandler sets the message handler
func (c *Client) SetHandler(handler MessageHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handler = handler
}

// Start starts listening for messages
func (c *Client) Start() error {
	c.mu.RLock()
	if !c.connected || c.conn == nil {
		c.mu.RUnlock()
		return fmt.Errorf("연결되지 않았습니다")
	}
	c.mu.RUnlock()

	go c.readLoop()
	return nil
}

// readLoop reads messages from the WebSocket
func (c *Client) readLoop() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		c.mu.RLock()
		conn := c.conn
		c.mu.RUnlock()

		if conn == nil {
			c.attemptReconnect()
			continue
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("메시지 수신 오류: %v", err)
			c.mu.Lock()
			c.connected = false
			c.mu.Unlock()
			c.attemptReconnect()
			continue
		}

		var msg RelayMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("메시지 파싱 오류: %v", err)
			continue
		}

		// Skip non-message types
		if msg.Type != "message" {
			if msg.Type == "error" {
				log.Printf("릴레이 서버 오류: %s", msg.Content)
			}
			continue
		}

		// Handle message
		c.mu.RLock()
		handler := c.handler
		c.mu.RUnlock()

		if handler != nil {
			go func(m RelayMessage) {
				response, err := handler(m)
				if err != nil {
					log.Printf("메시지 처리 오류: %v", err)
					return
				}
				c.SendResponse(m.From, response, m.SessionID)
			}(msg)
		}
	}
}

// attemptReconnect attempts to reconnect to the relay server
func (c *Client) attemptReconnect() {
	c.mu.Lock()
	if c.reconnecting {
		c.mu.Unlock()
		return
	}
	c.reconnecting = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.reconnecting = false
		c.mu.Unlock()
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		c.mu.Lock()
		c.reconnectCount++
		count := c.reconnectCount
		maxReconnect := c.maxReconnect
		c.mu.Unlock()

		if maxReconnect > 0 && count > maxReconnect {
			log.Printf("최대 재연결 시도 횟수 초과 (%d)", maxReconnect)
			return
		}

		log.Printf("재연결 시도 #%d...", count)

		c.mu.RLock()
		code := c.code
		c.mu.RUnlock()

		if err := c.Connect(code); err != nil {
			log.Printf("재연결 실패: %v", err)
			time.Sleep(c.reconnectDelay)
			continue
		}

		log.Println("재연결 성공!")
		return
	}
}

// SendResponse sends a response to the relay server
func (c *Client) SendResponse(to, content, sessionID string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected || c.conn == nil {
		return fmt.Errorf("연결되지 않았습니다")
	}

	response := RelayResponse{
		Type:      "response",
		To:        to,
		Content:   content,
		SessionID: sessionID,
		Code:      c.code,
	}

	return c.conn.WriteJSON(response)
}

// GetCode returns the current connection code
func (c *Client) GetCode() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.code
}

// IsConnected returns the connection status
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// Status represents relay client status
type Status struct {
	Connected bool   `json:"connected"`
	Code      string `json:"code"`
	URL       string `json:"url"`
}

// GetStatus returns the client status
func (c *Client) GetStatus() Status {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Status{
		Connected: c.connected,
		Code:      c.code,
		URL:       c.url,
	}
}
