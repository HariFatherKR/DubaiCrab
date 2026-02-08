<script lang="ts">
	import { page } from '$app/stores';
	
	interface NavItem {
		href: string;
		icon: string;
		label: string;
	}
	
	const navItems: NavItem[] = [
		{ href: '/dashboard', icon: '📊', label: 'Dashboard' },
		{ href: '/', icon: '💬', label: '채팅' },
		{ href: '/settings', icon: '⚙️', label: '설정' }
	];
	
	let currentPath = $derived($page.url.pathname);
</script>

<aside class="sidebar">
	<div class="sidebar-header">
		<div class="logo">
			<span class="logo-icon">🦞</span>
			<span class="logo-text">OpenKlaw</span>
		</div>
	</div>
	
	<nav class="sidebar-nav">
		{#each navItems as item}
			<a 
				href={item.href} 
				class="nav-item" 
				class:active={currentPath === item.href || (item.href === '/' && currentPath === '/')}
			>
				<span class="nav-icon">{item.icon}</span>
				<span class="nav-label">{item.label}</span>
			</a>
		{/each}
	</nav>
	
	<div class="sidebar-footer">
		<div class="version-badge">
			<span>v0.1.0</span>
		</div>
	</div>
</aside>

<style>
	.sidebar {
		width: 220px;
		min-height: 100vh;
		/* Light glassmorphism */
		background: rgba(255, 255, 255, 0.5);
		backdrop-filter: blur(20px);
		-webkit-backdrop-filter: blur(20px);
		border-right: 1px solid rgba(255, 255, 255, 0.8);
		display: flex;
		flex-direction: column;
		padding: 1.5rem 1rem;
		box-shadow: 4px 0 24px rgba(0, 0, 0, 0.04);
	}
	
	.sidebar-header {
		margin-bottom: 2rem;
	}
	
	.logo {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.5rem;
	}
	
	.logo-icon {
		font-size: 2rem;
	}
	
	.logo-text {
		font-size: 1.25rem;
		font-weight: 700;
		background: linear-gradient(135deg, #14b8a6, #06b6d4);
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		background-clip: text;
	}
	
	.sidebar-nav {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		flex: 1;
	}
	
	.nav-item {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.875rem 1rem;
		border-radius: 12px;
		text-decoration: none;
		color: #64748b;
		font-size: 0.95rem;
		font-weight: 500;
		transition: all 0.2s ease;
	}
	
	.nav-item:hover {
		background: rgba(255, 255, 255, 0.6);
		color: #1e293b;
	}
	
	.nav-item.active {
		background: linear-gradient(135deg, rgba(20, 184, 166, 0.15), rgba(6, 182, 212, 0.1));
		color: #0d9488;
		box-shadow: 0 2px 8px rgba(20, 184, 166, 0.15);
		border: 1px solid rgba(20, 184, 166, 0.2);
	}
	
	.nav-icon {
		font-size: 1.25rem;
	}
	
	.nav-label {
		letter-spacing: 0.01em;
	}
	
	.sidebar-footer {
		padding-top: 1rem;
		border-top: 1px solid rgba(0, 0, 0, 0.06);
	}
	
	.version-badge {
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 0.5rem;
		background: rgba(255, 255, 255, 0.5);
		border-radius: 8px;
		font-size: 0.8rem;
		color: #94a3b8;
	}
</style>
