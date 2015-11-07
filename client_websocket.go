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

//Majority of this code is heavy credits to github.com/gdraynz/go-discord/

//Listen to websocket connection
func (c *Client) Listen() (err error) {
	if !c.IsLoggedIn() {
		err = errors.New("You must be logged in")
		return
	}

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
	log.Println("dialing", c.GatewayURL)
	conn, _, err := websocket.DefaultDialer.Dial(c.GatewayURL, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	c.wsConn = conn

	err = c.initWebsocket()
	if err != nil {
		return
	}
	log.Println("Listening for Messages")
	for {
		var msgType int
		var msgData []byte
		msgType, msgData, err = c.wsConn.ReadMessage()
		if err != nil {
			return
		}
		go c.handleEvent(msgType, msgData)
	}
}

type Event struct {
}

func (c *Client) handleEvent(msgType int, msgData []byte) {
	var event interface{}
	err := bson.Unmarshal(msgData, &event)
	if err != nil {
		return
	}

	eventType := event.(map[string]interface{})["t"].(string)

	switch eventType {
	//case "READY":
	//	c.handleReady(msgData)
	case "MESSAGE_CREATE":
		c.handleMessageCreate(msgData)
	/*case "MESSAGE_ACK":
		c.handleMessageAck(msgData)
	case "MESSAGE_UPDATE":
		c.handleMessageUpdate(msgData)
	case "MESSAGE_DELETE":
		c.handleMessageDelete(msgData)
	case "TYPING_START":
		c.handleTypingStart(msgData)
	case "PRESENCE_UPDATE":
		c.handlePresenceUpdate(msgData)
	case "CHANNEL_CREATE":
		c.handleChannelCreate(msgData)
	case "CHANNEL_UPDATE":
		c.handleChannelUpdate(msgData)
	case "CHANNEL_DELETE":
		c.handleChannelDelete(msgData)
	case "GUILD_CREATE":
		c.handleGuildCreate(msgData)
	case "GUILD_DELETE":
		c.handleGuildDelete(msgData)
	*/
	default:

		log.Printf("Ignoring %s", eventType)
		log.Printf("event dump: %d %s", msgType, string(msgData[:]))
	}
}

// Ready is received when the websocket connection is made and helps set up everything
type Ready struct {
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
	Guilds            []Guild
	PrivateChannels   []PrivateChannel
}

type readyEvent struct {
	OpCode int    `json:"op"`
	Type   string `json:"t"`
	Data   Ready  `json:"d"`
}

func (c *Client) handleReady(eventStr []byte) {

	log.Println("handleReady")
	var ready readyEvent
	if err := json.Unmarshal(eventStr, &ready); err != nil {
		log.Printf("handleReady: %s", err)
		return
	}

	// WebSocket keepalive
	go func() {
		ticker := time.NewTicker(ready.Data.HeartbeatInterval * time.Millisecond)
		for range ticker.C {
			timestamp := int(time.Now().Unix())
			log.Printf("Sending keepalive with timestamp %d", timestamp)
			c.wsConn.WriteJSON(map[string]int{
				"op": 1,
				"d":  timestamp,
			})
		}
	}()

	//c.User = ready.Data.User
	//c.initServers(ready.Data)

	/*if c.OnReady == nil {
		log.Print("No handler for READY")
	} else {
		log.Print("Client ready, calling OnReady handler")
		c.OnReady(ready.Data)
	}*/
}

type messageEvent struct {
	OpCode int     `json:"op"`
	Type   string  `json:"t"`
	Data   Message `json:"d"`
}

func (c *Client) handleMessageCreate(msgData []byte) {
	if c.OnMessageCreate == nil {
		log.Print("No handler for MESSAGE_CREATE")
		return
	}

	var message messageEvent
	if err := json.Unmarshal(msgData, &message); err != nil {
		log.Printf("messageCreate: %s", err)
		return
	}

	if message.Data.Author.ID != c.ID {
		c.OnMessageCreate(message.Data)
	} else {
		log.Print("Ignoring message from self")
	}
}

//Initialize the websocket
func (c *Client) initWebsocket() (err error) {
	log.Println("Initialziing websocket")
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
