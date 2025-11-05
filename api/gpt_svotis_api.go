package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"myvoicego/model"
	"net/http"
)

// TTS API特有常量
const (
	// TTS API URLs
	DefaultTTSURL           = "http://127.0.0.1:9880/tts"
	DefaultSoVITSWeightsURL = "http://127.0.0.1:9880/set_sovits_weights"
	DefaultGPTWeightsURL    = "http://127.0.0.1:9880/set_gpt_weights"
)

// GPTSvotisAPI 结构体用于封装 TTS API 的配置
type GPTSvotisAPI struct {
	TTSURL           string
	SoVITSWeightsURL string
	GPTWeightsURL    string
	client           *http.Client // HTTP客户端，支持超时设置
}

// NewGPTSvotisAPI 创建一个新的 GPTSvotisAPI 实例
func NewGPTSvotisAPI() *GPTSvotisAPI {
	return &GPTSvotisAPI{
		TTSURL:           DefaultTTSURL,
		SoVITSWeightsURL: DefaultSoVITSWeightsURL,
		GPTWeightsURL:    DefaultGPTWeightsURL,
		client: &http.Client{
			Timeout: RequestTimeout,
		},
	}
}

// GenerateTTS 生成TTS语音
func (api *GPTSvotisAPI) GenerateTTS(req model.TTSRequest) ([]byte, error) {
	// 验证输入
	if req.Text == "" {
		return nil, fmt.Errorf("文本内容不能为空")
	}
	if req.TextLang == "" {
		return nil, fmt.Errorf("文本语言不能为空")
	}
	if req.RefAudioPath == "" {
		return nil, fmt.Errorf("参考音频路径不能为空")
	}
	if req.PromptLang == "" {
		return nil, fmt.Errorf("提示文本语言不能为空")
	}

	// 将请求体转换为 JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建 HTTP 请求
	httpReq, err := http.NewRequest("POST", api.TTSURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
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

// SetSoVITSWeights 设置SoVITS模型权重
func (api *GPTSvotisAPI) SetSoVITSWeights(weightsPath string) ([]byte, error) {
	// 验证输入
	if weightsPath == "" {
		return nil, fmt.Errorf("权重路径不能为空")
	}

	// 构建请求URL
	url := fmt.Sprintf("%s?weights_path=%s", api.SoVITSWeightsURL, weightsPath)

	// 创建 HTTP 请求
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

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

// SetGPTWeights 设置GPT模型权重
func (api *GPTSvotisAPI) SetGPTWeights(weightsPath string) ([]byte, error) {
	// 验证输入
	if weightsPath == "" {
		return nil, fmt.Errorf("权重路径不能为空")
	}

	// 构建请求URL
	url := fmt.Sprintf("%s?weights_path=%s", api.GPTWeightsURL, weightsPath)

	// 创建 HTTP 请求
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

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
