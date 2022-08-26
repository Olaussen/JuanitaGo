package juanitacore

import (
	"fmt"
	"juanitaGo/structs"

	"github.com/bwmarrin/discordgo"
)

type JuanitaPlayer struct {
	VoiceConnection *discordgo.VoiceConnection
	CurrentTrack    *structs.JuanitaSearch
	isLoop          bool
	isLocked        bool
	isRepeat        bool
	volume          uint
}

func NewJuanitaPlayer() *JuanitaPlayer {
	player := new(JuanitaPlayer)
	player.isLoop = false
	player.isLocked = false
	player.isRepeat = false
	player.volume = 100
	return player
}

func NewJuanitaPlayerWithQueue(audioPlayer string, voiceConnection *discordgo.VoiceConnection, currentTrack string, isLoop bool, isLocked bool, isRepeat bool, volume uint, queue []structs.JuanitaSearch) JuanitaPlayer {
	return JuanitaPlayer{VoiceConnection: voiceConnection, CurrentTrack: nil, isLoop: isLoop, isLocked: isLocked, isRepeat: isRepeat, volume: volume}
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

func (player *JuanitaPlayer) JoinVoiceChannel(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	voiceChannel, err := session.UserChannelCreate(interaction.Member.User.ID)
	if err != nil {
		return err
	}
	voiceConnection, err := session.ChannelVoiceJoin(interaction.GuildID, voiceChannel.ID, false, true)
	fmt.Printf("%v", voiceConnection)
	if err != nil {
		return err
	}

	player.VoiceConnection = voiceConnection
	return nil
}

func (player *JuanitaPlayer) LeaveVoiceChannel(session *discordgo.Session) error {
	err := player.VoiceConnection.Disconnect()
	if err != nil {
		return err
	}
	return nil
}

/*func (player *JuanitaPlayer) AddSongToBackOfQueue(search structs.JuanitaSearch) {
	player.Queue.EnqueueBack(search)
}

func (player *JuanitaPlayer) AddSongFirstToQueue(session *discordgo.Session, interaction *discordgo.InteractionCreate, search structs.JuanitaSearch) {
	player.Queue.EnqueueFirst(search)
}

// play song in discord channel if queue is not empty
func (player *JuanitaPlayer) Play(session *discordgo.Session, interaction *discordgo.InteractionCreate, song structs.JuanitaSearch) error {
	if player.VoiceConnection == nil {
		err := player.JoinVoiceChannel(session, interaction)
		if err != nil {
			return err
		}
	}
	player.AddSongToBackOfQueue(song)
	currentTrack := player.Queue.Dequeue()
	player.CurrentTrack = currentTrack
	audioStream := currentTrack.Song.AudioStream()
	dgvoice.PlayAudioFile(player.VoiceConnection, audioStream, make(<-chan bool))
	return nil
}

func (player *JuanitaPlayer) Echo(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	err := player.JoinVoiceChannel(session, interaction)
	if err != nil {
		return err
	}

	recv := make(chan *discordgo.Packet, 2)
	go dgvoice.ReceivePCM(player.VoiceConnection, recv)

	send := make(chan []int16, 2)
	go dgvoice.SendPCM(player.VoiceConnection, send)

	player.VoiceConnection.Speaking(true)
	defer player.VoiceConnection.Speaking(false)

	for {

		p, ok := <-recv
		if !ok {
			return nil
		}

		send <- p.PCM
	}

}*/
