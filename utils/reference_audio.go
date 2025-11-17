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

// 支持的音频文件扩展名
var SupportedAudioExtensions = []string{".mp3", ".wav", ".ogg", ".flac", ".m4a", ".aac"}

type ReferenceAudioPath struct {
	Path string `json:"path"`
}

type ReferenceAudioTone struct {
	Tone  string               `json:"tone"`
	Paths []ReferenceAudioPath `json:"paths"`
}

type ReferenceAudioModel struct {
	ModelID string               `json:"model_id"`
	Model   string               `json:"model"`
	Tones   []ReferenceAudioTone `json:"tones"`
}

type ReferenceAudioList struct {
	Models []ReferenceAudioModel `json:"models"`
}

// normalizePath 统一路径格式
func normalizePath(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

// isAudioFile 检查文件是否为支持的音频格式
func isAudioFile(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	for _, supportedExt := range SupportedAudioExtensions {
		if ext == supportedExt {
			return true
		}
	}
	return false
}

// readDir 安全地读取目录内容
func readDir(path string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败 %s: %w", path, err)
	}
	return entries, nil
}

// getAudioFilesInDir 获取指定目录下的所有音频文件
func getAudioFilesInDir(dirPath string) ([]ReferenceAudioPath, error) {
	entries, err := readDir(dirPath)
	if err != nil {
		return nil, err
	}

	audioPaths := make([]ReferenceAudioPath, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if !isAudioFile(fileName) {
			continue
		}

		filePath := filepath.Join(dirPath, fileName)
		relPath, err := filepath.Rel(dirPath, filePath)
		if err != nil {
			relPath = filePath
		}

		audioPaths = append(audioPaths, ReferenceAudioPath{
			Path: normalizePath(relPath),
		})
	}

	return audioPaths, nil
}

// hasSubModels 检查是否存在子模型结构
func hasSubModels(entries []os.DirEntry) bool {
	for _, entry := range entries {
		if entry.IsDir() && (entry.Name() == SubModel0 || entry.Name() == SubModel1) {
			return true
		}
	}
	return false
}

// processSubModels 处理子模型结构
func processSubModels(modelPath string) ([]ReferenceAudioTone, error) {
	entries, err := readDir(modelPath)
	if err != nil {
		return nil, err
	}

	var tones []ReferenceAudioTone
	for _, entry := range entries {
		if !entry.IsDir() || (entry.Name() != SubModel0 && entry.Name() != SubModel1) {
			continue
		}

		subModelPath := filepath.Join(modelPath, entry.Name())
		subModelEntries, err := readDir(subModelPath)
		if err != nil {
			continue
		}

		for _, subEntry := range subModelEntries {
			if !subEntry.IsDir() {
				continue
			}

			toneName := subEntry.Name()
			tonePath := filepath.Join(subModelPath, toneName)

			audioPaths, err := getAudioFilesInDir(tonePath)
			if err != nil || len(audioPaths) == 0 {
				continue
			}

			tones = append(tones, ReferenceAudioTone{
				Tone:  toneName,
				Paths: audioPaths,
			})
		}
	}

	return tones, nil
}

// processStandardModel 处理标准模型结构，将其作为"0"子模型处理
func processStandardModel(modelPath string) ([]ReferenceAudioTone, error) {
	entries, err := readDir(modelPath)
	if err != nil {
		return nil, err
	}

	tones := make([]ReferenceAudioTone, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		toneName := entry.Name()
		tonePath := filepath.Join(modelPath, toneName)

		audioPaths, err := getAudioFilesInDir(tonePath)
		if err != nil || len(audioPaths) == 0 {
			continue
		}

		tones = append(tones, ReferenceAudioTone{
			Tone:  toneName,
			Paths: audioPaths,
		})
	}

	return tones, nil
}

// getTonesInModel 获取指定模型下的所有语气
func getTonesInModel(rootDir, modelName string) (string, []ReferenceAudioTone, error) {
	modelPath := filepath.Join(rootDir, modelName)
	entries, err := readDir(modelPath)
	if err != nil {
		return "", nil, err
	}

	var modelID string
	var tones []ReferenceAudioTone

	if hasSubModels(entries) {
		modelID = SubModel0 // 默认使用0作为模型序号
		tones, err = processSubModels(modelPath)
	} else {
		modelID = SubModel0 // 标准结构也使用0作为模型序号
		tones, err = processStandardModel(modelPath)
	}

	if err != nil {
		return "", nil, err
	}

	return modelID, tones, nil
}

// BuildReferenceAudioList 构建参考音频列表结构
func BuildReferenceAudioList(rootDir string) ([]ReferenceAudioModel, error) {
	entries, err := readDir(rootDir)
	if err != nil {
		return nil, err
	}

	models := make([]ReferenceAudioModel, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		modelName := entry.Name()
		modelID, tones, err := getTonesInModel(rootDir, modelName)
		if err != nil {
			continue
		}

		if len(tones) > 0 {
			models = append(models, ReferenceAudioModel{
				ModelID: modelID,
				Model:   modelName,
				Tones:   tones,
			})
		}
	}

	return models, nil
}

// ListReferenceAudioFiles 遍历参考音频文件夹并将其内容序列化为JSON格式
func ListReferenceAudioFiles(rootDir string) (string, error) {
	models, err := BuildReferenceAudioList(rootDir)
	if err != nil {
		return "", err
	}

	audioList := ReferenceAudioList{
		Models: models,
	}

	jsonData, err := json.MarshalIndent(audioList, "", "  ")
	if err != nil {
		return "", fmt.Errorf("序列化JSON失败: %w", err)
	}

	return string(jsonData), nil
}

// SaveReferenceAudioListToFile 将参考音频列表保存到文件
func SaveReferenceAudioListToFile(rootDir, outputPath string) error {
	jsonData, err := ListReferenceAudioFiles(rootDir)
	if err != nil {
		return err
	}

	// 确保输出目录存在
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, DirPermission); err != nil {
		return fmt.Errorf("无法创建输出目录: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(outputPath, []byte(jsonData), FilePermission); err != nil {
		return fmt.Errorf("无法写入文件: %w", err)
	}

	return nil
}
