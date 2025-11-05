package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"myvoicego/model"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Coze API特有常量
const (
	// DefaultAPIURL API相关常量
	DefaultAPIURL = "https://api.coze.cn/v3/chat"

	// HeaderAuth 请求头常量
	HeaderAuth = "Authorization"

	// RoleUser 消息类型常量
	RoleUser        = "user"
	TypeToolOutput  = "tool_output"
	ContentTypeText = "text"

	// EventCompleted 事件类型常量
	EventCompleted = "conversation.message.completed"
	TypeAnswer     = "answer"
)

// CozeAPI 结构体用于封装 Coze API 的配置
type CozeAPI struct {
	BaseURL     string
	BearerToken string
	BotID       string
	UserID      string
	client      *http.Client // HTTP客户端，支持超时设置
}

// NewCozeAPI 创建一个新的 CozeAPI 实例
func NewCozeAPI(bearerToken, botID, userID string) *CozeAPI {
	return &CozeAPI{
		BaseURL:     DefaultAPIURL,
		BearerToken: bearerToken,
		BotID:       botID,
		UserID:      userID,
		client: &http.Client{
			Timeout: RequestTimeout,
		},
	}
}

// Config 表示Coze API的配置结构
type Config struct {
	APIURL          string `json:"api_url"`
	Token           string `json:"token"`
	BotID           string `json:"bot_id"`
	UserID          string `json:"user_id"`
	Stream          bool   `json:"stream"`
	AutoSaveHistory bool   `json:"auto_save_history"`
}

// NewCozeAPIFromConfig 从配置文件创建 CozeAPI 实例
func NewCozeAPIFromConfig(configPath string) (*CozeAPI, error) {
	// 读取配置文件
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析配置文件
	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %v", err)
	}

	// 使用默认API URL（如果未指定）
	apiURL := config.APIURL
	if apiURL == "" {
		apiURL = DefaultAPIURL
	}

	return &CozeAPI{
		BaseURL:     apiURL,
		BearerToken: config.Token,
		BotID:       config.BotID,
		UserID:      config.UserID,
		client: &http.Client{
			Timeout: RequestTimeout,
		},
	}, nil
}

// validateConfig 验证配置是否有效
func validateConfig(config *Config) error {
	if config.Token == "" {
		return fmt.Errorf("token 不能为空")
	}
	if config.BotID == "" {
		return fmt.Errorf("bot_id 不能为空")
	}
	if config.UserID == "" {
		return fmt.Errorf("user_id 不能为空")
	}
	return nil
}

// Message 表示发送给 Coze API 的消息结构
type Message struct {
	Role        string `json:"role"`
	Type        string `json:"type"`
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
}

// Request 表示发送给 Coze API 的请求结构
type Request struct {
	BotID              string    `json:"bot_id"`
	Stream             bool      `json:"stream"`
	AutoSaveHistory    bool      `json:"auto_save_history"`
	AdditionalMessages []Message `json:"additional_messages"`
	UserID             string    `json:"user_id"`
}

// SendDialogueToCoze 发送对话内容到 Coze API 进行情感解析
func (api *CozeAPI) SendDialogueToCoze(dialogueJSON string) ([]byte, error) {
	// 验证输入
	if dialogueJSON == "" {
		return nil, fmt.Errorf("对话内容不能为空")
	}

	// 创建请求体
	req := Request{
		BotID:           api.BotID,
		Stream:          true,
		AutoSaveHistory: false,
		AdditionalMessages: []Message{
			{
				Role:        RoleUser,
				Type:        TypeToolOutput,
				ContentType: ContentTypeText,
				Content:     dialogueJSON,
			},
		},
		UserID: api.UserID,
	}

	// 将请求体转换为 JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建带上下文的 HTTP 请求
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", api.BaseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	httpReq.Header.Set(HeaderAuth, "Bearer "+api.BearerToken)
	httpReq.Header.Set(HeaderContentType, ContentTypeJSON)

	// 发送请求
	resp, err := api.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 限制响应体大小
	limitedReader := io.LimitReader(resp.Body, MaxResponseSize)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 返回错误状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// ReadDialogueFromFile 从指定目录读取对话文件
func ReadDialogueFromFile(figuresDir, filename string) (string, error) {
	// 验证输入参数
	if figuresDir == "" {
		return "", fmt.Errorf("目录路径不能为空")
	}
	if filename == "" {
		return "", fmt.Errorf("文件名不能为空")
	}

	// 确保文件名有.json扩展名
	if filepath.Ext(filename) != ".json" {
		filename += ".json"
	}

	// 构建完整文件路径
	filePath := filepath.Join(figuresDir, filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("文件不存在: %s", filePath)
	}

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	return string(content), nil
}

// ListAvailableDialogues 列出可用的对话文件
func ListAvailableDialogues(figuresDir string) ([]string, error) {
	// 验证输入参数
	if figuresDir == "" {
		return nil, fmt.Errorf("目录路径不能为空")
	}

	// 确保目录存在
	if _, err := os.Stat(figuresDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("目录不存在: %s", figuresDir)
	}

	// 读取目录内容
	files, err := os.ReadDir(figuresDir)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %v", err)
	}

	// 预分配切片容量，提高性能
	dialogues := make([]string, 0, len(files))

	// 过滤出.json文件
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			// 安全地移除.json扩展名
			name := strings.TrimSuffix(file.Name(), ".json")
			dialogues = append(dialogues, name)
		}
	}

	return dialogues, nil
}

