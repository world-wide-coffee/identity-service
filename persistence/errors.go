package persistence

import "fmt"

type ErrorWrongPassword struct{ error }
type ErrorWrongUsername struct{ error }
type ErrorTooShortPassword struct{ error }
type ErrorEmptyUsername struct{ error }
type ErrorNotFound struct{ error }
type ErrorInternal struct{ error }

func NewErrorInternal(format string, a ...interface{}) ErrorInternal {
	return ErrorInternal{error: fmt.Errorf(format, a...)}
}
