package main

import (
	"flag"
	"juanitaGo/juanitacore"
	"juanitaGo/utils"
	"juanitaGo/youtube"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", utils.GetEnvironmentVariableByKey("BOT_TOKEN"), "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var session *discordgo.Session
var guildManager = juanitacore.NewGuildManager()
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
	"basic-command": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
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
		var user = interaction.Interaction.Member.User
		var searchString = messageArguments[0].(string)
		var searchResult = youtubeSearcher.Search(searchString, user)

		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			// Ignore type for now, they will be discussed in "responses"
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: utils.PlayEmbed(*searchResult),
			},
		})
	},

	"followups": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		// Followup messages are basically regular messages (you can create as many of them as you wish)
		// but work as they are created by webhooks and their functionality
		// is for handling additional messages after sending a response.

		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// Note: this isn't documented, but you can use that if you want to.
				// This flag just allows you to create messages visible only for the caller of the command
				// (user who triggered the command)
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Surprise!",
			},
		})
		msg, err := session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: "Followup message has been created, after 5 seconds it will be edited",
		})
		if err != nil {
			session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
				Content: "Something went wrong",
			})
			return
		}
		time.Sleep(time.Second * 5)

		content := "Now the original message is gone and after 10 seconds this message will ~~self-destruct~~ be deleted."
		session.FollowupMessageEdit(interaction.Interaction, msg.ID, &discordgo.WebhookEdit{
			Content: &content,
		})

		time.Sleep(time.Second * 10)

		session.FollowupMessageDelete(interaction.Interaction, msg.ID)

		session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: "For those, who didn't skip anything and followed tutorial along fairly, " +
				"take a unicorn :unicorn: as reward!\n" +
				"Also, as bonus... look at the original interaction response :D",
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
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, v := range registeredCommands {
			err := session.ApplicationCommandDelete(session.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
