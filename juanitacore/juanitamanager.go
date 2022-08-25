package juanitacore

type JuanitaManager struct {
	Guilds map[string]JuanitaGuild
}

func NewJuanitaManager() JuanitaManager {
	return JuanitaManager{Guilds: make(map[string]JuanitaGuild)}
}

func NewJuanitaManagerWithGuilds(guilds map[string]JuanitaGuild) JuanitaManager {
	return JuanitaManager{Guilds: guilds}
}

func (manager *JuanitaManager) GetOrAddGuild(guildID string) *JuanitaGuild {
	guild, ok := manager.Guilds[guildID]
	if !ok {
		guild = NewJuanitaGuild(guildID)
		manager.Guilds[guildID] = guild
	}
	return &guild
}

func (manager *JuanitaManager) GetGuild(guildID string) *JuanitaGuild {
	guild, ok := manager.Guilds[guildID]
	if !ok {
		return nil
	}
	return &guild
}

func (manager *JuanitaManager) AddGuild(guild JuanitaGuild) {
	manager.Guilds[guild.Id] = guild
}
