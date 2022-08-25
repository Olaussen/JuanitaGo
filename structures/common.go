package structures

import (
	"math"
	"time"

	"github.com/kkdai/youtube/v2"
)

type JuanitaRequestor struct {
	Tag string
	Id  string
}

func NewJuanitaRequestor(tag string, id string) JuanitaRequestor {
	return JuanitaRequestor{Tag: tag, Id: id}
}

func DefaultRequestor() JuanitaRequestor {
	return NewJuanitaRequestor("Unknown", "Unknown")
}

type JuanitaSong struct {
	Title   string
	Url     string
	Seconds uint
	Video   *youtube.Video
}

func NewJuanitaSong(title string, url string, seconds uint, video *youtube.Video) JuanitaSong {
	return JuanitaSong{Title: title, Url: url, Seconds: seconds / uint(math.Pow10(9)), Video: video}
}

func DefaultSong() JuanitaSong {
	return NewJuanitaSong("TestSong", "https://www.youtube.com/watch?v=9bZkp7q19f0", 0, nil)
}

type JuanitaSearch struct {
	Song      JuanitaSong
	Requestor JuanitaRequestor
	Date      time.Time
}

func NewJuanitaSearch(song JuanitaSong, requestor JuanitaRequestor, date time.Time) JuanitaSearch {
	return JuanitaSearch{Song: song, Requestor: requestor, Date: date}
}

func DefaultSearch() JuanitaSearch {

	return NewJuanitaSearch(DefaultSong(), DefaultRequestor(), time.Now())
}
