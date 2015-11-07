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

type Guild struct {
	Afk_timeout      int
	Joined_at        string
	Afk_channel_id   int `json:",string,omitempty"`
	Id               int `json:",string,omitempty"`
	Icon             int `json:"omitempty"`
	Name             string
	Roles            []Role
	Region           string
	Embed_channel_id int `json:",string,omitempty"`
	Embed_enabled    bool
	Owner_id         int `json:",string,omitempty"`
}

type Role struct {
	Managed     bool
	Name        string
	Color       int
	Hoist       bool
	Position    int
	Id          int `json:",string"`
	Permissions int
}

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
