package juanitacore

import (
	"juanitaGo/structs"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

type JuanitaPlayer struct {
	VoiceConnection *discordgo.VoiceConnection
	CurrentTrack    *structs.JuanitaSearch
	isLoop          bool
	isLocked        bool
	isRepeat        bool
	volume          uint
	Queue           JuanitaQueue
}

func NewJuanitaPlayer() JuanitaPlayer {
	return JuanitaPlayer{VoiceConnection: nil, CurrentTrack: nil, isLoop: false, isLocked: false, isRepeat: false, volume: 100, Queue: NewJuanitaQueue()}
}

func NewJuanitaPlayerWithQueue(audioPlayer string, voiceConnection *discordgo.VoiceConnection, currentTrack string, isLoop bool, isLocked bool, isRepeat bool, volume uint, queue []structs.JuanitaSearch) JuanitaPlayer {
	return JuanitaPlayer{VoiceConnection: voiceConnection, CurrentTrack: nil, isLoop: isLoop, isLocked: isLocked, isRepeat: isRepeat, volume: volume, Queue: NewJuanitaQueueWithTracks(queue)}
}

func (player *JuanitaPlayer) ToggleLoop() {
	player.isLoop = !player.isLoop
}

func (player *JuanitaPlayer) ToggleLock() {
	player.isLocked = !player.isLocked
}

func (player *JuanitaPlayer) ToggleRepeat() {
	player.isRepeat = !player.isRepeat
}

func (player *JuanitaPlayer) SetVolume(volume uint) {
	player.volume = volume
}

func (player *JuanitaPlayer) JoinVoiceChannel(session *discordgo.Session, interaction *discordgo.InteractionCreate) bool {
	voiceChannel, err := session.UserChannelCreate(interaction.Member.User.ID)
	if err != nil {
		return false
	}

	voiceConnection, err := session.ChannelVoiceJoin(interaction.GuildID, voiceChannel.ID, false, true)
	if err != nil {
		return false
	}

	player.VoiceConnection = voiceConnection
	return true
}

func (player *JuanitaPlayer) LeaveVoiceChannel(session *discordgo.Session) bool {
	err := player.VoiceConnection.Disconnect()
	if err != nil {
		return false
	}
	return true
}

func (player *JuanitaPlayer) AddSongToBackOfQueue(search structs.JuanitaSearch) {
	player.Queue.EnqueueBack(search)
}

func (player *JuanitaPlayer) AddSongFirstToQueue(session *discordgo.Session, interaction *discordgo.InteractionCreate, search structs.JuanitaSearch) {
	player.Queue.EnqueueFirst(search)
}

// play song in discord channel if queue is not empty
func (player *JuanitaPlayer) Play(session *discordgo.Session, interaction *discordgo.InteractionCreate, song structs.JuanitaSearch) error {
	if player.Queue.IsEmpty() {
		return nil
	}
	if player.CurrentTrack != nil {
		// add song to queue
		player.AddSongToBackOfQueue(song)
		return nil
	}
	if player.VoiceConnection == nil {
		ok := player.JoinVoiceChannel(session, interaction)
		if ok {
			currentTrack := player.Queue.Dequeue()
			audioStream := currentTrack.Song.AudioStream()
			dgvoice.PlayAudioFile(player.VoiceConnection, audioStream, make(<-chan bool))
		}
	}
	return nil
}

func (player *JuanitaPlayer) Echo(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	ok := player.JoinVoiceChannel(session, interaction)
	if ok {
		recv := make(chan *discordgo.Packet, 2)
		go dgvoice.ReceivePCM(player.VoiceConnection, recv)

		send := make(chan []int16, 2)
		go dgvoice.SendPCM(player.VoiceConnection, send)

		player.VoiceConnection.Speaking(true)
		defer player.VoiceConnection.Speaking(false)

		for {

			p, ok := <-recv
			if !ok {
				return
			}

			send <- p.PCM
		}
	}
}
