package structs

type JuanitaRequestor struct {
	Tag string
	Id  string
}

func NewJuanitaRequestor(tag string, id string) JuanitaRequestor {
	return JuanitaRequestor{Tag: tag, Id: id}
}

func (jr JuanitaRequestor) Mention() string {
	return "<@" + jr.Id + ">"
}
