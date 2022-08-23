package entity

type Message struct {
	Text               string `json:"text"`
	Telegram_chat_id   int64  `json:"telegram-chat-id"`
	Discord_channel_id string `json:"discord-channel-id"`
}
