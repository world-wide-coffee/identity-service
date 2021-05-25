package logic

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/world-wide-coffee/identity-service/persistence"
)

func (l *Logic) SignUp(ctx context.Context, username, password string) (_ uuid.UUID, err error) {
	if len(username) == 0 {
		err = ErrorEmptyUsername{}
		return
	}

	if len(password) < passwordLength {
		err = ErrorTooShortPassword{}
		return
	}

	_, err = l.persister.IdentityByUsername(ctx, username)
	if err == nil {
		switch {
		case errors.Is(err, persistence.ErrorWrongUsername{}):
			err = ErrorWrongUsername{}
		}

		return
	}

	ident := persistence.Identity{ID: uuid.New(), Username: username, Password: password}

	if err = l.persister.PersistIdentity(ctx, ident); err != nil {
		return
	}

	return ident.ID, nil
}
