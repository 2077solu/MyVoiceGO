package main

import (
	"fmt"
	"os"
	"path/filepath"

	"myvoicego/parser"
)

func main() {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取工作目录失败: %v\n", err)
		return
	}

	// 构建测试文件路径
	testFilePath := filepath.Join(wd, "test.txt")

	// 创建默认文件读取器
	fileReader := &parser.DefaultFileReader{}

	// 创建PreDialogue实例，注入依赖
	preDialogue := parser.NewPreDialogue(fileReader)

	// 读取并打印测试数据
	fmt.Printf("正在读取文件: %s\n", testFilePath)
	err = preDialogue.ReadAndPrintTestData(testFilePath)
	if err != nil {
		fmt.Printf("读取测试数据失败: %v\n", err)
		return
	}
}
