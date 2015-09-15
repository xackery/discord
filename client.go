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

var discordUrl = "http://discordapp.com/api"

type RestClient struct {
	Url     string
	Session *Session
	client  *http.Client
}

type Session struct {
	Id       string
	Email    string
	Password string
	Token    string
}

type Guild struct {
	Afk_timeout      int
	Joined_at        string
	Afk_channel_id   int `json:",string"`
	Id               int `json:",string"`
	Icon             int
	Name             string
	Roles            []Role
	Region           string
	Embed_channel_id int `json:",string"`
	Embed_enabled    bool
	Owner_id         int `json:",string"`
}

type Role struct {
	Permissions int
	Id          int `json:",string"`
	Name        string
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

// Create takes an email and password then prepares a RestClient with the given data,
// which is a simple object used for future requests.
func Create(email string, password string) (restClient *RestClient, err error) {
	if len(email) < 3 {
		err = errors.New("email too short")
		return
	}
	if len(password) < 3 {
		err = errors.New("password too short")
		return
	}
	session := &Session{"0", email, password, ""}
	httpClient := &http.Client{Timeout: (20 * time.Second)}
	restClient = &RestClient{discordUrl, session, httpClient}
	restClient.Session.Token, err = requestToken(restClient)
	if err != nil {
		return
	}
	restClient.Session.Id, err = requestSelf(restClient)
	if err != nil {
		return
	}
	return
}

// RequestToken asks the Rest server for a token by provided email/password
func requestToken(restClient *RestClient) (token string, err error) {

	if restClient == nil {
		err = errors.New("Empty restClient, Create() one first")
		return
	}

	if restClient.Session == nil || len(restClient.Session.Email) == 0 || len(restClient.Session.Password) == 0 {
		err = errors.New("Empty restClient.Session data, Create() to set email/password")
		return
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", restClient.Url, "auth/login"), bytes.NewBuffer([]byte(fmt.Sprintf(`{"email":"%s", "password":"%s"}`, restClient.Session.Email, restClient.Session.Password))))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := restClient.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("StatusCode: %d, %s", resp.StatusCode, string(body)))
		return
	}
	session := &Session{}
	err = json.Unmarshal(body, &session)
	token = session.Token
	return
}

// Identify user himself
func requestSelf(restClient *RestClient) (clientId string, err error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", restClient.Url, "users/@me"), bytes.NewBuffer([]byte(fmt.Sprintf(``))))
	if err != nil {
		return
	}
	req.Header.Set("authorization", restClient.Session.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := restClient.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("StatusCode: %d, %s", resp.StatusCode, string(body)))
		return
	}
	session := &Session{}
	err = json.Unmarshal(body, &session)
	clientId = session.Id
	return
}

func ListGuilds(restClient *RestClient) (guilds []Guild, err error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", restClient.Url, fmt.Sprintf("users/%s/guilds", restClient.Session.Id)), bytes.NewBuffer([]byte(fmt.Sprintf(``))))
	if err != nil {
		return
	}
	req.Header.Set("authorization", restClient.Session.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := restClient.client.Do(req)
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

	err = json.Unmarshal(body, &guilds)
	return
}

func ListChannels(restClient *RestClient, guildId int) (channels []Channel, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", restClient.Url, fmt.Sprintf("guilds/%d/channels", guildId)), bytes.NewBuffer([]byte(fmt.Sprintf(``))))
	if err != nil {
		return
	}
	req.Header.Set("authorization", restClient.Session.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := restClient.client.Do(req)
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

func SendMessage(restClient *RestClient, channelId int, message string) (responseMessage Message, err error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", restClient.Url, fmt.Sprintf("channels/%d/messages", channelId)), bytes.NewBuffer([]byte(fmt.Sprintf(`{"content":"%s"}`, message))))
	if err != nil {
		return
	}
	req.Header.Set("authorization", restClient.Session.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := restClient.client.Do(req)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	resp.Body.Close()

	fmt.Println(string(body))
	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("StatusCode: %d, %s", resp.StatusCode, string(body)))
		return
	}

	err = json.Unmarshal(body, &responseMessage)
	return
}

/*
func Track(restClient *RestClient) (e error) {
	req, e := http.NewRequest("POST", fmt.Sprintf("%s/%s", restClient.Url, fmt.Sprintf("track")), bytes.NewBuffer([]byte(fmt.Sprintf(`{"event":"Launch Game", "properties": {"Game":"Heroes of the Storm"}}`))))
	if e != nil {
		return
	}
	req.Header.Set("authorization", restClient.Session.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, e := restClient.client.Do(req)
	if e != nil {
		return
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		e = errors.New(fmt.Sprintf("StatusCode: %d, %s", resp.StatusCode, string(body)))
		return
	}

	//e = json.Unmarshal(body, &responseMessage)

	fmt.Println(string(body))
	return
}*/
