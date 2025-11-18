package ui

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func MainView() {
	// 创建应用
	myApp := app.New()
	myApp.Settings().SetTheme(myApp.Settings().Theme())
	window := myApp.NewWindow("MyVoiceGO - JSON查看器")
	window.Resize(fyne.NewSize(1200, 800))

	// 检查figures_output目录是否存在
	figuresPath := filepath.Join("..", "figures_output")
	if _, err := os.Stat(figuresPath); os.IsNotExist(err) {
		// 如果目录不存在，显示错误信息
		errorLabel := widget.NewLabel("错误: 找不到figures_output目录")
		window.SetContent(container.NewVBox(
			widget.NewCard("", "错误", errorLabel),
		))
		window.ShowAndRun()
		return
	}

	// 创建欢迎界面
	welcomeLabel := widget.NewLabelWithStyle("欢迎使用MyVoiceGO JSON查看器", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	descriptionLabel := widget.NewLabel("此工具用于查看和编辑figures_output文件夹中的JSON文件")
	startButton := widget.NewButton("启动JSON查看器", func() {
		// 启动JSON查看器
		JSONViewer()
		window.Close()
	})

	// 创建主界面
	content := container.NewVBox(
		widget.NewCard("", "欢迎", container.NewVBox(
			welcomeLabel,
			descriptionLabel,
			startButton,
		)),
	)

	window.SetContent(content)
	window.ShowAndRun()
}
