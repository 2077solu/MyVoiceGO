package parser

import (
	"fmt"
	"os"
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
