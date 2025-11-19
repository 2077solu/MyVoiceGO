package views

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"myvoicego/ui/components"
	"myvoicego/ui/models"
)

// ProcessorView 表示处理器视图
type ProcessorView struct {
	app           fyne.App
	window        fyne.Window
	statusBar     *components.StatusBar
	settings      *models.AppSettings
	filePath      string
	dialogueItems []models.DialogueItem
	resultText    *widget.Entry
	progressBar   *widget.ProgressBar
}

// NewProcessorView 创建新的处理器视图实例
func NewProcessorView(app fyne.App, window fyne.Window) *ProcessorView {
	settings := models.NewAppSettings()
	err := settings.LoadFromFile("settings.json")
	if err != nil {
		fmt.Printf("加载设置失败: %v\n", err)
		// 使用默认设置
	}

	view := &ProcessorView{
		app:         app,
		window:      window,
		settings:    settings,
		resultText:  widget.NewMultiLineEntry(),
		progressBar: widget.NewProgressBar(),
	}

	view.resultText.SetText("处理结果将显示在这里...")

	return view
}

// createUI 创建UI界面
func (v *ProcessorView) createUI() fyne.CanvasObject {
	// 创建状态栏
	v.statusBar = components.NewStatusBar()

	// 创建文件选择区域
	fileSelection := v.createFileSelection()

	// 创建处理选项区域
	processingOptions := v.createProcessingOptions()

	// 创建处理按钮区域
	processingButtons := v.createProcessingButtons()

	// 创建结果显示区域
	resultDisplay := v.createResultDisplay()

	// 创建主布局
	content := container.NewVBox(
		fileSelection,
		processingOptions,
		processingButtons,
		resultDisplay,
	)

	return container.NewBorder(
		nil,
		v.statusBar.Container(),
		nil,
		nil,
		content,
	)
}

// createFileSelection 创建文件选择区域
func (v *ProcessorView) createFileSelection() fyne.CanvasObject {
	fileLabel := widget.NewLabel("未选择文件")

	return widget.NewCard("文件选择", "", container.NewVBox(
		fileLabel,
		widget.NewButtonWithIcon("选择对话文件", theme.FileIcon(), func() {
			dialog := components.NewOpenFileDialog(func(reader fyne.URIReadCloser, err error) {
				if err != nil {
					dialog.ShowError(err, v.window)
					return
				}
				if reader != nil {
					defer reader.Close()
					v.filePath = reader.URI().Path()
					fileLabel.SetText("已选择: " + filepath.Base(v.filePath))
					v.statusBar.SetStatus("已选择文件: " + filepath.Base(v.filePath))

					// 自动加载文件
					v.loadDialogueFile()
				}
			}, v.window)
			dialog.ShowOpen(v.window)
		}),
	))
}

// createProcessingOptions 创建处理选项区域
func (v *ProcessorView) createProcessingOptions() fyne.CanvasObject {
	emotionCheck := widget.NewCheck("情绪分析", func(checked bool) {
		v.settings.EnableEmotionAnalysis = checked
		v.settings.SaveToFile("settings.json")
	})
	emotionCheck.SetChecked(v.settings.EnableEmotionAnalysis)

	keywordCheck := widget.NewCheck("关键词提取", func(checked bool) {
		v.settings.EnableKeywordExtraction = checked
		v.settings.SaveToFile("settings.json")
	})
	keywordCheck.SetChecked(v.settings.EnableKeywordExtraction)

	summaryCheck := widget.NewCheck("对话摘要", func(checked bool) {
		v.settings.EnableDialogueSummary = checked
		v.settings.SaveToFile("settings.json")
	})
	summaryCheck.SetChecked(v.settings.EnableDialogueSummary)

	return widget.NewCard("处理选项", "", container.NewVBox(
		emotionCheck,
		keywordCheck,
		summaryCheck,
	))
}

// createProcessingButtons 创建处理按钮区域
func (v *ProcessorView) createProcessingButtons() fyne.CanvasObject {
	return container.NewHBox(
		widget.NewButtonWithIcon("处理对话", theme.ConfirmIcon(), func() {
			v.processDialogue()
		}),
		widget.NewButtonWithIcon("保存结果", theme.DocumentSaveIcon(), func() {
			v.saveResults()
		}),
		widget.NewButtonWithIcon("清空结果", theme.DeleteIcon(), func() {
			v.clearResults()
		}),
	)
}

// createResultDisplay 创建结果显示区域
func (v *ProcessorView) createResultDisplay() fyne.CanvasObject {
	v.resultText.SetPlaceHolder("处理结果将显示在这里...")

	return widget.NewCard("处理结果", "", container.NewVBox(
		v.progressBar,
		v.resultText,
	))
}

// loadDialogueFile 加载对话文件
func (v *ProcessorView) loadDialogueFile() {
	if v.filePath == "" {
		dialog.ShowError(fmt.Errorf("未选择文件"), v.window)
		return
	}

	data, err := os.ReadFile(v.filePath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("读取文件失败: %w", err), v.window)
		return
	}

	// 尝试解析JSON
	var dialogueItems []models.DialogueItem
	err = json.Unmarshal(data, &dialogueItems)
	if err != nil {
		dialog.ShowError(fmt.Errorf("解析JSON失败: %w", err), v.window)
		return
	}

	v.dialogueItems = dialogueItems
	v.statusBar.SetStatusWithTimestamp(fmt.Sprintf("已加载 %d 条对话", len(dialogueItems)))
}

