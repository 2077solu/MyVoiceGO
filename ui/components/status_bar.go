package components

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// StatusBar 表示应用程序状态栏
type StatusBar struct {
	container    *fyne.Container
	statusLabel  *widget.Label
	timeLabel    *widget.Label
	progressBar  *widget.ProgressBar
	showProgress bool
}

// NewStatusBar 创建一个新的状态栏
func NewStatusBar() *StatusBar {
	statusLabel := widget.NewLabel("就绪")
	timeLabel := widget.NewLabel("")
	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	// 创建状态栏容器
	statusBar := container.NewHBox(
		statusLabel,
		widget.NewSeparator(),
		timeLabel,
		container.NewHBox(
			widget.NewLabel("进度:"),
			progressBar,
		),
	)

	return &StatusBar{
		container:    statusBar,
		statusLabel:  statusLabel,
		timeLabel:    timeLabel,
		progressBar:  progressBar,
		showProgress: false,
	}
}

// Container 返回状态栏容器
func (s *StatusBar) Container() *fyne.Container {
	return s.container
}

// SetStatus 设置状态文本
func (s *StatusBar) SetStatus(text string) {
	s.statusLabel.SetText(text)
}

// SetStatusWithTimestamp 设置带时间戳的状态文本
func (s *StatusBar) SetStatusWithTimestamp(text string) {
	timestamp := time.Now().Format("15:04:05")
	s.statusLabel.SetText(text + " - " + timestamp)
}

// ShowProgress 显示进度条
func (s *StatusBar) ShowProgress() {
	s.progressBar.Show()
	s.showProgress = true
}

// HideProgress 隐藏进度条
func (s *StatusBar) HideProgress() {
	s.progressBar.Hide()
	s.showProgress = false
}

// SetProgress 设置进度值 (0.0 到 1.0)
func (s *StatusBar) SetProgress(value float64) {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	s.progressBar.SetValue(value)
}

// UpdateTime 更新时间显示
func (s *StatusBar) UpdateTime() {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	s.timeLabel.SetText(currentTime)
}

// Refresh 刷新状态栏
func (s *StatusBar) Refresh() {
	s.UpdateTime()
	s.container.Refresh()
}
