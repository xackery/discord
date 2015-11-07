package message

import (
	"github.com/xackery/discord"
	"log"
)

var client discord.Client

func OnMessageCreateExample(event discord.Event, message discord.Message) {
	if message.Author.ID == client.ID {
		log.Println("Ignoring myself")
		return
	}
	log.Println("Got message!", message.ChannelID, client.Token, message.Author.Username, message.Content)

	message, err := client.ChannelMessageSend(message.ChannelID, "Hi "+message.Author.Username+", I see you.")
	if err != nil {
		log.Println("Err message sending:", err.Error())
	}
}

func ReplyOnMessageExample(email string, password string) {
	err := client.Login(email, password)
	if err != nil {
		log.Println("Error Logging in:", err.Error())
		return
	}
	client.OnMessageCreate = OnMessageCreateExample
	err = client.Listen()
	if err != nil {
		log.Println("Error Listening:", err.Error())
	}
}
