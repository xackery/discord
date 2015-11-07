package discord

import (
	"errors"
	"github.com/gorilla/websocket"
)

//Client wraps all Discord methods
type Client struct {
	User
	Token                     string
	GatewayURL                string `json:"url"`
	wsConn                    *websocket.Conn
	Guilds                    []Guild
	PrivateChannel            []PrivateChannel
	OnReady                   func(Event, Ready)
	OnTypingStart             func(Event, TypingEvent)
	OnMessageCreate           func(Event, Message)
	OnMessageUpdate           func(Event, Message)
	OnMessageDelete           func(Event, Message)
	OnPresenceStart           func(Event, Presence)
	OnPresenceUpdate          func(Event, Presence)
	OnUserSettingsUpdate      func(Event, UserSettings)
	OnGuildCreate             func(Event, GuildCreateEvent)
	OnGuildUpdate             func(Event, Guild)
	OnGuildDelete             func(Event, GuildDeleteEvent)
	OnGuildIntegrationsUpdate func(Event, GuildMemberEvent)
	OnGuildMemberAdd          func(Event, GuildMemberEvent)
	OnGuildMemberUpdate       func(Event, GuildMemberEvent)
	OnGuildMemberRemove       func(Event, GuildMemberEvent)
	OnGuildRoleDelete         func(Event, GuildRoleDeleteEvent)
	OnGuildRoleUpdate         func(Event, Guild)
	OnGuildRoleCreate         func(Event, GuildRoleEvent)
	OnVoiceStateUpdate        func(Event, VoiceState)
	isListening               bool
}

//Login Method for Discord
func (c *Client) Login(email string, pass string) (err error) {

	//Basic validation
	if len(email) < 3 {
		err = errors.New("email too short")
		return
	}
	if len(pass) < 3 {
		err = errors.New("password too short")
		return
	}
	//If user is not logged in, authorize and fetch a token
	if !c.IsLoggedIn() {
		err = c.authLogin(email, pass)
		if err != nil {
			return
		}
	}
	//sets User data on client
	c.UserMe()
	return
}

//IsLoggedIn returns if client is logged in to discord
func (c *Client) IsLoggedIn() (isLoggedIn bool) {
	return c.Token != ""
}
