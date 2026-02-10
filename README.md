# Dubai Crab Go 🦀

한국 사무직을 위한 로컬 AI 비서 - Wails (Go + Svelte) 버전

## 특징

- 🔒 **완전 로컬 실행**: 데이터가 외부로 전송되지 않음
- 🤖 **Ollama 통합**: 다양한 LLM 모델 지원
- 💬 **카카오톡 연동**: 오픈빌더 스킬서버 내장
- 📄 **HWP 지원**: 한글 문서 텍스트 추출
- 🎨 **다크 초콜릿 테마**: 편안한 다크 UI

## 요구사항

- macOS 10.15+, Windows 10+, 또는 Linux
- [Go 1.21+](https://go.dev/dl/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)
- [Ollama](https://ollama.ai/) (로컬 LLM 실행)
- Node.js 20+

## 설치

```bash
# Wails CLI 설치
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 프로젝트 클론
git clone https://github.com/HariFatherKR/DubaiCrab-Go.git
cd DubaiCrab-Go

# 의존성 설치 및 빌드
wails build
```

## 개발

```bash
# 개발 모드 실행
wails dev
```

## 프로젝트 구조

```
DubaiCrab-Go/
├── app.go                  # 메인 앱 (Wails 바인딩)
├── main.go                 # 엔트리 포인트
├── internal/
│   ├── agent/              # AI 에이전트 루프
│   ├── config/             # 설정 관리
│   ├── kakao/              # 카카오톡 웹훅 서버
│   ├── ollama/             # Ollama 클라이언트
│   └── tools/              # 도구 시스템
└── frontend/
    ├── src/
    │   ├── lib/
    │   │   ├── components/ # Svelte 컴포넌트
    │   │   └── stores/     # 상태 관리
    │   └── App.svelte      # 메인 앱
    └── wailsjs/            # 자동 생성 바인딩
```

## 기능

### 채팅
- Ollama 모델과 대화
- 세션 기반 대화 관리
- 컨텍스트 유지

### 카카오톡 연동
- 오픈빌더 스킬서버 웹훅
- 비동기 콜백 지원
- 접근 제어 (allowlist)

### 도구 시스템
- 시스템 정보 조회
- 클립보드 복사/붙여넣기
- URL 열기
- HWP 파일 파싱

## 설정

설정 파일: `~/.config/dubai-crab/config.json`

```json
{
  "ollamaUrl": "http://localhost:11434",
  "ollamaModel": "qwen2.5:3b",
  "kakaoEnabled": true,
  "kakaoPort": 3847,
  "kakaoWebhookPath": "/kakao/webhook"
}
```

## 라이선스

MIT License

## 기여

이슈와 PR을 환영합니다!
