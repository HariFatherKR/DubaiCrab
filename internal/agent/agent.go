package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"DubaiCrab-Go/internal/ollama"
	"DubaiCrab-Go/internal/tools"
)

// Message represents a conversation message
type Message struct {
	Role      string    `json:"role"`      // "user", "assistant", "system", "tool"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	ToolName  string    `json:"toolName,omitempty"`
	ToolCall  string    `json:"toolCall,omitempty"`
}

// Session represents a conversation session
type Session struct {
	ID        string    `json:"id"`
	Messages  []Message `json:"messages"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Agent orchestrates the AI conversation loop
type Agent struct {
	ollama       *ollama.Manager
	toolRegistry *tools.Registry
	sessions     map[string]*Session
	mu           sync.RWMutex
	
	// Configuration
	model        string
	systemPrompt string
	maxTokens    int
}

// Config holds agent configuration
type Config struct {
	Model        string
	SystemPrompt string
	MaxTokens    int
}

// NewAgent creates a new agent
func NewAgent(ollamaManager *ollama.Manager, registry *tools.Registry) *Agent {
	return &Agent{
		ollama:       ollamaManager,
		toolRegistry: registry,
		sessions:     make(map[string]*Session),
		model:        "qwen2.5:3b",
		systemPrompt: defaultSystemPrompt,
		maxTokens:    4096,
	}
}

const defaultSystemPrompt = `당신은 Dubai Crab, 한국 사무직을 위한 친절하고 유능한 AI 비서입니다.

## 역할
- 업무 관련 질문에 정확하고 간결하게 답변
- 문서 작성, 데이터 분석, 일정 관리 지원
- 한국어로 자연스럽게 대화

## 원칙
- 항상 존댓말 사용
- 불확실한 정보는 솔직히 모른다고 답변
- 요청에 따라 필요한 도구 사용
`

// Configure updates agent configuration
func (a *Agent) Configure(cfg Config) {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	if cfg.Model != "" {
		a.model = cfg.Model
	}
	if cfg.SystemPrompt != "" {
		a.systemPrompt = cfg.SystemPrompt
	}
	if cfg.MaxTokens > 0 {
		a.maxTokens = cfg.MaxTokens
	}
}

// GetOrCreateSession gets or creates a session
func (a *Agent) GetOrCreateSession(sessionID string) *Session {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	if session, ok := a.sessions[sessionID]; ok {
		return session
	}
	
	session := &Session{
		ID:        sessionID,
		Messages:  []Message{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	a.sessions[sessionID] = session
	return session
}

// ClearSession clears a session
func (a *Agent) ClearSession(sessionID string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.sessions, sessionID)
}

// ProcessMessage processes a user message and returns the response
func (a *Agent) ProcessMessage(ctx context.Context, sessionID, userMessage string) (string, error) {
	session := a.GetOrCreateSession(sessionID)
	
	// Add user message
	a.addMessage(session, Message{
		Role:      "user",
		Content:   userMessage,
		Timestamp: time.Now(),
	})
	
	// Build context
	messages := a.buildContext(session)
	
	// Get response from LLM
	response, err := a.ollama.Chat(ctx, a.model, a.formatMessages(messages), &a.systemPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to get response: %w", err)
	}
	
	// Check for tool calls (simple pattern matching for now)
	if toolCall := a.extractToolCall(response); toolCall != nil {
		toolResult, err := a.executeToolCall(ctx, toolCall)
		if err != nil {
			log.Printf("[agent] Tool execution failed: %v", err)
			// Add tool error to context and retry
			a.addMessage(session, Message{
				Role:      "tool",
				Content:   fmt.Sprintf("Error: %v", err),
				Timestamp: time.Now(),
				ToolName:  toolCall.Name,
			})
		} else {
			// Add tool result and get final response
			a.addMessage(session, Message{
				Role:      "tool",
				Content:   toolResult,
				Timestamp: time.Now(),
				ToolName:  toolCall.Name,
			})
			
			// Get final response incorporating tool result
			messages = a.buildContext(session)
			response, err = a.ollama.Chat(ctx, a.model, a.formatMessages(messages), &a.systemPrompt)
			if err != nil {
				return "", fmt.Errorf("failed to get final response: %w", err)
			}
		}
	}
	
	// Add assistant response
	a.addMessage(session, Message{
		Role:      "assistant",
		Content:   response,
		Timestamp: time.Now(),
	})
	
	return response, nil
}

// addMessage adds a message to a session
func (a *Agent) addMessage(session *Session, msg Message) {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	session.Messages = append(session.Messages, msg)
	session.UpdatedAt = time.Now()
	
	// Prune old messages if too many
	if len(session.Messages) > 100 {
		session.Messages = session.Messages[len(session.Messages)-50:]
	}
}

// buildContext builds the context for the LLM
func (a *Agent) buildContext(session *Session) []Message {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	// Get recent messages (context window management)
	messages := session.Messages
	if len(messages) > 20 {
		messages = messages[len(messages)-20:]
	}
	
	return messages
}

// formatMessages formats messages for the LLM
func (a *Agent) formatMessages(messages []Message) string {
	var parts []string
	for _, msg := range messages {
		prefix := ""
		switch msg.Role {
		case "user":
			prefix = "사용자: "
		case "assistant":
			prefix = "비서: "
		case "tool":
			prefix = fmt.Sprintf("[%s 결과]: ", msg.ToolName)
		}
		parts = append(parts, prefix+msg.Content)
	}
	return strings.Join(parts, "\n\n")
}

// ToolCall represents a tool call
type ToolCall struct {
	Name   string                 `json:"name"`
	Params map[string]interface{} `json:"params"`
}

// extractToolCall extracts a tool call from the response (simple pattern matching)
func (a *Agent) extractToolCall(response string) *ToolCall {
	// Look for tool call patterns like: @tool_name({"param": "value"})
	if !strings.Contains(response, "@") {
		return nil
	}
	
	// Simple extraction - can be enhanced with proper parsing
	idx := strings.Index(response, "@")
	if idx == -1 {
		return nil
	}
	
	rest := response[idx+1:]
	parenIdx := strings.Index(rest, "(")
	if parenIdx == -1 {
		return nil
	}
	
	toolName := rest[:parenIdx]
	if !a.toolRegistry.Has(toolName) {
		return nil
	}
	
	// Find matching closing paren
	closeIdx := strings.LastIndex(rest, ")")
	if closeIdx <= parenIdx {
		return nil
	}
	
	paramsStr := rest[parenIdx+1 : closeIdx]
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(paramsStr), &params); err != nil {
		log.Printf("[agent] Failed to parse tool params: %v", err)
		return nil
	}
	
	return &ToolCall{
		Name:   toolName,
		Params: params,
	}
}

// executeToolCall executes a tool call
func (a *Agent) executeToolCall(ctx context.Context, call *ToolCall) (string, error) {
	tool, ok := a.toolRegistry.Get(call.Name)
	if !ok {
		return "", fmt.Errorf("unknown tool: %s", call.Name)
	}
	
	result, err := tool.Execute(ctx, call.Params)
	if err != nil {
		return "", err
	}
	
	return result, nil
}

// GetSessionHistory returns the message history for a session
func (a *Agent) GetSessionHistory(sessionID string) []Message {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	session, ok := a.sessions[sessionID]
	if !ok {
		return []Message{}
	}
	
	// Return a copy
	messages := make([]Message, len(session.Messages))
	copy(messages, session.Messages)
	return messages
}

// ListSessions returns all session IDs
func (a *Agent) ListSessions() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	ids := make([]string, 0, len(a.sessions))
	for id := range a.sessions {
		ids = append(ids, id)
	}
	return ids
}
