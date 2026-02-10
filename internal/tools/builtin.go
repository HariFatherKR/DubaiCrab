package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// SystemInfoTool returns system information
type SystemInfoTool struct{}

func (t *SystemInfoTool) Name() string        { return "system_info" }
func (t *SystemInfoTool) Description() string { return "시스템 정보를 조회합니다" }
func (t *SystemInfoTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
}

func (t *SystemInfoTool) Execute(ctx context.Context, params map[string]interface{}) (string, error) {
	info := map[string]string{
		"os":      runtime.GOOS,
		"arch":    runtime.GOARCH,
		"cpus":    fmt.Sprintf("%d", runtime.NumCPU()),
		"version": runtime.Version(),
	}
	
	// Get hostname
	if hostname, err := os.Hostname(); err == nil {
		info["hostname"] = hostname
	}
	
	// Get home directory
	if home, err := os.UserHomeDir(); err == nil {
		info["home"] = home
	}
	
	result, _ := json.MarshalIndent(info, "", "  ")
	return string(result), nil
}

// ClipboardTool handles clipboard operations
type ClipboardTool struct{}

func (t *ClipboardTool) Name() string        { return "clipboard" }
func (t *ClipboardTool) Description() string { return "클립보드에 텍스트를 복사하거나 읽습니다" }
func (t *ClipboardTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"description": "copy 또는 paste",
				"enum":        []string{"copy", "paste"},
			},
			"text": map[string]interface{}{
				"type":        "string",
				"description": "복사할 텍스트 (action=copy일 때)",
			},
		},
		"required": []string{"action"},
	}
}

func (t *ClipboardTool) Execute(ctx context.Context, params map[string]interface{}) (string, error) {
	action, _ := params["action"].(string)
	
	switch action {
	case "copy":
		text, _ := params["text"].(string)
		if text == "" {
			return "", fmt.Errorf("복사할 텍스트가 필요합니다")
		}
		return t.copyToClipboard(text)
	case "paste":
		return t.pasteFromClipboard()
	default:
		return "", fmt.Errorf("알 수 없는 action: %s", action)
	}
}

func (t *ClipboardTool) copyToClipboard(text string) (string, error) {
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case "windows":
		cmd = exec.Command("powershell", "-Command", fmt.Sprintf("Set-Clipboard -Value '%s'", strings.ReplaceAll(text, "'", "''")))
		if err := cmd.Run(); err != nil {
			return "", err
		}
		return "클립보드에 복사되었습니다", nil
	default:
		return "", fmt.Errorf("지원되지 않는 운영체제입니다")
	}
	
	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	
	return "클립보드에 복사되었습니다", nil
}

func (t *ClipboardTool) pasteFromClipboard() (string, error) {
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbpaste")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard", "-o")
	case "windows":
		cmd = exec.Command("powershell", "-Command", "Get-Clipboard")
	default:
		return "", fmt.Errorf("지원되지 않는 운영체제입니다")
	}
	
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(string(output)), nil
}

// OpenURLTool opens URLs in the default browser
type OpenURLTool struct{}

func (t *OpenURLTool) Name() string        { return "open_url" }
func (t *OpenURLTool) Description() string { return "URL을 기본 브라우저에서 엽니다" }
func (t *OpenURLTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"url": map[string]interface{}{
				"type":        "string",
				"description": "열 URL",
			},
		},
		"required": []string{"url"},
	}
}

func (t *OpenURLTool) Execute(ctx context.Context, params map[string]interface{}) (string, error) {
	url, _ := params["url"].(string)
	if url == "" {
		return "", fmt.Errorf("URL이 필요합니다")
	}
	
	// Security check
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "", fmt.Errorf("허용되지 않는 URL 프로토콜입니다")
	}
	
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return "", fmt.Errorf("지원되지 않는 운영체제입니다")
	}
	
	if err := cmd.Start(); err != nil {
		return "", err
	}
	
	return fmt.Sprintf("URL을 열었습니다: %s", url), nil
}

// HWPParserTool parses HWP files
type HWPParserTool struct{}

func (t *HWPParserTool) Name() string        { return "parse_hwp" }
func (t *HWPParserTool) Description() string { return "HWP 파일에서 텍스트를 추출합니다" }
func (t *HWPParserTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "HWP 파일 경로",
			},
			"include_tables": map[string]interface{}{
				"type":        "boolean",
				"description": "표 포함 여부",
				"default":     false,
			},
		},
		"required": []string{"path"},
	}
}

