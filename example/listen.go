package main

import (
	"github.com/xackery/discord"
	"log"
)

var client discord.Client

func OnMessageCreate(event discord.Event, message discord.Message) {
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

//Listens
func main() {

	err := client.Login("email", "password")
	if err != nil {
		log.Println("Error Logging in:", err.Error())
		return
	}
	client.OnMessageCreate = OnMessageCreate
	err = client.Listen()
	if err != nil {
		log.Println("Error Listening:", err.Error())
	}
}
