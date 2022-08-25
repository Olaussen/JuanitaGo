package juanitacore

type JuanitaManager struct {
	Guilds map[string]JuanitaGuild
}

func NewGuildManager() JuanitaManager {
	return JuanitaManager{Guilds: make(map[string]JuanitaGuild)}
}

func NewGuildManagerWithGuilds(guilds map[string]JuanitaGuild) JuanitaManager {
	return JuanitaManager{Guilds: guilds}
}
