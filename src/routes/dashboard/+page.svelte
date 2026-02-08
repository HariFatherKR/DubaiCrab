<script lang="ts">
	import { onMount } from 'svelte';
	import { loadStats, type Stats } from '$lib/stores/stats-store';
	import Header from '$lib/components/Header.svelte';
	import GlassCard from '$lib/components/GlassCard.svelte';
	import StatCard from '$lib/components/StatCard.svelte';
	
	// PM 의사결정 데이터
	const integrationTools = [
		{ name: 'HWP', priority: 'P0', usage: 100, difficulty: 'Medium', status: 'Done' },
		{ name: '이메일', priority: 'P0', usage: 95, difficulty: 'Low', status: 'Done' },
		{ name: '엑셀', priority: 'P1', usage: 80, difficulty: 'Medium', status: 'Planned' },
		{ name: '캘린더', priority: 'P2', usage: 70, difficulty: 'High', status: 'Future' },
		{ name: '카카오톡', priority: 'P3', usage: 90, difficulty: 'Critical', status: 'Blocked' }
	];
	
	const features = [
		{ id: 'F01', name: '원클릭 설치', priority: 'P0', status: 'Done', effort: 1 },
		{ id: 'F02', name: 'Ollama 자동 설정', priority: 'P0', status: 'Done', effort: 0.5 },
		{ id: 'F03', name: '기본 채팅', priority: 'P0', status: 'Done', effort: 1 },
		{ id: 'F04', name: 'HWP 요약', priority: 'P0', status: 'Done', effort: 1 },
		{ id: 'F05', name: '이메일 작성', priority: 'P0', status: 'Done', effort: 0.5 },
		{ id: 'F06', name: '대화 저장', priority: 'P0', status: 'Done', effort: 0.5 },
		{ id: 'F07', name: '시스템 트레이', priority: 'P0', status: 'Done', effort: 0.5 },
		{ id: 'F08', name: '엑셀/CSV 분석', priority: 'P1', status: 'Done', effort: 1 },
		{ id: 'F09', name: '보고서 템플릿', priority: 'P1', status: 'Planned', effort: 0.5 },
		{ id: 'F10', name: '빠른 작업 버튼', priority: 'P1', status: 'Planned', effort: 0.5 },
		{ id: 'F11', name: '전역 단축키', priority: 'P1', status: 'Planned', effort: 0.5 },
		{ id: 'F12', name: '다크 모드', priority: 'P1', status: 'Done', effort: 0.3 }
	];
	
	const decisions = [
		{ date: '2026-02-08', decision: '엑셀/CSV 분석 구현 완료', reason: 'SheetJS + PapaParse 채택' },
		{ date: '2026-02-07', decision: 'HWP P0 확정', reason: '유일한 차별점, 경쟁사 전무' },
		{ date: '2026-02-07', decision: '이메일 P0 확정', reason: '높은 수요, 낮은 복잡도' },
		{ date: '2026-02-07', decision: '엑셀 P1로 조정', reason: 'MVP 범위 축소, Phase 2' },
		{ date: '2026-02-07', decision: '캘린더 P2로 결정', reason: 'OAuth 복잡성, 로컬 철학' },
		{ date: '2026-02-07', decision: '카카오톡 P3 (보류)', reason: '공식 API 없음, 법적 리스크' },
		{ date: '2026-02-07', decision: 'Qwen2.5-3B 기본 모델', reason: '한국어 성능/크기 밸런스' },
		{ date: '2026-02-07', decision: 'Tauri 선택', reason: 'Electron 대비 10x 경량' }
	];
	
	// 통계 (localStorage에서 가져오기)
	let stats = $state<Stats>({
		totalChats: 0,
		totalMessages: 0,
		hwpProcessed: 0,
		emailsGenerated: 0,
		lastUsed: ''
	});
	
	onMount(() => {
		stats = loadStats();
	});
	
	// 진행률 계산
	const doneCount = $derived(features.filter(f => f.status === 'Done').length);
	const totalCount = features.length;
	const progressPercent = $derived(Math.round((doneCount / totalCount) * 100));
	
	function getPriorityClass(priority: string): string {
		switch (priority) {
			case 'P0': return 'priority-p0';
			case 'P1': return 'priority-p1';
			case 'P2': return 'priority-p2';
			case 'P3': return 'priority-p3';
			default: return 'priority-p3';
		}
	}
	
	function getStatusClass(status: string): string {
		switch (status) {
			case 'Done': return 'status-done';
			case 'Planned': return 'status-planned';
			case 'Future': return 'status-future';
			case 'Blocked': return 'status-blocked';
			default: return '';
		}
	}
	
	function getDifficultyPercent(difficulty: string): number {
		switch (difficulty) {
			case 'Low': return 25;
			case 'Medium': return 50;
			case 'High': return 75;
			case 'Critical': return 100;
			default: return 25;
		}
	}
