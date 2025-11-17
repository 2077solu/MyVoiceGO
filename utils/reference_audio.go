package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	SubModel0      = "0"
	SubModel1      = "1"
	DirPermission  = 0755
	FilePermission = 0644
)

// 支持的音频文件扩展名映射，用于快速查找
var supportedAudioExts = map[string]struct{}{
	".mp3":  {},
	".wav":  {},
	".ogg":  {},
	".flac": {},
	".m4a":  {},
	".aac":  {},
}

// AudioRefPath 音频文件路径结构
type AudioRefPath struct {
	Path string `json:"path"`
}

// AudioRefTone 语气对应的音频文件列表
type AudioRefTone struct {
	Tone  string         `json:"tone"`
	Paths []AudioRefPath `json:"paths"`
}

// AudioRefModel 模型对应的语气列表
type AudioRefModel struct {
	ModelID string         `json:"model_id"`
	Model   string         `json:"model"`
	Tones   []AudioRefTone `json:"tones"`
}

// AudioRefList 所有模型的音频参考列表
type AudioRefList struct {
	Models []AudioRefModel `json:"models"`
}

// normalizePath 统一路径格式为Unix风格
func normalizePath(path string) string {
	return filepath.ToSlash(path)
}

// isAudioFile 检查文件是否为支持的音频格式
func isAudioFile(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	_, supported := supportedAudioExts[ext]
	return supported
}

// readDirectory 安全读取目录内容并返回详细错误信息
func readDirectory(path string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败 %q: %w", path, err)
	}
	return entries, nil
}

// getAudioFilesInDir 获取指定目录下的所有音频文件路径
func getAudioFilesInDir(dirPath string) ([]AudioRefPath, error) {
	entries, err := readDirectory(dirPath)
	if err != nil {
		return nil, err
	}

	audioPaths := make([]AudioRefPath, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue // 跳过子目录
		}

		fileName := entry.Name()
		if !isAudioFile(fileName) {
			continue // 跳过非音频文件
		}

		audioPaths = append(audioPaths, AudioRefPath{
			Path: normalizePath(fileName),
		})
	}

	return audioPaths, nil
}

// hasSubModels 检查是否存在子模型目录(0或1)
func hasSubModels(entries []os.DirEntry) bool {
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			if name == SubModel0 || name == SubModel1 {
				return true
			}
		}
	}
	return false
}

// collectTones 从基础目录收集语气信息(子目录作为语气名)
func collectTones(baseDir string) ([]AudioRefTone, error) {
	entries, err := readDirectory(baseDir)
	if err != nil {
		return nil, err
	}

	var tones []AudioRefTone
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		toneName := entry.Name()
		tonePath := filepath.Join(baseDir, toneName)

		audioPaths, err := getAudioFilesInDir(tonePath)
		if err != nil || len(audioPaths) == 0 {
			continue // 忽略无音频或读取失败的语气目录
		}

		tones = append(tones, AudioRefTone{
			Tone:  toneName,
			Paths: audioPaths,
		})
	}

	return tones, nil
}

// processSubModels 处理包含子模型(0/1)的目录结构
func processSubModels(modelPath string) ([]AudioRefTone, error) {
	entries, err := readDirectory(modelPath)
	if err != nil {
		return nil, err
	}

	var allTones []AudioRefTone
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		subModelName := entry.Name()
		if subModelName != SubModel0 && subModelName != SubModel1 {
			continue // 只处理指定的子模型目录
		}

		subModelPath := filepath.Join(modelPath, subModelName)
		tones, err := collectTones(subModelPath)
		if err != nil {
			continue // 忽略单个子模型的错误
		}
		allTones = append(allTones, tones...)
	}

	return allTones, nil
}

// getTonesInModel 获取指定模型目录下的所有语气信息
func getTonesInModel(rootDir, modelName string) (string, []AudioRefTone, error) {
	modelPath := filepath.Join(rootDir, modelName)
	entries, err := readDirectory(modelPath)
	if err != nil {
		return "", nil, err
	}

	// 无论哪种结构，默认使用SubModel0作为模型ID
	modelID := SubModel0
	var tones []AudioRefTone

	if hasSubModels(entries) {
		tones, err = processSubModels(modelPath)
	} else {
		tones, err = collectTones(modelPath)
	}

	if err != nil {
		return "", nil, err
	}

	return modelID, tones, nil
}

// BuildReferenceAudioList 构建所有模型的音频参考列表
func BuildReferenceAudioList(rootDir string) ([]AudioRefModel, error) {
	entries, err := readDirectory(rootDir)
	if err != nil {
		return nil, err
	}

	models := make([]AudioRefModel, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue // 只处理目录作为模型
		}

		modelName := entry.Name()
		modelID, tones, err := getTonesInModel(rootDir, modelName)
		if err != nil {
			continue // 忽略处理失败的模型
		}

		if len(tones) > 0 {
			models = append(models, AudioRefModel{
				ModelID: modelID,
				Model:   modelName,
				Tones:   tones,
			})
		}
	}

	return models, nil
}

// ListReferenceAudioFiles 将音频参考列表序列化为JSON
func ListReferenceAudioFiles(rootDir string) (string, error) {
	models, err := BuildReferenceAudioList(rootDir)
	if err != nil {
		return "", err
	}

	audioList := AudioRefList{
		Models: models,
	}

	jsonData, err := json.MarshalIndent(audioList, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON序列化失败: %w", err)
	}

	return string(jsonData), nil
}

// SaveReferenceAudioListToFile 将音频参考列表保存到指定文件
func SaveReferenceAudioListToFile(rootDir, outputPath string) error {
	jsonData, err := ListReferenceAudioFiles(rootDir)
	if err != nil {
		return err
	}

	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(outputPath), DirPermission); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(outputPath, []byte(jsonData), FilePermission); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}
