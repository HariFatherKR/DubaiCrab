<script lang="ts">
	import { onMount } from 'svelte';
	import { chat, isHealthy, type Message, DEFAULT_MODEL, SYSTEM_PROMPT } from '$lib/ollama';
	import { incrementMessages, incrementChats, incrementEmails, incrementHwp } from '$lib/stores/stats-store';
	import QuickActions from './QuickActions.svelte';
	import ReportTemplateModal from './ReportTemplateModal.svelte';
	import type { QuickAction } from './QuickActions.svelte';
	
	interface ChatMessage {
		id: string;
		role: 'user' | 'assistant';
		content: string;
		timestamp: Date;
		file?: { name: string; type: string; size: number };
	}
	
	let messages = $state<ChatMessage[]>([]);
	let inputValue = $state('');
	let isLoading = $state(false);
	let isConnected = $state(false);
	let chatContainer: HTMLDivElement;
	let showReportModal = $state(false);
	let textareaRef: HTMLTextAreaElement;
	let isDragging = $state(false);
	
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
				content: '안녕하세요! 🦀 Dubai Crab입니다.\n\n' +
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
	
	function handleQuickAction(action: QuickAction) {
		if (action.prompt) {
			// AI 응답으로 프롬프트 표시
			const assistantMessage: ChatMessage = {
				id: crypto.randomUUID(),
				role: 'assistant',
				content: action.prompt,
				timestamp: new Date()
			};
			messages = [...messages, assistantMessage];
			
			// 입력창 포커스
			setTimeout(() => {
				textareaRef?.focus();
			}, 100);
		}
	}
	
	function handleReportGenerate(content: string) {
		// 보고서를 AI 응답으로 추가
		const assistantMessage: ChatMessage = {
			id: crypto.randomUUID(),
			role: 'assistant',
			content: content,
			timestamp: new Date()
		};
		messages = [...messages, assistantMessage];
		scrollToBottom();
	}
	
	function openReportModal() {
		showReportModal = true;
	}
	
	function closeReportModal() {
		showReportModal = false;
	}
	
	// Drag and Drop handlers
	function handleDragEnter(e: DragEvent) {
		e.preventDefault();
		e.stopPropagation();
		isDragging = true;
	}
	
	function handleDragOver(e: DragEvent) {
		e.preventDefault();
		e.stopPropagation();
		isDragging = true;
	}
	
	function handleDragLeave(e: DragEvent) {
		e.preventDefault();
		e.stopPropagation();
		isDragging = false;
	}
	
	async function handleDrop(e: DragEvent) {
		e.preventDefault();
		e.stopPropagation();
		isDragging = false;
		
		if (!isConnected || isLoading) return;
		
		const files = e.dataTransfer?.files;
		if (!files || files.length === 0) return;
		
		const file = files[0];
		const fileName = file.name;
		const fileExt = fileName.split('.').pop()?.toLowerCase() || '';
		
		// 지원하는 파일 타입 체크
		const supportedTypes = ['txt', 'md', 'csv', 'json', 'hwp', 'hwpx', 'xlsx', 'xls'];
		if (!supportedTypes.includes(fileExt)) {
			const errorMessage: ChatMessage = {
				id: crypto.randomUUID(),
				role: 'assistant',
				content: `⚠️ 지원하지 않는 파일 형식입니다: .${fileExt}\n\n지원 형식: ${supportedTypes.join(', ')}`,
				timestamp: new Date()
			};
			messages = [...messages, errorMessage];
			return;
		}
		
		// 파일 읽기
		try {
			let fileContent = '';
			
			if (['txt', 'md', 'csv', 'json'].includes(fileExt)) {
				// 텍스트 파일 직접 읽기
				fileContent = await file.text();
			} else if (['hwp', 'hwpx'].includes(fileExt)) {
				// HWP 파일은 서버/백엔드 처리 필요
				fileContent = `[HWP 파일: ${fileName}]\n\n⚠️ HWP 파일 파싱은 현재 개발 중입니다.\n파일 크기: ${(file.size / 1024).toFixed(1)}KB`;
				incrementHwp();
			} else if (['xlsx', 'xls'].includes(fileExt)) {
				// Excel 파일
				fileContent = `[Excel 파일: ${fileName}]\n\n⚠️ Excel 파일 파싱은 현재 개발 중입니다.\n파일 크기: ${(file.size / 1024).toFixed(1)}KB`;
			}
			
			// 사용자 메시지로 파일 내용 추가
			const userMessage: ChatMessage = {
				id: crypto.randomUUID(),
				role: 'user',
				content: `📎 파일: ${fileName}\n\n${fileContent.slice(0, 3000)}${fileContent.length > 3000 ? '\n\n...(내용 생략)' : ''}`,
				timestamp: new Date(),
				file: { name: fileName, type: fileExt, size: file.size }
			};
			
			messages = [...messages, userMessage];
			incrementMessages();
			scrollToBottom();
			
			// 파일 분석 요청
			isLoading = true;
			
			const assistantMessage: ChatMessage = {
				id: crypto.randomUUID(),
				role: 'assistant',
				content: '',
				timestamp: new Date()
			};
			messages = [...messages, assistantMessage];
			
			const apiMessages: Message[] = [
				{ role: 'system', content: SYSTEM_PROMPT },
				{ role: 'user', content: `다음 파일의 내용을 분석하고 요약해주세요:\n\n파일명: ${fileName}\n\n${fileContent.slice(0, 5000)}` }
			];
			
			for await (const chunk of chat(DEFAULT_MODEL, apiMessages)) {
				if (chunk.message?.content) {
					assistantMessage.content += chunk.message.content;
					messages = [...messages.slice(0, -1), { ...assistantMessage }];
					scrollToBottom();
				}
			}
		} catch (error) {
			console.error('File read error:', error);
			const errorMessage: ChatMessage = {
				id: crypto.randomUUID(),
				role: 'assistant',
				content: `❌ 파일을 읽는 중 오류가 발생했습니다: ${error}`,
				timestamp: new Date()
			};
			messages = [...messages, errorMessage];
		} finally {
			isLoading = false;
			scrollToBottom();
		}
	}
