package juanitacore

import (
	"github.com/bwmarrin/discordgo"
)

type (
	JuanitaSession struct {
		Queue              *SongQueue
		guildId, ChannelId string
		connection         *JuanitaConnection
	}

	JuanitaSessionManager struct {
		sessions map[string]*JuanitaSession
	}

	JoinProperties struct {
		Muted    bool
		Deafened bool
	}
)

func newSession(guildId, channelId string, connection *JuanitaConnection) *JuanitaSession {
	session := new(JuanitaSession)
	session.Queue = newSongQueue()
	session.guildId = guildId
	session.ChannelId = channelId
	session.connection = connection
	return session
}

func (sess JuanitaSession) Play(song Song) error {
	return sess.connection.Play(song.Ffmpeg())
}

func (sess *JuanitaSession) Stop() {
	sess.connection.Stop()
}

func NewSessionManager() *JuanitaSessionManager {
	return &JuanitaSessionManager{make(map[string]*JuanitaSession)}
}

func (manager JuanitaSessionManager) GetByGuild(guildId string) *JuanitaSession {
	for _, sess := range manager.sessions {
		if sess.guildId == guildId {
			return sess
		}
	}
	return nil
}

func (manager JuanitaSessionManager) GetByChannel(channelId string) (*JuanitaSession, bool) {
	sess, found := manager.sessions[channelId]
	return sess, found
}

func (manager *JuanitaSessionManager) Join(discord *discordgo.Session, guildId, channelId string,
	properties JoinProperties) (*JuanitaSession, error) {
	vc, err := discord.ChannelVoiceJoin(guildId, channelId, properties.Muted, properties.Deafened)
	if err != nil {
		return nil, err
	}
	sess := newSession(guildId, channelId, NewJuanitaConnection(vc))
	manager.sessions[channelId] = sess
	return sess, nil
}

func (manager *JuanitaSessionManager) Leave(discord *discordgo.Session, session JuanitaSession) {
	session.connection.Stop()
	session.connection.Disconnect()
	delete(manager.sessions, session.ChannelId)
}
