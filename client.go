package discord

import (
	"errors"
	"net/http"
	"time"
)

type Response struct {
}

type Request struct {
}

var discordUrl = "http://discordapp.com"

type RestClient struct {
	url    string
	client *http.Client
}

func Create(email string, password string) (client *RestClient, e error) {
	if len(email) < 3 {
		e = errors.New("email too short")
		return
	}
	if len(password) < 3 {
		e = errors.New("password too short")
		return
	}

	httpClient := &http.Client{Timeout: (20 * time.Second)}
	client = &RestClient{discordUrl, httpClient}

	return
}
