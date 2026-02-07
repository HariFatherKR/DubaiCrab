/**
 * 사용 통계 저장소
 * localStorage에 저장되는 앱 사용 통계
 */

const STATS_KEY = 'openklaw_stats';

export interface Stats {
	totalChats: number;
	totalMessages: number;
	hwpProcessed: number;
	emailsGenerated: number;
	lastUsed: string;
}

const defaultStats: Stats = {
	totalChats: 0,
	totalMessages: 0,
	hwpProcessed: 0,
	emailsGenerated: 0,
	lastUsed: new Date().toISOString()
};

/**
 * 통계 로드
 */
export function loadStats(): Stats {
	if (typeof localStorage === 'undefined') {
		return defaultStats;
	}
	
	try {
		const saved = localStorage.getItem(STATS_KEY);
		if (saved) {
			return { ...defaultStats, ...JSON.parse(saved) };
		}
	} catch {
		// 파싱 실패 시 기본값 반환
	}
	
	return defaultStats;
}

/**
 * 통계 저장
 */
export function saveStats(stats: Stats): void {
	if (typeof localStorage === 'undefined') {
		return;
	}
	
	try {
		stats.lastUsed = new Date().toISOString();
		localStorage.setItem(STATS_KEY, JSON.stringify(stats));
	} catch {
		// 저장 실패 무시
	}
}

/**
 * 메시지 카운트 증가
 */
export function incrementMessages(): void {
	const stats = loadStats();
	stats.totalMessages++;
	saveStats(stats);
}

/**
 * 새 대화 시작
 */
export function incrementChats(): void {
	const stats = loadStats();
	stats.totalChats++;
	saveStats(stats);
}

/**
 * HWP 처리 카운트 증가
 */
export function incrementHwp(): void {
	const stats = loadStats();
	stats.hwpProcessed++;
	saveStats(stats);
}

/**
 * 이메일 생성 카운트 증가
 */
export function incrementEmails(): void {
	const stats = loadStats();
	stats.emailsGenerated++;
	saveStats(stats);
}
