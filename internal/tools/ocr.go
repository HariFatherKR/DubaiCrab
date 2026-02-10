package tools

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// OcrBox represents a detected text box
type OcrBox struct {
	Text       string      `json:"text"`
	Confidence float64     `json:"confidence"`
	BoxCoords  [][]float64 `json:"box_coords"`
}

// OcrResult represents OCR result
type OcrResult struct {
	Success   bool      `json:"success"`
	Text      *string   `json:"text,omitempty"`
	Error     *string   `json:"error,omitempty"`
	Boxes     []OcrBox  `json:"boxes,omitempty"`
	LineCount *int      `json:"line_count,omitempty"`
}

// Allowed image extensions for security
var allowedImageExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".webp": true,
	".bmp":  true,
	".tiff": true,
	".tif":  true,
}

// Maximum base64 size (10MB)
const maxBase64Size = 10 * 1024 * 1024

// OcrTool performs OCR on images using tesseract
type OcrTool struct{}

func (t *OcrTool) Name() string        { return "ocr" }
func (t *OcrTool) Description() string { return "이미지에서 텍스트를 추출합니다 (OCR)" }
func (t *OcrTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "이미지 파일 경로",
			},
			"lang": map[string]interface{}{
				"type":        "string",
				"description": "OCR 언어 (기본: kor+eng)",
				"default":     "kor+eng",
			},
		},
		"required": []string{"path"},
	}
}

func (t *OcrTool) Execute(ctx context.Context, params map[string]interface{}) (string, error) {
	path, _ := params["path"].(string)
	lang, _ := params["lang"].(string)

	if lang == "" {
		lang = "kor+eng"
	}

	result := OcrFromFile(path, lang)
	if !result.Success {
		if result.Error != nil {
			return "", fmt.Errorf(*result.Error)
		}
		return "", fmt.Errorf("OCR 실패")
	}

	if result.Text != nil {
		return *result.Text, nil
	}
	return "", nil
}

// OcrFromFile performs OCR on an image file
func OcrFromFile(imagePath, lang string) OcrResult {
	// Validate path
	validPath, err := validateImagePath(imagePath)
	if err != nil {
		errStr := err.Error()
		return OcrResult{Success: false, Error: &errStr}
	}

	if lang == "" {
		lang = "kor+eng"
	}

	// Check if tesseract is available
	_, err = exec.LookPath("tesseract")
	if err != nil {
		errStr := "tesseract가 설치되지 않았습니다. brew install tesseract tesseract-lang 으로 설치하세요"
		return OcrResult{Success: false, Error: &errStr}
	}

	// Run tesseract
	cmd := exec.Command("tesseract", validPath, "stdout", "-l", lang)
	output, err := cmd.Output()
	if err != nil {
		errStr := fmt.Sprintf("OCR 실행 실패: %v", err)
		return OcrResult{Success: false, Error: &errStr}
	}

	text := strings.TrimSpace(string(output))
	lines := strings.Split(text, "\n")
	lineCount := len(lines)

	return OcrResult{
		Success:   true,
		Text:      &text,
		LineCount: &lineCount,
	}
}

// OcrFromBase64 performs OCR on a base64-encoded image
func OcrFromBase64(base64Data, lang string) OcrResult {
	// Size limit check
	if len(base64Data) > maxBase64Size {
		errStr := fmt.Sprintf("이미지 크기가 제한을 초과합니다 (최대 10MB)")
		return OcrResult{Success: false, Error: &errStr}
	}

	// Strip data URL prefix if present
	if strings.HasPrefix(base64Data, "data:") {
		parts := strings.SplitN(base64Data, ",", 2)
		if len(parts) == 2 {
			base64Data = parts[1]
		}
	}

	// Decode base64
	imgData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		errStr := fmt.Sprintf("Base64 디코딩 실패: %v", err)
		return OcrResult{Success: false, Error: &errStr}
	}

	// Create temp file
	tmpFile, err := os.CreateTemp("", "ocr_*.png")
	if err != nil {
		errStr := fmt.Sprintf("임시 파일 생성 실패: %v", err)
		return OcrResult{Success: false, Error: &errStr}
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write image data
	if _, err := tmpFile.Write(imgData); err != nil {
		errStr := fmt.Sprintf("이미지 저장 실패: %v", err)
		return OcrResult{Success: false, Error: &errStr}
	}
	tmpFile.Close()

	return OcrFromFile(tmpFile.Name(), lang)
}

// validateImagePath validates the image path for security
func validateImagePath(imagePath string) (string, error) {
	// Expand home directory
	if strings.HasPrefix(imagePath, "~") {
		home, _ := os.UserHomeDir()
		imagePath = filepath.Join(home, imagePath[1:])
	}

	// Get absolute path
	absPath, err := filepath.Abs(imagePath)
	if err != nil {
		return "", fmt.Errorf("잘못된 경로입니다: %v", err)
	}

	// Check file exists
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("파일이 존재하지 않습니다: %s", absPath)
		}
		return "", fmt.Errorf("파일 접근 오류: %v", err)
	}

	if info.IsDir() {
		return "", fmt.Errorf("디렉토리가 아닌 파일이 필요합니다")
	}

	// Validate extension
	ext := strings.ToLower(filepath.Ext(absPath))
	if !allowedImageExtensions[ext] {
		return "", fmt.Errorf("지원하지 않는 이미지 형식입니다: %s", ext)
	}

	// Validate path is in allowed directories
	home, _ := os.UserHomeDir()
	allowedPrefixes := []string{
		home,
		"/tmp",
		"/var/folders",
		"/private/tmp",
	}

	allowed := false
	for _, prefix := range allowedPrefixes {
		if strings.HasPrefix(absPath, prefix) {
			allowed = true
			break
		}
	}

	if !allowed {
		return "", fmt.Errorf("허용되지 않은 디렉토리의 파일입니다")
	}

	return absPath, nil
}

// MarshalOcrResult converts OcrResult to JSON string
func MarshalOcrResult(result OcrResult) string {
	data, _ := json.Marshal(result)
	return string(data)
}

func init() {
	// This will be registered via RegisterOcrTools
}

// RegisterOcrTools registers OCR tools
func RegisterOcrTools(registry *Registry) {
	registry.Register(&OcrTool{})
}
