package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"myvoicego/parser"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取工作目录失败: %v\n", err)
		return
	}

	testFilePath := filepath.Join(wd, "test.txt")
	fmt.Printf("正在读取文件: %s\n", testFilePath)

	// 创建解析器实例
	p := parser.NewDialogueParser()

	// 读取文件内容
	content, err := os.ReadFile(testFilePath)
	if err != nil {
		fmt.Printf("读取文件失败: %v\n", err)
		return
	}

	// 按行分割内容
	lines := strings.Split(string(content), "\n")

	// 逐行解析
	for step, line := range lines {
		p.ParseDialogue(line, step)
	}

	// 保存解析结果到JSON
	//if err := p.SaveToJSON(filepath.Join(wd, "output")); err != nil {
	//    fmt.Printf("保存JSON失败: %v\n", err)
	//}
	//p.PrintFigures()
	
	// 导出figures到JSON文件
	outputDir := filepath.Join(wd, "figures_output")
	if err := p.ExportFiguresToJSON(outputDir); err != nil {
		fmt.Printf("导出figures到JSON失败: %v\n", err)
		return
	}
	fmt.Printf("成功导出所有figures到目录: %s\n", outputDir)
}
