package httpUtils

import (
	"time"
)

type Session struct {
	IP           string
	URL, Method  string
	ErrorMessage string
	ResponseCode string
	RequestTime  time.Time
}

func NewSessionRequest() *Session {
	return &Session{
		RequestTime: time.Now(),
	}
}