</script>

<div 
	class="chat-wrapper"
	class:dragging={isDragging}
	ondragenter={handleDragEnter}
	ondragover={handleDragOver}
	ondragleave={handleDragLeave}
	ondrop={handleDrop}
	role="region"
	aria-label="채팅 영역"
>
	{#if isDragging}
		<div class="drop-overlay">
			<div class="drop-icon">📎</div>
			<p>파일을 놓아주세요</p>
			<span>txt, md, csv, json, hwp, xlsx 지원</span>
		</div>
	{/if}
	<div class="chat-container" bind:this={chatContainer}>
		{#each messages as message (message.id)}
			<div class="message {message.role}">
				<div class="message-avatar">
					{message.role === 'user' ? '👤' : '🦀'}
				</div>
				<div class="message-content">
					<pre>{message.content}</pre>
				</div>
			</div>
		{/each}
		
		{#if isLoading}
			<div class="message assistant">
				<div class="message-avatar">🦀</div>
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
		<div class="input-row">
			<QuickActions 
				onAction={handleQuickAction}
				onOpenTemplates={openReportModal}
			/>
			<div class="input-wrapper">
				<textarea
					bind:this={textareaRef}
					bind:value={inputValue}
					onkeydown={handleKeyDown}
					placeholder={isConnected ? "메시지를 입력하세요..." : "Ollama 연결 대기 중..."}
					disabled={!isConnected || isLoading}
					rows="1"
				></textarea>
				<button type="submit" disabled={!inputValue.trim() || isLoading || !isConnected}>
					{isLoading ? '⏳' : '📤'}
				</button>
			</div>
		</div>
	</form>
</div>

<ReportTemplateModal 
	isOpen={showReportModal}
	onClose={closeReportModal}
	onGenerate={handleReportGenerate}
/>

<style>
	.chat-wrapper {
		flex: 1;
		display: flex;
		flex-direction: column;
		/* Dark pistachio glassmorphism */
		background: rgba(74, 124, 89, 0.15);
		backdrop-filter: blur(20px);
		-webkit-backdrop-filter: blur(20px);
		border: 1px solid rgba(74, 124, 89, 0.25);
		border-radius: 20px;
		overflow: hidden;
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
		position: relative;
		transition: all 0.3s ease;
	}
	
	.chat-wrapper.dragging {
		border-color: rgba(212, 165, 116, 0.6);
		box-shadow: 0 0 0 3px rgba(212, 165, 116, 0.2), 0 8px 32px rgba(0, 0, 0, 0.2);
	}
	
	.drop-overlay {
		position: absolute;
		inset: 0;
		background: rgba(45, 90, 63, 0.95);
		backdrop-filter: blur(10px);
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		z-index: 100;
		border-radius: 20px;
		animation: fadeIn 0.2s ease-out;
	}
	
	.drop-overlay .drop-icon {
		font-size: 4rem;
		margin-bottom: 1rem;
		animation: bounce 1s infinite;
	}
	
	.drop-overlay p {
		font-size: 1.5rem;
		font-weight: 600;
		color: #FFF8E1;
		margin: 0 0 0.5rem;
	}
	
	.drop-overlay span {
		font-size: 0.9rem;
		color: #BCAAA4;
	}
	
	.chat-container {
		flex: 1;
		overflow-y: auto;
		padding: 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}
	
	.message {
		display: flex;
		gap: 0.75rem;
		max-width: 85%;
		animation: fadeIn 0.3s ease-out;
	}
	
	@keyframes fadeIn {
		from { opacity: 0; transform: translateY(10px); }
		to { opacity: 1; transform: translateY(0); }
	}
	
	.message.user {
		align-self: flex-end;
		flex-direction: row-reverse;
	}
	
	.message.assistant {
		align-self: flex-start;
	}
	
	.message-avatar {
		width: 40px;
		height: 40px;
		border-radius: 12px;
		background: rgba(74, 124, 89, 0.2);
		border: 1px solid rgba(74, 124, 89, 0.3);
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.25rem;
		flex-shrink: 0;
	}
	
	.message.user .message-avatar {
		background: linear-gradient(135deg, rgba(212, 165, 116, 0.3), rgba(184, 134, 11, 0.2));
		border-color: rgba(212, 165, 116, 0.4);
	}
	
	.message-content {
		background: rgba(74, 124, 89, 0.15);
		padding: 1rem 1.25rem;
		border-radius: 16px;
		border: 1px solid rgba(74, 124, 89, 0.2);
		max-width: 100%;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
	}
	
	.message.user .message-content {
		background: linear-gradient(135deg, rgba(212, 165, 116, 0.2), rgba(184, 134, 11, 0.15));
		border-color: rgba(212, 165, 116, 0.3);
	}
	
	.message-content pre {
		margin: 0;
		white-space: pre-wrap;
		word-wrap: break-word;
		font-family: inherit;
		font-size: 0.95rem;
		line-height: 1.6;
		color: #FFF8E1;
	}
	
	.chat-input-container {
		padding: 1rem 1.5rem 1.5rem;
		background: rgba(45, 90, 63, 0.2);
		border-top: 1px solid rgba(74, 124, 89, 0.2);
	}
	
	.input-row {
		display: flex;
		gap: 0.75rem;
		align-items: flex-end;
	}
	
	.input-wrapper {
		flex: 1;
		display: flex;
		gap: 0.75rem;
		background: rgba(74, 124, 89, 0.15);
		border: 1px solid rgba(74, 124, 89, 0.25);
		border-radius: 16px;
		padding: 0.5rem;
		transition: all 0.2s ease;
	}
	
	.input-wrapper:focus-within {
		border-color: rgba(212, 165, 116, 0.5);
		box-shadow: 0 0 0 3px rgba(212, 165, 116, 0.1);
	}
	
	.chat-input-container textarea {
		flex: 1;
		padding: 0.75rem 1rem;
		border: none;
		border-radius: 12px;
		background: transparent;
		color: #FFF8E1;
		font-size: 0.95rem;
		resize: none;
		outline: none;
		font-family: inherit;
	}
	
	.chat-input-container textarea::placeholder {
		color: #BCAAA4;
	}
	
	.chat-input-container textarea:disabled {
		opacity: 0.5;
	}
	
	.chat-input-container button {
		padding: 0.75rem 1.25rem;
		border: none;
		border-radius: 12px;
		background: linear-gradient(135deg, #D4A574, #B8860B);
		color: #3E2723;
		font-size: 1.25rem;
		cursor: pointer;
		transition: all 0.2s ease;
		display: flex;
		align-items: center;
		justify-content: center;
		box-shadow: 0 4px 12px rgba(212, 165, 116, 0.3);
	}
	
	.chat-input-container button:hover:not(:disabled) {
		background: linear-gradient(135deg, #B8860B, #D4A574);
		box-shadow: 0 6px 20px rgba(212, 165, 116, 0.4);
		transform: translateY(-1px);
	}
	
	.chat-input-container button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
	
	.typing-indicator {
		display: flex;
		gap: 4px;
		padding: 0.25rem;
	}
	
	.typing-indicator span {
		width: 8px;
		height: 8px;
		background: #D4A574;
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
