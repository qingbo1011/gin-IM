package ws

type Trainer struct {
	Content   string `json:"content"`    // 内容
	StartTime int64  `json:"start_time"` // 创建时间
	EndTime   int64  `json:"end_time"`   // 过期时间
	Read      bool   `json:"read"`       // 是否已读
}

type Result struct {
	StartTime int64
	Msg       string
	Content   any
	From      string
}
