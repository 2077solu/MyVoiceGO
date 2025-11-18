package ui

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
)

// DialogueItem 表示JSON中的单个对话项
type DialogueItem struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	Text       string `json:"text"`
	Step       int    `json:"step"`
	Motion     string `json:"motion"`
	Expression string `json:"expression"`
	Model      string `json:"model"`
	Tone       string `json:"tone"`
}

// JSONViewer 创建并显示JSON查看器GUI
func JSONViewer() {
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	window := a.NewWindow("JSON查看器")
	window.Resize(fyne.NewSize(1000, 700))

	// 获取figures_output文件夹路径
	figuresPath := filepath.Join("..", "figures_output")

	// 创建文件列表
	fileList := widget.NewList(
		func() int {
			entries, _ := os.ReadDir(figuresPath)
			count := 0
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
					count++
				}
			}
			return count
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("文件名")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			entries, _ := os.ReadDir(figuresPath)
			jsonFiles := make([]string, 0)
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
					jsonFiles = append(jsonFiles, entry.Name())
				}
			}
			sort.Strings(jsonFiles)
			if id < len(jsonFiles) {
				obj.(*widget.Label).SetText(jsonFiles[id])
			}
		},
	)

	// 创建内容区域
	contentContainer := container.NewVBox()
	noFileSelected := widget.NewLabel("请从左侧选择一个JSON文件")
	contentContainer.Add(noFileSelected)

	// 当选择文件时加载其内容
	fileList.OnSelected = func(id widget.ListItemID) {
		entries, _ := os.ReadDir(figuresPath)
		jsonFiles := make([]string, 0)
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
				jsonFiles = append(jsonFiles, entry.Name())
			}
		}
		sort.Strings(jsonFiles)

		if id < len(jsonFiles) {
			filePath := filepath.Join(figuresPath, jsonFiles[id])
			loadJSONContent(filePath, contentContainer, window)
		}
	}

	// 创建保存按钮
	saveButton := widget.NewButton("保存修改", func() {
		dialog.ShowInformation("提示", "请在每个对话项的编辑面板中单独保存", window)
	})

	// 创建分割容器
	split := container.NewHSplit(fileList, contentContainer)
	split.Resize(fyne.NewSize(1000, 650))
	split.SetOffset(0.25) // 左侧占25%宽度

	// 创建主容器
	mainContainer := container.NewVBox(
		widget.NewCard("", "JSON文件查看器", split),
		container.NewHBox(
			saveButton,
		),
	)

	window.SetContent(mainContainer)
	window.ShowAndRun()
}

// loadJSONContent 加载并显示JSON文件内容
func loadJSONContent(filePath string, contentContainer *fyne.Container, window fyne.Window) {
	// 清空容器
	contentContainer.Objects = nil

	// 读取JSON文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		contentContainer.Add(widget.NewLabel(fmt.Sprintf("读取文件失败: %v", err)))
		return
	}

	// 解析JSON
	var dialogues []DialogueItem
	err = json.Unmarshal(data, &dialogues)
	if err != nil {
		contentContainer.Add(widget.NewLabel(fmt.Sprintf("解析JSON失败: %v", err)))
		return
	}

	// 创建文件名标签
	fileName := filepath.Base(filePath)
	fileLabel := widget.NewLabelWithStyle(fileName, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	contentContainer.Add(fileLabel)

	// 创建内容容器用于存放所有对话卡片
	cardContainer := container.NewVBox()

	// 为每个对话项创建卡片
	for i, dialogue := range dialogues {
		dialogueCard := createDialogueCard(i, dialogue, filePath, window)
		cardContainer.Add(dialogueCard)
	}

	// 创建滚动容器
	scrollContainer := container.NewScroll(cardContainer)

	// 添加滚动容器到主容器
	contentContainer.Add(scrollContainer)
	contentContainer.Refresh()
}

