<script lang="ts">
	import Chat from '$lib/components/Chat.svelte';
</script>

<svelte:head>
	<title>OpenKlaw - 로컬 AI 비서</title>
</svelte:head>

<main class="flex-1 flex flex-col">
	<header class="bg-surface px-4 py-3 flex items-center justify-between border-b border-gray-700">
		<div class="flex items-center gap-2">
			<span class="text-2xl">🦞</span>
			<h1 class="text-xl font-bold text-primary">OpenKlaw</h1>
		</div>
		<div class="flex items-center gap-4">
			<a href="/dashboard" class="dashboard-link">📊 대시보드</a>
			<span class="text-sm text-muted">로컬 AI 비서</span>
		</div>
	</header>
	
	<Chat />
</main>

<style>
	.bg-surface {
		background-color: var(--color-surface);
	}
	.text-primary {
		color: var(--color-primary);
	}
	.text-muted {
		color: var(--color-text-muted);
	}
	.dashboard-link {
		color: var(--color-text-muted);
		text-decoration: none;
		font-size: 0.9rem;
		padding: 0.4rem 0.75rem;
		border-radius: 0.5rem;
		background: rgba(255, 255, 255, 0.05);
		transition: all 0.2s;
	}
	.dashboard-link:hover {
		background: rgba(255, 255, 255, 0.1);
		color: var(--color-primary);
	}
</style>
