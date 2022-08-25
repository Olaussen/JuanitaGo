package juanitacore

import (
	"juanitaGo/structs"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

type JuanitaPlayer struct {
	AudioPlayer     string
	VoiceConnection *discordgo.VoiceConnection
	CurrentTrack    string
	isLoop          bool
	isLocked        bool
	isRepeat        bool
	volume          uint
	Queue           []structs.JuanitaSearch
}

func NewJuanitaPlayer() JuanitaPlayer {
	return JuanitaPlayer{AudioPlayer: "", VoiceConnection: nil, CurrentTrack: "", isLoop: false, isLocked: false, isRepeat: false, volume: 100, Queue: make([]structs.JuanitaSearch, 0)}
}

func NewJuanitaPlayerWithQueue(audioPlayer string, voiceConnection *discordgo.VoiceConnection, currentTrack string, isLoop bool, isLocked bool, isRepeat bool, volume uint, queue []structs.JuanitaSearch) JuanitaPlayer {
	return JuanitaPlayer{AudioPlayer: audioPlayer, VoiceConnection: voiceConnection, CurrentTrack: currentTrack, isLoop: isLoop, isLocked: isLocked, isRepeat: isRepeat, volume: volume, Queue: queue}
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

func (player *JuanitaPlayer) Shuffle() {
	for i := range player.Queue {
		j := rand.Intn(i + 1)
		player.Queue[i], player.Queue[j] = player.Queue[j], player.Queue[i]
	}
}

func (player *JuanitaPlayer) JoinVoiceChannel(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	voiceChannel, err := session.UserChannelCreate(interaction.Member.User.ID)
	if err != nil {
		return err
	}

	voiceConnection, err := session.ChannelVoiceJoin(interaction.GuildID, voiceChannel.ID, false, true)
	if err != nil {
		return err
	}

	player.VoiceConnection = voiceConnection
	return nil
}

func (player *JuanitaPlayer) AddSongToBackOfQueue(session *discordgo.Session, interaction *discordgo.InteractionCreate, search structs.JuanitaSearch) error {
	if player.VoiceConnection == nil {
		err := player.JoinVoiceChannel(session, interaction)
		if err != nil {
			return err
		}
	}

	player.Queue = append(player.Queue, search)
	return nil
}

func (player *JuanitaPlayer) AddSongFirstToQueue(session *discordgo.Session, interaction *discordgo.InteractionCreate, search structs.JuanitaSearch) error {
	if player.VoiceConnection == nil {
		err := player.JoinVoiceChannel(session, interaction)
		if err != nil {
			return err
		}
	}

	player.Queue = append([]structs.JuanitaSearch{search}, player.Queue...)
	return nil
}
