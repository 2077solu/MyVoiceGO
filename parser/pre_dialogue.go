package parser

import (
	"fmt"
	"myvoicego/model"
	"strings"
)


// LineType 定义行类型
type LineType int

const (
	UnknownLine LineType = iota
	FigureChangeLine
	DialogueLine
)

type DialogueParser struct {
	figures    []model.PreDialogue
	tempFigure map[string]model.PreDialogue
}

// NewDialogueParser 创建新的对话解析器
func NewDialogueParser() *DialogueParser {
	return &DialogueParser{
		figures:    make([]model.PreDialogue, 0),
		tempFigure: make(map[string]model.PreDialogue),
	}
}

// ParseDialogue 解析对话
func (p *DialogueParser) ParseDialogue(line string, step int) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	lineType := p.detectLineType(line)
	switch lineType {
	case FigureChangeLine:
		p.parseFigureChange(line, step)
	case DialogueLine:
		p.parseDialogueLine(line, step)
	}
}

// detectLineType 检测行类型
func (p *DialogueParser) detectLineType(line string) LineType {
	if strings.HasPrefix(line, "changeFigure:") {
		return FigureChangeLine
	}
	if strings.Contains(line, "-figureId=") {
		return DialogueLine
	}
	return UnknownLine
}

// parseFigureChange 解析角色变更
func (p *DialogueParser) parseFigureChange(line string, step int) {
	if !p.isValidFigureChangeLine(line) {
		return
	}

	figure := model.PreDialogue{Step: step}
	parts := p.lineToParts(line)

	for _, part := range parts {
		p.parseFigurePart(&figure, part)
	}

	if figure.Id != "" {
		p.tempFigure[figure.Id] = figure
	}
}

// isValidFigureChangeLine 检查是否为有效的角色变更行
func (p *DialogueParser) isValidFigureChangeLine(line string) bool {
	return strings.Contains(line, "-motion=") || strings.Contains(line, "-expression=")
}

// parseFigurePart 解析角色部分
func (p *DialogueParser) parseFigurePart(figure *model.PreDialogue, part string) {
	switch {
	case strings.HasPrefix(part, "changeFigure:"):
		figure.Model = strings.TrimPrefix(part, "changeFigure:")
	case strings.HasPrefix(part, "-id="):
		figure.Id = strings.TrimPrefix(part, "-id=")
	case strings.HasPrefix(part, "-motion="):
		figure.Motion = strings.TrimPrefix(part, "-motion=")
	case strings.HasPrefix(part, "-expression="):
		figure.Expression = strings.TrimPrefix(part, "-expression=")
	}
}
// parseDialogue 解析对话内容
func (p *DialogueParser) parseDialogueLine(line string, step int) {
	dialogueText, figureId := p.extractDialogueInfo(line)
	if dialogueText == "" || figureId == "" {
		return
	}

	name, text := p.extractNameAndText(dialogueText)

	if figure, exists := p.tempFigure[figureId]; exists {
		updatedFigure := p.updateFigure(figure, text, name, step)
		p.addOrUpdateFigure(updatedFigure, figureId)
	}
}

// extractDialogueInfo 提取对话信息
func (p *DialogueParser) extractDialogueInfo(line string) (string, string) {
	parts := p.lineToParts(line)
	if len(parts) <= 1 {
		return "", ""
	}

	dialogueText := parts[0]
	figureId := p.extractFigureId(parts)

	return dialogueText, figureId
}

// extractFigureId 提取角色ID
func (p *DialogueParser) extractFigureId(parts []string) string {
	for _, part := range parts {
		if strings.HasPrefix(part, "-figureId=") {
			return strings.TrimPrefix(part, "-figureId=")
		}
	}
	return ""
}

// extractNameAndText 提取名称和文本
func (p *DialogueParser) extractNameAndText(dialogueText string) (string, string) {
	if !strings.Contains(dialogueText, ":") {
		return "", dialogueText
	}

	parts := strings.Split(dialogueText, ":")
	if len(parts) < 2 {
		return "", dialogueText
	}

	return parts[0], parts[1]
}

// updateFigure 更新角色信息
func (p *DialogueParser) updateFigure(figure model.PreDialogue, text, name string, step int) model.PreDialogue {
	figure.Text = text
	figure.Step = step
	if name != "" {
		figure.Name = name
	}
	return figure
}

// addOrUpdateFigure 添加或更新角色
func (p *DialogueParser) addOrUpdateFigure(figure model.PreDialogue, figureId string) {
	// 更新临时存储
	if figure.Name != "" {
		p.tempFigure[figureId] = figure
	}

	// 添加到结果列表
	p.figures = append(p.figures, figure)

	// 更新列表中已有的相同ID的元素
	for i := len(p.figures) - 1; i >= 0; i-- {
		if p.figures[i].Id == figureId {
			p.figures[i] = figure
			break
		}
	}
}

// lineToParts 行裁剪器
func (p *DialogueParser) lineToParts(line string) []string {
	line = strings.TrimSuffix(line, ";")
	content := strings.TrimSpace(line)
	return strings.Fields(content)
}

// PrintFigures 打印所有figure信息
func (p *DialogueParser) PrintFigures() {
	fmt.Println("\n=== 打印所有Figure信息 ===")
	for _, figure := range p.figures {
		fmt.Printf("Step: %d\n", figure.Step)
		fmt.Printf("ID: %s\n", figure.Id)
		fmt.Printf("Name: %s\n", figure.Name)
		fmt.Printf("Text: %s\n", figure.Text)
		fmt.Printf("Motion: %s\n", figure.Motion)
		fmt.Printf("Expression: %s\n", figure.Expression)
		fmt.Println("------------------------")
	}
	fmt.Println("=== Figure信息打印完成 ===\n")
}
