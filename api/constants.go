package api

import "time"

// API通用常量
const (
	// 请求配置
	RequestTimeout  = 120 * time.Second // 增加到120秒，以处理大量对话
	MaxResponseSize = 10 << 20          // 10MB

	// 请求头常量
	HeaderContentType = "Content-Type"
	ContentTypeJSON   = "application/json"
)
