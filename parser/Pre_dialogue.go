package parser

import (
	"bufio"
	"fmt"
	"os"
)

// FileReader 创建文件读取依赖
type FileReader interface {
	ReadFile(path string) (string, error)
}

type DefaultFileReader struct{}

// ReadFile DefaultFileReader 结构体的 ReadFile 方法，用于读取文件内容
// 参数：
//   path - 要读取的文件路径
// 返回值：
//   string - 文件内容
//   error - 错误信息
func (d *DefaultFileReader) ReadFile(path string) (string, error) {
	// 打开指定路径的文件
	file, err := os.Open(path)
	// 如果打开文件出错，返回错误信息
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 用于存储文件内容的字符串变量
	var content string
	// 创建一个扫描器，用于逐行读取文件内容
	scanner := bufio.NewScanner(file)
	// 循环扫描每一行
	for scanner.Scan() {
		// 将读取的行内容添加到 content 中，并加上换行符
		content += scanner.Text() + "\n"
	}

	// 检查扫描过程中是否有错误发生
	if err := scanner.Err(); err != nil {
		return "", err
	}

	// 返回文件内容和 nil 错误
	return content, nil
}

// PreDialogue 结构体包含文件读取依赖
type PreDialogue struct {
	fileReader FileReader
}

// NewPreDialogue 创建一个新的PreDialogue实例，并注入依赖
func NewPreDialogue(fileReader FileReader) *PreDialogue {
	return &PreDialogue{
		fileReader: fileReader,
	}
}

// ReadTestData 读取测试数据并存到starttxt变量
func (p *PreDialogue) ReadTestData(path string) (string, error) {
	starttxt, err := p.fileReader.ReadFile(path)
	if err != nil {
		return "", err
	}
	return starttxt, nil
}

// ReadAndPrintTestData 读取测试数据并直接打印到控制台
func (p *PreDialogue) ReadAndPrintTestData(path string) error {
	starttxt, err := p.ReadTestData(path)
	if err != nil {
		return err
	}
	fmt.Println(starttxt)
	return nil
}
