<script lang="ts">
  import { onMount } from 'svelte';
  import { Chat as ChatAPI, GetChatHistory, ClearChatHistory } from '../../../wailsjs/go/main/App';

  interface Props {
    ollamaReady?: boolean;
  }

  interface Message {
    role: 'user' | 'assistant' | 'system';
    content: string;
    timestamp: number;
  }

  let { ollamaReady = false }: Props = $props();
  
  let messages = $state<Message[]>([]);
  let inputValue = $state('');
  let isLoading = $state(false);
  let messagesContainer: HTMLDivElement;
  let sessionId = 'default';

  onMount(() => {
    // Load chat history
    GetChatHistory(sessionId).then(history => {
      messages = history || [];
    });
  });

  async function sendMessage() {
    const content = inputValue.trim();
    if (!content || isLoading || !ollamaReady) return;

    // Add user message
    const userMessage: Message = {
      role: 'user',
      content,
      timestamp: Date.now()
    };
    messages = [...messages, userMessage];
    inputValue = '';
    isLoading = true;

    // Scroll to bottom
    setTimeout(() => {
      if (messagesContainer) {
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
      }
    }, 0);

    try {
      const response = await ChatAPI(sessionId, content);
      
      const assistantMessage: Message = {
        role: 'assistant',
        content: response,
        timestamp: Date.now()
      };
      messages = [...messages, assistantMessage];
    } catch (error) {
      const errorMessage: Message = {
        role: 'assistant',
        content: `오류가 발생했습니다: ${error}`,
        timestamp: Date.now()
      };
      messages = [...messages, errorMessage];
    } finally {
      isLoading = false;
      setTimeout(() => {
        if (messagesContainer) {
          messagesContainer.scrollTop = messagesContainer.scrollHeight;
        }
      }, 0);
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  }

  function clearChat() {
    ClearChatHistory(sessionId);
    messages = [];
  }
</script>

<div class="chat-wrapper glass">
  <div class="chat-header">
    <span>대화</span>
    <button class="clear-btn" onclick={clearChat} title="대화 지우기">
      🗑️
    </button>
  </div>
  
  <div class="messages" bind:this={messagesContainer}>
    {#if messages.length === 0}
      <div class="empty-state">
        <span class="empty-icon">🦀</span>
        <p>안녕하세요! 무엇이든 물어보세요.</p>
        <p class="empty-hint">예: "오늘 회의 안건 정리해줘" 또는 "이메일 초안 작성해줘"</p>
      </div>
    {:else}
      {#each messages as message}
        <div class="message {message.role}">
          <div class="message-content">
            {message.content}
          </div>
        </div>
      {/each}
      {#if isLoading}
        <div class="message assistant loading">
          <div class="message-content">
            <span class="typing-indicator">
              <span></span>
              <span></span>
              <span></span>
            </span>
          </div>
        </div>
      {/if}
    {/if}
  </div>
  
  <div class="input-area">
    <textarea
      bind:value={inputValue}
      onkeydown={handleKeydown}
      placeholder={ollamaReady ? "메시지를 입력하세요..." : "Ollama 연결 대기 중..."}
      disabled={!ollamaReady || isLoading}
      rows={1}
    ></textarea>
    <button 
      class="send-btn" 
      onclick={sendMessage}
      disabled={!inputValue.trim() || !ollamaReady || isLoading}
    >
      ➤
    </button>
  </div>
</div>

<style>
  .chat-wrapper {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
  
  .chat-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem 1rem;
    border-bottom: 1px solid rgba(239, 235, 233, 0.1);
    font-weight: 500;
  }
  
  .clear-btn {
    background: none;
    border: none;
    cursor: pointer;
    font-size: 1rem;
    opacity: 0.6;
    transition: opacity 0.2s;
  }
  
  .clear-btn:hover {
    opacity: 1;
  }
  
  .messages {
    flex: 1;
    overflow-y: auto;
    padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .empty-state {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    color: var(--crab-muted);
    text-align: center;
  }
  
  .empty-icon {
    font-size: 3rem;
    margin-bottom: 1rem;
  }
  
  .empty-hint {
    font-size: 0.85rem;
    opacity: 0.7;
  }
  
  .message {
    max-width: 80%;
    animation: fadeIn 0.3s ease-out;
  }
  
  .message.user {
    align-self: flex-end;
  }
  
  .message.assistant {
    align-self: flex-start;
  }
  
  .message-content {
    padding: 0.75rem 1rem;
    border-radius: 12px;
    line-height: 1.5;
    white-space: pre-wrap;
    word-break: break-word;
  }
  
  .message.user .message-content {
    background: var(--crab-orange);
    color: white;
    border-bottom-right-radius: 4px;
  }
  
  .message.assistant .message-content {
    background: var(--crab-medium);
    border-bottom-left-radius: 4px;
  }
  
  .typing-indicator {
    display: flex;
    gap: 4px;
  }
  
  .typing-indicator span {
    width: 8px;
    height: 8px;
    background: var(--crab-muted);
    border-radius: 50%;
    animation: bounce 1.4s infinite ease-in-out;
  }
  
  .typing-indicator span:nth-child(1) { animation-delay: -0.32s; }
  .typing-indicator span:nth-child(2) { animation-delay: -0.16s; }
  
  @keyframes bounce {
    0%, 80%, 100% { transform: scale(0); }
    40% { transform: scale(1); }
  }
  
  .input-area {
    display: flex;
    gap: 0.75rem;
    padding: 1rem;
    border-top: 1px solid rgba(239, 235, 233, 0.1);
  }
  
  .input-area textarea {
    flex: 1;
    resize: none;
    min-height: 44px;
    max-height: 120px;
    background: var(--crab-dark);
    border: 1px solid var(--crab-accent);
    border-radius: 12px;
    padding: 0.75rem 1rem;
    color: var(--crab-text);
    font-size: 0.95rem;
  }
  
  .input-area textarea:focus {
    outline: none;
    border-color: var(--crab-orange);
  }
  
  .input-area textarea:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  .send-btn {
    width: 44px;
    height: 44px;
    border-radius: 12px;
    background: var(--crab-orange);
    color: white;
    border: none;
    cursor: pointer;
    font-size: 1.1rem;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
  }
  
  .send-btn:hover:not(:disabled) {
    background: #FF7043;
    transform: scale(1.05);
  }
  
  .send-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }
</style>
