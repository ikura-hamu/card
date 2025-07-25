package merrors

import "fmt"

type Msg struct {
	err error
}

func New(err error) Msg {
	return Msg{err: err}
}

func (m Msg) Error() string {
	if m.err == nil {
		return "no error"
	}
	return fmt.Sprintf("msg error: %v", m.err)
}

func (m Msg) Unwrap() error {
	return m.err
}
