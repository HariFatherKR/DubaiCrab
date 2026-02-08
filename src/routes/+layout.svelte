<script lang="ts">
	import '../app.css';
	import Sidebar from '$lib/components/Sidebar.svelte';
	
	let { children } = $props();
</script>

<div class="app-layout">
	<Sidebar />
	<main class="main-content">
		{@render children()}
	</main>
</div>

<style>
	.app-layout {
		min-height: 100vh;
		display: flex;
		/* Light mint/teal gradient background */
		background: linear-gradient(
			135deg,
			#e6f5f3 0%,
			#d9ece9 25%,
			#cce4e7 50%,
			#c5e1ec 75%,
			#d4eaf5 100%
		);
		background-attachment: fixed;
	}
	
	.main-content {
		flex: 1;
		display: flex;
		flex-direction: column;
		overflow-x: hidden;
	}
</style>