// EmotionResponse 表示Coze API返回的情绪分析响应
type EmotionResponse struct {
	Emotions []EmotionResult `json:"emotions"`
}

// EmotionResult 表示单个情绪分析结果
type EmotionResult struct {
	Index   int    `json:"index"`
	Emotion string `json:"emotion"`
}

// AnalyzeEmotions 分析对话的情绪
func (api *CozeAPI) AnalyzeEmotions(dialogues []model.PreDialogue) ([]model.PreDialogue, error) {
	// 验证输入
	if len(dialogues) == 0 {
		return nil, fmt.Errorf("对话列表不能为空")
	}

	// 将对话转换为JSON
	dialogueJSON, err := json.Marshal(dialogues)
	if err != nil {
		return nil, fmt.Errorf("序列化对话失败: %v", err)
	}

	// 发送请求到Coze API
	response, err := api.SendDialogueToCoze(string(dialogueJSON))
	if err != nil {
		return nil, fmt.Errorf("发送情绪分析请求失败: %v", err)
	}

	// 解析流式响应
	responseStr := string(response)
	lines := strings.Split(responseStr, "\n")

	// 查找最终结果
	finalContent, err := extractFinalContent(lines)
	if err != nil {
		return nil, err
	}

	// 解析最终内容为带有情绪的对话数组
	var emotionDialogues []model.PreDialogue
	if err := json.Unmarshal([]byte(finalContent), &emotionDialogues); err != nil {
		return nil, fmt.Errorf("解析情绪分析结果失败: %v", err)
	}

	// 创建step到emotion的映射
	stepToEmotion := make(map[int]string, len(emotionDialogues))
	for _, dialogue := range emotionDialogues {
		if dialogue.Emotion != "" {
			stepToEmotion[dialogue.Step] = dialogue.Emotion
		}
	}

	// 更新原始对话的情绪
	for i := range dialogues {
		if emotion, exists := stepToEmotion[dialogues[i].Step]; exists {
			dialogues[i].Emotion = emotion
		}
	}

	return dialogues, nil
}

// extractFinalContent 从流式响应中提取最终内容
func extractFinalContent(lines []string) (string, error) {
	for i, line := range lines {
		if !strings.HasPrefix(line, "event:") {
			continue
		}

		event := strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		if event != EventCompleted {
			continue
		}

		if i+1 >= len(lines) {
			continue
		}

		nextLine := lines[i+1]
		if !strings.HasPrefix(nextLine, "data:") {
			continue
		}

		dataStr := strings.TrimSpace(strings.TrimPrefix(nextLine, "data:"))
		content, err := parseDataContent(dataStr)
		if err != nil {
			continue
		}

		return content, nil
	}

	return "", fmt.Errorf("未找到有效的情绪分析结果")
}

// parseDataContent 解析数据内容
func parseDataContent(dataStr string) (string, error) {
	var data struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	}

	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return "", err
	}

	if data.Type != TypeAnswer {
		return "", fmt.Errorf("不是答案类型")
	}

	return data.Content, nil
}
