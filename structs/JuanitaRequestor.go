package structs

type JuanitaRequestor struct {
	Tag string
	Id  string
}

func NewJuanitaRequestor(tag string, id string) *JuanitaRequestor {
	requestor := new(JuanitaRequestor)
	requestor.Tag = tag
	requestor.Id = id
	return requestor
}

func (jr JuanitaRequestor) Mention() string {
	return "<@" + jr.Id + ">"
}
