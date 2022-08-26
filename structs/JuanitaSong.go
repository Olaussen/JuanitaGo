package structs

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/kkdai/youtube/v2"
)

type JuanitaSong struct {
	Video  *youtube.Video
	Stream io.ReadCloser
}

func NewJuanitaSong(video *youtube.Video, stream io.ReadCloser) *JuanitaSong {
	song := new(JuanitaSong)
	song.Video = video
	song.Stream = stream
	return song
}

// get title of song
func (song JuanitaSong) Title() string {
	return song.Video.Title
}

// get url of song
func (song JuanitaSong) Url() string {
	return "https://www.youtube.com/watch?v=" + song.Video.ID
}

func (song JuanitaSong) Duration() string {

	seconds := int(math.Round(song.Video.Duration.Seconds()))
	minutes := seconds / 60
	seconds = seconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func (song JuanitaSong) AudioStream() string {
	file, err := os.Create(fmt.Sprintf("%v.mp3", song.Video.ID))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = io.Copy(file, song.Stream)
	return fmt.Sprintf("%v.mp3", song.Video.ID)
}
