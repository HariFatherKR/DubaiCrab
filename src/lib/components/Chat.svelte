<script lang="ts">
	import { onMount } from 'svelte';
	import { chat, isHealthy, type Message, DEFAULT_MODEL, SYSTEM_PROMPT } from '$lib/ollama';
	import { incrementMessages, incrementChats, incrementEmails } from '$lib/stores/stats-store';
	
	interface ChatMessage {
		id: string;
		role: 'user' | 'assistant';
		content: string;
		timestamp: Date;
	}
	
	let messages = $state<ChatMessage[]>([]);
	let inputValue = $state('');
	let isLoading = $state(false);
	let isConnected = $state(false);
	let chatContainer: HTMLDivElement;
	
	onMount(async () => {
		isConnected = await isHealthy();
		if (!isConnected) {
			messages = [{
				id: 'welcome',
				role: 'assistant',
				content: '⚠️ Ollama 서버에 연결할 수 없습니다.\n\n' +
					'Ollama가 실행 중인지 확인해주세요:\n' +
					'```\nollama serve\n```\n\n' +
					'모델이 설치되어 있지 않다면:\n' +
					'```\nollama pull qwen2.5:3b-instruct\n```',
				timestamp: new Date()
			}];
		} else {
			// 새 대화 시작 통계
			incrementChats();
			messages = [{
				id: 'welcome',
				role: 'assistant',
				content: '안녕하세요! 👋 OpenKlaw입니다.\n\n' +
					'무엇을 도와드릴까요?\n\n' +
					'- 📧 이메일 작성: "~에게 ~내용으로 이메일 써줘"\n' +
					'- 📄 문서 요약: 파일을 드래그하거나 텍스트를 붙여넣기\n' +
					'- 💬 일반 질문: 무엇이든 물어보세요!',
				timestamp: new Date()
			}];
		}
	});
	
	function scrollToBottom() {
		if (chatContainer) {
			chatContainer.scrollTop = chatContainer.scrollHeight;
		}
	}
	
	async function handleSubmit(e: Event) {
		e.preventDefault();
		
		if (!inputValue.trim() || isLoading || !isConnected) return;
		
		const userMessage: ChatMessage = {
			id: crypto.randomUUID(),
			role: 'user',
			content: inputValue.trim(),
			timestamp: new Date()
		};
		
		messages = [...messages, userMessage];
		inputValue = '';
		isLoading = true;
		
		// 메시지 통계 업데이트
		incrementMessages();
		
		// 이메일 관련 요청인지 확인
		const isEmailRequest = userMessage.content.includes('이메일') || 
			userMessage.content.includes('메일') ||
			userMessage.content.includes('email');
		
		// 스크롤
		setTimeout(scrollToBottom, 0);
		
		// AI 응답 메시지 준비
		const assistantMessage: ChatMessage = {
			id: crypto.randomUUID(),
			role: 'assistant',
			content: '',
			timestamp: new Date()
		};
		messages = [...messages, assistantMessage];
		
		try {
			// Ollama API 호출을 위한 메시지 형식 변환
			const apiMessages: Message[] = [
				{ role: 'system', content: SYSTEM_PROMPT },
				...messages.slice(0, -1).map(m => ({
					role: m.role as 'user' | 'assistant',
					content: m.content
				}))
			];
			
			// 스트리밍 응답 처리
			for await (const chunk of chat(DEFAULT_MODEL, apiMessages)) {
				if (chunk.message?.content) {
					assistantMessage.content += chunk.message.content;
					messages = [...messages.slice(0, -1), { ...assistantMessage }];
					scrollToBottom();
				}
			}
			
			// 이메일 생성 통계
			if (isEmailRequest && assistantMessage.content.length > 100) {
				incrementEmails();
			}
		} catch (error) {
			console.error('Chat error:', error);
			assistantMessage.content = '❌ 응답 생성 중 오류가 발생했습니다. 다시 시도해주세요.';
			messages = [...messages.slice(0, -1), { ...assistantMessage }];
		} finally {
			isLoading = false;
			scrollToBottom();
		}
	}
	
	function handleKeyDown(e: KeyboardEvent) {
		if (e.key === 'Enter' && !e.shiftKey) {
			e.preventDefault();
			handleSubmit(e);
		}
	}
