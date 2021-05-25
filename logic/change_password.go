package logic

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/world-wide-coffee/identity-service/persistence"
)

func (l *Logic) ChangePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) (err error) {
	_, err = l.persister.IdentityByIDAndPassword(ctx, id, oldPassword)
	if err != nil {
		switch {
		case errors.Is(err, persistence.ErrorNotFound{}):
			err = ErrorNotFound{}
		case errors.Is(err, persistence.ErrorWrongPassword{}):
			err = ErrorWrongPassword{}
		}

		return
	}

	if len(newPassword) < passwordLength {
		return ErrorTooShortPassword{}
	}

	err = l.persister.ChangePassword(ctx, id, newPassword)

	return
}
