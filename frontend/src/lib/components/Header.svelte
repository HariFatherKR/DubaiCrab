<script lang="ts">
  interface Props {
    title: string;
    subtitle?: string;
  }

  let { title, subtitle }: Props = $props();
</script>

<header class="header">
  <div class="header-content">
    <h1 class="header-title">{title}</h1>
    {#if subtitle}
      <p class="header-subtitle">{subtitle}</p>
    {/if}
  </div>
</header>

<style>
  .header {
    padding: 1.5rem 2rem;
    /* Add padding for macOS titlebar */
    padding-top: 2.5rem;
  }
  
  .header-content {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }
  
  .header-title {
    font-size: 1.75rem;
    font-weight: 600;
    color: var(--crab-text);
    margin: 0;
  }
  
  .header-subtitle {
    font-size: 0.9rem;
    color: var(--crab-muted);
    margin: 0;
  }
</style>
