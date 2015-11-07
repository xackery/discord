package discord

import (
	"errors"
)

type Client struct {
	Id    string
	Token string
}

//Log in to Discord
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

//Check if client is logged in to discord
func (c Client) IsLoggedIn() (isLoggedIn bool) {
	return c.Token != ""
}
