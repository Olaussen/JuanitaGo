package juanitacore

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

type JuanitaConnection struct {
	VoiceConnection *discordgo.VoiceConnection
	Send            chan []int16
	Lock            sync.Mutex
	Sendpcm         bool
	StopRunning     bool
	Playing         bool
}

func NewJuanitaConnection(voiceConnection *discordgo.VoiceConnection) *JuanitaConnection {
	connection := new(JuanitaConnection)
	connection.VoiceConnection = voiceConnection
	return connection
}
func (connection JuanitaConnection) Disconnect() {
	connection.VoiceConnection.Disconnect()
}
