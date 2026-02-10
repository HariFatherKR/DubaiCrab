# Dubai Crab Go - 아키텍처 문서

## 개요

Dubai Crab Go는 Wails 프레임워크를 사용하여 Go 백엔드와 Svelte 프론트엔드를 결합한 데스크톱 애플리케이션입니다.

## 기술 스택

### 백엔드 (Go)
- **Wails v2**: 데스크톱 앱 프레임워크
- **net/http**: 카카오 웹훅 서버
- **encoding/json**: JSON 처리

### 프론트엔드 (Svelte 5)
- **Svelte 5**: 반응형 UI
- **TypeScript**: 타입 안전성
- **Tailwind CSS**: 스타일링
- **Vite**: 빌드 도구

## 핵심 컴포넌트

### 1. Ollama Manager (`internal/ollama/`)

Ollama 서버의 생명주기와 API 호출을 관리합니다.

```go
type Manager struct {
    baseURL string
    process *exec.Cmd
    mu      sync.Mutex
    client  *http.Client
}
```

**주요 기능:**
- 서버 시작/종료
- 모델 목록 조회
- 채팅 API 호출
- 모델 다운로드

### 2. Kakao Server (`internal/kakao/`)

카카오톡 오픈빌더 스킬서버 웹훅을 처리합니다.

```go
type Server struct {
    ollama    *ollama.Manager
    config    *Config
    server    *http.Server
    running   bool
    mu        sync.RWMutex
    cancelFn  context.CancelFunc
}
```

**엔드포인트:**
- `POST /kakao/webhook`: 카카오 스킬 요청 처리
- `GET /health`: 헬스체크

**응답 형식:**
```json
{
  "version": "2.0",
  "template": {
    "outputs": [
      {"simpleText": {"text": "응답 메시지"}}
    ]
  }
}
```

### 3. Agent (`internal/agent/`)

AI 대화 루프를 관리합니다.

```go
type Agent struct {
    ollama       *ollama.Manager
    toolRegistry *tools.Registry
    sessions     map[string]*Session
    mu           sync.RWMutex
    model        string
    systemPrompt string
    maxTokens    int
}
```

**기능:**
- 세션 기반 대화 관리
- 컨텍스트 유지
- 도구 호출 지원

### 4. Tools Registry (`internal/tools/`)

확장 가능한 도구 시스템입니다.

```go
type Tool interface {
    Name() string
    Description() string
    Execute(ctx context.Context, params map[string]interface{}) (string, error)
    Schema() map[string]interface{}
}
```

**내장 도구:**
- `system_info`: 시스템 정보 조회
- `clipboard`: 클립보드 복사/붙여넣기
- `open_url`: URL 열기
- `parse_hwp`: HWP 파일 파싱

## 데이터 흐름

```
┌─────────────┐     Wails Binding      ┌─────────────┐
│   Frontend  │◀─────────────────────▶ │   app.go    │
│   (Svelte)  │                        │  (Go Core)  │
└─────────────┘                        └──────┬──────┘
                                              │
                    ┌─────────────────────────┼─────────────────────────┐
                    │                         │                         │
              ┌─────▼─────┐           ┌──────▼──────┐           ┌──────▼──────┐
              │  Ollama   │           │   Agent     │           │   Kakao     │
              │  Manager  │           │   Loop      │           │   Server    │
              └─────┬─────┘           └──────┬──────┘           └──────┬──────┘
                    │                        │                         │
                    │                  ┌─────▼─────┐                   │
                    │                  │   Tools   │                   │
                    │                  │  Registry │                   │
                    │                  └───────────┘                   │
                    ▼                                                  ▼
              ┌───────────┐                                    ┌───────────────┐
              │  Ollama   │                                    │  카카오톡      │
              │  Server   │                                    │  오픈빌더      │
              └───────────┘                                    └───────────────┘
```

## 설정 관리

설정은 JSON 형식으로 저장됩니다:

```
~/.config/dubai-crab/config.json
```

### Config 구조

```go
type Config struct {
    AppName          string   `json:"appName"`
    Version          string   `json:"version"`
    OllamaURL        string   `json:"ollamaUrl"`
    OllamaModel      string   `json:"ollamaModel"`
    KakaoEnabled     bool     `json:"kakaoEnabled"`
    KakaoPort        int      `json:"kakaoPort"`
    KakaoWebhookPath string   `json:"kakaoWebhookPath"`
    KakaoDMPolicy    string   `json:"kakaoDmPolicy"`
    KakaoAllowFrom   []string `json:"kakaoAllowFrom"`
    KakaoSystemPrompt string  `json:"kakaoSystemPrompt"`
    KakaoModel       string   `json:"kakaoModel"`
    RelayURL         string   `json:"relayUrl"`
    RelayToken       string   `json:"relayToken"`
}
```

## 보안 고려사항

1. **로컬 전용**: 모든 AI 처리는 로컬에서 수행
2. **접근 제어**: 카카오 서버에 allowlist 기반 접근 제어
3. **URL 검증**: http/https 프로토콜만 허용
4. **파일 경로**: 홈 디렉토리 내 파일만 접근 허용

## 빌드 및 배포

### 개발

```bash
wails dev
```

### 프로덕션 빌드

```bash
wails build
```

빌드 결과:
- macOS: `build/bin/DubaiCrab-Go.app`
- Windows: `build/bin/DubaiCrab-Go.exe`
- Linux: `build/bin/DubaiCrab-Go`

## 확장 포인트

### 새 도구 추가

1. `internal/tools/` 에 새 파일 생성
2. `Tool` 인터페이스 구현
3. `RegisterBuiltinTools()` 에 등록

### 새 LLM 프로바이더 추가

1. `internal/` 에 새 프로바이더 패키지 생성
2. `Manager` 인터페이스 구현
3. `app.go` 에서 초기화

## 성능 최적화

- 동시 요청 처리를 위한 goroutine 사용
- 세션별 메시지 히스토리 제한 (100개)
- 컨텍스트 윈도우 관리 (20개 메시지)
- HTTP 클라이언트 재사용
