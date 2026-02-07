#!/bin/bash
# OpenKlaw - Ollama 자동 설치 스크립트 (macOS/Linux)

set -e

echo "🦞 OpenKlaw - Ollama 설치 확인 중..."

# Ollama 설치 확인
if command -v ollama &> /dev/null; then
    echo "✅ Ollama가 이미 설치되어 있습니다."
    ollama --version
else
    echo "📦 Ollama 설치 중..."
    
    # OS 감지
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command -v brew &> /dev/null; then
            echo "🍺 Homebrew로 Ollama 설치 중..."
            brew install ollama
        else
            echo "⚠️ Homebrew가 설치되어 있지 않습니다."
            echo "📥 공식 스크립트로 설치 중..."
            curl -fsSL https://ollama.ai/install.sh | sh
        fi
    else
        # Linux
        echo "🐧 Linux에 Ollama 설치 중..."
        curl -fsSL https://ollama.ai/install.sh | sh
    fi
fi

# Ollama 서비스 시작
echo "🚀 Ollama 서비스 시작 중..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS: brew services 사용
    if command -v brew &> /dev/null; then
        brew services start ollama 2>/dev/null || true
    else
        ollama serve &
    fi
else
    # Linux: systemd 사용
    sudo systemctl start ollama 2>/dev/null || ollama serve &
fi

# 서비스 시작 대기
sleep 3

# 연결 확인
if curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
    echo "✅ Ollama 서버가 실행 중입니다."
else
    echo "⚠️ Ollama 서버 시작 중... 잠시 기다려주세요."
    sleep 5
fi

# 모델 확인 및 다운로드
MODEL="qwen2.5:3b-instruct"
echo ""
echo "🤖 기본 모델 확인 중: $MODEL"

if ollama list | grep -q "qwen2.5:3b"; then
    echo "✅ 모델이 이미 설치되어 있습니다."
else
    echo "📥 모델 다운로드 중... (약 2GB, 네트워크 상태에 따라 5-15분 소요)"
    ollama pull $MODEL
    echo "✅ 모델 다운로드 완료!"
fi

echo ""
echo "🎉 OpenKlaw 설정 완료!"
echo "   openklaw 명령어로 실행하세요."
