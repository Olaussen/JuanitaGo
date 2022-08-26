package structs

import "time"

type JuanitaSearch struct {
	Song      JuanitaSong
	Requestor JuanitaRequestor
	Date      time.Time
}

func NewJuanitaSearch(song JuanitaSong, requestor JuanitaRequestor, date time.Time) *JuanitaSearch {
	search := new(JuanitaSearch)
	search.Song = song
	search.Requestor = requestor
	search.Date = date
	return search
}
