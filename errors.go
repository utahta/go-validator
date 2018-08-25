package validator

type (
	// Error represents a validation error
	Error struct {
	}

	Errors []Error
)

func ToErrors(err error) (Errors, bool) {
	es, ok := err.(Errors)
	return es, ok
}

func (e Error) Error() string {
	return ""
}

func (es Errors) Error() string {
	return ""
}
