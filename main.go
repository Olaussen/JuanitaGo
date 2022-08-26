package main

import (
	"fmt"
	cmd "juanitaGo/commands"
	core "juanitaGo/juanitacore"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	conf       *core.JuanitaConfig
	CmdHandler *core.JuanitaCommandHandler
	Sessions   *core.JuanitaSessionManager
	youtube    *core.Youtube
	botId      string
	PREFIX     string
)

func init() {
	conf = core.LoadConfig("config.json")
	PREFIX = conf.Prefix

}

func main() {
	CmdHandler = core.NewCommandHandler()
	registerCommands()
	Sessions = core.NewSessionManager()
	youtube = &core.Youtube{Conf: conf}
	discord, err := discordgo.New("Bot " + conf.BotToken)
	if err != nil {
		fmt.Println("Error creating discord session,", err)
		return
	}
	if conf.UseSharding {
		discord.ShardID = conf.ShardId
		discord.ShardCount = conf.ShardCount
	}

	discord.AddHandler(commandHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		discord.UpdateGameStatus(0, conf.DefaultStatus)
		guilds := discord.State.Guilds
		fmt.Println("Ready with", len(guilds), "guilds.")
	})
	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}
	fmt.Println("Started")
	<-make(chan struct{})
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == botId || user.Bot {
		return
	}
	content := message.Content
	if len(content) <= len(PREFIX) {
		return
	}
	if content[:len(PREFIX)] != PREFIX {
		return
	}
	content = content[len(PREFIX):]
	if len(content) < 1 {
		return
	}
	args := strings.Fields(content)
	name := strings.ToLower(args[0])
	command, found := CmdHandler.Get(name)
	if !found {
		return
	}
	channel, err := discord.State.Channel(message.ChannelID)
	if err != nil {
		fmt.Println("Error getting channel,", err)
		return
	}
	guild, err := discord.State.Guild(channel.GuildID)
	if err != nil {
		fmt.Println("Error getting guild,", err)
		return
	}
	ctx := core.NewContext(discord, guild, channel, user, message, conf, CmdHandler, Sessions, youtube)
	ctx.Args = args[1:]
	c := *command
	c(*ctx)
}

func registerCommands() {
	CmdHandler.Register("play", cmd.PlayCommand, "Plays whats in the queue")
	CmdHandler.Register("join", cmd.JoinCommand, "Join a voice channel !join attic")
	CmdHandler.Register("add", cmd.AddCommand, "Add a song to the queue !add <youtube-link>")
	// ??? means I haven't dug in
	// TODO: Consistant order?
	/*CmdHandler.Register("help", cmd.HelpCommand, "Gives you this help message!")
	CmdHandler.Register("admin", cmd.AdminCommand, "???")
	CmdHandler.Register("leave", cmd.LeaveCommand, "Leaves current voice channel")*/
	/*CmdHandler.Register("stop", cmd.StopCommand, "Stops the music")
	CmdHandler.Register("info", cmd.InfoCommand, "???")
	CmdHandler.Register("skip", cmd.SkipCommand, "Skip")
	CmdHandler.Register("queue", cmd.QueueCommand, "Print queue???")
	CmdHandler.Register("eval", cmd.EvalCommand, "???")
	CmdHandler.Register("debug", cmd.DebugCommand, "???")
	CmdHandler.Register("clear", cmd.ClearCommand, "empty queue???")
	CmdHandler.Register("current", cmd.CurrentCommand, "Name current song???")
	CmdHandler.Register("youtube", cmd.YoutubeCommand, "???")
	CmdHandler.Register("shuffle", cmd.ShuffleCommand, "Shuffle queue???")
	CmdHandler.Register("pausequeue", cmd.PauseCommand, "Pause song in place???")
	CmdHandler.Register("pick", cmd.PickCommand, "???")*/
}
