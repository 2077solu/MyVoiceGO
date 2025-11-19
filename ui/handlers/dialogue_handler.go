package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"myvoicego/ui/models"
)

// DialogueHandler 处理对话数据
type DialogueHandler struct {
	settings *models.AppSettings
}

// NewDialogueHandler 创建新的对话处理器
func NewDialogueHandler(settings *models.AppSettings) *DialogueHandler {
	return &DialogueHandler{
		settings: settings,
	}
}

// LoadDialogueFromFile 从文件加载对话数据
func (h *DialogueHandler) LoadDialogueFromFile(filePath string) ([]models.DialogueItem, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	var dialogueItems []models.DialogueItem
	err = json.Unmarshal(data, &dialogueItems)
	if err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	return dialogueItems, nil
}

// SaveDialogueToFile 将对话数据保存到文件
func (h *DialogueHandler) SaveDialogueToFile(dialogueItems []models.DialogueItem, filePath string) error {
	data, err := json.MarshalIndent(dialogueItems, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化JSON失败: %w", err)
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// ProcessDialogue 处理对话数据
func (h *DialogueHandler) ProcessDialogue(dialogueItems []models.DialogueItem) ([]models.ProcessedDialogueItem, error) {
	var processedItems []models.ProcessedDialogueItem

	for _, item := range dialogueItems {
		processedItem := models.ProcessedDialogueItem{
			ID:        item.ID,
			Timestamp: item.Timestamp,
			Speaker:   item.Speaker,
			Content:   item.Content,
		}

		// 情绪分析
		if h.settings.EnableEmotionAnalysis {
			emotion := h.analyzeEmotion(item.Content)
			processedItem.Emotion = emotion
		}

		// 关键词提取
		if h.settings.EnableKeywordExtraction {
			keywords := h.extractKeywords(item.Content)
			processedItem.Keywords = keywords
		}

		// 对话摘要
		if h.settings.EnableDialogueSummary {
			summary := h.generateSummary(item.Content)
			processedItem.Summary = summary
		}

		processedItems = append(processedItems, processedItem)
	}

	return processedItems, nil
}

// analyzeEmotion 分析情绪
func (h *DialogueHandler) analyzeEmotion(content string) string {
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
func (h *DialogueHandler) extractKeywords(content string) []string {
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

// generateSummary 生成摘要
func (h *DialogueHandler) generateSummary(content string) string {
	// 这里实现摘要生成逻辑
	// 可以调用API或使用本地模型

	// 简单实现：截取前50个字符作为摘要
	if len(content) > 50 {
		return content[:50] + "..."
	}
	return content
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

// ExportProcessedDialogue 导出处理后的对话数据
func (h *DialogueHandler) ExportProcessedDialogue(processedItems []models.ProcessedDialogueItem, outputPath string, format string) error {
	switch format {
	case "json":
		return h.exportAsJSON(processedItems, outputPath)
	case "csv":
		return h.exportAsCSV(processedItems, outputPath)
	case "txt":
		return h.exportAsText(processedItems, outputPath)
	default:
		return fmt.Errorf("不支持的导出格式: %s", format)
	}
}

// exportAsJSON 导出为JSON格式
func (h *DialogueHandler) exportAsJSON(processedItems []models.ProcessedDialogueItem, outputPath string) error {
	data, err := json.MarshalIndent(processedItems, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化JSON失败: %w", err)
	}

	err = ioutil.WriteFile(outputPath, data, 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// exportAsCSV 导出为CSV格式
func (h *DialogueHandler) exportAsCSV(processedItems []models.ProcessedDialogueItem, outputPath string) error {
	// 创建CSV内容
	var csvContent strings.Builder

	// 写入标题行
	csvContent.WriteString("ID,Timestamp,Speaker,Content,Emotion,Keywords,Summary\n")

	// 写入数据行
	for _, item := range processedItems {
		csvContent.WriteString(fmt.Sprintf("%d,%s,%s,\"%s\",%s,\"%s\",\"%s\"\n",
			item.ID,
			item.Timestamp,
			item.Speaker,
			strings.ReplaceAll(item.Content, "\"", "\"\""),
			item.Emotion,
			strings.Join(item.Keywords, ";"),
			strings.ReplaceAll(item.Summary, "\"", "\"\""),
		))
	}

	err := ioutil.WriteFile(outputPath, []byte(csvContent.String()), 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// exportAsText 导出为文本格式
func (h *DialogueHandler) exportAsText(processedItems []models.ProcessedDialogueItem, outputPath string) error {
	var textContent strings.Builder

	for _, item := range processedItems {
		textContent.WriteString(fmt.Sprintf("对话 %d:\n", item.ID))
		textContent.WriteString(fmt.Sprintf("时间: %s\n", item.Timestamp))
		textContent.WriteString(fmt.Sprintf("说话人: %s\n", item.Speaker))
		textContent.WriteString(fmt.Sprintf("内容: %s\n", item.Content))

		if h.settings.EnableEmotionAnalysis {
			textContent.WriteString(fmt.Sprintf("情绪: %s\n", item.Emotion))
		}

		if h.settings.EnableKeywordExtraction {
			textContent.WriteString(fmt.Sprintf("关键词: %s\n", strings.Join(item.Keywords, ", ")))
		}

		if h.settings.EnableDialogueSummary {
			textContent.WriteString(fmt.Sprintf("摘要: %s\n", item.Summary))
		}

		textContent.WriteString("\n" + strings.Repeat("-", 50) + "\n\n")
	}

	err := ioutil.WriteFile(outputPath, []byte(textContent.String()), 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}
