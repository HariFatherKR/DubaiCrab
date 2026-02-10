<script lang="ts">
  import { onMount } from 'svelte';
  import Sidebar from '$lib/components/Sidebar.svelte';
  import Chat from '$lib/components/Chat.svelte';
  import Header from '$lib/components/Header.svelte';
  import SettingsModal from '$lib/components/SettingsModal.svelte';
  import { appState } from '$lib/stores/app';
  import { EventsOn } from '../wailsjs/runtime/runtime';
  import { CheckOllama, GetOllamaModels } from '../wailsjs/go/main/App';

  let showSettings = $state(false);
  let ollamaReady = $state(false);
  let models = $state<string[]>([]);

  onMount(() => {
    // Check Ollama status
    CheckOllama().then(ready => {
      ollamaReady = ready;
      if (ready) {
        GetOllamaModels().then(m => models = m);
      }
    });

    // Listen for Ollama events
    EventsOn('ollama:ready', (ready: boolean) => {
      ollamaReady = ready;
    });

    EventsOn('ollama:models', (m: string[]) => {
      models = m;
    });

    EventsOn('ollama:error', (err: string) => {
      console.error('Ollama error:', err);
    });
  });

  function openSettings() {
    showSettings = true;
  }

  function closeSettings() {
    showSettings = false;
  }
</script>

<div class="app-layout">
  <Sidebar onSettings={openSettings} />
  
  <main class="main-content">
    <Header 
      title="채팅" 
      subtitle={ollamaReady ? `${models.length}개 모델 사용 가능` : 'Ollama 연결 중...'} 
    />
    
    <div class="chat-container">
      <Chat {ollamaReady} />
    </div>
  </main>
</div>

{#if showSettings}
  <SettingsModal onClose={closeSettings} />
{/if}

<style>
  .app-layout {
    height: 100%;
    display: flex;
  }
  
  .main-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
  
  .chat-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    padding: 0 2rem 2rem 2rem;
    overflow: hidden;
  }
</style>
