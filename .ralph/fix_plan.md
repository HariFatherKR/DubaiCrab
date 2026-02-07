# OpenKlaw - 태스크 리스트

> 자동 업데이트됨. 각 에이전트가 완료 시 체크.

---

## 🔄 현재 Phase: 완료 🎉

---

## Phase 1: PM - 문서 업데이트 ✅ 완료

### 목표: 설치 경험 목표 문서화
- [x] docs/PRD.md 업데이트 - Windows 원클릭 설치 요구사항 추가
- [x] docs/PRD.md 업데이트 - macOS npm 설치 요구사항 추가
- [x] docs/INSTALL.md 생성 - 설치 가이드 초안
- [x] docs/ARCHITECTURE.md 생성 - 시스템 아키텍처

**완료 일시**: 2025-02-07

---

## Phase 2: Research - 사무직 워크플로우 리서치 ✅ 완료

### 목표: 타겟 사용자 업무 패턴 파악
- [x] 한국 사무직 일과 리서치 (블라인드, 긱뉴스)
- [x] AI 도구 사용 현황 조사
- [x] 경쟁 제품 분석 (Ollama, LM Studio, Jan.ai, GPT4All)
- [x] .ralph/specs/research-workflow.md 작성
- [x] git commit

**완료 일시**: 2025-02-07

---

## Phase 3: PM - 의사결정 ✅ 완료

### 목표: 리서치 기반 우선순위 결정
- [x] 리서치 결과 vs PRD 비교 분석
- [x] 핵심 기능 우선순위 재정의
- [x] fix_plan.md 개발 태스크 업데이트
- [x] docs/DECISIONS.md 생성 - 의사결정 기록
- [x] git commit
- [x] 추가 도구 분석 (캘린더, 엑셀, 카카오톡)
- [x] docs/pm-decisions.md 생성 - 상세 PM 의사결정 문서
- [x] 연동 도구 우선순위 확정 (HWP > 이메일 > 엑셀 > 캘린더 > 카카오톡)
- [x] MVP 스코프 및 Go/No-Go 기준 문서화

**완료 일시**: 2026-02-07

---

## Phase 4: Developer - 프로젝트 초기화 ✅ 완료

### 목표: 개발 환경 및 기본 구조
- [x] 프로젝트 디렉토리 구조 생성
- [x] 루트 package.json 설정 (pnpm workspace)
- [x] Tauri 프로젝트 초기화 (src-tauri/)
- [x] SvelteKit 프론트엔드 설정 (src/)
- [x] Ollama 클라이언트 래퍼 구현 (src/lib/ollama/)
- [x] 기본 채팅 UI 구현 (입력 → 스트리밍 응답)
- [x] pnpm run validate 통과
- [x] git commit & push

**완료 일시**: 2025-02-07

---

## Phase 5: Developer - 핵심 기능 ✅ 완료

### 목표: MVP 기능 구현
- [x] 시스템 트레이 구현 (Tauri lib.rs에 포함)
- [x] HWP 파싱 연동 (pyhwp subprocess via Rust command)
- [x] 이메일 작성 스킬 (src/lib/skills/email.ts)
- [x] 대화 저장/로드 (localStorage, SQLite 예정)
- [x] pnpm run validate 통과
- [x] git commit & push

**완료 일시**: 2025-02-07

---

## Phase 6: Developer - 설치 패키지 ✅ 완료

### 목표: 배포 준비
- [x] Windows NSIS 설치 패키지 설정 (tauri.conf.json)
- [x] macOS npm global 설치 설정 (package.json)
- [x] Ollama 자동 설치 스크립트 (scripts/install-ollama.*)
- [x] 모델 자동 다운로드 로직 (bin/openklaw.js)
- [x] bin/openklaw.js CLI 진입점
- [x] README.md 작성
- [x] pnpm run validate 통과
- [x] git commit & push

**완료 일시**: 2025-02-07

---

## Phase 7: QA - 기능 테스트 ✅ 완료

### 목표: 전체 기능 검증
- [x] pnpm test 통과 확인 (8개 테스트)
- [x] pnpm build 빌드 테스트 (SvelteKit 정상)
- [x] pnpm run validate 전체 검증 통과
- [ ] Ollama 연결 테스트 (환경 필요)
- [ ] HWP 파싱 테스트 (pyhwp 필요)
- [ ] Windows 빌드 테스트 (환경 필요)
- [x] docs/bugs/ 버그 리포트 작성
- [x] git commit & push

**완료 일시**: 2025-02-07
**참고**: 전체 통합 테스트는 환경 구성 후 진행 예정

---

## Phase 8: Developer - PM 대시보드 ✅ 완료

### 목표: PM 의사결정 시각화
- [x] 대시보드 라우트 생성 (/dashboard)
- [x] PM 의사결정 시각화
- [x] 기능 우선순위 차트
- [x] 프로젝트 진행 상황 표시
- [x] 연동 도구 우선순위 매트릭스
- [x] 사용 통계 UI (대화 수, 토큰 등)
- [x] 메인 페이지에서 대시보드 링크
- [x] pnpm run check/build/lint 통과
- [x] git commit

**완료 일시**: 2026-02-07

---

## 🐛 버그/이슈 (QA 발견)

> QA에서 발견된 버그는 여기에 추가

---

## ✅ 완료

- [x] 프로젝트 생성
- [x] GitHub 레포지토리 생성
- [x] 초기 문서 작성 (BRAINSTORM.md, PRD.md)
- [x] Ralph 구조 설정
