package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ExpandableNode 表示可展开/折叠的节点
type ExpandableNode struct {
	Title      string
	Content    fyne.CanvasObject
	IsExpanded bool
	OnClick    func(*ExpandableNode)
}

// NewExpandableNode 创建一个新的可展开节点
func NewExpandableNode(title string, content fyne.CanvasObject) *ExpandableNode {
	return &ExpandableNode{
		Title:      title,
		Content:    content,
		IsExpanded: false,
	}
}

// BuildUI 构建节点的UI组件
func (n *ExpandableNode) BuildUI() fyne.CanvasObject {
	// 创建展开/折叠图标
	expandIcon := widget.NewIcon(theme.MenuDropDownIcon())
	if n.IsExpanded {
		expandIcon.SetResource(theme.MenuDropUpIcon())
	}

	// 创建标题标签
	titleLabel := widget.NewLabel(n.Title)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Wrapping = fyne.TextWrapBreak

	// 创建可点击的标题区域
	titleButton := widget.NewButton("", func() {
		n.IsExpanded = !n.IsExpanded
		if n.IsExpanded {
			expandIcon.SetResource(theme.MenuDropUpIcon())
		} else {
			expandIcon.SetResource(theme.MenuDropDownIcon())
		}
		if n.OnClick != nil {
			n.OnClick(n)
		}
	})
	titleButton.Importance = widget.LowImportance

	// 创建标题行
	titleRow := container.NewHBox(expandIcon, titleLabel)

	// 创建内容容器（初始时隐藏）
	contentContainer := container.NewVBox()
	if n.IsExpanded && n.Content != nil {
		contentContainer.Add(n.Content)
	}

	// 创建节点容器
	nodeContainer := container.NewVBox(
		titleRow,
		container.NewHBox(
			widget.NewLabel("  "), // 缩进
			contentContainer,
		),
	)

	// 设置点击回调，用于显示/隐藏内容
	n.OnClick = func(node *ExpandableNode) {
		contentContainer.Objects = nil
		if node.IsExpanded && node.Content != nil {
			contentContainer.Add(node.Content)
		}
		contentContainer.Refresh()
	}

	return nodeContainer
}

// Toggle 切换节点的展开/折叠状态
func (n *ExpandableNode) Toggle() {
	n.IsExpanded = !n.IsExpanded
}

// Expand 展开节点
func (n *ExpandableNode) Expand() {
	n.IsExpanded = true
}

// Collapse 折叠节点
func (n *ExpandableNode) Collapse() {
	n.IsExpanded = false
}
