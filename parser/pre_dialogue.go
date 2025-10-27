package parser

import (
	"fmt"
	"myvoicego/model"
	"os"
	"strings"
)

// ReadAndPrintTestData 读取并打印测试数据
func ReadAndPrintTestData(path string) error {
	startTxt, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}
	fmt.Println(string(startTxt))
	return nil
}

type DialogueParser struct {
	figures    []model.PreDialogue
	tempFigure map[string]model.PreDialogue
}

// ParseDialogue 初始化
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

	switch {
	case strings.HasPrefix(line, "changeFigure:"):
		p.parserFigureChange(line, step)
	case strings.Contains(line, "-figureId="):
		p.parserDialogue(line, step)
	}
}

func (p *DialogueParser) parserFigureChange(line string, step int) {
	line = strings.TrimSuffix(line, ";")
	content := strings.TrimSpace(line)
	if strings.Contains(content, "-motion=") || strings.Contains(content, "-expression=") {
		parts := strings.Fields(content)

		figure := model.PreDialogue{
			Step: step,
		}
		for _, part := range parts {
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
		if figure.Id != "" {
			p.tempFigure[figure.Id] = figure
			p.figures = append(p.figures, figure)
		}
	}else {
		return
	}
}
func (p *DialogueParser) parserDialogue(line string, step int) {
	line = strings.TrimSuffix(line, ";")
	content := strings.TrimSpace(line)
    parts := strings.Fields(content)

	var dialogueText string
	var Id string
	if len(parts) > 1 {
		dialogueText = parts[0]
		for _, part := range parts {
			if strings.HasPrefix(part, "-figureId=") {
				Id = strings.TrimPrefix(part, "-figureId=")
				break
			}
		}
	}else {
		return
	}
	var name string
	var text string
	if strings.Contains(dialogueText, ":") {
		ingredients := strings.Split(dialogueText, ":")
		name = ingredients[0]
		text = ingredients[1]
	}else {
		text = dialogueText
	}
	//TODO name和text未正确导入
	if Id != "" {
		if figure, exists := p.tempFigure[Id]; exists {
			figure.Name = name
			figure.Text = text
			p.tempFigure[Id] = figure
			for i := len(p.figures) - 1; i >= 0; i-- {
				if p.figures[i].Id == Id {
					p.figures[i] = figure
					break
				}
			}
		}
	}
}
//TODO figure貌似只能按轮次存储数据
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

