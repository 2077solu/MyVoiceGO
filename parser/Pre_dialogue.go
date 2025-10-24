package parser

import (
	"fmt"
	"os"
)

// FileReader 文件读取接口
type FileReader interface {
	ReadFile(path string) (string, error)
}

// DefaultFileReader 默认文件读取实现
type DefaultFileReader struct{}

// ReadFile 读取文件内容
func (d *DefaultFileReader) ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}
	return string(content), nil
}

// TestReader 测试数据读取器
type TestReader struct {
	fileReader FileReader
}

// NewTestReader 创建新的测试读取器
func NewTestReader(fileReader FileReader) *TestReader {
	if fileReader == nil {
		fileReader = &DefaultFileReader{}
	}
	return &TestReader{fileReader: fileReader}
}

// ReadAndPrintTestData 读取并打印测试数据
func (t *TestReader) ReadAndPrintTestData(path string) error {
	content, err := t.fileReader.ReadFile(path)
	if err != nil {
		return err
	}
	fmt.Println(content)
	return nil
}
