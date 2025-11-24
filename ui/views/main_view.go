package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// MainView 表示主应用视图
type MainView struct {
	app    fyne.App
	window fyne.Window
	tabs   *container.AppTabs
}

// NewMainView 创建新的主视图实例
func NewMainView() *MainView {
	a := app.NewWithID("com.myvoicego.mainview")
	a.Settings().SetTheme(theme.DefaultTheme()) // 使用白色主题

	window := a.NewWindow("MyVoiceGO - 简化界面")
	window.Resize(fyne.NewSize(800, 600))

	view := &MainView{
		app:    a,
		window: window,
	}

	return view
}

// createUI 创建UI界面
func (v *MainView) createUI() {
	// 创建三个简单的页面，每个页面只显示"Hello"
	page1 := widget.NewLabel("Hello - Page 1")
	page1.Alignment = fyne.TextAlignCenter

	page2 := widget.NewLabel("Hello - Page 2")
	page2.Alignment = fyne.TextAlignCenter

	page3 := widget.NewLabel("Hello - Page 3")
	page3.Alignment = fyne.TextAlignCenter

	// 创建标签页
	v.tabs = container.NewAppTabs(
		container.NewTabItem("", container.NewCenter(page1)),
		container.NewTabItem("", container.NewCenter(page2)),
		container.NewTabItem("", container.NewCenter(page3)),
	)

	// 隐藏标签页的标题栏
	v.tabs.SetTabLocation(container.TabLocationBottom)

	// 创建页面切换按钮
	pageButtons := container.NewHBox(
		widget.NewButton("Page 1", func() {
			v.tabs.SelectIndex(0)
		}),
		widget.NewButton("Page 2", func() {
			v.tabs.SelectIndex(1)
		}),
		widget.NewButton("Page 3", func() {
			v.tabs.SelectIndex(2)
		}),
	)

	// 创建主布局，将页面切换按钮放在左下角
	content := container.NewBorder(
		nil, // 顶部
		nil, // 底部
		nil, // 左侧
		nil, // 右侧
		v.tabs,
	)

	// 创建底部容器，将页面切换按钮放在左侧
	bottomContainer := container.NewHBox(
		pageButtons,
		widget.NewLabel(""), // 占位符，将按钮推到左侧
	)

	// 使用Border布局，将底部容器放在底部
	mainContainer := container.NewBorder(
		nil,             // 顶部
		bottomContainer, // 底部
		nil,             // 左侧
		nil,             // 右侧
		content,         // 中心内容
	)

	v.window.SetContent(mainContainer)
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
