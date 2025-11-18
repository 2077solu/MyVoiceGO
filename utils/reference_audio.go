package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	DirPermission  = 0755
	FilePermission = 0644
)

// 支持的音频文件扩展名
var audioExts = map[string]struct{}{
	".mp3": {}, ".wav": {}, ".ogg": {}, ".flac": {}, ".m4a": {}, ".aac": {},
}

// AudioRefPath 音频文件路径
type AudioRefPath struct {
	Path string `json:"path"`
}

// AudioRefTone 语气对应的音频文件列表
type AudioRefTone struct {
	Tone  string         `json:"tone"`
	Paths []AudioRefPath `json:"paths"`
}

// AudioRefSubDir 子模型对应的语气列表
type AudioRefSubDir struct {
	AudioId string         `json:"audioid"`
	Tones   []AudioRefTone `json:"tones"`
}

// AudioRefModel 模型对应的语气列表
type AudioRefModel struct {
	Model   string           `json:"model"`
	SubDirs []AudioRefSubDir `json:"subdirs"`
}

// AudioRefList 所有模型的音频参考列表
type AudioRefList struct {
	Models []AudioRefModel `json:"models"`
}

// normalizePath 统一路径格式为Unix风格，并确保使用绝对路径
func normalizePath(path string) string {
	p, _ := filepath.Abs(path)
	return filepath.ToSlash(p)
}

// isAudioFile 检查文件是否为支持的音频格式
func isAudioFile(name string) bool {
	_, ok := audioExts[strings.ToLower(filepath.Ext(name))]
	return ok
}

// readDir 读取目录内容
func readDir(path string) ([]os.DirEntry, error) {
	d, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败 %q: %w", path, err)
	}
	return d, nil
}

// getAudioFiles 获取目录下的音频文件，使用绝对路径
func getAudioFiles(dir string) ([]AudioRefPath, error) {
	entries, err := readDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []AudioRefPath
	for _, entry := range entries {
		if entry.IsDir() || !isAudioFile(entry.Name()) {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		paths = append(paths, AudioRefPath{Path: normalizePath(path)})
	}
	return paths, nil
}

// collectTones 收集语气信息
func collectTones(baseDir string) ([]AudioRefTone, error) {
	entries, err := readDir(baseDir)
	if err != nil {
		return nil, err
	}

	var tones []AudioRefTone
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		path := filepath.Join(baseDir, name)
		paths, err := getAudioFiles(path)
		if err != nil || len(paths) == 0 {
			continue
		}

		tones = append(tones, AudioRefTone{Tone: name, Paths: paths})
	}
	return tones, nil
}

// getSubDirs 获取模型的子目录和语气
func getSubDirs(rootDir, modelName string) ([]AudioRefSubDir, error) {
	modelPath := filepath.Join(rootDir, modelName)
	entries, err := readDir(modelPath)
	if err != nil {
		return nil, err
	}

	var subDirs []AudioRefSubDir
	// 检查是否有子目录
	hasSubDirs := false

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		hasSubDirs = true
		name := entry.Name()
		path := filepath.Join(modelPath, name)
		tones, err := collectTones(path)
		if err != nil {
			continue
		}

		subDirs = append(subDirs, AudioRefSubDir{AudioId: name, Tones: tones})
	}

	// 没有子目录，直接从模型目录收集语气
	if !hasSubDirs {
		tones, err := collectTones(modelPath)
		if err != nil {
			return nil, err
		}
		subDirs = append(subDirs, AudioRefSubDir{AudioId: "", Tones: tones})
	}

	return subDirs, nil
}

// BuildReferenceAudioList 构建所有模型的音频参考列表
func BuildReferenceAudioList(rootDir string) ([]AudioRefModel, error) {
	entries, err := readDir(rootDir)
	if err != nil {
		return nil, err
	}

	var models []AudioRefModel
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		subDirs, err := getSubDirs(rootDir, name)
		if err != nil {
			continue
		}

		// 检查是否有语气数据
		hasTones := false
		for _, subDir := range subDirs {
			if len(subDir.Tones) > 0 {
				hasTones = true
				break
			}
		}

		if hasTones {
			models = append(models, AudioRefModel{Model: name, SubDirs: subDirs})
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

	data, err := json.MarshalIndent(AudioRefList{Models: models}, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON序列化失败: %w", err)
	}

	return string(data), nil
}

// SaveReferenceAudioListToFile 将音频参考列表保存到文件
func SaveReferenceAudioListToFile(rootDir, outputPath string) error {
	data, err := ListReferenceAudioFiles(rootDir)
	if err != nil {
		return err
	}

	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(outputPath), DirPermission); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(outputPath, []byte(data), FilePermission); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}