// processDialogue 处理对话
func (v *ProcessorView) processDialogue() {
	if len(v.dialogueItems) == 0 {
		dialog.ShowError(fmt.Errorf("没有可处理的对话数据"), v.window)
		return
	}

	v.progressBar.SetValue(0)
	v.statusBar.ShowProgress()
	v.statusBar.SetProgress(0.0)

	// 清空结果
	v.resultText.SetText("")

	// 处理每条对话
	total := len(v.dialogueItems)
	for i, item := range v.dialogueItems {
		// 更新进度
		progress := float64(i) / float64(total)
		v.progressBar.SetValue(progress)

		// 处理单条对话
		result := v.processSingleDialogue(item)

		// 添加到结果
		currentText := v.resultText.Text
		v.resultText.SetText(currentText + result + "\n\n")
	}

	// 完成处理
	v.progressBar.SetValue(1.0)
	v.statusBar.HideProgress()
	v.statusBar.SetStatusWithTimestamp("对话处理完成")
}

// processSingleDialogue 处理单条对话
func (v *ProcessorView) processSingleDialogue(item models.DialogueItem) string {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("对话 %s:\n", item.ID))
	result.WriteString(fmt.Sprintf("说话人: %s\n", item.Name))
	result.WriteString(fmt.Sprintf("内容: %s\n", item.Text))

	// 情绪分析
	if v.settings.EnableEmotionAnalysis {
		emotion := v.analyzeEmotion(item.Text)
		result.WriteString(fmt.Sprintf("情绪: %s\n", emotion))
	}

	// 关键词提取
	if v.settings.EnableKeywordExtraction {
		keywords := v.extractKeywords(item.Text)
		result.WriteString(fmt.Sprintf("关键词: %s\n", strings.Join(keywords, ", ")))
	}

	return result.String()
}

// analyzeEmotion 分析情绪
func (v *ProcessorView) analyzeEmotion(content string) string {
	// 这里实现情绪分析逻辑
	// 可以调用API或使用本地模型

	// 简单实现：基于关键词的情绪分析
	content = strings.ToLower(content)

	if strings.Contains(content, "开心") || strings.Contains(content, "高兴") || strings.Contains(content, "快乐") {
		return "开心"
	} else if strings.Contains(content, "难过") || strings.Contains(content, "伤心") || strings.Contains(content, "悲伤") {
		return "难过"
	} else if strings.Contains(content, "生气") || strings.Contains(content, "愤怒") || strings.Contains(content, "恼火") {
		return "生气"
	} else if strings.Contains(content, "害怕") || strings.Contains(content, "恐惧") || strings.Contains(content, "担心") {
		return "害怕"
	} else {
		return "中性"
	}
}

// extractKeywords 提取关键词
func (v *ProcessorView) extractKeywords(content string) []string {
	// 这里实现关键词提取逻辑
	// 可以调用API或使用本地模型

	// 简单实现：基于词频的关键词提取
	words := strings.Fields(content)
	wordCount := make(map[string]int)

	for _, word := range words {
		// 过滤掉常见的停用词
		if !isStopWord(word) {
			wordCount[word]++
		}
	}

	// 返回出现频率最高的几个词
	var keywords []string
	for word, count := range wordCount {
		if count >= 2 {
			keywords = append(keywords, word)
		}
	}

	return keywords
}

// isStopWord 检查是否是停用词
func isStopWord(word string) bool {
	stopWords := []string{
		"的", "了", "在", "是", "我", "有", "和", "就", "不", "人", "都", "一", "一个", "上", "也", "很", "到", "说", "要", "去", "你", "会", "着", "没有", "看", "好", "自己", "这",
		"the", "a", "an", "and", "or", "but", "in", "on", "at", "to", "for", "of", "with", "by", "is", "are", "was", "were", "be", "been", "being", "have", "has", "had", "do", "does", "did", "will", "would", "could", "should", "may", "might", "must", "can", "shall",
	}

	for _, stopWord := range stopWords {
		if word == stopWord {
			return true
		}
	}

	return false
}

// saveResults 保存结果
func (v *ProcessorView) saveResults() {
	if v.resultText.Text == "" {
		dialog.ShowError(fmt.Errorf("没有可保存的结果"), v.window)
		return
	}

	dialog := components.NewSaveFileDialog(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, v.window)
			return
		}
		if writer != nil {
			defer writer.Close()

			_, err = writer.Write([]byte(v.resultText.Text))
			if err != nil {
				dialog.ShowError(fmt.Errorf("保存文件失败: %w", err), v.window)
				return
			}

			v.statusBar.SetStatusWithTimestamp("结果已保存")
		}
	}, v.window)
	dialog.ShowSave(v.window)
}

// clearResults 清空结果
func (v *ProcessorView) clearResults() {
	confirmDialog := components.NewConfirmDialog(
		"确认清空",
		"确定要清空所有结果吗？",
		func(confirm bool) {
			if confirm {
				v.resultText.SetText("")
				v.progressBar.SetValue(0)
				v.statusBar.SetStatus("结果已清空")
			}
		},
		v.window)
	confirmDialog.Show()
}

// Show 显示处理器视图
func (v *ProcessorView) Show() {
	content := v.createUI()
	v.window.SetContent(content)
}
