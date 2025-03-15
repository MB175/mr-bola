package notes

// Note represents a simple note entity
type Note struct {
	ID      string `json:"id"`
	Owner   string `json:"owner"`
	Content string `json:"content"`
}
