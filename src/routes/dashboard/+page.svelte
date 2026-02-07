<script lang="ts">
	import { onMount } from 'svelte';
	import { loadStats, type Stats } from '$lib/stores/stats-store';
	
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
		{ id: 'F08', name: '엑셀/CSV 분석', priority: 'P1', status: 'Planned', effort: 1 },
		{ id: 'F09', name: '보고서 템플릿', priority: 'P1', status: 'Planned', effort: 0.5 },
		{ id: 'F10', name: '빠른 작업 버튼', priority: 'P1', status: 'Planned', effort: 0.5 },
		{ id: 'F11', name: '전역 단축키', priority: 'P1', status: 'Planned', effort: 0.5 },
		{ id: 'F12', name: '다크 모드', priority: 'P1', status: 'Done', effort: 0.3 }
	];
	
	const decisions = [
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
		// localStorage에서 통계 로드
		stats = loadStats();
	});
	
	// 진행률 계산
	const doneCount = $derived(features.filter(f => f.status === 'Done').length);
	const totalCount = features.length;
	const progressPercent = $derived(Math.round((doneCount / totalCount) * 100));
	
	function getPriorityColor(priority: string): string {
		switch (priority) {
			case 'P0': return 'bg-red-500';
			case 'P1': return 'bg-yellow-500';
			case 'P2': return 'bg-blue-500';
			case 'P3': return 'bg-gray-500';
			default: return 'bg-gray-500';
		}
	}
	
	function getStatusColor(status: string): string {
		switch (status) {
			case 'Done': return 'text-green-400';
			case 'Planned': return 'text-yellow-400';
			case 'Future': return 'text-blue-400';
			case 'Blocked': return 'text-red-400';
			default: return 'text-gray-400';
		}
	}
	
	function getDifficultyWidth(difficulty: string): string {
		switch (difficulty) {
			case 'Low': return 'w-1/4';
			case 'Medium': return 'w-2/4';
			case 'High': return 'w-3/4';
			case 'Critical': return 'w-full';
			default: return 'w-1/4';
		}
	}
</script>

<svelte:head>
	<title>OpenKlaw - 대시보드</title>
</svelte:head>

