package discord

type TypingEvent struct {
	UserID    int `json:"user_id,string"`
	Timestamp int
	ChannelID int `json:"channel_id,string"`
}
