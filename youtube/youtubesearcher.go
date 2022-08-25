package youtube

import (
	"fmt"
	"juanitaGo/structures"
	"juanitaGo/utils"
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	youtube "github.com/kkdai/youtube/v2"
	"google.golang.org/api/googleapi/transport"
	googleYoutube "google.golang.org/api/youtube/v3"
)

type YoutubeSearcher struct {
	service *googleYoutube.Service
	client  youtube.Client
}

// initialize new youtube client with api key
func NewYoutubeSearcher() YoutubeSearcher {
	httpClient := &http.Client{
		Transport: &transport.APIKey{Key: utils.GetEnvironmentVariableByKey("YOUTUBE_API_KEY")},
	}
	client := youtube.Client{}

	service, err := googleYoutube.New(httpClient)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}
	return YoutubeSearcher{service: service, client: client}
}

func (searcher YoutubeSearcher) GetVideoId(query string) string {
	var result, err = searcher.service.Search.List([]string{"id,snippet"}).SafeSearch("none").Q(query).MaxResults(1).Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
		return ""
	}
	var items = result.Items
	if len(items) > 0 {
		if items[0].Id.Kind == "youtube#video" {
			return items[0].Id.VideoId
		}
	}
	return ""
}

func (searcher YoutubeSearcher) Search(query string, user *discordgo.User) *structures.JuanitaSearch {
	var videoId = searcher.GetVideoId(query)
	if videoId == "" {
		return nil
	}
	video, err := searcher.client.GetVideo(videoId)
	if err != nil {
		log.Fatalf("Error getting video: %v", err)
		return nil
	}
	/*formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := searcher.client.GetStream(video, &formats[0])
	if err != nil {
		log.Fatalf("Error getting stream: %v", err)
		return nil
	}*/

	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoId)
	song := structures.NewJuanitaSong(video.Title, url, uint(video.Duration), video)
	requestor := structures.NewJuanitaRequestor(utils.ExtractUserTag(user), user.ID)
	search := structures.NewJuanitaSearch(song, requestor, time.Now())
	fmt.Printf("%v", search)
	return &search
}
