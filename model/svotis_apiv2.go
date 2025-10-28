package model

// TTSRequest 表示发送到TTS API的请求
type TTSRequest struct {
	Text              string   `json:"text"`                // (required) 要合成的文本
	TextLang          string   `json:"text_lang"`           // (required) 要合成文本的语言
	RefAudioPath      string   `json:"ref_audio_path"`      // (required) 参考音频路径
	AuxRefAudioPaths  []string `json:"aux_ref_audio_paths"` // (optional) 多说话人音色融合的辅助参考音频路径
	PromptText        string   `json:"prompt_text"`         // (optional) 参考音频的提示文本
	PromptLang        string   `json:"prompt_lang"`         // (required) 参考音频提示文本的语言
	TopK              int      `json:"top_k"`               // top k 采样
	TopP              float64  `json:"top_p"`               // top p 采样
	Temperature       float64  `json:"temperature"`         // 采样的温度
	TextSplitMethod   string   `json:"text_split_method"`   // 文本分割方法
	BatchSize         int      `json:"batch_size"`          // 推理的批大小
	BatchThreshold    float64  `json:"batch_threshold"`     // 批分割的阈值
	SplitBucket       bool     `json:"split_bucket"`        // 是否将批次分割到多个桶
	SpeedFactor       float64  `json:"speed_factor"`        // 控制合成音频的速度
	FragmentInterval  float64  `json:"fragment_interval"`   // 控制音频片段的间隔
	Seed              int      `json:"seed"`                // 用于可重复性的随机种子
	MediaType         string   `json:"media_type"`          // 输出音频的媒体类型
	StreamingMode     bool     `json:"streaming_mode"`      // 是否返回流式响应
	ParallelInfer     bool     `json:"parallel_infer"`      // 是否使用并行推理
	RepetitionPenalty float64  `json:"repetition_penalty"`  // T2S模型的重复惩罚
	SampleSteps       int      `json:"sample_steps"`        // VITS模型V3的采样步数
	SuperSampling     bool     `json:"super_sampling"`      // VITS模型V3时是否使用超采样
}

// DefaultTTSRequest 创建一个带有默认值的TTS请求
// 必需参数仍需提供：text, textLang, refAudioPath, promptLang
func DefaultTTSRequest(text, textLang, refAudioPath, promptLang string) TTSRequest {
	return TTSRequest{
		Text:              text,
		TextLang:          textLang,
		RefAudioPath:      refAudioPath,
		PromptLang:        promptLang,
		TopK:              5,
		TopP:              1.0,
		Temperature:       1.0,
		TextSplitMethod:   "cut5",
		BatchSize:         1,
		BatchThreshold:    0.75,
		SplitBucket:       true,
		SpeedFactor:       1.0,
		FragmentInterval:  0.3,
		Seed:              -1,
		MediaType:         "wav",
		StreamingMode:     false,
		ParallelInfer:     true,
		RepetitionPenalty: 1.35,
		SampleSteps:       32,
		SuperSampling:     false,
	}
}

// ControlRequest 表示控制API的请求
type ControlRequest struct {
	Command string `json:"command"` // 控制命令: "restart" 或 "exit"
}

// SetGPTWeightsRequest 表示设置GPT模型权重的请求
type SetGPTWeightsRequest struct {
	WeightsPath string `json:"weights_path"` // GPT模型权重文件路径
}

// SetSoVITSWeightsRequest 表示设置SoVITS模型权重的请求
type SetSoVITSWeightsRequest struct {
	WeightsPath string `json:"weights_path"` // SoVITS模型权重文件路径
}

// SetReferAudioRequest 表示设置参考音频的请求
type SetReferAudioRequest struct {
	ReferAudioPath string `json:"refer_audio_path"` // 参考音频路径
}

// APIError 表示API返回的错误信息
type APIError struct {
	Message   string `json:"message"`   // 错误信息
	Exception string `json:"Exception"` // 异常详情
}

// APIResponse 表示API的通用响应
type APIResponse struct {
	Message string `json:"message"` // 响应消息
}
