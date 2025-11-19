package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ToolbarButton 表示工具栏按钮
type ToolbarButton struct {
	Icon    fyne.Resource
	Text    string
	OnClick func()
}

// Toolbar 表示应用程序工具栏
type Toolbar struct {
	container *fyne.Container
	buttons   []*widget.Button
}

// NewToolbar 创建一个新的工具栏
func NewToolbar(buttons []ToolbarButton) *Toolbar {
	toolbarButtons := make([]*widget.Button, 0, len(buttons))

	for _, btn := range buttons {
		var button *widget.Button
		if btn.Icon != nil {
			button = widget.NewButtonWithIcon(btn.Text, btn.Icon, btn.OnClick)
		} else {
			button = widget.NewButton(btn.Text, btn.OnClick)
		}
		toolbarButtons = append(toolbarButtons, button)
	}

	// 创建工具栏容器
	toolbar := container.NewHBox()
	for _, btn := range toolbarButtons {
		toolbar.Add(btn)
	}

	return &Toolbar{
		container: toolbar,
		buttons:   toolbarButtons,
	}
}

// Container 返回工具栏容器
func (t *Toolbar) Container() *fyne.Container {
	return t.container
}

// AddButton 添加新按钮到工具栏
func (t *Toolbar) AddButton(button ToolbarButton) {
	var newButton *widget.Button
	if button.Icon != nil {
		newButton = widget.NewButtonWithIcon(button.Text, button.Icon, button.OnClick)
	} else {
		newButton = widget.NewButton(button.Text, button.OnClick)
	}

	t.buttons = append(t.buttons, newButton)
	t.container.Objects = append(t.container.Objects, newButton)
	t.container.Refresh()
}

// RemoveButton 从工具栏移除按钮
func (t *Toolbar) RemoveButton(index int) {
	if index < 0 || index >= len(t.buttons) {
		return
	}

	t.buttons = append(t.buttons[:index], t.buttons[index+1:]...)
	t.container.Objects = append(t.container.Objects[:index], t.container.Objects[index+1:]...)
	t.container.Refresh()
}

// EnableButton 启用指定索引的按钮
func (t *Toolbar) EnableButton(index int) {
	if index >= 0 && index < len(t.buttons) {
		t.buttons[index].Enable()
	}
}

// DisableButton 禁用指定索引的按钮
func (t *Toolbar) DisableButton(index int) {
	if index >= 0 && index < len(t.buttons) {
		t.buttons[index].Disable()
	}
}

// Refresh 刷新工具栏
func (t *Toolbar) Refresh() {
	t.container.Refresh()
}

// CreateDefaultFileToolbar 创建默认的文件操作工具栏
func CreateDefaultFileToolbar(onNew, onOpen, onSave, onSaveAs func()) *Toolbar {
	buttons := []ToolbarButton{
		{Icon: theme.FileIcon(), Text: "新建", OnClick: onNew},
		{Icon: theme.FileIcon(), Text: "打开", OnClick: onOpen},
		{Icon: theme.DocumentSaveIcon(), Text: "保存", OnClick: onSave},
		{Icon: theme.DocumentSaveIcon(), Text: "另存为", OnClick: onSaveAs},
	}
	return NewToolbar(buttons)
}

// CreateDefaultEditToolbar 创建默认的编辑操作工具栏
func CreateDefaultEditToolbar(onCut, onCopy, onPaste, onDelete func()) *Toolbar {
	buttons := []ToolbarButton{
		{Icon: theme.ContentCutIcon(), Text: "剪切", OnClick: onCut},
		{Icon: theme.ContentCopyIcon(), Text: "复制", OnClick: onCopy},
		{Icon: theme.ContentPasteIcon(), Text: "粘贴", OnClick: onPaste},
		{Icon: theme.DeleteIcon(), Text: "删除", OnClick: onDelete},
	}
	return NewToolbar(buttons)
}
