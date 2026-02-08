<script lang="ts">
	interface Props {
		title?: string;
		icon?: string;
		class?: string;
		fullWidth?: boolean;
		children: import('svelte').Snippet;
		headerAction?: import('svelte').Snippet;
	}
	
	let { title, icon, class: className = '', fullWidth = false, children, headerAction }: Props = $props();
</script>

<div class="glass-card {className}" class:full-width={fullWidth}>
	{#if title}
		<div class="card-header">
			<div class="card-title">
				{#if icon}
					<span class="card-icon">{icon}</span>
				{/if}
				<h3>{title}</h3>
			</div>
			{#if headerAction}
				<div class="card-action">
					{@render headerAction()}
				</div>
			{/if}
		</div>
	{/if}
	<div class="card-content">
		{@render children()}
	</div>
</div>

<style>
	.glass-card {
		/* Light glassmorphism - CoachPro style */
		background: rgba(255, 255, 255, 0.6);
		backdrop-filter: blur(20px);
		-webkit-backdrop-filter: blur(20px);
		border: 1px solid rgba(255, 255, 255, 0.8);
		border-radius: 20px;
		padding: 1.5rem;
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.06);
	}
	
	.glass-card.full-width {
		grid-column: 1 / -1;
	}
	
	.card-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1.25rem;
	}
	
	.card-title {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}
	
	.card-icon {
		font-size: 1.25rem;
	}
	
	.card-title h3 {
		font-size: 1rem;
		font-weight: 600;
		color: #1e293b;
		margin: 0;
	}
	
	.card-action {
		font-size: 0.85rem;
		color: #14b8a6;
	}
	
	.card-action :global(a) {
		color: inherit;
		text-decoration: none;
	}
	
	.card-action :global(a:hover) {
		color: #0d9488;
	}
	
	.card-content {
		color: #475569;
	}
</style>
