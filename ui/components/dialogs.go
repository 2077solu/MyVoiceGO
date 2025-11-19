package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ConfirmDialog 确认对话框
type ConfirmDialog struct {
	dialog    dialog.Dialog
	onConfirm func(bool)
}

// NewConfirmDialog 创建确认对话框
func NewConfirmDialog(title, message string, callback func(bool), parent fyne.Window) *ConfirmDialog {
	confirmDialog := &ConfirmDialog{
		onConfirm: callback,
	}

	confirmDialog.dialog = dialog.NewConfirm(title, message, func(confirm bool) {
		if callback != nil {
			callback(confirm)
		}
	}, parent)

	return confirmDialog
}

// Show 显示对话框
func (d *ConfirmDialog) Show() {
	d.dialog.Show()
}

// Hide 隐藏对话框
func (d *ConfirmDialog) Hide() {
	d.dialog.Hide()
}

// SetDismissText 设置关闭按钮文本
func (d *ConfirmDialog) SetDismissText(text string) {
	d.dialog.SetDismissText(text)
}

// InputDialog 输入对话框
type InputDialog struct {
	dialog   dialog.Dialog
	onSubmit func(string)
	entry    *widget.Entry
}

// NewInputDialog 创建输入对话框
func NewInputDialog(title, message string, callback func(string), parent fyne.Window) *InputDialog {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("请输入内容")

	inputDialog := &InputDialog{
		onSubmit: callback,
		entry:    entry,
	}

	inputDialog.dialog = dialog.NewForm(title, message, "确定", []*widget.FormItem{
		widget.NewFormItem("", entry),
	}, func(b bool) {
		if b && callback != nil {
			callback(entry.Text)
		}
	}, parent)

	return inputDialog
}

// Show 显示对话框
func (d *InputDialog) Show() {
	d.dialog.Show()
}

// Hide 隐藏对话框
func (d *InputDialog) Hide() {
	d.dialog.Hide()
}

// SetText 设置输入框默认文本
func (d *InputDialog) SetText(text string) {
	d.entry.SetText(text)
}

// GetText 获取输入框文本
func (d *InputDialog) GetText() string {
	return d.entry.Text
}

// SetPlaceHolder 设置输入框占位符
func (d *InputDialog) SetPlaceHolder(text string) {
	d.entry.SetPlaceHolder(text)
}

// ProgressDialog 进度对话框
type ProgressDialog struct {
	dialog      dialog.Dialog
	progressBar *widget.ProgressBar
	statusLabel *widget.Label
}

// NewProgressDialog 创建进度对话框
func NewProgressDialog(title string, parent fyne.Window) *ProgressDialog {
	progressBar := widget.NewProgressBar()
	statusLabel := widget.NewLabel("处理中...")

	content := container.NewVBox(
		statusLabel,
		progressBar,
	)

	progressDialog := &ProgressDialog{
		progressBar: progressBar,
		statusLabel: statusLabel,
	}

	progressDialog.dialog = dialog.NewCustom(title, "取消", content, parent)

	return progressDialog
}

// Show 显示对话框
func (d *ProgressDialog) Show() {
	d.dialog.Show()
}

// Hide 隐藏对话框
func (d *ProgressDialog) Hide() {
	d.dialog.Hide()
}

// SetProgress 设置进度值 (0.0 到 1.0)
func (d *ProgressDialog) SetProgress(value float64) {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	d.progressBar.SetValue(value)
}

// SetStatus 设置状态文本
func (d *ProgressDialog) SetStatus(text string) {
	d.statusLabel.SetText(text)
}

// FileDialog 文件对话框
type FileDialog struct {
	onOpen func(fyne.URIReadCloser, error)
	onSave func(fyne.URIWriteCloser, error)
}

// NewOpenFileDialog 创建打开文件对话框
func NewOpenFileDialog(callback func(fyne.URIReadCloser, error), parent fyne.Window) *FileDialog {
	return &FileDialog{onOpen: callback}
}

// NewSaveFileDialog 创建保存文件对话框
func NewSaveFileDialog(callback func(fyne.URIWriteCloser, error), parent fyne.Window) *FileDialog {
	return &FileDialog{onSave: callback}
}

// ShowOpen 显示打开文件对话框
func (d *FileDialog) ShowOpen(parent fyne.Window) {
	dialog.ShowFileOpen(d.onOpen, parent)
}

// ShowSave 显示保存文件对话框
func (d *FileDialog) ShowSave(parent fyne.Window) {
	if d.onSave != nil {
		dialog.ShowFileSave(d.onSave, parent)
	}
}
