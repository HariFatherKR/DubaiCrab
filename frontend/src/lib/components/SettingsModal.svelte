<script lang="ts">
  import { onMount } from 'svelte';
  import { 
    GetOllamaModels, 
    SetOllamaModel,
    GetKakaoStatus,
    GetKakaoConfig,
    UpdateKakaoConfig,
    StartKakaoServer,
    StopKakaoServer
  } from '../../../wailsjs/go/main/App';

  interface Props {
    onClose: () => void;
  }

  let { onClose }: Props = $props();
  
  // Ollama settings
  let models = $state<string[]>([]);
  let selectedModel = $state('qwen2.5:3b');
  
  // Kakao settings
  let kakaoEnabled = $state(true);
  let kakaoRunning = $state(false);
  let kakaoPort = $state(3847);
  let kakaoWebhookPath = $state('/kakao/webhook');
  let kakaoDmPolicy = $state('open');
  let kakaoSystemPrompt = $state('');
  let kakaoModel = $state('qwen2.5:3b');
  
  let activeTab = $state<'general' | 'kakao' | 'advanced'>('general');
  let saving = $state(false);

  onMount(async () => {
    try {
      // Load Ollama models
      models = await GetOllamaModels();
      
      // Load Kakao status
      const status = await GetKakaoStatus();
      kakaoRunning = status.running;
      
      // Load Kakao config
      const config = await GetKakaoConfig();
      kakaoEnabled = config.enabled;
      kakaoPort = config.port;
      kakaoWebhookPath = config.webhookPath;
      kakaoDmPolicy = config.dmPolicy;
      kakaoSystemPrompt = config.systemPrompt || '';
      kakaoModel = config.model || 'qwen2.5:3b';
    } catch (error) {
      console.error('Failed to load settings:', error);
    }
  });

  async function saveSettings() {
    saving = true;
    try {
      // Save Ollama model
      await SetOllamaModel(selectedModel);
      
      // Save Kakao config
      await UpdateKakaoConfig({
        enabled: kakaoEnabled,
        port: kakaoPort,
        webhookPath: kakaoWebhookPath,
        dmPolicy: kakaoDmPolicy,
        allowFrom: [],
        systemPrompt: kakaoSystemPrompt,
        model: kakaoModel
      });
      
      onClose();
    } catch (error) {
      console.error('Failed to save settings:', error);
      alert('설정 저장에 실패했습니다.');
    } finally {
      saving = false;
    }
  }

  async function toggleKakaoServer() {
    try {
      if (kakaoRunning) {
        await StopKakaoServer();
        kakaoRunning = false;
      } else {
        await StartKakaoServer();
        kakaoRunning = true;
      }
    } catch (error) {
      console.error('Failed to toggle Kakao server:', error);
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      onClose();
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="modal-overlay" onclick={onClose} role="presentation">
  <div class="modal" onclick={(e) => e.stopPropagation()} role="dialog" aria-modal="true">
    <header class="modal-header">
      <h2>설정</h2>
      <button class="close-btn" onclick={onClose}>✕</button>
    </header>
    
    <nav class="tabs">
      <button 
        class="tab" 
        class:active={activeTab === 'general'}
        onclick={() => activeTab = 'general'}
      >
        일반
      </button>
      <button 
        class="tab" 
        class:active={activeTab === 'kakao'}
        onclick={() => activeTab = 'kakao'}
      >
        카카오톡
      </button>
      <button 
        class="tab" 
        class:active={activeTab === 'advanced'}
        onclick={() => activeTab = 'advanced'}
      >
        고급
      </button>
    </nav>
    
    <div class="modal-content">
      {#if activeTab === 'general'}
        <div class="settings-section">
          <h3>AI 모델</h3>
          <div class="form-group">
            <label for="model">기본 모델</label>
            <select id="model" bind:value={selectedModel}>
              {#each models as model}
                <option value={model}>{model}</option>
              {/each}
            </select>
          </div>
        </div>
      {:else if activeTab === 'kakao'}
        <div class="settings-section">
          <h3>카카오톡 연동</h3>
          
          <div class="form-group row">
            <label for="kakao-enabled">활성화</label>
            <input 
              type="checkbox" 
              id="kakao-enabled" 
              bind:checked={kakaoEnabled}
            />
          </div>
          
          <div class="form-group row">
            <span>서버 상태</span>
            <span class="status" class:running={kakaoRunning}>
              {kakaoRunning ? '🟢 실행 중' : '🔴 중지됨'}
            </span>
            <button class="btn-small" onclick={toggleKakaoServer}>
              {kakaoRunning ? '중지' : '시작'}
            </button>
          </div>
          
          <div class="form-group">
            <label for="kakao-port">포트</label>
            <input 
              type="number" 
              id="kakao-port" 
              bind:value={kakaoPort}
              min={1024}
              max={65535}
            />
          </div>
          
          <div class="form-group">
            <label for="kakao-model">사용 모델</label>
            <select id="kakao-model" bind:value={kakaoModel}>
              {#each models as model}
                <option value={model}>{model}</option>
              {/each}
            </select>
          </div>
          
          <div class="form-group">
            <label for="kakao-prompt">시스템 프롬프트</label>
            <textarea 
              id="kakao-prompt" 
              bind:value={kakaoSystemPrompt}
              rows={3}
              placeholder="AI 비서의 역할과 성격을 정의하세요"
            ></textarea>
          </div>
        </div>
      {:else if activeTab === 'advanced'}
        <div class="settings-section">
          <h3>고급 설정</h3>
          <p class="hint">추후 업데이트 예정</p>
        </div>
      {/if}
    </div>
    
    <footer class="modal-footer">
      <button class="btn-secondary" onclick={onClose}>취소</button>
      <button class="btn-primary" onclick={saveSettings} disabled={saving}>
        {saving ? '저장 중...' : '저장'}
      </button>
    </footer>
  </div>
</div>

<style>
  .modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(4px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 100;
  }
  
  .modal {
    background: var(--crab-dark);
    border: 1px solid rgba(239, 235, 233, 0.2);
    border-radius: 16px;
    width: 90%;
    max-width: 560px;
    max-height: 80vh;
    display: flex;
    flex-direction: column;
    animation: slideUp 0.2s ease-out;
  }
  
  @keyframes slideUp {
    from { opacity: 0; transform: translateY(20px); }
    to { opacity: 1; transform: translateY(0); }
  }
  
  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.25rem 1.5rem;
    border-bottom: 1px solid rgba(239, 235, 233, 0.1);
  }
  
  .modal-header h2 {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 600;
  }
  
  .close-btn {
    background: none;
    border: none;
    color: var(--crab-muted);
    cursor: pointer;
    font-size: 1.25rem;
    padding: 0.25rem;
    transition: color 0.2s;
  }
  
  .close-btn:hover {
    color: var(--crab-text);
  }
  
  .tabs {
    display: flex;
    padding: 0 1.5rem;
    gap: 0.5rem;
    border-bottom: 1px solid rgba(239, 235, 233, 0.1);
  }
  
  .tab {
    padding: 0.75rem 1rem;
    background: none;
    border: none;
    color: var(--crab-muted);
    cursor: pointer;
    font-size: 0.9rem;
    border-bottom: 2px solid transparent;
    transition: all 0.2s;
    margin-bottom: -1px;
  }
  
  .tab:hover {
    color: var(--crab-text);
  }
  
  .tab.active {
    color: var(--crab-orange);
    border-bottom-color: var(--crab-orange);
  }
  
  .modal-content {
    flex: 1;
    overflow-y: auto;
    padding: 1.5rem;
  }
  
  .settings-section h3 {
    font-size: 1rem;
    font-weight: 600;
    margin: 0 0 1rem 0;
    color: var(--crab-text);
  }
  
  .form-group {
    margin-bottom: 1rem;
  }
  
  .form-group.row {
    display: flex;
    align-items: center;
    gap: 1rem;
  }
  
  .form-group label {
    display: block;
    margin-bottom: 0.5rem;
    color: var(--crab-muted);
    font-size: 0.9rem;
  }
  
  .form-group.row label {
    margin-bottom: 0;
    min-width: 80px;
  }
  
  .form-group input[type="text"],
  .form-group input[type="number"],
  .form-group select,
  .form-group textarea {
    width: 100%;
    padding: 0.625rem 0.875rem;
    background: var(--crab-medium);
    border: 1px solid var(--crab-accent);
    border-radius: 8px;
    color: var(--crab-text);
    font-size: 0.9rem;
  }
  
  .form-group input:focus,
  .form-group select:focus,
  .form-group textarea:focus {
    outline: none;
    border-color: var(--crab-orange);
  }
  
  .form-group input[type="checkbox"] {
    width: 18px;
    height: 18px;
    accent-color: var(--crab-orange);
  }
  
  .status {
    font-size: 0.85rem;
    flex: 1;
  }
  
  .btn-small {
    padding: 0.375rem 0.75rem;
    font-size: 0.8rem;
    background: var(--crab-accent);
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    transition: opacity 0.2s;
  }
  
  .btn-small:hover {
    opacity: 0.8;
  }
  
  .hint {
    color: var(--crab-muted);
    font-size: 0.9rem;
  }
  
  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
    padding: 1rem 1.5rem;
    border-top: 1px solid rgba(239, 235, 233, 0.1);
  }
  
  .btn-primary,
  .btn-secondary {
    padding: 0.625rem 1.25rem;
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    border: none;
  }
  
  .btn-primary {
    background: var(--crab-orange);
    color: white;
  }
  
  .btn-primary:hover:not(:disabled) {
    background: #FF7043;
  }
  
  .btn-primary:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  .btn-secondary {
    background: var(--crab-medium);
    color: var(--crab-text);
  }
  
  .btn-secondary:hover {
    background: var(--crab-light);
  }
</style>
