package views

import (
	"encoding/json"
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
)

// JSONViewer 表示JSON查看器应用
type JSONViewer struct {
	app               fyne.App
	window            fyne.Window
	fileList          *widget.List
	leftContainer     *fyne.Container
	rightPanel        *fyne.Container
	statusBar         *components.StatusBar
	jsonFiles         []string
	hasUnsavedChanges bool
}

// NewJSONViewer 创建新的JSON查看器实例
func NewJSONViewer() *JSONViewer {
	a := app.NewWithID("com.myvoicego.jsonviewer")
	a.Settings().SetTheme(theme.DarkTheme())

	window := a.NewWindow("MyVoiceGO - JSON查看器")
	window.Resize(fyne.NewSize(1200, 800))

	viewer := &JSONViewer{
		app:       a,
		window:    window,
		jsonFiles: make([]string, 0),
	}

	// 设置窗口关闭拦截
	window.SetCloseIntercept(func() {
		if viewer.hasUnsavedChanges {
			confirmDialog := components.NewConfirmDialog(
				"确认关闭",
				"有未保存的修改，确定要关闭吗？",
				func(confirm bool) {
					if confirm {
						window.Close()
					}
				},
				window)
			confirmDialog.Show()
		} else {
			window.Close()
		}
	})

	return viewer
}

// createUI 创建UI界面
func (v *JSONViewer) createUI() {
	// 创建状态栏
	v.statusBar = components.NewStatusBar()

	// 创建文件列表
	v.fileList = v.createFileList()

	// 创建内容区域（左侧主编辑窗口）
	v.leftContainer = container.NewVBox()
	v.leftContainer.Add(widget.NewLabel("请从左侧选择一个JSON文件"))

	// 创建右侧详情面板（初始隐藏）
	v.rightPanel = container.NewVBox(
		widget.NewLabelWithStyle("详情面板", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel("选择左侧项目查看详情"),
	)
	v.rightPanel.Hide()

	// 创建工具栏
	toolbar := components.CreateDefaultFileToolbar(
		func() { /* 新建 */ },
		func() { v.openFile() },
		func() { v.saveFile() },
		func() { v.saveFileAs() },
	)

	// 创建刷新按钮
	toolbar.AddButton(components.ToolbarButton{
		Icon:    theme.ViewRefreshIcon(),
		Text:    "刷新列表",
		OnClick: func() { v.refreshFileList() },
	})

	// 创建分割视图
	splitView := container.NewHSplit(
		container.NewVBox(
			widget.NewCard("文件列表", "", v.fileList),
		),
		container.NewHSplit(
			v.leftContainer,
			v.rightPanel,
		),
	)
	splitView.Offset = 0.7 // 默认左面板占70%宽度

	// 创建主布局
	content := container.NewBorder(
		toolbar.Container(),
		v.statusBar.Container(),
		nil,
		nil,
		splitView,
	)

	v.window.SetContent(content)
}

// createFileList 创建文件列表
func (v *JSONViewer) createFileList() *widget.List {
	// 获取figures_output文件夹路径
	figuresPath := filepath.Join("..", "figures_output")

	// 初始化jsonFiles
	v.loadJSONFiles(figuresPath)

	// 创建文件列表
	fileList := widget.NewList(
		func() int {
			return len(v.jsonFiles)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("文件名")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(v.jsonFiles) {
				obj.(*widget.Label).SetText(v.jsonFiles[id])
			}
		},
	)

	// 设置文件选择事件
	fileList.OnSelected = func(id widget.ListItemID) {
		v.loadSelectedFile(id)
	}

	return fileList
}

// loadJSONFiles 加载JSON文件列表
func (v *JSONViewer) loadJSONFiles(figuresPath string) {
	entries, err := os.ReadDir(figuresPath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("无法读取目录: %w", err), v.window)
		return
	}

	v.jsonFiles = make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			v.jsonFiles = append(v.jsonFiles, entry.Name())
		}
	}
	sort.Strings(v.jsonFiles)
}

// refreshFileList 刷新文件列表
func (v *JSONViewer) refreshFileList() {
	figuresPath := filepath.Join("..", "figures_output")
	v.loadJSONFiles(figuresPath)
	v.fileList.Refresh()
}

