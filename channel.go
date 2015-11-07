package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//Channels are sections inside guilds that messages are grouped by
type Channel struct {
	GuildID              int `json:",string,omitempty"`
	Name                 string
	PermissionOverWrites []PermissionOverWrites `json:"permission_overwrites,omitempty"`
	Topic                string
	Position             int `json:",omitempty"`
	LastMessageID        string
	Type                 string
	Id                   int `json:",string,omitempty"`
	IsPrivate            bool
}

//Channel setting regarding if a permission overwrites
type PermissionOverWrites struct {
}

//Posted message inside a channel
type Message struct {
	Attachments     []Attachment
	Tts             bool
	Embeds          []Embed
	Timestamp       string
	MentionEveryone bool
	ID              int `json:",string,omitempty"`
	EditedTimestamp string
	Author          *Author
	Content         string
	ChannelID       int `json:",string,omitempty"`
	Mentions        []Mention
}

//Mentions are @user inside a message
type Mention struct {
}

//Attachment array entry
type Attachment struct {
}

//Embedded media array entry
type Embed struct {
}

//An author of a message
type Author struct {
	Username      string
	Discriminator int `json:",string"`
	ID            int `json:",string"`
	Avatar        string
}

//Send a message to specified channel
func (c Client) ChannelMessageSend(channelId int, messageText string) (responseMessage Message, err error) {
	if !c.IsLoggedIn() {
		err = errors.New("You must be logged in")
		return
	}
	httpClient := &http.Client{Timeout: (20 * time.Second)}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", DISCORD_URL, fmt.Sprintf("channels/%d/messages", channelId)), bytes.NewBuffer([]byte(fmt.Sprintf(`{"content":"%s"}`, messageText))))
	if err != nil {
		return
	}
	req.Header.Set("authorization", c.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("StatusCode: %d, %s", resp.StatusCode, string(body)))
		return
	}

	err = json.Unmarshal(body, &responseMessage)
	return
}
