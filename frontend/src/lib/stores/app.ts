import { writable } from 'svelte/store';

interface AppState {
  ollamaReady: boolean;
  currentModel: string;
  sessionId: string;
  kakaoRunning: boolean;
}

const initialState: AppState = {
  ollamaReady: false,
  currentModel: 'qwen2.5:3b',
  sessionId: 'default',
  kakaoRunning: false
};

export const appState = writable<AppState>(initialState);

export function setOllamaReady(ready: boolean) {
  appState.update(s => ({ ...s, ollamaReady: ready }));
}

export function setCurrentModel(model: string) {
  appState.update(s => ({ ...s, currentModel: model }));
}

export function setSessionId(id: string) {
  appState.update(s => ({ ...s, sessionId: id }));
}

export function setKakaoRunning(running: boolean) {
  appState.update(s => ({ ...s, kakaoRunning: running }));
}
