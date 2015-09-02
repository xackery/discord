package main

import (
	"fmt"
	"github.com/xackery/discord"
)

func main() {
	//Create a new httpClient
	restClient, e := discord.Create("email", "password")
	if e != nil {
		panic(e.Error())
	}

	//Get list of guilds
	guilds, e := discord.ListGuilds(restClient)
	if e != nil {
		panic(e.Error())
	}

	guildId := "0"
	if len(guilds) < 1 {
		panic("No guilds returned..!")
	}

	//Just snag initial guild
	guildId = guilds[0].Id

	//Get list of channels
	channels, e := discord.ListChannels(restClient, guildId)
	channelId := "0"
	for _, channel := range channels {
		if channel.Name == "lineage2" {
			fmt.Println("Found lineage 2 channel")
			channelId = channel.Id
		}
	}

	//I'm looking for a specific channel, not found? error here.
	if channelId == "0" {
		panic("Fail!?")
	}

	//Send a message
	message, e := discord.SendMessage(restClient, channelId, "Capturing the struct of message responses..")

	//output example stuff.
	fmt.Println(message)
	fmt.Println("Done", guildId, channelId)
}
