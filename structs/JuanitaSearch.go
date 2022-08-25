package structs

import "time"

type JuanitaSearch struct {
	Song      JuanitaSong
	Requestor JuanitaRequestor
	Date      time.Time
}

func NewJuanitaSearch(song JuanitaSong, requestor JuanitaRequestor, date time.Time) JuanitaSearch {
	return JuanitaSearch{Song: song, Requestor: requestor, Date: date}
}
