package types

import (
	"time"

	"github.com/stellaraf/go-utils"
)

type State uint
type LogType uint

type Parser func(*Request) (Log, error)

const (
	UP State = iota + 1
	DOWN
)

const (
	ISIS LogType = iota + 1
	BGP
)

type Base struct {
	Type     LogType        `json:"type"`
	Extra    map[string]any `json:"extra"`
	Original string         `json:"original"`
}

type Log interface {
	Is(Log) bool
	ID() (string, error)
	Attrs() map[string]any
	Up() bool
	Down() bool
	LogType() LogType
}

type ISISLog struct {
	Base
	Local     string    `json:"local"`
	Remote    string    `json:"remote"`
	Timestamp time.Time `json:"timestamp"`
	State     State     `json:"state"`
	Interface string    `json:"interface"`
	Reason    string    `json:"reason"`
}

type BGPLog struct {
	Base
	Local     string    `json:"local"`
	Remote    string    `json:"remote"`
	Timestamp time.Time `json:"timestamp"`
	State     State     `json:"state"`
	RemoteAS  string    `json:"remote_as"`
	Table     string    `json:"table"`
}

// ISISLog Methods

func (l *ISISLog) Up() bool {
	return l.State == UP
}

func (l *ISISLog) Down() bool {
	return l.State == DOWN
}

func (l *ISISLog) Is(other Log) bool {
	return other.LogType() == ISIS
}

func (l *ISISLog) ID() (string, error) {
	return utils.HashFromStrings(l.Local, l.Remote, l.Interface)
}

func (l *ISISLog) LogType() LogType {
	return l.Type
}

func (l *ISISLog) Attrs() map[string]any {
	return map[string]any{
		"local":     l.Local,
		"remote":    l.Remote,
		"timestamp": l.Timestamp,
		"state":     l.State,
		"interface": l.Interface,
		"reason":    l.Reason,
		"type":      l.Type,
		"extra":     l.Extra,
		"original":  l.Original,
	}
}

// BGPLog Methods

func (l *BGPLog) Is(other Log) bool {
	return other.LogType() == BGP
}

func (l *BGPLog) LogType() LogType {
	return l.Type
}

func (l *BGPLog) ID() (string, error) {
	return utils.HashFromStrings(l.Local, l.Remote, l.RemoteAS, l.Table)
}

func (l *BGPLog) Up() bool {
	return l.State == UP
}

func (l *BGPLog) Down() bool {
	return l.State == DOWN
}

func (l *BGPLog) Attrs() map[string]any {
	return map[string]any{
		"local":     l.Local,
		"remote":    l.Remote,
		"timestamp": l.Timestamp,
		"state":     l.State,
		"remote_as": l.RemoteAS,
		"table":     l.Table,
		"type":      l.Type,
		"extra":     l.Extra,
		"original":  l.Original,
	}
}
