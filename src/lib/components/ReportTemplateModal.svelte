<script lang="ts">
	import { REPORT_TEMPLATES, generateReport, type ReportTemplate } from '$lib/skills/report';
	
	interface Props {
		isOpen: boolean;
		onClose: () => void;
		onGenerate: (content: string) => void;
	}
	
	let { isOpen, onClose, onGenerate }: Props = $props();
	
	let selectedTemplate = $state<ReportTemplate | null>(null);
	let fieldValues = $state<Record<string, string>>({});
	let isGenerating = $state(false);
	let step = $state<'select' | 'fill'>('select');
	
	function selectTemplate(template: ReportTemplate) {
		selectedTemplate = template;
		fieldValues = {};
		step = 'fill';
	}
	
	function goBack() {
		step = 'select';
		selectedTemplate = null;
		fieldValues = {};
	}
	
	async function handleGenerate() {
		if (!selectedTemplate) return;
		
		isGenerating = true;
		try {
			const result = await generateReport(selectedTemplate.id, fieldValues);
			onGenerate(`## ${result.title}\n\n${result.content}`);
			handleClose();
		} catch (error) {
			console.error('Report generation error:', error);
			onGenerate(`❌ 보고서 생성 중 오류가 발생했습니다: ${error}`);
		} finally {
			isGenerating = false;
		}
	}
	
	function handleClose() {
		step = 'select';
		selectedTemplate = null;
		fieldValues = {};
		onClose();
	}
	
	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) {
			handleClose();
		}
	}
</script>

