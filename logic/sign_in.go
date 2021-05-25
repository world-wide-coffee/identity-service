package logic

import (
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/world-wide-coffee/identity-service/persistence"
)

func (l *Logic) SignIn(ctx context.Context, username, password string) (_ string, err error) {
	ident, err := l.persister.IdentityByUsernameAndPassword(ctx, username, password)
	if err != nil {
		switch {
		case errors.Is(err, persistence.ErrorWrongUsername{}):
			err = ErrorWrongUsername{}
		case errors.Is(err, persistence.ErrorWrongPassword{}):
			err = ErrorWrongPassword{}
		}

		return
	}

	// Create the token
	token := jwt.New(jwt.GetSigningMethod("HS256"))

	// Set some claims
	issuedAt := time.Now()
	token.Claims = jwt.StandardClaims{
		Subject:   ident.ID.String(),
		IssuedAt:  issuedAt.Unix(),
		ExpiresAt: issuedAt.Add(time.Hour).Unix(),
	}

	token.Header["kid"] = hmacKeyID

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(hmacKey))
	if err != nil {
		return
	}

	return tokenString, nil
}
