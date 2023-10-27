package xerror

type Error struct {
	Code   int
	Status string
	Err    error
}

func New(code int, status string, err error) *Error {
	return &Error{
		Code:   code,
		Status: status,
		Err:    err,
	}
}

func (error *Error) String() string {
	if error.Err == nil {
		return ""
	}
	return error.Err.Error()
}
