# 🦀 Dubai Crab

> 한국 사무직을 위한 로컬 AI 비서

[![Website](https://img.shields.io/badge/Website-DubaiCrab.ai-D4A574)](https://dubaicrab.ai)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub](https://img.shields.io/github/stars/HariFatherKR/DubaiCrab?style=social)](https://github.com/HariFatherKR/DubaiCrab)

## ✨ 특징

- 🔒 **100% 로컬 실행** - 데이터가 외부로 나가지 않습니다
- 📄 **HWP 지원** - 한글 문서 요약 및 분석
- 📧 **이메일 작성** - 비즈니스 이메일 초안 생성
- 🇰🇷 **한국어 최적화** - 비즈니스 한국어 특화
- ⚡ **5분 설치** - 간편한 설치 경험

## 🚀 설치

### macOS / Linux

```bash
npm install -g dubai-crab
dubai-crab
```

### Windows

[Releases](https://github.com/HariFatherKR/DubaiCrab/releases)에서 설치 파일 다운로드

## 📋 시스템 요구사항

| 항목 | 최소 | 권장 |
|------|------|------|
| CPU | Intel i5 / Ryzen 5 | Intel i7 / Ryzen 7 |
| RAM | 8GB | 16GB |
| 저장공간 | 10GB | 20GB |
| OS | Windows 10/11, macOS 12+ | 동일 |

## 🎯 주요 기능

### 💬 AI 채팅
로컬 LLM(Qwen2.5-3B)을 사용한 한국어 대화

### 📄 HWP 요약
한글 문서를 드래그하면 자동 요약

### 📧 이메일 작성
"김과장님께 회의 일정 조율 이메일 써줘"

### 🗂️ 대화 기록
모든 대화가 로컬에 저장되어 나중에 검색 가능

## 🛠️ 개발

```bash
# 의존성 설치
pnpm install

# 개발 서버 실행
pnpm dev

# Tauri 앱 개발
pnpm tauri:dev

# 빌드
pnpm build
```

## 📁 프로젝트 구조

```
dubai-crab/
├── src/
│   ├── lib/
│   │   ├── components/   # Svelte 컴포넌트
│   │   ├── stores/       # 상태 관리
│   │   ├── skills/       # AI 스킬 모듈
│   │   └── ollama.ts     # Ollama API 클라이언트
│   └── routes/           # SvelteKit 라우트
├── src-tauri/            # Tauri 백엔드
├── static/               # 정적 파일
└── bin/                  # CLI 스크립트
```

## 🎨 테마

Dubai Crab은 두바이 쫀득 쿠키에서 영감을 받은 독특한 색상 팔레트를 사용합니다:

- **피스타치오 그린**: #4a7c59, #3d6b4f, #2d5a3f
- **초콜릿 브라운**: #5D4037, #4E342E, #3E2723
- **골든/카라멜**: #D4A574, #C49A6C, #B8860B

## 📄 라이선스

MIT License - 자유롭게 사용하세요!

## 🙏 기여

PR 환영합니다! [Contributing Guide](CONTRIBUTING.md)를 참고해주세요.

---

Made with ❤️ by [HariFatherKR](https://github.com/HariFatherKR)