// createDialogueCard 为单个对话项创建可编辑的卡片
func createDialogueCard(index int, dialogue DialogueItem, filePath string, window fyne.Window) *widget.Card {
	// 创建标签
	nameLabel := widget.NewLabel("角色:")
	idLabel := widget.NewLabel("ID:")
	textLabel := widget.NewLabel("文本:")
	stepLabel := widget.NewLabel("步骤:")
	motionLabel := widget.NewLabel("动作:")
	expressionLabel := widget.NewLabel("表情:")
	modelLabel := widget.NewLabel("模型:")
	toneLabel := widget.NewLabel("语气:")

	// 创建输入框
	nameEntry := widget.NewEntry()
	nameEntry.SetText(dialogue.Name)

	idEntry := widget.NewEntry()
	idEntry.SetText(dialogue.ID)

	textEntry := widget.NewMultiLineEntry()
	textEntry.SetText(dialogue.Text)
	textEntry.SetMinRowsVisible(2)

	stepEntry := widget.NewEntry()
	stepEntry.SetText(fmt.Sprintf("%d", dialogue.Step))

	motionEntry := widget.NewEntry()
	motionEntry.SetText(dialogue.Motion)

	expressionEntry := widget.NewEntry()
	expressionEntry.SetText(dialogue.Expression)

	modelEntry := widget.NewEntry()
	modelEntry.SetText(dialogue.Model)

	toneEntry := widget.NewEntry()
	toneEntry.SetText(dialogue.Tone)

	// 创建保存按钮
	saveButton := widget.NewButtonWithIcon("保存", theme.DocumentSaveIcon(), func() {
		// 读取整个JSON文件
		data, err := os.ReadFile(filePath)
		if err != nil {
			dialog.ShowError(fmt.Errorf("读取文件失败: %v", err), window)
			return
		}

		// 解析JSON
		var dialogues []DialogueItem
		err = json.Unmarshal(data, &dialogues)
		if err != nil {
			dialog.ShowError(fmt.Errorf("解析JSON失败: %v", err), window)
			return
		}

		// 更新当前项
		if index < len(dialogues) {
			dialogues[index].Name = nameEntry.Text
			dialogues[index].ID = idEntry.Text
			dialogues[index].Text = textEntry.Text

			// 解析步骤
			var step int
			_, err := fmt.Sscanf(stepEntry.Text, "%d", &step)
			if err == nil {
				dialogues[index].Step = step
			}

			dialogues[index].Motion = motionEntry.Text
			dialogues[index].Expression = expressionEntry.Text
			dialogues[index].Model = modelEntry.Text
			dialogues[index].Tone = toneEntry.Text

			// 重新编码为JSON
			updatedData, err := json.MarshalIndent(dialogues, "", "    ")
			if err != nil {
				dialog.ShowError(fmt.Errorf("编码JSON失败: %v", err), window)
				return
			}

			// 写入文件
			err = os.WriteFile(filePath, updatedData, 0644)
			if err != nil {
				dialog.ShowError(fmt.Errorf("保存文件失败: %v", err), window)
				return
			}

			dialog.ShowInformation("成功", "对话项已成功保存", window)
		}
	})

	// 创建表单布局
	form := container.NewVBox(
		container.NewGridWithColumns(2,
			nameLabel, nameEntry,
			idLabel, idEntry,
		),
		textLabel,
		textEntry,
		container.NewGridWithColumns(2,
			stepLabel, stepEntry,
		),
		container.NewGridWithColumns(2,
			motionLabel, motionEntry,
			expressionLabel, expressionEntry,
		),
		modelLabel,
		modelEntry,
		container.NewGridWithColumns(2,
			toneLabel, toneEntry,
		),
		container.NewHBox(
			saveButton,
		),
	)

	// 创建卡片并返回
	card := widget.NewCard(
		fmt.Sprintf("对话项 #%d", index+1),
		"",
		form,
	)

	return card
}
