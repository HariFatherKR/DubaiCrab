<script lang="ts">
  interface Props {
    onSettings?: () => void;
  }

  let { onSettings }: Props = $props();
</script>

<aside class="sidebar">
  <div class="logo">
    <span class="logo-icon">🦀</span>
    <span class="logo-text">Dubai Crab</span>
  </div>
  
  <nav class="nav">
    <a href="/" class="nav-item active">
      <span class="nav-icon">💬</span>
      <span>채팅</span>
    </a>
    <a href="/dashboard" class="nav-item">
      <span class="nav-icon">📊</span>
      <span>대시보드</span>
    </a>
  </nav>
  
  <div class="sidebar-footer">
    <button class="nav-item" onclick={onSettings}>
      <span class="nav-icon">⚙️</span>
      <span>설정</span>
    </button>
  </div>
</aside>

<style>
  .sidebar {
    width: 240px;
    display: flex;
    flex-direction: column;
    background: rgba(62, 39, 35, 0.8);
    backdrop-filter: blur(10px);
    border-right: 1px solid rgba(239, 235, 233, 0.1);
    padding: 1rem;
  }
  
  .logo {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 1rem 0.5rem;
    margin-bottom: 1rem;
  }
  
  .logo-icon {
    font-size: 1.75rem;
  }
  
  .logo-text {
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--crab-text);
  }
  
  .nav {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }
  
  .nav-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    border-radius: 8px;
    color: var(--crab-muted);
    text-decoration: none;
    transition: all 0.2s;
    cursor: pointer;
    background: none;
    border: none;
    font-size: 0.95rem;
    width: 100%;
    text-align: left;
  }
  
  .nav-item:hover {
    background: rgba(141, 110, 99, 0.3);
    color: var(--crab-text);
  }
  
  .nav-item.active {
    background: var(--crab-orange);
    color: white;
  }
  
  .nav-icon {
    font-size: 1.1rem;
  }
  
  .sidebar-footer {
    padding-top: 1rem;
    border-top: 1px solid rgba(239, 235, 233, 0.1);
    margin-top: 1rem;
  }
</style>
