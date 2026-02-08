/**
 * 보고서 템플릿 스킬
 * 주간 보고, 회의록, 제안서 등 비즈니스 문서 생성
 */

import { chat, DEFAULT_MODEL, type Message } from '$lib/ollama';

export interface ReportTemplate {
	id: string;
	name: string;
	icon: string;
	description: string;
	fields: TemplateField[];
	systemPrompt: string;
}

export interface TemplateField {
	id: string;
	label: string;
	type: 'text' | 'textarea' | 'date' | 'select';
	placeholder?: string;
	required?: boolean;
	options?: string[]; // for select type
}

export interface ReportResult {
	title: string;
	content: string;
	template: string;
}

/**
 * 보고서 템플릿 목록
 */
export const REPORT_TEMPLATES: ReportTemplate[] = [
	{
		id: 'weekly-report',
		name: '주간 업무 보고',
		icon: '📋',
		description: '주간 업무 진행 현황 및 계획',
		fields: [
			{ id: 'period', label: '보고 기간', type: 'text', placeholder: '예: 2024.01.01 ~ 01.07', required: true },
			{ id: 'completed', label: '완료된 업무', type: 'textarea', placeholder: '이번 주 완료한 업무들...', required: true },
			{ id: 'ongoing', label: '진행 중인 업무', type: 'textarea', placeholder: '현재 진행 중인 업무들...' },
			{ id: 'issues', label: '이슈/건의사항', type: 'textarea', placeholder: '발생한 이슈나 건의사항...' },
			{ id: 'nextWeek', label: '다음 주 계획', type: 'textarea', placeholder: '다음 주 예정된 업무들...' }
		],
		systemPrompt: `당신은 한국 기업 문서 작성 전문가입니다.

## 주간 업무 보고서 작성 규칙
1. 격식체 사용 (습니다/합니다 체)
2. 간결하고 명확한 문장
3. 핵심 성과 강조
4. 수치/데이터 포함 시 구체적으로
5. 문제점은 해결 방안과 함께 제시

## 보고서 구조
1. 보고 기간
2. 금주 완료 업무
3. 진행 중 업무
4. 이슈 및 건의사항
5. 차주 계획

마크다운 형식으로 깔끔하게 작성해주세요.`
	},
	{
		id: 'meeting-minutes',
		name: '회의록',
		icon: '📝',
		description: '회의 내용 정리 및 기록',
		fields: [
			{ id: 'title', label: '회의명', type: 'text', placeholder: '예: 프로젝트 킥오프 회의', required: true },
			{ id: 'datetime', label: '일시', type: 'text', placeholder: '예: 2024.01.15 14:00~15:30', required: true },
			{ id: 'attendees', label: '참석자', type: 'text', placeholder: '예: 김팀장, 이과장, 박대리' },
			{ id: 'agenda', label: '안건', type: 'textarea', placeholder: '논의된 안건들...', required: true },
			{ id: 'decisions', label: '결정사항', type: 'textarea', placeholder: '회의에서 결정된 사항들...' },
			{ id: 'actionItems', label: '후속 조치', type: 'textarea', placeholder: '담당자별 후속 조치 사항...' }
		],
		systemPrompt: `당신은 한국 기업 회의록 작성 전문가입니다.

## 회의록 작성 규칙
1. 객관적이고 사실적인 기록
2. 발언자와 발언 내용 구분
3. 결정사항은 명확하게
4. 후속 조치는 담당자와 기한 명시
5. 핵심 내용 중심으로 정리

## 회의록 구조
1. 회의 개요 (제목, 일시, 장소, 참석자)
2. 안건별 논의 내용
3. 결정 사항
4. Action Items (담당자, 기한)
5. 비고

마크다운 형식으로 체계적으로 작성해주세요.`
	},
	{
		id: 'proposal',
		name: '제안서',
		icon: '💡',
		description: '프로젝트/사업 제안서 초안',
		fields: [
			{ id: 'title', label: '제안 제목', type: 'text', placeholder: '예: 업무 자동화 시스템 도입 제안', required: true },
			{ id: 'background', label: '배경/현황', type: 'textarea', placeholder: '현재 상황과 문제점...', required: true },
			{ id: 'proposal', label: '제안 내용', type: 'textarea', placeholder: '제안하는 해결 방안...', required: true },
			{ id: 'benefits', label: '기대 효과', type: 'textarea', placeholder: '도입 시 기대되는 효과...' },
			{ id: 'timeline', label: '일정/예산', type: 'textarea', placeholder: '예상 일정 및 소요 예산...' }
		],
		systemPrompt: `당신은 한국 기업 제안서 작성 전문가입니다.

## 제안서 작성 규칙
1. 설득력 있는 논리 구조
2. 정량적 데이터로 효과 입증
3. 리스크와 대응 방안 포함
4. 단계별 실행 계획 제시
5. 경영진 관점에서 작성

## 제안서 구조
1. 제안 개요 (Executive Summary)
2. 현황 및 문제점
3. 제안 솔루션
4. 기대 효과 (ROI)
5. 실행 계획 및 일정
6. 예산
7. 결론

마크다운 형식으로 전문적으로 작성해주세요.`
	},
	{
		id: 'status-report',
		name: '프로젝트 현황 보고',
		icon: '📊',
		description: '프로젝트 진행 상황 보고',
		fields: [
			{ id: 'projectName', label: '프로젝트명', type: 'text', placeholder: '프로젝트 이름', required: true },
			{ id: 'reportDate', label: '보고일', type: 'text', placeholder: '예: 2024.01.15' },
			{ id: 'progress', label: '진행률', type: 'text', placeholder: '예: 65%' },
			{ id: 'completed', label: '완료 항목', type: 'textarea', placeholder: '완료된 마일스톤/태스크...' },
			{ id: 'inProgress', label: '진행 중 항목', type: 'textarea', placeholder: '현재 진행 중인 작업...' },
			{ id: 'risks', label: '리스크/이슈', type: 'textarea', placeholder: '발생한 리스크나 이슈...' }
		],
		systemPrompt: `당신은 한국 기업 프로젝트 관리 전문가입니다.

## 프로젝트 현황 보고서 작성 규칙
1. 객관적 진행 상황 제시
2. RAG 상태 표시 (정상/주의/위험)
3. 마일스톤 기준 진척도
4. 이슈는 영향도와 대응 방안 포함
5. 의사결정 필요 사항 명시

## 보고서 구조
1. 프로젝트 개요
2. 전체 진행 현황 (진행률, 일정)
3. 주요 완료 항목
4. 진행 중 업무
5. 리스크 및 이슈
6. 의사결정 요청 사항
7. 다음 단계

마크다운 형식으로 체계적으로 작성해주세요.`
	},
	{
		id: 'handover',
		name: '업무 인수인계서',
		icon: '🔄',
		description: '업무 이관 및 인수인계 문서',
		fields: [
			{ id: 'department', label: '부서/팀', type: 'text', placeholder: '예: 개발팀', required: true },
			{ id: 'fromTo', label: '인수인계자', type: 'text', placeholder: '예: 김철수 → 이영희' },
			{ id: 'period', label: '인수인계 기간', type: 'text', placeholder: '예: 2024.01.15 ~ 01.19' },
			{ id: 'tasks', label: '담당 업무', type: 'textarea', placeholder: '주요 담당 업무 목록...', required: true },
			{ id: 'contacts', label: '주요 연락처', type: 'textarea', placeholder: '업무 관련 주요 연락처...' },
			{ id: 'notes', label: '특이사항', type: 'textarea', placeholder: '업무 수행 시 유의사항...' }
		],
		systemPrompt: `당신은 한국 기업 업무 인수인계 전문가입니다.

## 인수인계서 작성 규칙
1. 업무별 상세 프로세스 기술
2. 주기적/비정기적 업무 구분
3. 필수 시스템/도구 접근 정보
4. 주요 이해관계자 연락처
5. 업무 수행 팁과 주의사항

## 인수인계서 구조
1. 인수인계 개요
2. 담당 업무 목록
3. 업무별 상세 프로세스
4. 주요 연락처/협업 부서
5. 시스템/도구 접근 정보
6. 진행 중 업무 현황
7. 특이사항 및 주의점

마크다운 형식으로 상세하게 작성해주세요.`
	}
];

