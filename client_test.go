package discord

import (
	"fmt"
	"github.com/xackery/discord"
	"testing"
)

func TestClient(t *testing.T) {
	restClient, err := discord.Create("test@hotmail.com", "test")
	if err != nil {
		t.Error(err)
		return
	}

	guilds, err := discord.ListGuilds(restClient)
	if err != nil {
		t.Error(err)
		return
	}

	if len(guilds) < 1 {
		t.Error("No guilds returned..")
	}

	guildId := guilds[0].Id

	channels, err := discord.ListChannels(restClient, guildId)
	channelId := 0

	for _, channel := range channels {
		if channel.Name == "test" {
			channelId = channel.Id
			break
		}
	}

	if channelId == 0 {
		t.Error("Failed to get a channel")
		return
	}

	//Send a message
	_, err = discord.SendMessage(restClient, channelId, "Test")
	if err != nil {
		t.Error(err)
		return
	}

	err = discord.Close(restClient)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("Done", guildId, channelId)
	return
}
