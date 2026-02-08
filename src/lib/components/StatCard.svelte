<script lang="ts">
	interface Props {
		icon: string;
		label: string;
		value: string | number;
		iconBg?: string;
	}
	
	let { icon, label, value, iconBg = 'bg-teal' }: Props = $props();
</script>

<div class="stat-card">
	<div class="stat-icon {iconBg}">
		{icon}
	</div>
	<div class="stat-info">
		<span class="stat-label">{label}</span>
		<span class="stat-value">{value}</span>
	</div>
</div>

<style>
	.stat-card {
		/* Light glassmorphism */
		background: rgba(255, 255, 255, 0.6);
		backdrop-filter: blur(20px);
		-webkit-backdrop-filter: blur(20px);
		border: 1px solid rgba(255, 255, 255, 0.8);
		border-radius: 16px;
		padding: 1.25rem;
		display: flex;
		align-items: center;
		gap: 1rem;
		box-shadow: 0 4px 16px rgba(0, 0, 0, 0.04);
	}
	
	.stat-icon {
		width: 48px;
		height: 48px;
		border-radius: 12px;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.5rem;
	}
	
	/* Light theme backgrounds */
	.stat-icon.bg-teal {
		background: linear-gradient(135deg, rgba(20, 184, 166, 0.15), rgba(6, 182, 212, 0.1));
	}
	
	.stat-icon.bg-blue {
		background: linear-gradient(135deg, rgba(59, 130, 246, 0.15), rgba(99, 102, 241, 0.1));
	}
	
	.stat-icon.bg-purple {
		background: linear-gradient(135deg, rgba(168, 85, 247, 0.15), rgba(192, 132, 252, 0.1));
	}
	
	.stat-icon.bg-amber {
		background: linear-gradient(135deg, rgba(245, 158, 11, 0.15), rgba(251, 191, 36, 0.1));
	}
	
	.stat-icon.bg-green {
		background: linear-gradient(135deg, rgba(34, 197, 94, 0.15), rgba(74, 222, 128, 0.1));
	}
	
	.stat-icon.bg-rose {
		background: linear-gradient(135deg, rgba(244, 63, 94, 0.15), rgba(251, 113, 133, 0.1));
	}
	
	.stat-info {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}
	
	.stat-label {
		font-size: 0.8rem;
		color: #64748b;
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}
	
	.stat-value {
		font-size: 1.5rem;
		font-weight: 700;
		color: #1e293b;
	}
</style>
