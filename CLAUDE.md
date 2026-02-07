# OpenKlaw - 로컬 AI 비서

항상 한글로 답변

## 프로젝트 개요

OpenClaw의 경량화 포크. 한국 사무직을 위한 로컬 LLM 기반 AI 비서.

**핵심 가치:**
- 🔒 100% 로컬 실행 (데이터 유출 없음)
- 🇰🇷 한국어 + 한국 문서 특화 (HWP)
- 💨 가벼움 (사무용 노트북에서 실행)

## 타겟 사용자

- ChatGPT 사용 제한된 대기업/공공기관 직장인
- 노트북: i5, 16GB RAM, 내장 GPU

## 기술 스택

| 구성요소 | 기술 |
|----------|------|
| Desktop App | Tauri (Rust + SvelteKit) |
| LLM Runtime | Ollama |
| Default Model | Qwen2.5-3B-Instruct |
| HWP Parser | pyhwp (Python) |

## 폴더 구조

```
OpenKlaw/
├── docs/               # 문서
│   ├── BRAINSTORM.md   # 브레인스토밍
│   └── PRD.md          # 제품 요구사항
├── src/
│   ├── tauri/          # Rust 백엔드
│   └── web/            # SvelteKit 프론트엔드
├── scripts/            # 설치/빌드 스크립트
└── tests/              # 테스트
```

## 핵심 기능 (MVP)

1. **기본 채팅**: 로컬 LLM과 대화
2. **HWP 요약**: 한글 문서 파싱 + 요약
3. **이메일 작성**: 비즈니스 이메일 초안
4. **시스템 트레이**: 빠른 접근

## 개발 명령어

```bash
# 개발 모드
cd src/web && pnpm dev

# Tauri 개발
cd src/tauri && cargo tauri dev

# 빌드
cargo tauri build
```

## 참고 문서

- `docs/BRAINSTORM.md` - 아이디어 정리
- `docs/PRD.md` - 제품 요구사항
- OpenClaw 소스: https://github.com/openclaw/openclaw
