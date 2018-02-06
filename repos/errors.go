package repos

type EmptyListError struct {
	message string
}

func (e EmptyListError) Error() string {
	return e.message
}

func NewEmptyListError(message string) EmptyListError {
	return EmptyListError{
		message: message,
	}
}
