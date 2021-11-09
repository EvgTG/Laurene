package unmarshal_export

import "time"

type Message struct {
	ID        int64
	Type      string // message animation sticker photo voice_message audio_file video_file video_message
	Date      time.Time
	FromID    string
	ReplyToID int64
	Text      string
	Hashtags  []string
	Raw       string
}

type messageForUnmarshal struct {
	ID           float64     `json:"id"`
	Type         string      `json:"type"`
	Date         string      `json:"date"`
	Edited       string      `json:"edited"`
	FromID       string      `json:"from_id"`
	ReplyToID    float64     `json:"reply_to_message_id"`
	Text         interface{} `json:"text"`
	MediaType    string      `json:"media_type"`
	StickerEmoji string      `json:"sticker_emoji"`
	Photo        string      `json:"photo"`
	Performer    string      `json:"performer"`
	Title        string      `json:"title"`
}
