package game

import (
	"othello/board"
	"othello/builtinai"
)

type Agent int

const (
	AgentNone Agent = iota
	AgentHuman
	AgentBuiltIn
	AgentExternal
)

type Parameter struct {
	BlackAgent   Agent
	WhiteAgent   Agent
	BlackPath    string
	WhitePath    string
	BlackAILevel builtinai.Level
	WhiteAILevel builtinai.Level
	GoesFirst    board.Color
}

func NewParam() Parameter {
	return Parameter{
		BlackAgent: AgentNone,
		BlackPath:  "",
		WhiteAgent: AgentNone,
		WhitePath:  "",
	}
}

func (params Parameter) AllSelected() bool {
	return params.BlackAgent != AgentNone && params.WhiteAgent != AgentNone
}
