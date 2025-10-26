package model

// Dialogue 对话初解析模型
// Dialogue 对话初解析模型
type PreDialogue struct {
    Name        string `json:"name"`
    Id          string `json:"id"`
    Text        string `json:"text"`
    Step        int    `json:"step"`
    Motion      string `json:"motion"`
    Expression  string `json:"expression"`
    Model       string `json:"model"`
}

