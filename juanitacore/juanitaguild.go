package juanitacore

type JuanitaGuild struct {
	Id     string
	Player JuanitaPlayer
}

func NewJuanitaGuild(id string) *JuanitaGuild {
	guild := new(JuanitaGuild)
	guild.Id = id
	return guild
}

func (guild JuanitaGuild) GetPlayer() JuanitaPlayer {
	return guild.Player
}
