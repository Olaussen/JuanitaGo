package messages

import (
	"fmt"
	"juanitaGo/structs"

	"github.com/bwmarrin/discordgo"
)

func PlayEmbed(songSearch structs.JuanitaSearch) []*discordgo.MessageEmbed {
	song := songSearch.Song
	requestor := songSearch.Requestor

	return []*discordgo.MessageEmbed{
		{
			Title: "Bruk knappene for å kontrollere musikken 🎶",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Nå spiller:",
					Value: fmt.Sprintf("[%v](%v) lagt til av <@%v>", song.Title(), song.Url(), requestor.Id),
				},
				{
					Name:  "Duration",
					Value: fmt.Sprintf("%v", song.Duration()),
				},
			},
		},
	}
}
