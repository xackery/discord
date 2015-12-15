package discord

import (
	"encoding/json"
	"log"
	"time"
)

type Event struct {
	Type      string          `json:"t"`
	State     int             `json:"s"`
	Operation int             `json:"o"`
	Direction int             `json:"dir"` //Direction of command, 0-received, 1-sent
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

	switch event.Type {
	case "READY":
		var ready Ready
		err := json.Unmarshal(event.RawData, &ready)
		if err != nil {
			log.Println("Error Ready Parse:", err.Error())
			return
		}
		//log.Println("Size of guilds:", len(ready.Guilds), "members:", len(ready.Guilds[0].Members))
		c.handleReady(event, ready)
	case "GUILD_CREATE":
		var guildCreateEvent GuildCreateEvent
		err := json.Unmarshal(event.RawData, &guildCreateEvent)
		if err != nil {
			log.Println("Error Ready Parse:", err.Error())
			return
		}
		if c.OnGuildCreate != nil {
			c.OnGuildCreate(event, guildCreateEvent)
		}
	case "GUILD_DELETE":
		var guildDeleteEvent GuildDeleteEvent
		err := json.Unmarshal(event.RawData, &guildDeleteEvent)
		if err != nil {
			log.Println("Error Ready Parse:", err.Error())
			return
		}
		if c.OnGuildDelete != nil {
			c.OnGuildDelete(event, guildDeleteEvent)
		}
	case "GUILD_INTEGRATIONS_UPDATE":
		var guildMemberEvent GuildMemberEvent
		err := json.Unmarshal(event.RawData, &guildMemberEvent)
		if err != nil {
			log.Println("Error Ready Parse:", err.Error())
			return
		}
		if c.OnGuildIntegrationsUpdate != nil {
			c.OnGuildIntegrationsUpdate(event, guildMemberEvent)
		}
	case "GUILD_MEMBER_ADD":
		var guildMemberEvent GuildMemberEvent
		err := json.Unmarshal(event.RawData, &guildMemberEvent)
		if err != nil {
			log.Println("Error Ready Parse:", err.Error())
			return
		}
		if c.OnGuildMemberAdd != nil {
			c.OnGuildMemberAdd(event, guildMemberEvent)
		}
	case "GUILD_MEMBER_UPDATE":
		var guildMemberEvent GuildMemberEvent
		err := json.Unmarshal(event.RawData, &guildMemberEvent)
		if err != nil {
			log.Println("Error Ready Parse:", err.Error())
			return
		}
		if c.OnGuildMemberUpdate != nil {
			c.OnGuildMemberUpdate(event, guildMemberEvent)
		}
	case "GUILD_MEMBER_REMOVE":
		var guildMemberEvent GuildMemberEvent
		err := json.Unmarshal(event.RawData, &guildMemberEvent)
		if err != nil {
			log.Println("Error Ready Parse:", err.Error())
			return
		}
		if c.OnGuildMemberRemove != nil {
			c.OnGuildMemberRemove(event, guildMemberEvent)
		}
	case "GUILD_ROLE_CREATE":
		var guildRoleEvent GuildRoleEvent
		err := json.Unmarshal(event.RawData, &guildRoleEvent)
		if err != nil {
			log.Println("Error Ready Parse:", err.Error())
			return
		}
		if c.OnGuildRoleCreate != nil {
			c.OnGuildRoleCreate(event, guildRoleEvent)
		}
	case "GUILD_ROLE_DELETE":
		var guildRoleDeleteEvent GuildRoleDeleteEvent
		err := json.Unmarshal(event.RawData, &guildRoleDeleteEvent)
		if err != nil {
			log.Println("Error Ready Parse:", err.Error())
			return
		}
		if c.OnGuildRoleDelete != nil {
			c.OnGuildRoleDelete(event, guildRoleDeleteEvent)
		}
	case "GUILD_ROLE_UPDATE":
		var guild Guild
		err := json.Unmarshal(event.RawData, &guild)
		if err != nil {
			log.Println("Error Ready Parse:", err.Error())
			return
		}
		if c.OnGuildRoleUpdate != nil {
			c.OnGuildRoleUpdate(event, guild)
		}
	case "GUILD_UPDATE":
		var guild Guild
		err := json.Unmarshal(event.RawData, &guild)
		if err != nil {
			log.Println("Error Ready Parse:", err.Error())
			return
		}
		if c.OnGuildUpdate != nil {
			c.OnGuildUpdate(event, guild)
		}
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
		if c.OnPresenceUpdate != nil {
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
	case "MESSAGE_UPDATE":
		var message Message
		err := json.Unmarshal(event.RawData, &message)
		if err != nil {
			log.Println("Error Message Parse:", err.Error())
			return
		}
		if c.OnMessageUpdate != nil {
			c.OnMessageUpdate(event, message)
		}
	case "MESSAGE_DELETE":
		var message Message
		err := json.Unmarshal(event.RawData, &message)
		if err != nil {
			log.Println("Error Message Parse:", err.Error())
			return
		}
		if c.OnMessageDelete != nil {
			c.OnMessageDelete(event, message)
		}
	case "USER_SETTINGS_UPDATE":
		var userSettings UserSettings
		err := json.Unmarshal(event.RawData, &userSettings)
		if err != nil {
			log.Println("Error Message Parse:", err.Error())
			return
		}
		if c.OnUserSettingsUpdate != nil {
			c.OnUserSettingsUpdate(event, userSettings)
		}
	case "VOICE_STATE_UPDATE":
		var voiceState VoiceState
		err := json.Unmarshal(event.RawData, &voiceState)
		if err != nil {
			log.Println("Error Message Parse:", err.Error())
			return
		}
		if c.OnVoiceStateUpdate != nil {
			c.OnVoiceStateUpdate(event, voiceState)
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
