<script lang="ts">
	import { 
		settingsStore, 
		AVAILABLE_MODELS, 
		THEME_OPTIONS, 
		ACCENT_COLORS,
		type AppSettings 
	} from '$lib/stores/settings-store';
	
	let settings = $state<AppSettings>($settingsStore);
	let newModel = $state('');
	let showResetConfirm = $state(false);
	
	// Store 구독
	$effect(() => {
		const unsubscribe = settingsStore.subscribe(value => {
			settings = value;
		});
		return unsubscribe;
	});
	
	function updateSetting<K extends keyof AppSettings>(key: K, value: AppSettings[K]) {
		settingsStore.updateSettings({ [key]: value });
	}
	
	function addModel() {
		if (newModel.trim()) {
			settingsStore.addCustomModel(newModel.trim());
			newModel = '';
		}
	}
	
	function removeModel(model: string) {
		settingsStore.removeCustomModel(model);
	}
	
	function resetSettings() {
		settingsStore.resetSettings();
		showResetConfirm = false;
	}
</script>

<div class="settings-page">
	<div class="settings-header">
		<h1>⚙️ 설정</h1>
		<p>OpenKlaw 앱 설정을 관리합니다</p>
	</div>
	
	<div class="settings-content">
		<!-- 모델 설정 -->
		<section class="settings-section">
			<h2>🤖 AI 모델</h2>
			
			<div class="setting-item">
				<label for="model-select">기본 모델</label>
				<select 
					id="model-select" 
					value={settings.model}
					onchange={(e) => updateSetting('model', e.currentTarget.value)}
				>
					{#each AVAILABLE_MODELS as model}
						<option value={model.id}>
							{model.name} ({model.size})
						</option>
					{/each}
					{#each settings.customModels as model}
						<option value={model}>{model} (커스텀)</option>
					{/each}
				</select>
				<span class="setting-hint">Ollama에 설치된 모델을 선택하세요</span>
			</div>
			
			<div class="setting-item">
				<label>커스텀 모델 추가</label>
				<div class="inline-form">
					<input 
						type="text" 
						bind:value={newModel}
						placeholder="예: codellama:7b"
						onkeydown={(e) => e.key === 'Enter' && addModel()}
					/>
					<button onclick={addModel} disabled={!newModel.trim()}>추가</button>
				</div>
				
				{#if settings.customModels.length > 0}
					<div class="custom-models">
						{#each settings.customModels as model}
							<span class="model-tag">
								{model}
								<button onclick={() => removeModel(model)}>✕</button>
							</span>
						{/each}
					</div>
				{/if}
			</div>
		</section>
		
		<!-- 테마 설정 -->
		<section class="settings-section">
			<h2>🎨 테마</h2>
			
			<div class="setting-item">
				<label>테마 모드</label>
				<div class="theme-options">
					{#each THEME_OPTIONS as theme}
						<button 
							class="theme-btn"
							class:active={settings.theme === theme.id}
							onclick={() => updateSetting('theme', theme.id as 'dark' | 'light' | 'system')}
						>
							<span class="theme-icon">{theme.icon}</span>
							<span>{theme.name}</span>
						</button>
					{/each}
				</div>
			</div>
			
			<div class="setting-item">
				<label>액센트 컬러</label>
				<div class="color-options">
					{#each ACCENT_COLORS as color}
						<button 
							class="color-btn"
							class:active={settings.accentColor === color.id}
							style="--color: {color.id}"
							onclick={() => updateSetting('accentColor', color.id)}
							title={color.name}
						></button>
					{/each}
				</div>
			</div>
		</section>
		
		<!-- 단축키 설정 -->
		<section class="settings-section">
			<h2>⌨️ 단축키</h2>
			
			<div class="setting-item">
				<label>앱 열기/숨기기</label>
				<div class="shortcut-display">
					<kbd>⌘</kbd> + <kbd>⇧</kbd> + <kbd>O</kbd>
				</div>
				<span class="setting-hint">전역 단축키 (변경 불가)</span>
			</div>
			
			<div class="setting-item">
				<label>Enter로 전송</label>
				<div class="toggle-switch">
					<input 
						type="checkbox" 
						id="send-on-enter"
						checked={settings.sendOnEnter}
						onchange={(e) => updateSetting('sendOnEnter', e.currentTarget.checked)}
					/>
					<label for="send-on-enter" class="toggle-label"></label>
				</div>
				<span class="setting-hint">
					{settings.sendOnEnter ? 'Enter로 전송, Shift+Enter로 줄바꿈' : 'Cmd+Enter로 전송'}
				</span>
			</div>
		</section>
		
		<!-- 데이터 설정 -->
		<section class="settings-section">
			<h2>💾 데이터</h2>
			
			<div class="setting-item">
				<label>자동 저장</label>
				<div class="toggle-switch">
					<input 
						type="checkbox" 
						id="auto-save"
						checked={settings.autoSave}
						onchange={(e) => updateSetting('autoSave', e.currentTarget.checked)}
					/>
					<label for="auto-save" class="toggle-label"></label>
				</div>
				<span class="setting-hint">대화 내용을 자동으로 저장합니다</span>
			</div>
			
			<div class="setting-item">
				<label>데이터 저장 위치</label>
				<input 
					type="text" 
					value={settings.dataPath}
					onchange={(e) => updateSetting('dataPath', e.currentTarget.value)}
					placeholder="~/.openklaw"
				/>
			</div>
		</section>
		
		<!-- 초기화 -->
		<section class="settings-section danger-zone">
			<h2>⚠️ 초기화</h2>
			
			<div class="setting-item">
				{#if showResetConfirm}
					<div class="confirm-box">
						<p>정말 모든 설정을 초기화하시겠습니까?</p>
						<div class="confirm-actions">
							<button class="cancel-btn" onclick={() => showResetConfirm = false}>취소</button>
							<button class="danger-btn" onclick={resetSettings}>초기화</button>
						</div>
					</div>
				{:else}
					<button class="reset-btn" onclick={() => showResetConfirm = true}>
						설정 초기화
					</button>
				{/if}
			</div>
		</section>
	</div>
	
	<div class="settings-footer">
		<div class="version-info">
			<span>OpenKlaw v0.1.0</span>
			<span>•</span>
			<a href="https://github.com/openklaw/openklaw" target="_blank">GitHub</a>
		</div>
	</div>
</div>

<style>
	.settings-page {
		padding: 2rem;
		max-width: 700px;
		margin: 0 auto;
	}
	
	.settings-header {
		margin-bottom: 2rem;
	}
	
	.settings-header h1 {
		font-size: 1.75rem;
		font-weight: 700;
		margin: 0 0 0.5rem;
		color: rgba(255, 255, 255, 0.95);
	}
	
	.settings-header p {
		color: rgba(255, 255, 255, 0.6);
		margin: 0;
	}
	
	.settings-content {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}
	
	.settings-section {
		background: rgba(255, 255, 255, 0.06);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 16px;
		padding: 1.5rem;
	}
	
	.settings-section h2 {
		font-size: 1.1rem;
		font-weight: 600;
		margin: 0 0 1.25rem;
		color: rgba(255, 255, 255, 0.9);
	}
	
	.setting-item {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		margin-bottom: 1.25rem;
	}
	
	.setting-item:last-child {
		margin-bottom: 0;
	}
	
	.setting-item > label {
		font-size: 0.95rem;
		font-weight: 500;
		color: rgba(255, 255, 255, 0.8);
	}
	
	.setting-hint {
		font-size: 0.8rem;
		color: rgba(255, 255, 255, 0.4);
	}
	
	select, input[type="text"] {
		padding: 0.75rem 1rem;
		background: rgba(255, 255, 255, 0.06);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 10px;
		color: rgba(255, 255, 255, 0.95);
		font-size: 0.95rem;
		font-family: inherit;
	}
	
	select:focus, input:focus {
		outline: none;
		border-color: rgba(20, 184, 166, 0.5);
	}
	
	.inline-form {
		display: flex;
		gap: 0.5rem;
	}
	
	.inline-form input {
		flex: 1;
	}
	
	.inline-form button {
		padding: 0.75rem 1.25rem;
		background: linear-gradient(135deg, #14b8a6, #0d9488);
		border: none;
		border-radius: 10px;
		color: white;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
	}
	
	.inline-form button:hover:not(:disabled) {
		background: linear-gradient(135deg, #0d9488, #0f766e);
	}
	
	.inline-form button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
	
	.custom-models {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
		margin-top: 0.5rem;
	}
	
	.model-tag {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.375rem 0.75rem;
		background: rgba(255, 255, 255, 0.1);
		border-radius: 6px;
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.8);
	}
	
	.model-tag button {
		background: none;
		border: none;
		color: rgba(255, 255, 255, 0.5);
		cursor: pointer;
		padding: 0;
		font-size: 0.8rem;
	}
	
	.model-tag button:hover {
		color: #f87171;
	}
	
	.theme-options {
		display: flex;
		gap: 0.75rem;
	}
	
	.theme-btn {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
		padding: 1rem;
		background: rgba(255, 255, 255, 0.06);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 12px;
		color: rgba(255, 255, 255, 0.7);
		cursor: pointer;
		transition: all 0.2s;
	}
	
	.theme-btn:hover {
		background: rgba(255, 255, 255, 0.1);
	}
	
	.theme-btn.active {
		background: rgba(20, 184, 166, 0.15);
		border-color: rgba(20, 184, 166, 0.4);
		color: #5eead4;
	}
	
	.theme-icon {
		font-size: 1.5rem;
	}
	
	.color-options {
		display: flex;
		gap: 0.75rem;
	}
	
	.color-btn {
		width: 40px;
		height: 40px;
		border-radius: 50%;
		background: var(--color);
		border: 3px solid transparent;
		cursor: pointer;
		transition: all 0.2s;
	}
	
	.color-btn:hover {
		transform: scale(1.1);
	}
	
	.color-btn.active {
		border-color: white;
		box-shadow: 0 0 0 3px rgba(255, 255, 255, 0.2);
	}
	
	.shortcut-display {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}
	
	kbd {
		padding: 0.375rem 0.625rem;
		background: rgba(255, 255, 255, 0.1);
		border: 1px solid rgba(255, 255, 255, 0.15);
		border-radius: 6px;
		font-family: inherit;
		font-size: 0.9rem;
		color: rgba(255, 255, 255, 0.8);
	}
	
	.toggle-switch {
		position: relative;
		width: 52px;
		height: 28px;
	}
	
	.toggle-switch input {
		opacity: 0;
		width: 0;
		height: 0;
	}
	
	.toggle-label {
		position: absolute;
		cursor: pointer;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(255, 255, 255, 0.1);
		border-radius: 14px;
		transition: 0.3s;
	}
	
	.toggle-label::before {
		position: absolute;
		content: "";
		height: 20px;
		width: 20px;
		left: 4px;
		bottom: 4px;
		background: white;
		border-radius: 50%;
		transition: 0.3s;
	}
	
	.toggle-switch input:checked + .toggle-label {
		background: #14b8a6;
	}
	
	.toggle-switch input:checked + .toggle-label::before {
		transform: translateX(24px);
	}
	
	.danger-zone {
		border-color: rgba(239, 68, 68, 0.3);
	}
	
	.danger-zone h2 {
		color: #fca5a5;
	}
	
	.reset-btn {
		padding: 0.75rem 1.5rem;
		background: rgba(239, 68, 68, 0.15);
		border: 1px solid rgba(239, 68, 68, 0.3);
		border-radius: 10px;
		color: #fca5a5;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
	}
	
	.reset-btn:hover {
		background: rgba(239, 68, 68, 0.25);
	}
	
	.confirm-box {
		background: rgba(239, 68, 68, 0.1);
		border: 1px solid rgba(239, 68, 68, 0.2);
		border-radius: 12px;
		padding: 1rem;
	}
	
	.confirm-box p {
		margin: 0 0 1rem;
		color: rgba(255, 255, 255, 0.8);
	}
	
	.confirm-actions {
		display: flex;
		gap: 0.75rem;
	}
	
	.cancel-btn, .danger-btn {
		flex: 1;
		padding: 0.625rem 1rem;
		border-radius: 8px;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
	}
	
	.cancel-btn {
		background: rgba(255, 255, 255, 0.1);
		border: 1px solid rgba(255, 255, 255, 0.1);
		color: rgba(255, 255, 255, 0.8);
	}
	
	.danger-btn {
		background: #ef4444;
		border: none;
		color: white;
	}
	
	.danger-btn:hover {
		background: #dc2626;
	}
	
	.settings-footer {
		margin-top: 2rem;
		padding-top: 1.5rem;
		border-top: 1px solid rgba(255, 255, 255, 0.1);
		text-align: center;
	}
	
	.version-info {
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 0.75rem;
		color: rgba(255, 255, 255, 0.4);
		font-size: 0.85rem;
	}
	
	.version-info a {
		color: #5eead4;
		text-decoration: none;
	}
	
	.version-info a:hover {
		text-decoration: underline;
	}
</style>
