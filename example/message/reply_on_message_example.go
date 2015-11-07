package message

import (
	"github.com/xackery/discord"
	"log"
)

var client discord.Client

//Whenever a message creation event triggers, this function is notified
func OnMessageCreateExample(event discord.Event, message discord.Message) {
	if message.Author.ID == client.ID {
		log.Println("Ignoring myself")
		return
	}
	log.Println("Got message!", message.ChannelID, client.Token, message.Author.Username, message.Content)

	//Reply to the event with a message
	message, err := client.ChannelMessageSend(message.ChannelID, "Hi "+message.Author.Username+", I see you.")
	if err != nil {
		log.Println("Err message sending:", err.Error())
	}
}

//This is the equivalent of main()
func ReplyOnMessageExample(email string, password string) {
	//Log in to discord
	err := client.Login(email, password)
	if err != nil {
		log.Println("Error Logging in:", err.Error())
		return
	}
	//Register an event handler
	client.OnMessageCreate = OnMessageCreateExample

	//Listen for events.
	//This code is blocked and will never end until client.StopListen() is called.
	err = client.Listen()
	if err != nil {
		log.Println("Error Listening:", err.Error())
	}
}
