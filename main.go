package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"myvoicego/api"
	"myvoicego/model"
	"myvoicego/parser"
	"myvoicego/utils"
)

func main() {
	// 获取工作目录
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取工作目录失败: %v\n", err)
		return
	}

	// 确保基础目录存在
	referenceAudioPath := filepath.Join(wd, "reference_audio")
	if _, err := os.Stat(referenceAudioPath); os.IsNotExist(err) {
		fmt.Printf("reference_audio 目录不存在: %v\n", err)
		return
	}

	// 调用方法获取JSON格式的音频列表
	jsonData, err := utils.ListReferenceAudioFiles(referenceAudioPath)
	if err != nil {
		fmt.Printf("获取音频列表失败: %v\n", err)
		return
	}

	// 打印结果
	fmt.Println(jsonData)

	// 在main函数中添加
	if err := utils.StartGPTSvits(); err != nil {
		fmt.Printf("启动GPT-SoVITS失败: %v\n", err)
		return
	}

	// 确保输出目录存在
	storagePath := filepath.Join(wd, "reference_audio_storage")
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		fmt.Printf("创建输出目录失败: %v\n", err)
		return
	}

	// 保存到文件
	outputPath := filepath.Join(storagePath, "爱音.json")
	err = utils.SaveReferenceAudioListToFile(referenceAudioPath, outputPath)
	if err != nil {
		fmt.Printf("保存文件时出错: %v\n", err)
		return
	}

	// 工作目录已在上面获取
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

	// 导出figures到JSON文件
	outputDir := filepath.Join(wd, "figures_output")
	if err := p.ExportFiguresToJSON(outputDir); err != nil {
		fmt.Printf("导出figures到JSON失败: %v\n", err)
		return
	}
	fmt.Printf("成功导出所有figures到目录: %s\n", outputDir)

	// 读取导出的JSON文件并进行情绪分析
	figuresFiles, err := api.ListAvailableDialogues(outputDir)
	if err != nil {
		fmt.Printf("获取figures文件列表失败: %v\n", err)
		return
	}

	// 创建Coze API实例
	cozeAPI, err := api.NewCozeAPIFromConfig(filepath.Join(wd, "config", "coze_config.json"))
	if err != nil {
		fmt.Printf("创建Coze API实例失败: %v\n", err)
		return
	}

	// 对每个角色的对话进行情绪分析
	for _, figureFile := range figuresFiles {
		// 读取对话文件
		dialogueJSON, err := api.ReadDialogueFromFile(outputDir, figureFile)
		if err != nil {
			fmt.Printf("读取对话文件失败 (%s): %v\n", figureFile, err)
			continue
		}

		// 解析对话
		var dialogues []model.PreDialogue
		if err := json.Unmarshal([]byte(dialogueJSON), &dialogues); err != nil {
			fmt.Printf("解析对话JSON失败 (%s): %v\n", figureFile, err)
			continue
		}

		// 进行语气分析
		updatedDialogues, err := cozeAPI.AnalyzeTones(dialogues)
		if err != nil {
			fmt.Printf("语气分析失败 (%s): %v\n", figureFile, err)
			continue
		}

		// 将带有语气的对话重新序列化为JSON
		updatedJSON, err := json.MarshalIndent(updatedDialogues, "", "    ")
		if err != nil {
			fmt.Printf("序列化更新后的对话失败 (%s): %v\n", figureFile, err)
			continue
		}

		// 保存更新后的对话
		filePath := filepath.Join(outputDir, fmt.Sprintf("%s.json", figureFile))
		if err := os.WriteFile(filePath, updatedJSON, 0644); err != nil {
			fmt.Printf("保存更新后的对话失败 (%s): %v\n", figureFile, err)
			continue
		}

		fmt.Printf("成功更新角色 %s 的情绪分析结果\n", figureFile)
	}
}
