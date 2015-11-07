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

type Mention struct {
}

type Attachment struct {
}

type Embed struct {
}

type Author struct {
	Username      string
	Discriminator int `json:",string"`
	Id            int `json:",string"`
	Avatar        string
}

type Channel struct {
	Guild_id              int `json:",string"`
	Name                  string
	Permission_overwrites string
	Position              int `json:",string"`
	Last_message_id       string
	Type                  string
	Id                    int `json:",string"`
	Is_private            string
}

type Message struct {
	Attachments      []Attachment
	Tts              bool
	Embeds           []Embed
	Timestamp        string
	Mention_everyone bool
	Id               int `json:",string"`
	Edited_timestamp string
	Author           *Author
	Content          string
	Channel_id       int `json:",string"`
	Mentions         []Mention
}

func (c Client) SendMessage(channelId int, messageText string) (responseMessage Message, err error) {
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
		err = errors.New(fmt.Sprintf("StatusCode: %d, %s", resp.StatusCode, string(body)))
		return
	}

	err = json.Unmarshal(body, &responseMessage)
	return
}
