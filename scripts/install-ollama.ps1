# OpenKlaw - Ollama 자동 설치 스크립트 (Windows PowerShell)

Write-Host "🦞 OpenKlaw - Ollama 설치 확인 중..." -ForegroundColor Cyan

# Ollama 설치 확인
$ollamaPath = Get-Command ollama -ErrorAction SilentlyContinue

if ($ollamaPath) {
    Write-Host "✅ Ollama가 이미 설치되어 있습니다." -ForegroundColor Green
    ollama --version
} else {
    Write-Host "📦 Ollama 설치 중..." -ForegroundColor Yellow
    
    # Ollama Windows 설치 파일 다운로드
    $installerUrl = "https://ollama.ai/download/OllamaSetup.exe"
    $installerPath = "$env:TEMP\OllamaSetup.exe"
    
    Write-Host "📥 설치 파일 다운로드 중..."
    Invoke-WebRequest -Uri $installerUrl -OutFile $installerPath
    
    Write-Host "🔧 Ollama 설치 중... (관리자 권한이 필요할 수 있습니다)"
    Start-Process -FilePath $installerPath -Wait
    
    # 환경 변수 새로고침
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
}

# Ollama 서비스 시작 확인
Write-Host ""
Write-Host "🚀 Ollama 서비스 시작 중..." -ForegroundColor Cyan

# 서비스가 실행 중인지 확인
$running = $false
for ($i = 0; $i -lt 10; $i++) {
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:11434/api/tags" -UseBasicParsing -TimeoutSec 2
        if ($response.StatusCode -eq 200) {
            $running = $true
            break
        }
    } catch {
        Start-Sleep -Seconds 2
    }
}

if (-not $running) {
    Write-Host "⚠️ Ollama 서비스를 수동으로 시작합니다..."
    Start-Process ollama -ArgumentList "serve" -WindowStyle Hidden
    Start-Sleep -Seconds 5
}

# 연결 확인
try {
    $response = Invoke-WebRequest -Uri "http://localhost:11434/api/tags" -UseBasicParsing -TimeoutSec 5
    Write-Host "✅ Ollama 서버가 실행 중입니다." -ForegroundColor Green
} catch {
    Write-Host "⚠️ Ollama 서버에 연결할 수 없습니다. 시스템을 재시작해보세요." -ForegroundColor Yellow
}

# 모델 확인 및 다운로드
$model = "qwen2.5:3b-instruct"
Write-Host ""
Write-Host "🤖 기본 모델 확인 중: $model" -ForegroundColor Cyan

$modelList = ollama list 2>&1
if ($modelList -match "qwen2.5:3b") {
    Write-Host "✅ 모델이 이미 설치되어 있습니다." -ForegroundColor Green
} else {
    Write-Host "📥 모델 다운로드 중... (약 2GB, 네트워크 상태에 따라 5-15분 소요)" -ForegroundColor Yellow
    ollama pull $model
    Write-Host "✅ 모델 다운로드 완료!" -ForegroundColor Green
}

Write-Host ""
Write-Host "🎉 OpenKlaw 설정 완료!" -ForegroundColor Cyan
Write-Host "   OpenKlaw 앱을 실행하세요." -ForegroundColor White
