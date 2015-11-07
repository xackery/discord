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

//A guild are basically servers on Discord
type Guild struct {
	VoiceStates    []VoiceState //Used by Websocket
	JoinedAt       string
	AfkChannelID   int `json:"afk_channel_id,string,omitempty"`
	AfkTimeout     int
	ID             int `json:"id,string,omitempty"`
	Icon           int `json:"omitempty"`
	Name           string
	Roles          []Role
	Region         string
	Presences      []Presence //Used by Websocket
	EmbedChannelID int        `json:"embed_channel_id,string,omitempty"`
	EmbedEnabled   bool
	OwnerID        int `json:"owner_id,string,omitempty"`
	Members        []Member
	Large          bool      //Used by Websocket
	Channels       []Channel //Used by Websocket
}

//Roles are permission groupings
type Role struct {
	Managed     bool
	Name        string
	Color       int
	Hoist       bool
	Position    int
	ID          int `json:"id,string"`
	Permissions int
}

type VoiceState struct {
	UserID    int `json:"string"`
	Suppress  bool
	SessionID string `json:"session_id"`
	SelfMute  bool
	SelfDeaf  bool
	Mute      bool
	Deaf      bool
	ChannelID int `json:"channel_id,string"`
}

type Member struct {
	User     User
	Roles    []string
	Mute     bool
	JoinedAt string
	Deaf     bool
}

//List channels found on given guildId
func (c *Client) GuildChannels(guildId int) (channels []Channel, err error) {
	if !c.IsLoggedIn() {
		err = errors.New("You must be logged in")
		return
	}
	httpClient := &http.Client{Timeout: (20 * time.Second)}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", DISCORD_URL, fmt.Sprintf("guilds/%d/channels", guildId)), bytes.NewBuffer([]byte(fmt.Sprintf(``))))
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
		err = errors.New(fmt.Sprintf("StatusCode: %d, %s", resp.StatusCode, string(body)))
		return
	}

	err = json.Unmarshal(body, &channels)
	return
}
