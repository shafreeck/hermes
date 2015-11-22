package session

import (
	"github.com/shafreeck/hermes/hermes"
)

type Session interface {
	Get(sessionid string)
	Set(sessionid string, s session)
}

type session struct {
	sessionid string
}
type Client struct {
	clientID string
	cursors  map[string]hermes.Cursor
}
