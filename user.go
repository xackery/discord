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

//Fetch information about user logged in (me)
func (c *Client) UserMe() (err error) {
	if !c.IsLoggedIn() {
		err = errors.New("You must be logged in")
		return
	}

	httpClient := &http.Client{Timeout: (20 * time.Second)}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", DISCORD_URL, "users/@me"), bytes.NewBuffer([]byte(fmt.Sprintf(``))))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", c.Token)
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("StatusCode: %d, %s", resp.StatusCode, string(body))
		return
	}
	err = json.Unmarshal(body, &c)
	return
}

//Get a list of guilds for current logged in user
func (c Client) UserGuilds() (guilds []Guild, err error) {
	if !c.IsLoggedIn() {
		err = errors.New("You must be logged in")
		return
	}
	httpClient := &http.Client{Timeout: (20 * time.Second)}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", DISCORD_URL, fmt.Sprintf("users/%s/guilds", c.Id)), bytes.NewBuffer([]byte(fmt.Sprintf(``))))
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
	resp.Body.Close()

	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("StatusCode: %d, %s", resp.StatusCode, string(body))
		return
	}

	err = json.Unmarshal(body, &guilds)
	return
}
