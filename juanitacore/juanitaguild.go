package juanitacore

type JuanitaGuild struct {
	Id     string
	Player JuanitaPlayer
}

func NewJuanitaGuild(id string) JuanitaGuild {
	return JuanitaGuild{Id: id, Player: NewJuanitaPlayer()}
}

func NewJuanitaGuildWithQueue(id string, player JuanitaPlayer) JuanitaGuild {
	return JuanitaGuild{Id: id, Player: player}
}

func (guild *JuanitaGuild) GetPlayer() *JuanitaPlayer {
	return &guild.Player
}
