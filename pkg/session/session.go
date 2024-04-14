package session

import (
	"avito_hr/pkg/user"
	"time"
)

type Session struct {
	Sub user.Role
	Exp time.Time
}

type SessionsPacker interface {
	Pack(role user.Role) (*string, error)
	Unpack(inToken string) (*Session, error)
}
