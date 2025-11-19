package main

import (
	"myvoicego/ui/views"
)

func main() {
	// 创建主视图
	mainView := views.NewMainView()

	// 运行应用
	mainView.Run()
}
