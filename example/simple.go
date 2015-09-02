package main

import (
	"fmt"
	"github.com/xackery/discord"
)

func main() {
	c, e := discord.Create("email", "password")
	fmt.Println(c.Session.Id)
	if e != nil {
		panic(e.Error())
	}
	guilds, e := discord.ListGuilds(c)
	/*if e != nil {
		panic(e.Error())
	}*/
	guildId := "0"
	if len(guilds) == 1 {
		guildId = guilds[0].Id
	}

	channels, e := discord.ListChannels(c, guildId)
	channelId := "0"
	for _, channel := range channels {
		if channel.Name == "lineage2" {
			fmt.Println("Found lineage 2 channel")
			channelId = channel.Id
		}
	}
	if channelId == "0" {
		fmt.Println("Fail!")
		return
	}
	message, e := discord.SendMessage(c, channelId, "Capturing the struct of message responses..")
	fmt.Println(message)

	fmt.Println("Done", guildId, channelId)
}
