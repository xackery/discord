package discord

type Presence struct {
	User    User
	Status  string
	Roles   []string
	GuildID int `json:"guild_id,string"`
	GameID  int `json:"game_id"`
}
