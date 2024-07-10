package v1

import "errors"

const (
	Service = "service"
)

var (
	ErrServiceNotFound     = errors.New("service not found")
	ErrTuningParamNotFound = errors.New("tuning parameter not found")
)

type Err struct {
	Occurred bool
	Type     string
	Msg      string
	Raw      error
}

func (e Err) Error() string {
	return e.Msg
}

func ErrService(err error) Err {
	return Err{
		Occurred: true,
		Type:     Service,
		Msg:      "configuration operation failure",
		Raw:      err,
	}
}
