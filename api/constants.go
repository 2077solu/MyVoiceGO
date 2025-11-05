package api

import "time"

// API通用常量
const (
	// 请求配置
	RequestTimeout  = 30 * time.Second
	MaxResponseSize = 10 << 20 // 10MB

	// 请求头常量
	HeaderContentType = "Content-Type"
	ContentTypeJSON   = "application/json"
)
