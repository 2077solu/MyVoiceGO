package models

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ModelConfig 模型配置结构体
type ModelConfig struct {
	Name     string `json:"name"`
	Alias    string `json:"alias"`
	ModelIDs []struct {
		ModelPath string `json:"model_path"`
	} `json:"models"`
}

// DialogueItem 表示JSON中的单个对话项
type DialogueItem struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	Text       string `json:"text"`
	Step       int    `json:"step"`
	Motion     string `json:"motion"`
	Expression string `json:"expression"`
	Model      string `json:"model"`
	Tone       string `json:"tone"`
}

// ProcessedDialogueItem 表示处理后的对话项
type ProcessedDialogueItem struct {
	ID        int      `json:"id"`
	Timestamp string   `json:"timestamp"`
	Speaker   string   `json:"speaker"`
	Content   string   `json:"content"`
	Emotion   string   `json:"emotion,omitempty"`
	Keywords  []string `json:"keywords,omitempty"`
	Summary   string   `json:"summary,omitempty"`
}

// AppSettings 应用程序设置
type AppSettings struct {
	Theme                   string `json:"theme"`
	Language                string `json:"language"`
	LastOpenFile            string `json:"last_open_file"`
	WindowWidth             int    `json:"window_width"`
	WindowHeight            int    `json:"window_height"`
	EnableEmotionAnalysis   bool   `json:"enable_emotion_analysis"`
	EnableKeywordExtraction bool   `json:"enable_keyword_extraction"`
	EnableDialogueSummary   bool   `json:"enable_dialogue_summary"`
}

// NewAppSettings 创建默认应用设置
func NewAppSettings() *AppSettings {
	return &AppSettings{
		Theme:                   "dark",
		Language:                "zh-CN",
		LastOpenFile:            "",
		WindowWidth:             1200,
		WindowHeight:            800,
		EnableEmotionAnalysis:   true,
		EnableKeywordExtraction: true,
		EnableDialogueSummary:   false,
	}
}

// SaveToFile 保存设置到文件
func (s *AppSettings) SaveToFile(filePath string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return os.WriteFile(filePath, data, 0644)
}

// LoadFromFile 从文件加载设置
func (s *AppSettings) LoadFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, s)
}

// LoadModelConfigs 加载模型配置
func LoadModelConfigs(configPath string) ([]ModelConfig, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		// 如果读取失败，返回默认模型
		return []ModelConfig{
			{
				Name:  "默认模型",
				Alias: "default",
				ModelIDs: []struct {
					ModelPath string `json:"model_path"`
				}{
					{ModelPath: "models/default"},
				},
			},
		}, nil
	}

	var configs []ModelConfig
	err = json.Unmarshal(file, &configs)
	if err != nil {
		// 如果解析失败，返回默认模型
		return []ModelConfig{
			{
				Name:  "默认模型",
				Alias: "default",
				ModelIDs: []struct {
					ModelPath string `json:"model_path"`
				}{
					{ModelPath: "models/default"},
				},
			},
		}, nil
	}

	return configs, nil
}

// SaveModelConfigs 保存模型配置
func SaveModelConfigs(configs []ModelConfig, configPath string) error {
	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return err
	}

	// 确保目录存在
	dir := filepath.Dir(configPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return os.WriteFile(configPath, data, 0644)
}
