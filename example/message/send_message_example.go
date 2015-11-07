package message

import (
	"log"
)

//This is the equivalent of main()
func SendMessageExample(email string, password string) {
	//Log in to discord
	err := client.Login(email, password)
	if err != nil {
		log.Println("Error Logging in:", err.Error())
		return
	}

	//Get a list of guilds (servers)
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

	//Get a list of channels (rooms)
	channels, err := client.GuildChannels(guild.ID)
	if err != nil {
		log.Println("Error getting channels:", err.Error())
		return
	}
	log.Println("Channel size:", len(channels))

	//Find the channel "test"
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
}