{#if isOpen}
	<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
	<div class="modal-backdrop" onclick={handleBackdropClick}>
		<div class="modal">
			<div class="modal-header">
				{#if step === 'fill'}
					<button class="back-btn" onclick={goBack}>←</button>
				{/if}
				<h2>
					{#if step === 'select'}
						📋 보고서 템플릿
					{:else if selectedTemplate}
						{selectedTemplate.icon} {selectedTemplate.name}
					{/if}
				</h2>
				<button class="close-btn" onclick={handleClose}>✕</button>
			</div>
			
			<div class="modal-content">
				{#if step === 'select'}
					<div class="template-grid">
						{#each REPORT_TEMPLATES as template}
							<button 
								class="template-card" 
								onclick={() => selectTemplate(template)}
							>
								<span class="template-icon">{template.icon}</span>
								<span class="template-name">{template.name}</span>
								<span class="template-desc">{template.description}</span>
							</button>
						{/each}
					</div>
				{:else if selectedTemplate}
					<div class="field-form">
						{#each selectedTemplate.fields as field}
							<div class="field-group">
								<label for={field.id}>
									{field.label}
									{#if field.required}<span class="required">*</span>{/if}
								</label>
								{#if field.type === 'textarea'}
									<textarea
										id={field.id}
										bind:value={fieldValues[field.id]}
										placeholder={field.placeholder}
										rows="3"
									></textarea>
								{:else if field.type === 'select' && field.options}
									<select id={field.id} bind:value={fieldValues[field.id]}>
										<option value="">선택하세요</option>
										{#each field.options as option}
											<option value={option}>{option}</option>
										{/each}
									</select>
								{:else}
									<input
										type="text"
										id={field.id}
										bind:value={fieldValues[field.id]}
										placeholder={field.placeholder}
									/>
								{/if}
							</div>
						{/each}
					</div>
				{/if}
			</div>
			
			{#if step === 'fill'}
				<div class="modal-footer">
					<button 
						class="generate-btn" 
						onclick={handleGenerate}
						disabled={isGenerating}
					>
						{#if isGenerating}
							⏳ 생성 중...
						{:else}
							✨ 보고서 생성
						{/if}
					</button>
				</div>
			{/if}
		</div>
	</div>
{/if}

<style>
	.modal-backdrop {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.3);
		backdrop-filter: blur(4px);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		animation: fadeIn 0.2s ease-out;
	}
	
	@keyframes fadeIn {
		from { opacity: 0; }
		to { opacity: 1; }
	}
	
	.modal {
		/* Light glassmorphism */
		background: rgba(255, 255, 255, 0.9);
		backdrop-filter: blur(20px);
		-webkit-backdrop-filter: blur(20px);
		border: 1px solid rgba(255, 255, 255, 0.8);
		border-radius: 20px;
		width: 90%;
		max-width: 600px;
		max-height: 80vh;
		display: flex;
		flex-direction: column;
		animation: slideUp 0.3s ease-out;
		box-shadow: 0 20px 60px rgba(0, 0, 0, 0.12);
	}
	
	@keyframes slideUp {
		from { transform: translateY(20px); opacity: 0; }
		to { transform: translateY(0); opacity: 1; }
	}
	
	.modal-header {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 1.25rem 1.5rem;
		border-bottom: 1px solid rgba(0, 0, 0, 0.06);
	}
	
	.modal-header h2 {
		flex: 1;
		margin: 0;
		font-size: 1.25rem;
		font-weight: 600;
		color: #1e293b;
	}
	
	.back-btn, .close-btn {
		background: rgba(0, 0, 0, 0.04);
		border: none;
		width: 32px;
		height: 32px;
		border-radius: 8px;
		color: #64748b;
		cursor: pointer;
		transition: all 0.2s;
	}
	
	.back-btn:hover, .close-btn:hover {
		background: rgba(0, 0, 0, 0.08);
		color: #1e293b;
	}
	
	.modal-content {
		flex: 1;
		overflow-y: auto;
		padding: 1.5rem;
	}
	
	.template-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
		gap: 1rem;
	}
	
	.template-card {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
		padding: 1.25rem 1rem;
		background: rgba(255, 255, 255, 0.7);
		border: 1px solid rgba(0, 0, 0, 0.06);
		border-radius: 12px;
		cursor: pointer;
		transition: all 0.2s;
		text-align: center;
	}
	
	.template-card:hover {
		background: rgba(20, 184, 166, 0.08);
		border-color: rgba(20, 184, 166, 0.2);
		transform: translateY(-2px);
		box-shadow: 0 8px 24px rgba(20, 184, 166, 0.1);
	}
	
	.template-icon {
		font-size: 2rem;
	}
	
	.template-name {
		font-weight: 600;
		color: #1e293b;
		font-size: 0.95rem;
	}
	
	.template-desc {
		font-size: 0.8rem;
		color: #64748b;
	}
	
	.field-form {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}
	
	.field-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}
	
	.field-group label {
		font-size: 0.9rem;
		font-weight: 500;
		color: #475569;
	}
	
	.required {
		color: #dc2626;
		margin-left: 0.25rem;
	}
	
	.field-group input,
	.field-group textarea,
	.field-group select {
		padding: 0.75rem 1rem;
		background: rgba(255, 255, 255, 0.8);
		border: 1px solid rgba(0, 0, 0, 0.1);
		border-radius: 10px;
		color: #1e293b;
		font-size: 0.95rem;
		font-family: inherit;
		transition: all 0.2s;
	}
	
	.field-group input:focus,
	.field-group textarea:focus,
	.field-group select:focus {
		outline: none;
		border-color: rgba(20, 184, 166, 0.4);
		box-shadow: 0 0 0 3px rgba(20, 184, 166, 0.08);
	}
	
	.field-group textarea {
		resize: vertical;
		min-height: 80px;
	}
	
	.field-group input::placeholder,
	.field-group textarea::placeholder {
		color: #94a3b8;
	}
	
	.modal-footer {
		padding: 1rem 1.5rem 1.5rem;
		border-top: 1px solid rgba(0, 0, 0, 0.06);
	}
	
	.generate-btn {
		width: 100%;
		padding: 0.875rem 1.5rem;
		background: linear-gradient(135deg, #14b8a6, #0d9488);
		border: none;
		border-radius: 12px;
		color: white;
		font-size: 1rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
		box-shadow: 0 4px 12px rgba(20, 184, 166, 0.25);
	}
	
	.generate-btn:hover:not(:disabled) {
		background: linear-gradient(135deg, #0d9488, #0f766e);
		box-shadow: 0 6px 20px rgba(20, 184, 166, 0.35);
	}
	
	.generate-btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}
</style>