/**
 * 템플릿 기반 보고서 생성
 */
export async function generateReport(
	templateId: string,
	fieldValues: Record<string, string>
): Promise<ReportResult> {
	const template = REPORT_TEMPLATES.find(t => t.id === templateId);
	if (!template) {
		throw new Error(`템플릿을 찾을 수 없습니다: ${templateId}`);
	}

	const userPrompt = buildReportPrompt(template, fieldValues);
	
	const messages: Message[] = [
		{ role: 'system', content: template.systemPrompt },
		{ role: 'user', content: userPrompt }
	];

	let response = '';
	for await (const chunk of chat(DEFAULT_MODEL, messages, { temperature: 0.7 })) {
		if (chunk.message?.content) {
			response += chunk.message.content;
		}
	}

	return {
		title: `${template.name} - ${fieldValues['title'] || fieldValues['period'] || new Date().toLocaleDateString('ko-KR')}`,
		content: response,
		template: template.id
	};
}

/**
 * 템플릿 필드값으로 프롬프트 생성
 */
function buildReportPrompt(template: ReportTemplate, values: Record<string, string>): string {
	let prompt = `다음 정보를 바탕으로 ${template.name}를 작성해주세요:\n\n`;

	for (const field of template.fields) {
		const value = values[field.id];
		if (value && value.trim()) {
			prompt += `### ${field.label}\n${value}\n\n`;
		}
	}

	prompt += `\n위 정보를 바탕으로 완성된 ${template.name}를 작성해주세요.`;
	return prompt;
}

/**
 * 템플릿 ID로 템플릿 조회
 */
export function getTemplate(templateId: string): ReportTemplate | undefined {
	return REPORT_TEMPLATES.find(t => t.id === templateId);
}

/**
 * 보고서 관련 요청인지 감지
 */
export function isReportRequest(input: string): boolean {
	const keywords = [
		'보고서', '회의록', '제안서', '인수인계', '현황', '주간보고',
		'월간보고', '일일보고', 'report', '작성해', '써줘'
	];
	const lowerInput = input.toLowerCase();
	return keywords.some(kw => lowerInput.includes(kw));
}
