# OpenKlaw - 태스크 리스트

> 자동 업데이트됨. 각 에이전트가 완료 시 체크.

---

## 🔄 현재 Phase: Developer - 프로젝트 초기화 (Phase 4)

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

**완료 일시**: 2025-02-07

---

## Phase 4: Developer - 프로젝트 초기화 ⏳ 진행 중

### 목표: 개발 환경 및 기본 구조
- [ ] 프로젝트 디렉토리 구조 생성
- [ ] 루트 package.json 설정 (pnpm workspace)
- [ ] Tauri 프로젝트 초기화 (src-tauri/)
- [ ] SvelteKit 프론트엔드 설정 (src/)
- [ ] Ollama 클라이언트 래퍼 구현 (src/lib/ollama/)
- [ ] 기본 채팅 UI 구현 (입력 → 스트리밍 응답)
- [ ] pnpm run validate 통과
- [ ] git commit & push

---

## Phase 5: Developer - 핵심 기능 ⏸️

### 목표: MVP 기능 구현
- [ ] 시스템 트레이 구현
- [ ] HWP 파싱 연동 (pyhwp)
- [ ] 이메일 작성 스킬
- [ ] 대화 저장/로드
- [ ] git commit

---

## Phase 6: Developer - 설치 패키지 ⏸️

### 목표: 배포 준비
- [ ] Windows: NSIS 또는 MSI 설치 패키지
- [ ] macOS: npm global 설치 테스트
- [ ] 자동 Ollama 설치 스크립트
- [ ] 모델 자동 다운로드 로직
- [ ] git commit

---

## Phase 7: QA - 기능 테스트 ⏸️

### 목표: 전체 기능 검증
- [ ] 설치 테스트 (Windows)
- [ ] 설치 테스트 (macOS npm)
- [ ] 채팅 기능 테스트
- [ ] HWP 파싱 테스트
- [ ] 이메일 작성 테스트
- [ ] 시스템 트레이 테스트
- [ ] docs/bugs/ 버그 리포트 작성
- [ ] git commit

---

## 🐛 버그/이슈 (QA 발견)

> QA에서 발견된 버그는 여기에 추가

---

## ✅ 완료

- [x] 프로젝트 생성
- [x] GitHub 레포지토리 생성
- [x] 초기 문서 작성 (BRAINSTORM.md, PRD.md)
- [x] Ralph 구조 설정
