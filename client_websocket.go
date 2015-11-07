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
	Type      string          `json:"t"`
	State     int             `json:"s"`
	Operation int             `json:"o"`
	RawData   json.RawMessage `json:"d"`
}

type Ready struct {
	Version           int `json:"v"`
	User              User
	SessionID         string
	ReadState         []ReadState
	PrivateChannels   []PrivateChannel
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
	Guilds            []Guild
}

type ReadState struct {
	MentionCount  int
	LastMessageID int `json:"string"`
	ID            int `json:"string"`
}

func (c *Client) handleEvent(msgType int, msgData []byte) {
	var event Event
	err := json.Unmarshal(msgData, &event)
	if err != nil {
		log.Println("Err HandleEvent:", err.Error())
		return
	}

	//eventType := event.(map[string]interface{})["t"].(string)
	switch event.Type {
	case "READY":
		var ready Ready
		err := json.Unmarshal(event.RawData, &ready)
		if err != nil {
			log.Println("Error Ready Parse:", err.Error())
			return
		}
		log.Println("Size of guilds:", len(ready.Guilds), "members:", len(ready.Guilds[0].Members))
		c.handleReady(event, ready)
	case "TYPING_START":
		var typingEvent TypingEvent
		err := json.Unmarshal(event.RawData, &typingEvent)
		if err != nil {
			log.Println("Error Typing Parse:", err.Error())
			return
		}
		if c.OnTypingStart != nil {
			c.OnTypingStart(event, typingEvent)
		}
	case "PRESENCE_START":
		var presence Presence
		err := json.Unmarshal(event.RawData, &presence)
		if err != nil {
			log.Println("Error Presence Parse:", err.Error())
			return
		}
		if c.OnPresenceStart != nil {
			c.OnPresenceStart(event, presence)
		}
	case "PRESENCE_UPDATE":
		var presence Presence
		err := json.Unmarshal(event.RawData, &presence)
		if err != nil {
			log.Println("Error Presence Parse:", err.Error())
			return
		}
		if c.OnPresenceStart != nil {
			c.OnPresenceUpdate(event, presence)
		}
	case "MESSAGE_CREATE":
		var message Message
		err := json.Unmarshal(event.RawData, &message)
		if err != nil {
			log.Println("Error Message Parse:", err.Error())
			return
		}
		if c.OnMessageCreate != nil {
			c.OnMessageCreate(event, message)
		}
	default:
		log.Printf("Ignoring %s", event.Type)
		log.Printf("event dump: %d %s", msgType, string(msgData[:]))
	}
}

func (c *Client) handleReady(event Event, ready Ready) {

	// WebSocket keepalive
	go func() {
		ticker := time.NewTicker(ready.HeartbeatInterval * time.Millisecond)
		for range ticker.C {
			timestamp := int(time.Now().Unix())
			log.Printf("Sending keepalive with timestamp %d", timestamp)
			c.wsConn.WriteJSON(map[string]int{
				"op": 1,
				"d":  timestamp,
			})
		}
	}()

	if c.OnReady != nil {
		c.OnReady(event, ready)
	}
}

type messageEvent struct {
	OpCode int     `json:"op"`
	Type   string  `json:"t"`
	Data   Message `json:"d"`
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
