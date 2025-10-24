package model

// Dialogue 对话初解析模型
type FirstDialogue struct {
	Name    string `json:"name"`
	Id      string `json:"id"`
	Text    string `json:"text"`
	step    int    `json:"step"`
	motion  string `json:"motion"`
	expression string `json:"expression"`
	model 	string `json:"model"`
}