// loadSelectedFile 加载选中的文件
func (v *JSONViewer) loadSelectedFile(id widget.ListItemID) {
	if v.hasUnsavedChanges {
		confirmDialog := components.NewConfirmDialog(
			"确认切换文件",
			"有未保存的修改，确定要切换文件吗？",
			func(confirm bool) {
				if confirm {
					v.doLoadSelectedFile(id)
					v.hasUnsavedChanges = false
				}
			},
			v.window)
		confirmDialog.Show()
	} else {
		v.doLoadSelectedFile(id)
	}
}

// doLoadSelectedFile 实际加载选中文件的逻辑
func (v *JSONViewer) doLoadSelectedFile(id widget.ListItemID) {
	if id >= len(v.jsonFiles) {
		return
	}

	filePath := filepath.Join("..", "figures_output", v.jsonFiles[id])
	v.loadJSONContentWithExpandableView(filePath)
	v.statusBar.SetStatus("已加载: " + v.jsonFiles[id])
}

// loadJSONContentWithExpandableView 加载JSON内容并以可展开视图显示
func (v *JSONViewer) loadJSONContentWithExpandableView(filePath string) {
	// 读取JSON文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("无法读取文件: %w", err), v.window)
		return
	}

	// 解析JSON
	var jsonData interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		dialog.ShowError(fmt.Errorf("JSON解析错误: %w", err), v.window)
		return
	}

	// 清空左侧容器
	v.leftContainer.RemoveAll()

	// 创建根节点
	tree := v.createJSONTree(jsonData, func(key string, value interface{}) {
		// 更新右侧详情面板
		v.updateDetailPanel(key, value)
	})

	// 添加到左侧容器
	v.leftContainer.Add(tree)
	v.leftContainer.Refresh()
}

// createJSONTree 创建JSON树形结构
func (v *JSONViewer) createJSONTree(data interface{}, onClick func(string, interface{})) fyne.CanvasObject {
	switch data := data.(type) {
	case map[string]interface{}:
		// 处理对象
		nodes := make([]fyne.CanvasObject, 0, len(data))
		for key, value := range data {
			node := components.NewExpandableNode(key, v.createJSONTree(value, onClick))
			node.OnClick = func(n *components.ExpandableNode) {
				if onClick != nil {
					onClick(n.Title, data)
				}
			}
			nodes = append(nodes, node.BuildUI())
		}
		return container.NewVBox(nodes...)

	case []interface{}:
		// 处理数组
		nodes := make([]fyne.CanvasObject, 0, len(data))
		for i, value := range data {
			node := components.NewExpandableNode(fmt.Sprintf("[%d]", i), v.createJSONTree(value, onClick))
			node.OnClick = func(n *components.ExpandableNode) {
				if onClick != nil {
					onClick(n.Title, data)
				}
			}
			nodes = append(nodes, node.BuildUI())
		}
		return container.NewVBox(nodes...)

	default:
		// 处理基本类型
		return widget.NewLabel(fmt.Sprintf("%v", data))
	}
}

// updateDetailPanel 更新详情面板
func (v *JSONViewer) updateDetailPanel(key string, value interface{}) {
	v.rightPanel.RemoveAll()
	v.rightPanel.Add(widget.NewLabelWithStyle("详情: "+key, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	v.rightPanel.Add(widget.NewLabel(fmt.Sprintf("值: %v", value)))

	// 添加编辑按钮
	editButton := widget.NewButton("编辑", func() {
		// 这里可以添加编辑功能
	})
	v.rightPanel.Add(editButton)

	// 添加保存按钮
	saveButton := widget.NewButton("保存", func() {
		v.hasUnsavedChanges = true
		v.statusBar.SetStatus("有未保存的修改")
	})
	v.rightPanel.Add(saveButton)

	v.rightPanel.Show()
}

// openFile 打开文件
func (v *JSONViewer) openFile() {
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
func (v *JSONViewer) saveFile() {
	if v.hasUnsavedChanges {
		v.statusBar.SetStatus("保存中...")
		v.statusBar.Refresh()

		// 这里可以实现保存逻辑
		v.statusBar.SetStatusWithTimestamp("保存完成")
		v.hasUnsavedChanges = false
	} else {
		v.statusBar.SetStatus("没有需要保存的修改")
	}
}

// saveFileAs 另存为文件
func (v *JSONViewer) saveFileAs() {
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

// Show 显示JSON查看器
func (v *JSONViewer) Show() {
	v.createUI()
	v.window.Show()
}

// Run 运行JSON查看器应用
func (v *JSONViewer) Run() {
	v.Show()
	v.app.Run()
}
