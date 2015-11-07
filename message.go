package discord

//Posted message inside a channel
type Message struct {
	Attachments     []Attachment
	Tts             bool
	Embeds          []Embed
	Timestamp       string
	MentionEveryone bool
	Nonce           string //Used by websocket
	ID              int    `json:"id,string,omitempty"`
	EditedTimestamp string
	Author          User
	Content         string
	ChannelID       int `json:"channel_id,string,omitempty"`
	Mentions        []User
}

//Attachment array entry
type Attachment struct {
	Width    int
	Url      string
	Size     int
	ProxyURL string
	ID       int `json:"id,string"`
	Filename string
}

//Embedded media array entry
type Embed struct {
}

//An author of a message
type Author struct {
	Username      string
	Discriminator int `json:",string"`
	ID            int `json:"id,string"`
	Avatar        string
}
