#!/usr/bin/env node
/**
 * OpenKlaw CLI 진입점
 * npm install -g openklaw 후 openklaw 명령어로 실행
 */

import { spawn, exec } from 'child_process';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';
import { existsSync } from 'fs';
import { promisify } from 'util';

const execAsync = promisify(exec);
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const OLLAMA_URL = 'http://localhost:11434';
const DEFAULT_MODEL = 'qwen2.5:3b-instruct';

async function checkOllama() {
	try {
		const response = await fetch(`${OLLAMA_URL}/api/tags`);
		return response.ok;
	} catch {
		return false;
	}
}

async function checkModel() {
	try {
		const response = await fetch(`${OLLAMA_URL}/api/tags`);
		const data = await response.json();
		return data.models?.some((m) => m.name.includes('qwen2.5'));
	} catch {
		return false;
	}
}

async function installOllama() {
	console.log('📦 Ollama 설치 중...');

	const platform = process.platform;

	if (platform === 'darwin') {
		// macOS
		try {
			await execAsync('brew install ollama');
			console.log('✅ Ollama 설치 완료 (Homebrew)');
		} catch {
			console.log('📥 공식 스크립트로 설치 중...');
			await execAsync('curl -fsSL https://ollama.ai/install.sh | sh');
		}
	} else if (platform === 'linux') {
		// Linux
		await execAsync('curl -fsSL https://ollama.ai/install.sh | sh');
		console.log('✅ Ollama 설치 완료');
	} else if (platform === 'win32') {
		// Windows - PowerShell 스크립트 실행
		const scriptPath = join(__dirname, '..', 'scripts', 'install-ollama.ps1');
		spawn('powershell', ['-ExecutionPolicy', 'Bypass', '-File', scriptPath], {
			stdio: 'inherit'
		});
		return;
	}
}

async function pullModel() {
	console.log(`📥 모델 다운로드 중: ${DEFAULT_MODEL}`);
	console.log('   (약 2GB, 5-15분 소요)');

	return new Promise((resolve, reject) => {
		const pull = spawn('ollama', ['pull', DEFAULT_MODEL], {
			stdio: 'inherit'
		});

		pull.on('close', (code) => {
			if (code === 0) {
				console.log('✅ 모델 다운로드 완료!');
				resolve();
			} else {
				reject(new Error('모델 다운로드 실패'));
			}
		});
	});
}

async function startApp() {
	// Tauri 앱 경로 찾기
	const appPaths = {
		darwin: join(__dirname, '..', 'dist', 'OpenKlaw.app', 'Contents', 'MacOS', 'OpenKlaw'),
		linux: join(__dirname, '..', 'dist', 'openklaw'),
		win32: join(__dirname, '..', 'dist', 'OpenKlaw.exe')
	};

	const appPath = appPaths[process.platform];

	if (existsSync(appPath)) {
		console.log('🚀 OpenKlaw 실행 중...');
		spawn(appPath, [], {
			detached: true,
			stdio: 'ignore'
		}).unref();
	} else {
		// 앱이 없으면 개발 모드로 실행
		console.log('🔧 개발 모드로 실행 중...');
		console.log('   http://localhost:1420 에서 UI 확인');

		spawn('pnpm', ['dev'], {
			stdio: 'inherit',
			cwd: join(__dirname, '..')
		});
	}
}

async function main() {
	console.log('');
	console.log('🦞 OpenKlaw - 로컬 AI 비서');
	console.log('');

	// 1. Ollama 확인
	console.log('🔍 Ollama 확인 중...');
	const ollamaRunning = await checkOllama();

	if (!ollamaRunning) {
		console.log('⚠️ Ollama가 실행되지 않았습니다.');

		try {
			await execAsync('which ollama || where ollama');
			console.log('🚀 Ollama 서비스 시작 중...');
			spawn('ollama', ['serve'], {
				detached: true,
				stdio: 'ignore'
			}).unref();

			// 시작 대기
			await new Promise((resolve) => setTimeout(resolve, 3000));
		} catch {
			await installOllama();
		}
	} else {
		console.log('✅ Ollama 실행 중');
	}

	// 2. 모델 확인
	console.log('🔍 모델 확인 중...');
	const modelExists = await checkModel();

	if (!modelExists) {
		await pullModel();
	} else {
		console.log('✅ 모델 준비 완료');
	}

	// 3. 앱 시작
	console.log('');
	await startApp();
}

// 명령어 처리
const args = process.argv.slice(2);

if (args.includes('--help') || args.includes('-h')) {
	console.log(`
🦞 OpenKlaw - 한국 사무직을 위한 로컬 AI 비서

사용법:
  openklaw              앱 실행
  openklaw --setup      Ollama + 모델 설치만
  openklaw --version    버전 표시
  openklaw --help       도움말 표시

문서: https://github.com/HariFatherKR/OpenKlaw
`);
	process.exit(0);
}

if (args.includes('--version') || args.includes('-v')) {
	console.log('OpenKlaw v0.1.0');
	process.exit(0);
}

main().catch((err) => {
	console.error('❌ 오류:', err.message);
	process.exit(1);
});
