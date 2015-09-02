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

type Response struct {
}

type Request struct {
}

var discordUrl = "http://discordapp.com/api"

type RestClient struct {
	Url     string
	Session *Session
	client  *http.Client
}

type Session struct {
	Email    string
	Password string
	Token    string
}

// Create takes an email and password then prepares a RestClient with the given data,
// which is a simple object used for future requests.
func Create(email string, password string) (restClient *RestClient, e error) {
	if len(email) < 3 {
		e = errors.New("email too short")
		return
	}
	if len(password) < 3 {
		e = errors.New("password too short")
		return
	}
	session := &Session{email, password, ""}
	httpClient := &http.Client{Timeout: (20 * time.Second)}
	restClient = &RestClient{discordUrl, session, httpClient}
	restClient.Session.Token, e = requestToken(restClient)
	if e != nil {
		return
	}

	return
}

// RequestToken asks the Rest server for a token by provided email/password
func requestToken(restClient *RestClient) (token string, e error) {

	if restClient == nil {
		e = errors.New("Empty restClient, Create() one first")
		return
	}

	if restClient.Session == nil || len(restClient.Session.Email) == 0 || len(restClient.Session.Password) == 0 {
		e = errors.New("Empty restClient.Session data, Create() to set email/password")
		return
	}

	req, e := http.NewRequest("POST", fmt.Sprintf("%s/%s", restClient.Url, "auth/login"), bytes.NewBuffer([]byte(fmt.Sprintf(`{"email":"%s", "password":"%s"}`, restClient.Session.Email, restClient.Session.Password))))
	if e != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, e := restClient.client.Do(req)
	if e != nil {
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		e = errors.New(fmt.Sprintf("StatusCode: %d, %s", resp.StatusCode, string(body)))
		return
	}
	session := &Session{}
	e = json.Unmarshal(body, &session)
	token = session.Token
	return
}
