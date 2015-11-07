package discord

import (
	"errors"
)

//Client wrapper for Discord
type Client struct {
	Id    string
	Token string
}

//Login Method for Discord
func (c *Client) Login(email string, pass string) (err error) {
	if len(email) < 3 {
		err = errors.New("email too short")
		return
	}
	if len(pass) < 3 {
		err = errors.New("password too short")
		return
	}
	if !c.IsLoggedIn() {
		err = c.authLogin(email, pass)
		if err != nil {
			return
		}
	}
	c.UserMe()
	return
}

//Checks if client IsLoggedIn to discord
func (c Client) IsLoggedIn() (isLoggedIn bool) {
	return c.Token != ""
}
