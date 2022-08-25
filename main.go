package main

import (
	"flag"
	"juanitaGo/juanitacore"
	"juanitaGo/messages"
	"juanitaGo/utils"
	"juanitaGo/youtube"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	GuildID        = flag.String("guild", "", "Leave blank to register global commands")
	BotToken       = flag.String("token", utils.GetEnvironmentVariableByKey("BOT_TOKEN"), "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var session *discordgo.Session
var guildManager = juanitacore.NewJuanitaManager()
var youtubeSearcher = youtube.NewYoutubeSearcher()

func init() { flag.Parse() }

func init() {
	var err error
	session, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

const integerOptionMinValue = 1.0

var commands = GetCommandConfig()

var commandHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate){
	"echo": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		guild := guildManager.GetOrAddGuild(interaction.Interaction.GuildID)
		guild.Player.Echo(session, interaction)
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hey there! Congratulations, you just executed your first slash command",
			},
		})
	},
	"test": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		options := utils.ExtractOptions(interaction)
		messageArguments := make([]interface{}, 0, len(options))
		if option, ok := options["sangnavn"]; ok {
			messageArguments = append(messageArguments, option.StringValue())
		}
		user := interaction.Interaction.Member.User
		searchString := messageArguments[0].(string)
		searchResult := youtubeSearcher.Search(searchString, user)
		guild := guildManager.GetOrAddGuild(interaction.Interaction.GuildID)
		guild.Player.Play(session, interaction, *searchResult)

		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: messages.PlayEmbed(*searchResult),
			},
		})
	},
}

func init() {
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	session.AddHandler(func(session *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", session.State.User.Username, session.State.User.Discriminator)
	})
	err := session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")

		for _, v := range registeredCommands {
			err := session.ApplicationCommandDelete(session.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Juanita shutting down.")
}