</script>

<svelte:head>
	<title>OpenKlaw - Dashboard</title>
</svelte:head>

<Header title="Dashboard" subtitle="프로젝트 현황 및 통계" userName="PM" />

<div class="dashboard-content">
	<!-- 프로젝트 진행 상황 -->
	<GlassCard title="프로젝트 진행 상황" icon="📈" fullWidth>
		<div class="progress-section">
			<div class="progress-bar-container">
				<div class="progress-bar" style="width: {progressPercent}%"></div>
			</div>
			<div class="progress-info">
				<span class="progress-text">{doneCount} / {totalCount} 기능 완료</span>
				<span class="progress-percent">{progressPercent}%</span>
			</div>
			<div class="milestone-badge">
				<span class="milestone-label">현재 Phase</span>
				<span class="milestone-value">Phase 10 - P1 기능 구현</span>
			</div>
		</div>
	</GlassCard>
	
	<!-- 사용 통계 카드 그리드 -->
	<div class="stats-grid">
		<StatCard icon="💬" label="총 대화" value={stats.totalChats} iconBg="bg-teal" />
		<StatCard icon="📝" label="메시지" value={stats.totalMessages} iconBg="bg-blue" />
		<StatCard icon="📄" label="HWP 처리" value={stats.hwpProcessed} iconBg="bg-purple" />
		<StatCard icon="📧" label="이메일 생성" value={stats.emailsGenerated} iconBg="bg-amber" />
	</div>
	
	<!-- 연동 도구 우선순위 -->
	<GlassCard title="연동 도구 우선순위" icon="🔌">
		{#snippet headerAction()}
			<a href="#all">View all</a>
		{/snippet}
		
		<table class="glass-table">
			<thead>
				<tr>
					<th>도구</th>
					<th>우선순위</th>
					<th>활용도</th>
					<th>난이도</th>
					<th>상태</th>
				</tr>
			</thead>
			<tbody>
				{#each integrationTools as tool}
					<tr>
						<td class="tool-name">{tool.name}</td>
						<td>
							<span class="priority-badge {getPriorityClass(tool.priority)}">{tool.priority}</span>
						</td>
						<td>
							<div class="bar-container">
								<div class="bar usage-bar" style="width: {tool.usage}%"></div>
							</div>
						</td>
						<td>
							<div class="bar-container">
								<div class="bar difficulty-bar" style="width: {getDifficultyPercent(tool.difficulty)}%"></div>
							</div>
						</td>
						<td class={getStatusClass(tool.status)}>{tool.status}</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</GlassCard>
	
	<!-- 기능 목록 -->
	<GlassCard title="기능 목록" icon="✅">
		<div class="features-list">
			{#each features as feature}
				<div class="feature-item">
					<span class="feature-id">{feature.id}</span>
					<span class="feature-name">{feature.name}</span>
					<span class="priority-badge small {getPriorityClass(feature.priority)}">{feature.priority}</span>
					<span class="feature-effort">{feature.effort}d</span>
					<span class="feature-status {getStatusClass(feature.status)}">
						{#if feature.status === 'Done'}✓{:else if feature.status === 'Planned'}○{:else}◇{/if}
					</span>
				</div>
			{/each}
		</div>
	</GlassCard>
	
	<!-- PM 의사결정 로그 -->
	<GlassCard title="PM 의사결정 로그" icon="📝" fullWidth>
		<div class="decisions-list">
			{#each decisions as decision}
				<div class="decision-item">
					<span class="decision-date">{decision.date}</span>
					<span class="decision-text">{decision.decision}</span>
					<span class="decision-reason">{decision.reason}</span>
				</div>
			{/each}
		</div>
	</GlassCard>
	
	<!-- CTA 카드 -->
	<div class="cta-card">
		<div class="cta-content">
			<div class="cta-decoration">
				<div class="cta-icon">🦞</div>
			</div>
			<div class="cta-text">
				<span class="cta-label">Don't forget!</span>
				<h3>채팅으로 테스트하기</h3>
			</div>
		</div>
		<a href="/" class="cta-button">Go to Chat</a>
	</div>
</div>

<style>
	.dashboard-content {
		padding: 0 2rem 2rem;
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 1.5rem;
		max-width: 1400px;
	}
	
	/* Stats Grid */
	.stats-grid {
		grid-column: 1 / -1;
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 1rem;
	}
	
	/* Progress Section */
	.progress-section {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}
	
	.progress-bar-container {
		height: 12px;
		background: rgba(0, 0, 0, 0.06);
		border-radius: 6px;
		overflow: hidden;
	}
	
	.progress-bar {
		height: 100%;
		background: linear-gradient(90deg, #14b8a6, #06b6d4);
		border-radius: 6px;
		transition: width 0.5s ease;
	}
	
	.progress-info {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}
	
	.progress-text {
		font-size: 0.9rem;
		color: #64748b;
	}
	
	.progress-percent {
		font-size: 1.25rem;
		font-weight: 700;
		color: #0d9488;
	}
	
	.milestone-badge {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.75rem 1rem;
		background: rgba(20, 184, 166, 0.08);
		border: 1px solid rgba(20, 184, 166, 0.15);
		border-radius: 12px;
		width: fit-content;
	}
	
	.milestone-label {
		font-size: 0.8rem;
		color: #64748b;
	}
	
	.milestone-value {
		font-size: 0.9rem;
		font-weight: 600;
		color: #0d9488;
	}
	
	/* Table styles - Light theme */
	.glass-table {
		width: 100%;
		border-collapse: separate;
		border-spacing: 0 0.5rem;
	}
	
	.glass-table th {
		font-size: 0.7rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: #94a3b8;
		text-align: left;
		padding: 0.5rem 1rem;
	}
	
	.glass-table td {
		background: rgba(255, 255, 255, 0.5);
		padding: 0.875rem 1rem;
		font-size: 0.9rem;
		color: #475569;
	}
	
	.glass-table tr td:first-child {
		border-radius: 12px 0 0 12px;
	}
	
	.glass-table tr td:last-child {
		border-radius: 0 12px 12px 0;
	}
	
	.tool-name {
		font-weight: 500;
		color: #1e293b;
	}
	
	/* Priority badges - Light theme */
	.priority-badge {
		display: inline-block;
		padding: 0.25rem 0.5rem;
		border-radius: 6px;
		font-size: 0.75rem;
		font-weight: 600;
	}
	
	.priority-badge.small {
		padding: 0.15rem 0.4rem;
		font-size: 0.7rem;
	}
	
	.priority-p0 { background: #fef2f2; color: #dc2626; border: 1px solid #fecaca; }
	.priority-p1 { background: #fffbeb; color: #d97706; border: 1px solid #fde68a; }
	.priority-p2 { background: #eff6ff; color: #2563eb; border: 1px solid #bfdbfe; }
	.priority-p3 { background: #f3f4f6; color: #6b7280; border: 1px solid #e5e7eb; }
	
	/* Status colors - Light theme */
	.status-done { color: #059669; font-weight: 600; }
	.status-planned { color: #d97706; font-weight: 600; }
	.status-future { color: #2563eb; font-weight: 600; }
	.status-blocked { color: #dc2626; font-weight: 600; }
	
	/* Bar containers */
	.bar-container {
		height: 6px;
		background: rgba(0, 0, 0, 0.06);
		border-radius: 3px;
		overflow: hidden;
		min-width: 60px;
	}
	
	.bar {
		height: 100%;
		border-radius: 3px;
	}
	
	.usage-bar {
		background: linear-gradient(90deg, #14b8a6, #06b6d4);
	}
	
	.difficulty-bar {
		background: linear-gradient(90deg, #22c55e, #eab308, #ef4444);
	}
	
	/* Features list */
	.features-list {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		max-height: 350px;
		overflow-y: auto;
	}
	
	.feature-item {
		display: grid;
		grid-template-columns: 3rem 1fr auto 2.5rem 1.5rem;
		gap: 0.75rem;
		align-items: center;
		padding: 0.625rem 0.875rem;
		background: rgba(255, 255, 255, 0.5);
		border-radius: 10px;
		font-size: 0.875rem;
	}
	
	.feature-id {
		font-family: 'SF Mono', Monaco, monospace;
		color: #94a3b8;
		font-size: 0.75rem;
	}
	
	.feature-name {
		color: #1e293b;
	}
	
	.feature-effort {
		color: #94a3b8;
		font-size: 0.75rem;
		text-align: right;
	}
	
	.feature-status {
		font-size: 1rem;
		text-align: center;
	}
	
	/* Decisions list */
	.decisions-list {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
		gap: 0.75rem;
	}
	
	.decision-item {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		padding: 1rem;
		background: rgba(255, 255, 255, 0.5);
		border-radius: 12px;
	}
	
	.decision-date {
		font-family: 'SF Mono', Monaco, monospace;
		font-size: 0.75rem;
		color: #94a3b8;
	}
	
	.decision-text {
		font-weight: 500;
		color: #1e293b;
		font-size: 0.9rem;
	}
	
	.decision-reason {
		font-size: 0.8rem;
		color: #64748b;
	}
	
	/* CTA Card - Light theme with accent */
	.cta-card {
		grid-column: 1 / -1;
		background: linear-gradient(135deg, rgba(20, 184, 166, 0.12), rgba(6, 182, 212, 0.08));
		backdrop-filter: blur(20px);
		-webkit-backdrop-filter: blur(20px);
		border: 1px solid rgba(20, 184, 166, 0.2);
		border-radius: 20px;
		padding: 1.5rem 2rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
		box-shadow: 0 8px 32px rgba(20, 184, 166, 0.08);
	}
	
	.cta-content {
		display: flex;
		align-items: center;
		gap: 1.5rem;
	}
	
	.cta-decoration {
		position: relative;
	}
	
	.cta-icon {
		font-size: 3rem;
		filter: drop-shadow(0 4px 12px rgba(20, 184, 166, 0.3));
	}
	
	.cta-text {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}
	
	.cta-label {
		font-size: 0.8rem;
		text-transform: uppercase;
		letter-spacing: 0.1em;
		color: #64748b;
	}
	
	.cta-text h3 {
		font-size: 1.25rem;
		font-weight: 600;
		color: #1e293b;
		margin: 0;
	}
	
	.cta-button {
		padding: 0.75rem 1.5rem;
		background: linear-gradient(135deg, #14b8a6, #0d9488);
		border: none;
		border-radius: 12px;
		color: white;
		text-decoration: none;
		font-weight: 500;
		font-size: 0.9rem;
		transition: all 0.2s ease;
		box-shadow: 0 4px 12px rgba(20, 184, 166, 0.25);
	}
	
	.cta-button:hover {
		background: linear-gradient(135deg, #0d9488, #0f766e);
		transform: translateY(-2px);
		box-shadow: 0 6px 20px rgba(20, 184, 166, 0.35);
	}
	
	/* Responsive */
	@media (max-width: 1024px) {
		.stats-grid {
			grid-template-columns: repeat(2, 1fr);
		}
	}
	
	@media (max-width: 768px) {
		.dashboard-content {
			grid-template-columns: 1fr;
			padding: 0 1rem 1rem;
		}
		
		.stats-grid {
			grid-template-columns: 1fr 1fr;
		}
		
		.decisions-list {
			grid-template-columns: 1fr;
		}
		
		.cta-card {
			flex-direction: column;
			gap: 1rem;
			text-align: center;
		}
		
		.cta-content {
			flex-direction: column;
		}
	}
</style>
