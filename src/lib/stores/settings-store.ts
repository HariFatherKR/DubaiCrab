/**
 * 설정 스토어
 * 앱 설정 상태 관리
 */

import { writable, get } from 'svelte/store';
import { browser } from '$app/environment';

export interface AppSettings {
	// 모델 설정
	model: string;
	customModels: string[];
	
	// 테마 설정
	theme: 'dark' | 'light' | 'system';
	accentColor: string;
	
	// 단축키 설정
	shortcuts: {
		toggleApp: string;
		focusInput: string;
		newChat: string;
	};
	
	// 데이터 설정
	dataPath: string;
	autoSave: boolean;
	
	// 기타
	language: 'ko' | 'en';
	sendOnEnter: boolean;
}

const DEFAULT_SETTINGS: AppSettings = {
	model: 'qwen2.5:3b-instruct',
	customModels: [],
	theme: 'dark',
	accentColor: '#14b8a6',
	shortcuts: {
		toggleApp: 'Cmd+Shift+O',
		focusInput: 'Cmd+/',
		newChat: 'Cmd+N'
	},
	dataPath: '~/.openklaw',
	autoSave: true,
	language: 'ko',
	sendOnEnter: true
};

const STORAGE_KEY = 'openklaw_settings';

function loadSettings(): AppSettings {
	if (!browser) return DEFAULT_SETTINGS;
	
	try {
		const stored = localStorage.getItem(STORAGE_KEY);
		if (stored) {
			return { ...DEFAULT_SETTINGS, ...JSON.parse(stored) };
		}
	} catch (e) {
		console.error('Failed to load settings:', e);
	}
	return DEFAULT_SETTINGS;
}

function createSettingsStore() {
	const { subscribe, set, update } = writable<AppSettings>(loadSettings());
	
	return {
		subscribe,
		
		/**
		 * 설정 업데이트
		 */
		updateSettings: (partial: Partial<AppSettings>) => {
			update(current => {
				const updated = { ...current, ...partial };
				if (browser) {
					localStorage.setItem(STORAGE_KEY, JSON.stringify(updated));
				}
				return updated;
			});
		},
		
		/**
		 * 설정 초기화
		 */
		resetSettings: () => {
			set(DEFAULT_SETTINGS);
			if (browser) {
				localStorage.setItem(STORAGE_KEY, JSON.stringify(DEFAULT_SETTINGS));
			}
		},
		
		/**
		 * 현재 설정 가져오기
		 */
		getSettings: (): AppSettings => {
			return get({ subscribe });
		},
		
		/**
		 * 커스텀 모델 추가
		 */
		addCustomModel: (modelName: string) => {
			update(current => {
				if (!current.customModels.includes(modelName)) {
					const updated = {
						...current,
						customModels: [...current.customModels, modelName]
					};
					if (browser) {
						localStorage.setItem(STORAGE_KEY, JSON.stringify(updated));
					}
					return updated;
				}
				return current;
			});
		},
		
		/**
		 * 커스텀 모델 제거
		 */
		removeCustomModel: (modelName: string) => {
			update(current => {
				const updated = {
					...current,
					customModels: current.customModels.filter(m => m !== modelName)
				};
				if (browser) {
					localStorage.setItem(STORAGE_KEY, JSON.stringify(updated));
				}
				return updated;
			});
		}
	};
}

export const settingsStore = createSettingsStore();

// 사용 가능한 기본 모델 목록
export const AVAILABLE_MODELS = [
	{ id: 'qwen2.5:3b-instruct', name: 'Qwen 2.5 3B (권장)', size: '~2GB' },
	{ id: 'qwen2.5:7b-instruct', name: 'Qwen 2.5 7B', size: '~4.5GB' },
	{ id: 'llama3.2:3b', name: 'Llama 3.2 3B', size: '~2GB' },
	{ id: 'gemma2:2b', name: 'Gemma 2 2B', size: '~1.6GB' },
	{ id: 'phi3:mini', name: 'Phi-3 Mini', size: '~2.3GB' },
	{ id: 'mistral:7b', name: 'Mistral 7B', size: '~4.1GB' }
];

// 테마 옵션
export const THEME_OPTIONS = [
	{ id: 'dark', name: '다크', icon: '🌙' },
	{ id: 'light', name: '라이트', icon: '☀️' },
	{ id: 'system', name: '시스템', icon: '💻' }
];

// 액센트 컬러 옵션
export const ACCENT_COLORS = [
	{ id: '#14b8a6', name: '틸' },
	{ id: '#3b82f6', name: '블루' },
	{ id: '#8b5cf6', name: '퍼플' },
	{ id: '#ec4899', name: '핑크' },
	{ id: '#f97316', name: '오렌지' },
	{ id: '#22c55e', name: '그린' }
];
