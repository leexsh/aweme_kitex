package errno

import (
	"errors"
	"fmt"
)

type ErrMsg struct {
	ErrCode int64
	ErrMsg  string
}

func (e ErrMsg) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.ErrCode, e.ErrMsg)
}

func NewErr(code int64, msg string) ErrMsg {
	return ErrMsg{code, msg}
}

func (e ErrMsg) WithMessage(msg string) ErrMsg {
	e.ErrMsg = msg
	return e
}

// ConvertErr convert error to Errno
func ConvertErr(err error) ErrMsg {
	Err := ErrMsg{}
	if errors.As(err, &Err) {
		return Err
	}

	s := ServiceErr
	s.ErrMsg = err.Error()
	return s
}
