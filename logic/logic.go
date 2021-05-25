package logic

import (
	"github.com/google/uuid"

	"github.com/world-wide-coffee/identity-service/persistence"
)

const (
	hmacKey   = "monkeybar"
	hmacKeyID = "woopwoop"
)

type Logic struct {
	persister persistence.Persister
}

func NewLogic() *Logic {
	return &Logic{}
}

func (l *Logic) SetPersister(p persistence.Persister) *Logic {
	l.persister = p
	return l
}

const passwordLength = 10

type Identity struct {
	ID       uuid.UUID
	Username string
}