func (t *HWPParserTool) Execute(ctx context.Context, params map[string]interface{}) (string, error) {
	path, _ := params["path"].(string)
	if path == "" {
		return "", fmt.Errorf("파일 경로가 필요합니다")
	}
	
	// Expand home directory
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[1:])
	}
	
	// Check file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("파일이 존재하지 않습니다: %s", path)
	}
	
	includeTables, _ := params["include_tables"].(bool)
	
	subcmd := "text"
	if includeTables {
		subcmd = "rich-text"
	}
	
	cmd := exec.CommandContext(ctx, "hwpparser", subcmd, path)
	output, err := cmd.Output()
	if err != nil {
		// Check if hwpparser is not installed
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "executable file not found") {
			return "", fmt.Errorf("hwpparser가 설치되지 않았습니다. pip install -e ~/Documents/snovium/hwp-parser 로 설치하세요")
		}
		return "", fmt.Errorf("HWP 파싱 실패: %v", err)
	}
	
	return string(output), nil
}

// HWPToPDFTool converts HWP files to PDF
type HWPToPDFTool struct{}

func (t *HWPToPDFTool) Name() string        { return "hwp_to_pdf" }
func (t *HWPToPDFTool) Description() string { return "HWP 파일을 PDF로 변환합니다" }
func (t *HWPToPDFTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"input_path": map[string]interface{}{
				"type":        "string",
				"description": "입력 HWP 파일 경로",
			},
			"output_path": map[string]interface{}{
				"type":        "string",
				"description": "출력 PDF 파일 경로",
			},
		},
		"required": []string{"input_path", "output_path"},
	}
}

func (t *HWPToPDFTool) Execute(ctx context.Context, params map[string]interface{}) (string, error) {
	inputPath, _ := params["input_path"].(string)
	outputPath, _ := params["output_path"].(string)

	if inputPath == "" || outputPath == "" {
		return "", fmt.Errorf("입력 경로와 출력 경로가 필요합니다")
	}

	// Expand home directory
	if strings.HasPrefix(inputPath, "~") {
		home, _ := os.UserHomeDir()
		inputPath = filepath.Join(home, inputPath[1:])
	}
	if strings.HasPrefix(outputPath, "~") {
		home, _ := os.UserHomeDir()
		outputPath = filepath.Join(home, outputPath[1:])
	}

	// Check input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return "", fmt.Errorf("입력 파일이 존재하지 않습니다: %s", inputPath)
	}

	// Try hwpparser convert first
	cmd := exec.CommandContext(ctx, "hwpparser", "convert", inputPath, outputPath)
	output, err := cmd.CombinedOutput()
	if err == nil {
		return fmt.Sprintf("PDF 변환 완료: %s", outputPath), nil
	}

	// Fallback to LibreOffice
	lofficePath := "libreoffice"
	if runtime.GOOS == "darwin" {
		lofficePath = "/Applications/LibreOffice.app/Contents/MacOS/soffice"
	}

	outputDir := filepath.Dir(outputPath)
	cmd = exec.CommandContext(ctx, lofficePath, "--headless", "--convert-to", "pdf", "--outdir", outputDir, inputPath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("PDF 변환 실패: %s\n%s", err.Error(), string(output))
	}

	// LibreOffice generates filename with .pdf extension
	expectedOutput := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath)) + ".pdf"
	expectedPath := filepath.Join(outputDir, expectedOutput)

	// Rename if output path is different
	if expectedPath != outputPath {
		if err := os.Rename(expectedPath, outputPath); err != nil {
			return fmt.Sprintf("PDF 생성됨 (경로: %s)", expectedPath), nil
		}
	}

	return fmt.Sprintf("PDF 변환 완료: %s", outputPath), nil
}

// RegisterBuiltinTools registers all built-in tools
func RegisterBuiltinTools(registry *Registry) {
	registry.Register(&SystemInfoTool{})
	registry.Register(&ClipboardTool{})
	registry.Register(&OpenURLTool{})
	registry.Register(&HWPParserTool{})
	registry.Register(&HWPToPDFTool{})
}
