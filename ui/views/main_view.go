package views

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"myvoicego/ui/components"
	"myvoicego/ui/models"
)

// MainView 表示主应用视图
type MainView struct {
	app         fyne.App
	window      fyne.Window
	tabs        *container.AppTabs
	statusBar   *components.StatusBar
	modelConfig *models.ModelConfig
	settings    *models.AppSettings
}

// NewMainView 创建新的主视图实例
func NewMainView() *MainView {
	a := app.NewWithID("com.myvoicego.mainview")
	a.Settings().SetTheme(theme.DarkTheme())

	window := a.NewWindow("MyVoiceGO - 主界面")
	window.Resize(fyne.NewSize(1200, 800))

	// 加载设置
	settings := models.NewAppSettings()
	err := settings.LoadFromFile("settings.json")
	if err != nil {
		fmt.Printf("加载设置失败: %v\n", err)
		// 使用默认设置
	}

	view := &MainView{
		app:      a,
		window:   window,
		settings: settings,
	}

	return view
}

// createUI 创建UI界面
func (v *MainView) createUI() {
	// 创建状态栏
	v.statusBar = components.NewStatusBar()

	// 创建工具栏
	toolbar := components.CreateDefaultFileToolbar(
		func() { v.newFile() },
		func() { v.openFile() },
		func() { v.saveFile() },
		func() { v.saveFileAs() },
	)

	// 创建标签页
	v.tabs = container.NewAppTabs()

	// 添加JSON查看器标签页
	jsonTab := container.NewTabItemWithIcon("JSON查看器", theme.FileIcon(), v.createJSONViewerTab())
	v.tabs.Append(jsonTab)

	// 添加处理器标签页
	processorView := NewProcessorView(v.app, v.window)
	processorTab := container.NewTabItemWithIcon("处理器", theme.ComputerIcon(), processorView.createUI())
	v.tabs.Append(processorTab)

	// 添加设置标签页
	settingsTab := container.NewTabItemWithIcon("设置", theme.SettingsIcon(), v.createSettingsTab())
	v.tabs.Append(settingsTab)

	// 创建主布局
	content := container.NewBorder(
		toolbar.Container(),
		v.statusBar.Container(),
		nil,
		nil,
		v.tabs,
	)

	v.window.SetContent(content)
}

// createTabContent 为JSON查看器创建标签页内容
func (v *MainView) createJSONViewerTab() fyne.CanvasObject {
	// 创建文件列表
	fileList := v.createFileList()

	// 创建内容区域（左侧主编辑窗口）
	leftContainer := container.NewVBox()
	leftContainer.Add(widget.NewLabel("请从左侧选择一个JSON文件"))

	// 创建右侧详情面板（初始隐藏）
	rightPanel := container.NewVBox(
		widget.NewLabelWithStyle("详情面板", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel("选择左侧项目查看详情"),
	)
	rightPanel.Hide()

	// 创建刷新按钮
	refreshButton := widget.NewButtonWithIcon("刷新列表", theme.ViewRefreshIcon(), func() {
		v.refreshFileList()
		v.statusBar.SetStatusWithTimestamp("列表已刷新")
	})

	// 创建分割视图
	splitView := container.NewHSplit(
		container.NewVBox(
			widget.NewCard("文件列表", "", fileList),
			refreshButton,
		),
		container.NewHSplit(
			leftContainer,
			rightPanel,
		),
	)
	splitView.Offset = 0.7 // 默认左面板占70%宽度

	return splitView
}

// createFileList 创建文件列表
func (v *MainView) createFileList() *widget.List {
	// 获取figures_output文件夹路径
	figuresPath := filepath.Join("..", "figures_output")

	// 初始化jsonFiles
	jsonFiles := make([]string, 0)
	entries, err := os.ReadDir(figuresPath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("无法读取目录: %w", err), v.window)
	} else {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
				jsonFiles = append(jsonFiles, entry.Name())
			}
		}
		sort.Strings(jsonFiles)
	}

	// 创建文件列表
	fileList := widget.NewList(
		func() int {
			return len(jsonFiles)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("文件名")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(jsonFiles) {
				obj.(*widget.Label).SetText(jsonFiles[id])
			}
		},
	)

	// 设置文件选择事件
	fileList.OnSelected = func(id widget.ListItemID) {
		if id < len(jsonFiles) {
			filePath := filepath.Join("..", "figures_output", jsonFiles[id])
			v.loadJSONContent(filePath)
			v.statusBar.SetStatus("已加载: " + jsonFiles[id])
		}
	}

	return fileList
}

