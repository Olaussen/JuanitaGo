package juanitacore

type (
	JuanitaCommand func(Context)

	JuanitaCommandStruct struct {
		command JuanitaCommand
		help    string
	}

	JuanitaCmdMap map[string]JuanitaCommandStruct

	JuanitaCommandHandler struct {
		cmds JuanitaCmdMap
	}
)

func NewCommandHandler() *JuanitaCommandHandler {
	return &JuanitaCommandHandler{make(JuanitaCmdMap)}
}

func (handler JuanitaCommandHandler) GetCmds() JuanitaCmdMap {
	return handler.cmds
}

func (handler JuanitaCommandHandler) Get(name string) (*JuanitaCommand, bool) {
	cmd, found := handler.cmds[name]
	// For legacy reasons, lets just deliver the commnd
	// A new function can be made GetAll() ??
	return &cmd.command, found
}

func (handler JuanitaCommandHandler) Register(name string, command JuanitaCommand, helpmsg string) {
	// Massage the arguments into a "Full command"
	cmdstruct := JuanitaCommandStruct{command: command, help: helpmsg}
	handler.cmds[name] = cmdstruct
	if len(name) > 1 {
		handler.cmds[name[:1]] = cmdstruct
	}
}

func (command JuanitaCommandStruct) GetHelp() string {
	return command.help
}
