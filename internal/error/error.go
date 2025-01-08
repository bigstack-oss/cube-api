package error

import "errors"

const (
	Service = "service"
)

var (
	ServiceNotFound     = errors.New("service not found")
	TuningParamNotFound = errors.New("tuning parameter not found")
)

type Template struct {
	Occurred bool
	Type     string
	Msg      string
	Raw      error
}

func (t Template) Error() string {
	return t.Msg
}

func ErrService(err error) Template {
	return Template{
		Occurred: true,
		Type:     Service,
		Msg:      "configuration operation failure",
		Raw:      err,
	}
}
