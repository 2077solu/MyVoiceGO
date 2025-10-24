package main

import (
	"fmt"
	"os"
	"path/filepath"

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
	err = parser.ReadAndPrintTestData(testFilePath)
	if err != nil {
		fmt.Printf("读取测试数据失败: %v\n", err)
		return
	}
}
