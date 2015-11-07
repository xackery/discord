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
	GuildID              int `json:"guild_id,string,omitempty"`
	Name                 string
	PermissionOverwrites []PermissionOverwrites `json:"permission_overwrites,omitempty"`
	Topic                string
	Position             int    `json:",omitempty"`
	LastMessageID        string `json:"last_message_id"`
	Type                 string
	ID                   int `json:",string,omitempty"`
	IsPrivate            bool
}

//Private Channels between two users
type PrivateChannel struct {
	ID            string `json:"id"`
	Recipient     User   `json:"recipient"`
	LastMessageID string `json:"last_message_id"`
}

//Channel setting regarding if a permission overwrites
type PermissionOverwrites struct {
}

//Send a message to specified channel
func (c *Client) ChannelMessageSend(channelId int, messageText string) (responseMessage Message, err error) {
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
		err = fmt.Errorf("StatusCode: %d, %s", resp.StatusCode, string(body))
		return
	}

	err = json.Unmarshal(body, &responseMessage)
	return
}
