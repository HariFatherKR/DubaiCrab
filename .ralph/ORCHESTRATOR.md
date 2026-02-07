# OpenKlaw - Ralph Loop Orchestrator

## 현재 상태

| Phase | 에이전트 | 상태 | 시작 시간 |
|-------|----------|------|----------|
| Phase 1 | PM | 🟡 진행중 | 2025-02-07 23:15 |
| Phase 2 | Research | 🟡 진행중 | 2025-02-07 23:15 |
| Phase 3 | PM (의사결정) | ⏸️ 대기 | - |
| Phase 4 | Developer | ⏸️ 대기 | - |
| Phase 5 | Developer | ⏸️ 대기 | - |
| Phase 6 | Developer | ⏸️ 대기 | - |
| Phase 7 | QA | ⏸️ 대기 | - |

## 다음 Phase 트리거 조건

### Phase 3 시작 조건 (PM 의사결정)
- [x] Phase 1 PM 완료
- [x] Phase 2 Research 완료
- 리서치 문서 존재:
  - .ralph/specs/research-workflow.md
  - .ralph/specs/research-competitors.md
  - .ralph/specs/research-insights.md

### Phase 4-6 시작 조건 (Developer)
- [ ] Phase 3 PM 의사결정 완료
- [ ] fix_plan.md 개발 태스크 정의됨

### Phase 7 시작 조건 (QA)
- [ ] Phase 4-6 Developer 완료
- [ ] pnpm run validate 통과

## 루프 종료 조건

- [ ] fix_plan.md 모든 태스크 완료
- [ ] QA 100% 통과
- [ ] Windows 설치 테스트 통과
- [ ] macOS npm 설치 테스트 통과

## 에러 핸들링

### Developer 에러 시
1. 에러 로그 기록
2. 3회 연속 실패 시 PM 호출
3. PM이 방향성 수정 후 재시작

### QA 실패 시
1. 버그 리포트 작성 (docs/bugs/)
2. Developer로 태스크 전달
3. 수정 후 다시 QA

## 진행 로그

```
[2025-02-07 23:15] Ralph Loop 시작
[2025-02-07 23:15] PM Phase 1 시작
[2025-02-07 23:15] Research Phase 2 시작
```
