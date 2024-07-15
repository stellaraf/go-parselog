package types

import (
	"encoding/json"
	"time"
)

type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	str := ""
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	// 2024-07-13 21:57:59
	ts, err := time.Parse(time.DateTime, str)
	if err != nil {
		return err
	}
	t.Time = ts
	return nil
}

type Request struct {
	Message   string    `json:"message"`
	Platform  string    `json:"platform"`
	Source    string    `json:"source"`
	Timestamp Timestamp `json:"timestamp"`
}
