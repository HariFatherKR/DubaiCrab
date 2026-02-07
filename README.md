# 🦞 OpenKlaw

> 한국 사무직을 위한 로컬 AI 비서

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## ✨ 특징

- 🔒 **100% 로컬 실행** - 데이터가 외부로 나가지 않습니다
- 📄 **HWP 지원** - 한글 문서 요약 및 분석
- 📧 **이메일 작성** - 비즈니스 이메일 초안 생성
- 🇰🇷 **한국어 최적화** - 비즈니스 한국어 특화
- ⚡ **5분 설치** - 간편한 설치 경험

## 🚀 설치

### macOS / Linux

```bash
npm install -g openklaw
openklaw
```

### Windows

[Releases](https://github.com/HariFatherKR/OpenKlaw/releases)에서 설치 파일 다운로드

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

# 개발 서버
pnpm dev

# Tauri 개발 모드
pnpm tauri:dev

# 빌드
pnpm tauri:build

# 테스트
pnpm test

# 전체 검증
pnpm validate
```

## 📁 프로젝트 구조

```
OpenKlaw/
├── src/                  # SvelteKit 프론트엔드
│   ├── lib/
│   │   ├── components/   # UI 컴포넌트
│   │   ├── ollama/       # Ollama 클라이언트
│   │   ├── skills/       # AI 스킬 (이메일 등)
│   │   ├── tools/        # 도구 (HWP 파싱 등)
│   │   └── stores/       # 상태 관리
│   └── routes/           # 페이지
├── src-tauri/            # Tauri 백엔드 (Rust)
├── scripts/              # 설치 스크립트
├── bin/                  # CLI 진입점
└── docs/                 # 문서
```

## 🤝 기여

버그 리포트, 기능 제안, PR 모두 환영합니다!

1. Fork
2. Feature branch (`git checkout -b feature/amazing`)
3. Commit (`git commit -m 'feat: add amazing feature'`)
4. Push (`git push origin feature/amazing`)
5. Pull Request

## 📄 라이센스

MIT License - [LICENSE](LICENSE)

## 🔗 링크

- [문서](docs/)
- [이슈](https://github.com/HariFatherKR/OpenKlaw/issues)
- [릴리스](https://github.com/HariFatherKR/OpenKlaw/releases)