</script>

<div class="chat-container" bind:this={chatContainer}>
	{#each messages as message (message.id)}
		<div class="message {message.role}">
			<div class="message-avatar">
				{message.role === 'user' ? '👤' : '🦞'}
			</div>
			<div class="message-content">
				<pre>{message.content}</pre>
			</div>
		</div>
	{/each}
	
	{#if isLoading}
		<div class="message assistant">
			<div class="message-avatar">🦞</div>
			<div class="message-content">
				<span class="typing-indicator">
					<span></span>
					<span></span>
					<span></span>
				</span>
			</div>
		</div>
	{/if}
</div>

<form class="chat-input-container" onsubmit={handleSubmit}>
	<textarea
		bind:value={inputValue}
		onkeydown={handleKeyDown}
		placeholder={isConnected ? "메시지를 입력하세요..." : "Ollama 연결 대기 중..."}
		disabled={!isConnected || isLoading}
		rows="1"
	></textarea>
	<button type="submit" disabled={!inputValue.trim() || isLoading || !isConnected}>
		{isLoading ? '⏳' : '📤'}
	</button>
</form>

<style>
	.chat-container {
		flex: 1;
		overflow-y: auto;
		padding: 1rem;
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}
	
	.message {
		display: flex;
		gap: 0.75rem;
		max-width: 85%;
	}
	
	.message.user {
		align-self: flex-end;
		flex-direction: row-reverse;
	}
	
	.message.assistant {
		align-self: flex-start;
	}
	
	.message-avatar {
		font-size: 1.5rem;
		flex-shrink: 0;
	}
	
	.message-content {
		background: var(--color-surface);
		padding: 0.75rem 1rem;
		border-radius: 1rem;
		max-width: 100%;
	}
	
	.message.user .message-content {
		background: var(--color-primary);
	}
	
	.message-content pre {
		margin: 0;
		white-space: pre-wrap;
		word-wrap: break-word;
		font-family: inherit;
		font-size: 0.95rem;
		line-height: 1.5;
	}
	
	.chat-input-container {
		display: flex;
		gap: 0.5rem;
		padding: 1rem;
		background: var(--color-surface);
		border-top: 1px solid rgba(255, 255, 255, 0.1);
	}
	
	.chat-input-container textarea {
		flex: 1;
		padding: 0.75rem 1rem;
		border: none;
		border-radius: 1.5rem;
		background: var(--color-bg);
		color: var(--color-text);
		font-size: 1rem;
		resize: none;
		outline: none;
	}
	
	.chat-input-container textarea:disabled {
		opacity: 0.5;
	}
	
	.chat-input-container button {
		padding: 0.75rem 1rem;
		border: none;
		border-radius: 50%;
		background: var(--color-primary);
		color: white;
		font-size: 1.25rem;
		cursor: pointer;
		transition: background 0.2s;
	}
	
	.chat-input-container button:hover:not(:disabled) {
		background: var(--color-primary-hover);
	}
	
	.chat-input-container button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
	
	.typing-indicator {
		display: flex;
		gap: 4px;
	}
	
	.typing-indicator span {
		width: 8px;
		height: 8px;
		background: var(--color-text-muted);
		border-radius: 50%;
		animation: bounce 1.4s infinite;
	}
	
	.typing-indicator span:nth-child(2) {
		animation-delay: 0.2s;
	}
	
	.typing-indicator span:nth-child(3) {
		animation-delay: 0.4s;
	}
	
	@keyframes bounce {
		0%, 60%, 100% {
			transform: translateY(0);
		}
		30% {
			transform: translateY(-8px);
		}
	}
</style>
