package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//Listen to websocket connection
func (c *Client) Listen() (err error) {
	if !c.IsLoggedIn() {
		err = errors.New("You must be logged in")
		return
	}

	if c.isListening {
		err = errors.New("A websocket already exists")
		return
	}

	//Request API Getway URL
	httpClient := &http.Client{Timeout: (20 * time.Second)}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", DISCORD_URL, "gateway"), bytes.NewBuffer([]byte(fmt.Sprintf(``))))
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
	if err != nil {
		return
	}

	//Dial websocket and initialize connection
	//log.Println("dialing", c.GatewayURL)
	conn, _, err := websocket.DefaultDialer.Dial(c.GatewayURL, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	c.wsConn = conn
	//Sends token data and ensures websocket validates connection
	err = c.initWebsocket()
	if err != nil {
		return
	}
	c.isListening = true
	log.Println("Listening on", c.GatewayURL)

	//Message pool loops until isListening is set to false
	for c.isListening {
		var msgType int
		var msgData []byte
		msgType, msgData, err = c.wsConn.ReadMessage()
		if err != nil {
			return
		}
		go c.handleEvent(msgType, msgData)
	}
	return
}

//Stop listening for websocket data
func (c *Client) StopListen() {
	c.isListening = false
}

//Initialize the websocket
func (c *Client) initWebsocket() (err error) {
	err = c.wsConn.WriteJSON(bson.M{
		"op": 2,
		"d": bson.M{
			"token": c.Token,
			"properties": bson.M{
				"$os":               "linux",
				"$browser":          "discord",
				"$device":           "discord",
				"$referer":          "",
				"$referring_domain": "",
			},
			"v": 3},
	})
	if err != nil {
		return
	}
	return
}
