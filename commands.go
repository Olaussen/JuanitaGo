package main

import "github.com/bwmarrin/discordgo"

var CommandDescriptions = map[string]string{
	"test":     "Søk og spill av en sang",
	"first":    "Sniker en sang først inn i køen",
	"kys":      "Forlater kanalen og sletter køen",
	"lists":    "Lister opp alle lagrede Spotify spillelister",
	"spotify":  "Legger alle sanger i en definert Spotify spillesliste til køen",
	"remember": "Lagrer en eksisterende Spotify spilleliste med alias for å finne den lettere senere",
	"skip":     "Hopp over en sang i køen",
	"resume":   "Gjenoppta en pauset sang",
	"pause":    "Setter en sang på pause",
	"restart":  "Restarter Juanita",
}

func GetCommandConfig() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "basic-command",
			Description: "Basic command",
		},
		{
			Name:        "test",
			Description: CommandDescriptions["test"],
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "sangnavn",
					Description: "Sangnavn eller søkeord",
					Required:    true,
				},
			},
		},
		{
			Name:        "followups",
			Description: "Followup messages",
		},
	}
}
