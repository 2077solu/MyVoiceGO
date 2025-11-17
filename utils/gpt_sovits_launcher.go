package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// StartGPTSvits 启动GPT-SoVITS服务
func StartGPTSvits() error {
	// 读取配置文件获取路径
	configPath := filepath.Join("config", "gpt_svotis_path.json")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config struct {
		GSV_Path string `json:"GSV_Path"`
	}
	if err := json.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 构建命令
	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf("cd %s; ./runtime/python api_v2.py -a 127.0.0.1 -p 9880 -c GPT_SoVITS/configs/tts_infer.yaml", config.GSV_Path))

	// 启动命令
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动GPT-SoVITS失败: %v", err)
	}

	fmt.Println("GPT-SoVITS服务已启动")
	return nil
}
