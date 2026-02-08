# Dubai Crab 🦀 - 로컬 AI 비서

항상 한글로 답변

## 프로젝트 개요

OpenClaw의 경량화 포크. 한국 사무직을 위한 로컬 LLM 기반 AI 비서.

**핵심 가치:**
- 🔒 100% 로컬 실행 (데이터 유출 없음)
- 🇰🇷 한국어 + 한국 문서 특화 (HWP)
- 💨 가벼움 (사무용 노트북에서 실행)

## 링크

- **웹사이트**: https://dubaicrab.ai
- **GitHub**: https://github.com/HariFatherKR/DubaiCrab

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
DubaiCrab/
├── docs/               # 문서
│   ├── BRAINSTORM.md   # 브레인스토밍
│   └── PRD.md          # 제품 요구사항
├── src/                # SvelteKit 프론트엔드
├── src-tauri/          # Rust 백엔드
├── static/             # 정적 파일
└── bin/                # CLI 스크립트
```

## 핵심 기능 (MVP ✅ 100%)

1. **기본 채팅**: 로컬 LLM과 대화
2. **HWP 요약**: 한글 문서 파싱 + 요약
3. **이메일 작성**: 비즈니스 이메일 초안
4. **엑셀/CSV 분석**: 데이터 요약
5. **보고서 템플릿**: 5종 템플릿
6. **빠른 작업 버튼**: 원클릭 기능
7. **전역 단축키**: Ctrl+Shift+D
8. **시스템 트레이**: 빠른 접근
9. **설정 페이지**: 모델/단축키 설정

## 개발 명령어

```bash
# 의존성 설치
pnpm install

# 개발 모드
pnpm dev

# Tauri 앱 개발
pnpm tauri dev

# 빌드
pnpm build
pnpm tauri build
```

## 테마 (두바이 쫀득 쿠키 🍪)

- **피스타치오 그린**: #4a7c59
- **초콜릿 브라운**: #5D4037
- **골든/카라멜**: #D4A574

## 참고 문서

- `docs/PRD.md` - 제품 요구사항
- `docs/ARCHITECTURE.md` - 기술 아키텍처
- `docs/pm-decisions.md` - PM 의사결정 로그
