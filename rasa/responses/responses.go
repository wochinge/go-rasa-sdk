package responses

type Button struct {
	Title   string `json:"title"`
	PayLoad string `json:"payload"`
}

type BotMessage struct {
	Text string `json:"text"`
	Template string `json:"template,omitempty"`
	Elements     []interface{} `json:"elements,omitempty"`
	QuickReplies []interface{} `json:"quick_replies,omitempty"`
	Buttons      []Button      `json:"buttons,omitempty"`
	Attachment   interface{}   `json:"attachment,omitempty"`
	Image        string        `json:"image,omitempty"`
	Custom       interface{}   `json:"custom,omitempty"`
}

