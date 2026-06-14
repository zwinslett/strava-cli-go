package telegram

type Update struct {
	OK     bool     `json:"ok"`
	Result []Result `json:"result"`
}

type Result struct {
	UpdateID int64   `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageID int64  `json:"message_id"`
	Text      string `json:"text"`
}
