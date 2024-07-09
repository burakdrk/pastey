package main

import (
	"embed"
	"runtime"

	"github.com/burakdrk/pastey/pastey-wails/backend"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := backend.NewApp()

	isFrameless := true
	if runtime.GOOS == "darwin" {
		isFrameless = false
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:                    "pastey",
		Width:                    1024,
		Height:                   768,
		MinWidth:                 800,
		MinHeight:                600,
		EnableDefaultContextMenu: false,
		Frameless:                isFrameless,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 1},
		OnStartup:        app.Startup,
		OnShutdown:       app.Shutdown,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			IsZoomControlEnabled: false,
			DisablePinchZoom:     true,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   "pastey app",
				Message: "© 2024 Burak Duruk",
			},
		},
	})

	if err != nil {
		app.Logger.Log("Error: " + err.Error())
	}
}
