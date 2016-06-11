package helpers

//RecoverableError is an error for recoverable state
type RecoverableError struct {
	message string
}

//UnRecoverableError is an error for unrecoverable state
type UnRecoverableError struct {
	message string
}

//Error returns error message
func (err *RecoverableError) Error() string {
	return err.message
}

//Error returns error message
func (err *UnRecoverableError) Error() string {
	return err.message
}

//NewRecoverableError returns new recoverable error
func NewRecoverableError(message string) error {
	return &RecoverableError{message: message}
}

//NewUnRecoverableError returns new unrecoverable error
func NewUnRecoverableError(message string) error {
	return &UnRecoverableError{message: message}
}
