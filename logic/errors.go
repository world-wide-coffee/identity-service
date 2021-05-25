package logic

type ErrorWrongPassword struct{ error }
type ErrorWrongUsername struct{ error }
type ErrorTooShortPassword struct{ error }
type ErrorEmptyUsername struct{ error }
type ErrorNotFound struct{ error }
