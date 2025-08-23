package serviceError

type ErrApp interface {
	Error() string
	Code(code int) ErrApp
	GetCode() int
}

type errApp struct {
	message string
	code    int
}

func NewErr(message string) *errApp {
	return &errApp{
		message: message,
	}
}

func (e *errApp) Error() string {
	return e.message
}

func (e *errApp) Code(code int) ErrApp {
	e.code = code
	return e
}

func (e *errApp) GetCode() int {
	return e.code
}
