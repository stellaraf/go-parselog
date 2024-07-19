package parselog

import (
	"github.com/stellaraf/go-parselog/junos"
	"github.com/stellaraf/go-parselog/types"
)

type (
	Request = types.Request
	BGPLog  = types.BGPLog
	ISISLog = types.ISISLog
	Log     = types.Log
	State   = types.State
	LogType = types.LogType
)

var (
	ErrNoMatchingParser   = types.ErrNoMatchingParser
	ErrIncompleteMatch    = types.ErrIncompleteMatch
	ErrNoMatchingPlatform = types.ErrNoMatchingPlatform
	ISIS                  = types.ISIS
	BGP                   = types.BGP
	UP                    = types.UP
	DOWN                  = types.DOWN
	ISISLogType           = &types.ISISLog{Base: types.Base{Type: ISIS}}
	BGPLogType            = &types.BGPLog{Base: types.Base{Type: BGP}}
)

var parseMap = map[string]types.Parser{
	"junos": junos.Parse,
}

func Parse(request *Request) ([]Log, error) {
	parser, ok := parseMap[request.Platform]
	if !ok {
		return nil, types.ErrNoMatchingPlatform
	}
	return parser(request)
}