// refreshFileList 刷新文件列表
func (v *MainView) refreshFileList() {
	// 这里可以实现刷新逻辑
	v.statusBar.SetStatus("列表已刷新")
}

// loadJSONContent 加载JSON内容
func (v *MainView) loadJSONContent(filePath string) {
	// 这里可以实现JSON内容加载逻辑
	v.statusBar.SetStatus("已加载: " + filepath.Base(filePath))
}

// createSettingsTab 创建设置标签页
func (v *MainView) createSettingsTab() fyne.CanvasObject {
	// 创建设置界面
	return container.NewScroll(
		container.NewVBox(
			widget.NewLabelWithStyle("应用设置", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),

			// 模型配置
			widget.NewCard("模型配置", "", container.NewVBox(
				widget.NewLabel("API密钥:"),
				widget.NewPasswordEntry(),
				widget.NewLabel("模型名称:"),
				widget.NewSelect([]string{"gpt-3.5-turbo", "gpt-4", "gpt-4-turbo"}, nil),
				widget.NewLabel("温度:"),
				widget.NewSlider(0.0, 1.0),
			)),

			// 保存按钮
			widget.NewButtonWithIcon("保存设置", theme.DocumentSaveIcon(), func() {
				v.saveSettings()
			}),
		),
	)
}

// selectDialogueFile 选择对话文件
func (v *MainView) selectDialogueFile() {
	dialog := components.NewOpenFileDialog(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, v.window)
			return
		}
		if reader != nil {
			defer reader.Close()
			// 处理选择的文件
			v.statusBar.SetStatus("已选择文件: " + reader.URI().Name())
		}
	}, v.window)
	dialog.ShowOpen(v.window)
}

// processDialogue 处理对话
func (v *MainView) processDialogue() {
	v.statusBar.ShowProgress()
	v.statusBar.SetProgress(0.5)

	// 这里实现对话处理逻辑
	// 可以调用processor包中的函数

	v.statusBar.HideProgress()
	v.statusBar.SetStatusWithTimestamp("对话处理完成")
}

// saveSettings 保存设置
func (v *MainView) saveSettings() {
	err := v.settings.SaveToFile("settings.json")
	if err != nil {
		dialog.ShowError(fmt.Errorf("保存设置失败: %w", err), v.window)
		return
	}

	dialog.ShowInformation("保存成功", "设置已保存", v.window)
	v.statusBar.SetStatusWithTimestamp("设置已保存")
}

// newFile 新建文件
func (v *MainView) newFile() {
	dialog := components.NewConfirmDialog(
		"新建文件",
		"确定要新建文件吗？当前未保存的更改将丢失。",
		func(confirm bool) {
			if confirm {
				// 实现新建文件逻辑
				v.statusBar.SetStatus("已创建新文件")
			}
		},
		v.window)
	dialog.Show()
}

// openFile 打开文件
func (v *MainView) openFile() {
	dialog := components.NewOpenFileDialog(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, v.window)
			return
		}
		if reader != nil {
			defer reader.Close()
			// 处理打开的文件
			v.statusBar.SetStatus("已打开文件: " + reader.URI().Name())
		}
	}, v.window)
	dialog.ShowOpen(v.window)
}

// saveFile 保存文件
func (v *MainView) saveFile() {
	v.statusBar.SetStatus("保存中...")
	v.statusBar.Refresh()

	// 这里可以实现保存逻辑

	v.statusBar.SetStatusWithTimestamp("保存完成")
}

// saveFileAs 另存为文件
func (v *MainView) saveFileAs() {
	dialog := components.NewSaveFileDialog(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, v.window)
			return
		}
		if writer != nil {
			defer writer.Close()
			// 处理保存文件
			v.statusBar.SetStatus("文件已保存")
		}
	}, v.window)
	dialog.ShowSave(v.window)
}

// Show 显示主视图
func (v *MainView) Show() {
	v.createUI()
	v.window.Show()
}

// Run 运行主视图应用
func (v *MainView) Run() {
	v.Show()
	v.app.Run()
}
