package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application menu
	appMenu := menu.NewMenu()

	// File menu
	fileMenu := appMenu.AddSubmenu("파일")
	fileMenu.AddText("새 대화", keys.CmdOrCtrl("n"), func(_ *menu.CallbackData) {
		// Handle new conversation
	})
	fileMenu.AddSeparator()
	fileMenu.AddText("종료", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		// wails.Quit()
	})

	// Edit menu
	editMenu := appMenu.AddSubmenu("편집")
	editMenu.AddText("실행취소", keys.CmdOrCtrl("z"), nil)
	editMenu.AddText("다시실행", keys.CmdOrCtrl("shift+z"), nil)
	editMenu.AddSeparator()
	editMenu.AddText("잘라내기", keys.CmdOrCtrl("x"), nil)
	editMenu.AddText("복사", keys.CmdOrCtrl("c"), nil)
	editMenu.AddText("붙여넣기", keys.CmdOrCtrl("v"), nil)
	editMenu.AddText("전체선택", keys.CmdOrCtrl("a"), nil)

	// View menu
	viewMenu := appMenu.AddSubmenu("보기")
	viewMenu.AddText("새로고침", keys.CmdOrCtrl("r"), nil)
	viewMenu.AddSeparator()
	viewMenu.AddText("개발자 도구", keys.CmdOrCtrl("shift+i"), nil)

	// Help menu
	helpMenu := appMenu.AddSubmenu("도움말")
	helpMenu.AddText("Dubai Crab 정보", nil, func(_ *menu.CallbackData) {
		// Show about dialog
	})
	helpMenu.AddText("문서", nil, func(_ *menu.CallbackData) {
		// Open documentation
	})

	// Create application with options
	err := wails.Run(&options.App{
		Title:         "Dubai Crab",
		Width:         1200,
		Height:        800,
		MinWidth:      800,
		MinHeight:     600,
		DisableResize: false,
		Fullscreen:    false,
		Frameless:     false,
		StartHidden:   false,
		Menu:          appMenu,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		// Dark chocolate theme color
		BackgroundColour: &options.RGBA{R: 62, G: 39, B: 35, A: 255},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
		// macOS specific options
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            true,
				UseToolbar:                 false,
			},
			About: &mac.AboutInfo{
				Title:   "Dubai Crab",
				Message: "한국 사무직을 위한 로컬 AI 비서\n\n버전 0.1.0",
			},
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
		// Windows specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
	})

	if err != nil {
		log.Fatal("Error:", err.Error())
	}
}