<main class="dashboard">
	<header class="dashboard-header">
		<div class="header-left">
			<a href="/" class="back-link">← 채팅으로</a>
		</div>
		<div class="header-center">
			<span class="text-2xl">📊</span>
			<h1>PM 대시보드</h1>
		</div>
		<div class="header-right">
			<span class="version">v0.1.0</span>
		</div>
	</header>
	
	<div class="dashboard-content">
		<!-- 프로젝트 진행 상황 -->
		<section class="card progress-card">
			<h2>📈 프로젝트 진행 상황</h2>
			<div class="progress-bar-container">
				<div class="progress-bar" style="width: {progressPercent}%"></div>
			</div>
			<div class="progress-text">
				<span>{doneCount} / {totalCount} 기능 완료</span>
				<span class="progress-percent">{progressPercent}%</span>
			</div>
			<div class="milestone">
				<span class="milestone-label">현재 Phase:</span>
				<span class="milestone-value">Phase 6 - 설치 패키지</span>
			</div>
		</section>
		
		<!-- 사용 통계 -->
		<section class="card stats-card">
			<h2>📊 사용 통계</h2>
			<div class="stats-grid">
				<div class="stat-item">
					<span class="stat-value">{stats.totalChats}</span>
					<span class="stat-label">총 대화 수</span>
				</div>
				<div class="stat-item">
					<span class="stat-value">{stats.totalMessages}</span>
					<span class="stat-label">총 메시지</span>
				</div>
				<div class="stat-item">
					<span class="stat-value">{stats.hwpProcessed}</span>
					<span class="stat-label">HWP 처리</span>
				</div>
				<div class="stat-item">
					<span class="stat-value">{stats.emailsGenerated}</span>
					<span class="stat-label">이메일 생성</span>
				</div>
			</div>
		</section>
		
		<!-- 연동 도구 우선순위 -->
		<section class="card">
			<h2>🔌 연동 도구 우선순위</h2>
			<div class="tools-table">
				<div class="table-header">
					<span>도구</span>
					<span>우선순위</span>
					<span>활용도</span>
					<span>난이도</span>
					<span>상태</span>
				</div>
				{#each integrationTools as tool}
					<div class="table-row">
						<span class="tool-name">{tool.name}</span>
						<span class="tool-priority">
							<span class="priority-badge {getPriorityColor(tool.priority)}">{tool.priority}</span>
						</span>
						<span class="tool-usage">
							<div class="usage-bar-bg">
								<div class="usage-bar" style="width: {tool.usage}%"></div>
							</div>
							<span class="usage-text">{tool.usage}%</span>
						</span>
						<span class="tool-difficulty">
							<div class="difficulty-bar-bg">
								<div class="difficulty-bar {getDifficultyWidth(tool.difficulty)}"></div>
							</div>
						</span>
						<span class="tool-status {getStatusColor(tool.status)}">{tool.status}</span>
					</div>
				{/each}
			</div>
		</section>
		
		<!-- 기능 목록 -->
		<section class="card">
			<h2>✅ 기능 목록</h2>
			<div class="features-list">
				{#each features as feature}
					<div class="feature-item">
						<span class="feature-id">{feature.id}</span>
						<span class="feature-name">{feature.name}</span>
						<span class="priority-badge {getPriorityColor(feature.priority)}">{feature.priority}</span>
						<span class="feature-effort">{feature.effort}d</span>
						<span class="feature-status {getStatusColor(feature.status)}">
							{#if feature.status === 'Done'}✓{:else if feature.status === 'Planned'}○{:else}◇{/if}
						</span>
					</div>
				{/each}
			</div>
		</section>
		
		<!-- PM 의사결정 로그 -->
		<section class="card decisions-card">
			<h2>📝 PM 의사결정 로그</h2>
			<div class="decisions-list">
				{#each decisions as decision}
					<div class="decision-item">
						<span class="decision-date">{decision.date}</span>
						<span class="decision-text">{decision.decision}</span>
						<span class="decision-reason">{decision.reason}</span>
					</div>
				{/each}
			</div>
		</section>
		
		<!-- 우선순위 매트릭스 시각화 -->
		<section class="card matrix-card">
			<h2>📊 우선순위 매트릭스</h2>
			<div class="matrix">
				<div class="matrix-y-label">수요 ↑</div>
				<div class="matrix-x-label">구현 용이성 →</div>
				<div class="matrix-grid">
					<div class="matrix-quadrant q1">
						<span class="quadrant-label">높은 수요 / 쉬운 구현</span>
						<div class="matrix-item p0">이메일</div>
						<div class="matrix-item p0">HWP</div>
					</div>
					<div class="matrix-quadrant q2">
						<span class="quadrant-label">높은 수요 / 어려운 구현</span>
						<div class="matrix-item p1">엑셀</div>
					</div>
					<div class="matrix-quadrant q3">
						<span class="quadrant-label">낮은 수요 / 쉬운 구현</span>
					</div>
					<div class="matrix-quadrant q4">
						<span class="quadrant-label">낮은 수요 / 어려운 구현</span>
						<div class="matrix-item p2">캘린더</div>
						<div class="matrix-item p3">카카오톡</div>
					</div>
				</div>
			</div>
		</section>
	</div>
</main>

<style>
	.dashboard {
		min-height: 100vh;
		background: var(--color-bg);
		color: var(--color-text);
	}
	
	.dashboard-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 1rem 2rem;
		background: var(--color-surface);
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
	}
	
	.header-center {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}
	
	.header-center h1 {
		font-size: 1.25rem;
		font-weight: 600;
	}
	
	.back-link {
		color: var(--color-primary);
		text-decoration: none;
		font-size: 0.9rem;
	}
	
	.back-link:hover {
		text-decoration: underline;
	}
	
	.version {
		color: var(--color-text-muted);
		font-size: 0.8rem;
	}
	
	.dashboard-content {
		max-width: 1200px;
		margin: 0 auto;
		padding: 2rem;
		display: grid;
		gap: 1.5rem;
		grid-template-columns: repeat(2, 1fr);
	}
	
	.card {
		background: var(--color-surface);
		border-radius: 1rem;
		padding: 1.5rem;
	}
	
	.card h2 {
		font-size: 1.1rem;
		margin-bottom: 1rem;
		color: var(--color-text);
	}
	
	.progress-card, .decisions-card, .matrix-card {
		grid-column: span 2;
	}
	
	/* 진행 바 */
	.progress-bar-container {
		height: 1.5rem;
		background: rgba(255, 255, 255, 0.1);
		border-radius: 0.75rem;
		overflow: hidden;
		margin-bottom: 0.75rem;
	}
	
	.progress-bar {
		height: 100%;
		background: linear-gradient(90deg, var(--color-primary), #22c55e);
		border-radius: 0.75rem;
		transition: width 0.5s ease;
	}
	
	.progress-text {
		display: flex;
		justify-content: space-between;
		font-size: 0.9rem;
		color: var(--color-text-muted);
	}
	
	.progress-percent {
		font-weight: 600;
		color: var(--color-primary);
	}
	
	.milestone {
		margin-top: 1rem;
		padding-top: 1rem;
		border-top: 1px solid rgba(255, 255, 255, 0.1);
	}
	
	.milestone-label {
		color: var(--color-text-muted);
		margin-right: 0.5rem;
	}
	
	.milestone-value {
		color: var(--color-primary);
		font-weight: 500;
	}
	
	/* 통계 그리드 */
	.stats-grid {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 1rem;
	}
	
	.stat-item {
		text-align: center;
		padding: 1rem;
		background: rgba(255, 255, 255, 0.05);
		border-radius: 0.75rem;
	}
	
	.stat-value {
		display: block;
		font-size: 2rem;
		font-weight: 700;
		color: var(--color-primary);
	}
	
	.stat-label {
		font-size: 0.85rem;
		color: var(--color-text-muted);
	}
	
	/* 도구 테이블 */
	.tools-table {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}
	
	.table-header, .table-row {
		display: grid;
		grid-template-columns: 1fr 0.8fr 1.5fr 1fr 0.8fr;
		gap: 0.75rem;
		align-items: center;
		padding: 0.75rem;
	}
	
	.table-header {
		font-size: 0.8rem;
		color: var(--color-text-muted);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}
	
	.table-row {
		background: rgba(255, 255, 255, 0.03);
		border-radius: 0.5rem;
	}
	
	.tool-name {
		font-weight: 500;
	}
	
	.priority-badge {
		display: inline-block;
		padding: 0.2rem 0.5rem;
		border-radius: 0.25rem;
		font-size: 0.75rem;
		font-weight: 600;
		color: white;
	}
	
	.tool-usage {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}
	
	.usage-bar-bg, .difficulty-bar-bg {
		flex: 1;
		height: 0.5rem;
		background: rgba(255, 255, 255, 0.1);
		border-radius: 0.25rem;
		overflow: hidden;
	}
	
	.usage-bar {
		height: 100%;
		background: var(--color-primary);
		border-radius: 0.25rem;
	}
	
	.usage-text {
		font-size: 0.8rem;
		color: var(--color-text-muted);
		min-width: 35px;
	}
	
	.difficulty-bar {
		height: 100%;
		background: linear-gradient(90deg, #22c55e, #eab308, #ef4444);
		border-radius: 0.25rem;
	}
	
	/* 기능 목록 */
	.features-list {
		display: flex;
		flex-direction: column;
		gap: 0.4rem;
		max-height: 400px;
		overflow-y: auto;
	}
	
	.feature-item {
		display: grid;
		grid-template-columns: 3rem 1fr auto 2.5rem 1.5rem;
		gap: 0.5rem;
		align-items: center;
		padding: 0.5rem 0.75rem;
		background: rgba(255, 255, 255, 0.03);
		border-radius: 0.4rem;
		font-size: 0.9rem;
	}
	
	.feature-id {
		font-family: monospace;
		color: var(--color-text-muted);
		font-size: 0.8rem;
	}
	
	.feature-effort {
		color: var(--color-text-muted);
		font-size: 0.8rem;
		text-align: right;
	}
	
	.feature-status {
		font-size: 1rem;
		text-align: center;
	}
	
	/* 의사결정 목록 */
	.decisions-list {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}
	
	.decision-item {
		display: grid;
		grid-template-columns: 6rem 1fr 1fr;
		gap: 1rem;
		padding: 0.75rem;
		background: rgba(255, 255, 255, 0.03);
		border-radius: 0.5rem;
		font-size: 0.9rem;
	}
	
	.decision-date {
		color: var(--color-text-muted);
		font-family: monospace;
		font-size: 0.8rem;
	}
	
	.decision-text {
		font-weight: 500;
	}
	
	.decision-reason {
		color: var(--color-text-muted);
		font-size: 0.85rem;
	}
	
	/* 매트릭스 */
	.matrix {
		position: relative;
		padding: 2rem 0 0 2rem;
	}
	
	.matrix-y-label {
		position: absolute;
		left: 0;
		top: 50%;
		transform: rotate(-90deg) translateX(-50%);
		font-size: 0.8rem;
		color: var(--color-text-muted);
	}
	
	.matrix-x-label {
		text-align: center;
		font-size: 0.8rem;
		color: var(--color-text-muted);
		margin-top: 0.5rem;
	}
	
	.matrix-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		grid-template-rows: 1fr 1fr;
		gap: 0.5rem;
		min-height: 200px;
	}
	
	.matrix-quadrant {
		background: rgba(255, 255, 255, 0.03);
		border-radius: 0.5rem;
		padding: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}
	
	.quadrant-label {
		font-size: 0.7rem;
		color: var(--color-text-muted);
		margin-bottom: 0.5rem;
	}
	
	.matrix-item {
		display: inline-block;
		padding: 0.25rem 0.5rem;
		border-radius: 0.25rem;
		font-size: 0.8rem;
		font-weight: 500;
	}
	
	.matrix-item.p0 {
		background: rgba(239, 68, 68, 0.2);
		color: #ef4444;
	}
	
	.matrix-item.p1 {
		background: rgba(234, 179, 8, 0.2);
		color: #eab308;
	}
	
	.matrix-item.p2 {
		background: rgba(59, 130, 246, 0.2);
		color: #3b82f6;
	}
	
	.matrix-item.p3 {
		background: rgba(107, 114, 128, 0.2);
		color: #6b7280;
	}
	
	/* 색상 유틸리티 */
	.bg-red-500 { background-color: #ef4444; }
	.bg-yellow-500 { background-color: #eab308; }
	.bg-blue-500 { background-color: #3b82f6; }
	.bg-gray-500 { background-color: #6b7280; }
	
	.text-green-400 { color: #4ade80; }
	.text-yellow-400 { color: #facc15; }
	.text-blue-400 { color: #60a5fa; }
	.text-red-400 { color: #f87171; }
	.text-gray-400 { color: #9ca3af; }
	
	.w-1\/4 { width: 25%; }
	.w-2\/4 { width: 50%; }
	.w-3\/4 { width: 75%; }
	.w-full { width: 100%; }
	
	/* 반응형 */
	@media (max-width: 768px) {
		.dashboard-content {
			grid-template-columns: 1fr;
			padding: 1rem;
		}
		
		.progress-card, .decisions-card, .matrix-card {
			grid-column: span 1;
		}
		
		.table-header, .table-row {
			grid-template-columns: 1fr 0.6fr 1fr 0.6fr;
		}
		
		.tool-usage {
			display: none;
		}
		
		.decision-item {
			grid-template-columns: 1fr;
			gap: 0.25rem;
		}
	}
</style>
