package canalx

// Message 可以根据需要把其他字段加进来
// 这里用作canal->kafka
type Message[T any] struct {
	Data     []T    `json:"data"`
	Database string `json:"database"`
	Table    string `json:"table"`
	Type     string `json:"type"`
}
