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

//Request a token from web server
func (c *Client) authLogin(email string, pass string) (err error) {

	httpClient := &http.Client{Timeout: (20 * time.Second)}
	var resp *http.Response
	resp, err = httpClient.Post(DISCORD_URL+"auth/login", "application/json", bytes.NewBuffer([]byte(fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, pass))))
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("StatusCode: %d, %s", resp.StatusCode, string(body)))
		return
	}
	err = json.Unmarshal(body, &c)
	if err != nil {
		return
	}

	return
}

func (c *Client) authLogout() (err error) {
	if !c.IsLoggedIn() {
		err = errors.New("You must be logged in")
		return
	}
	httpClient := &http.Client{Timeout: (20 * time.Second)}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", DISCORD_URL, fmt.Sprintf("auth/logout")), bytes.NewBuffer([]byte(fmt.Sprintf(``))))
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

	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("StatusCode: %d, %s", resp.StatusCode, string(body)))
		return
	}
	return
}
