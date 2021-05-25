package persistence

import (
	"context"

	"github.com/google/uuid"
)

type Persister interface {
	PersistIdentity(_ context.Context, _ Identity) error
	IdentityByUsername(_ context.Context, username string) (Identity, error)
	IdentityByIDAndPassword(_ context.Context, id uuid.UUID, password string) (Identity, error)
	IdentityByUsernameAndPassword(_ context.Context, username, password string) (Identity, error)
	ChangePassword(_ context.Context, id uuid.UUID, password string) error
}

type Identity struct {
	ID       uuid.UUID
	Username string
	Password string
}
