package models

import "time"

type APIResponse struct {
	Status         string      `json:"status"`
	Data           interface{} `json:"data"`
	ErrorMessage   string      `json:"message,omitempty"`
}

type Timestamp struct {
	Unix int64 `json:"timestamp"`
}

func (t *Timestamp) Time() time.Time {
	return time.Unix(t.Unix/1000, (t.Unix%1000)*1000000)
}
