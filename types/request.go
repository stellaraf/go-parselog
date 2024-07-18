package types

import (
	"encoding/json"
	"strings"
	"time"
)

type Request struct {
	Messages  []string       `json:"message"`
	Platform  string         `json:"platform"`
	Source    string         `json:"source"`
	Timestamp time.Time      `json:"timestamp"`
	Extra     map[string]any `json:"extra"`
}

func (req *Request) UnmarshalJSON(b []byte) error {
	var initial map[string]any
	err := json.Unmarshal(b, &initial)
	if err != nil {
		return err
	}
	_platform, ok := initial["platform"]
	if !ok {
		return MissingFieldErr("platform")
	}
	platform, ok := _platform.(string)
	if !ok {
		return InvalidTypeErr("platform")
	}
	_source, ok := initial["source"]
	if !ok {
		return MissingFieldErr("source")
	}
	source, ok := _source.(string)
	if !ok {
		return InvalidTypeErr("source")
	}
	_tss, ok := initial["timestamp"]
	if !ok {
		return MissingFieldErr("timestamp")
	}
	_ts, ok := _tss.(string)
	if !ok {
		return InvalidTypeErr("timestamp")
	}
	ts, err := time.Parse(time.DateTime, _ts)
	if err != nil {
		return err
	}
	_extra, ok := initial["extra"]
	var extra map[string]any
	if !ok {
		extra = make(map[string]any, 0)
	} else {
		_extra, ok := _extra.(map[string]any)
		if !ok {
			return InvalidTypeErr("extra")
		}
		extra = _extra
	}
	_msg, ok := initial["message"]
	if !ok {
		return MissingFieldErr("message")
	}
	msg, ok := _msg.(string)
	if !ok {
		return InvalidTypeErr("message")
	}
	msgs := strings.Split(msg, "__")
	for i, msg := range msgs {
		msgs[i] = strings.TrimSpace(msg)
	}
	req.Messages = msgs
	req.Platform = platform
	req.Source = source
	req.Timestamp = ts
	req.Extra = extra
	return nil
}
