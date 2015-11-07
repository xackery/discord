# discord
A Discord API for Golang

[![godoc badge](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/Xackery/discord/discord) [![Go report](http://goreportcard.com/badge/xackery/discord)](http://goreportcard.com/report/xackery/discord) [![Build Status](https://travis-ci.org/Xackery/discord.svg)](https://travis-ci.org/Xackery/discord.svg?branch=master)
Features
---

* Log In/Out
* List Guilds
* List Channels
* Send Message to Guild/Channel
* Connect via Websocket, Trigger functions on Events

Example
---


```
	client := discord.Client{}
	err := client.Login(email, password)
	if err != nil {
		log.Println("Error Logging in:", err.Error())
		return
	}

	guilds, err := client.UserGuilds()
	if err != nil {
		log.Println("Error getting guilds:", err.Error())
	}
	log.Println("Guild size:", len(guilds))
	if len(guilds) < 1 {
		log.Println("no guilds, exiting")
		return
	}

	guild := guilds[0]
	channels, err := client.GuildChannels(guild.ID)
	if err != nil {
		log.Println("Error getting channels:", err.Error())
		return
	}
	log.Println("Channel size:", len(channels))
	for _, channel := range channels {
		if channel.Name == "test" {
			log.Println("Sending message to ", channel.ID)
			resp, err := client.ChannelMessageSend(channel.ID, "So this method works...")
			if err != nil {
				log.Println("Error messaging:", err.Error())
				return
			}
			log.Println("Done", resp)
			return
		}
	}
```