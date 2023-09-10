package errors

import (
	"encoding/json"
	"fmt"

	"github.com/wfrodriguez/mimir/eris"
)

func NewError(cls, msg string) error {
	return eris.New(fmt.Sprintf("[%s] %s", cls, msg))
}

func NewWrapError(err error, cls, msg string) error {
	return eris.Wrap(err, fmt.Sprintf("[%s] %s", cls, msg))
}

func TraceString(err error) string {
	format := eris.NewDefaultStringFormat(eris.FormatOptions{
		InvertOutput: true,
		WithTrace:    true,
		InvertTrace:  true,
	})

	// format the error to a string and return
	return eris.ToCustomString(err, format)
}

func ToJsonString(err error, indent bool) string {
	value := eris.ToJson(err, true)

	dat, _ := json.Marshal(value)
	return string(dat)
}

func ToJsonMap(err error, indent bool) map[string]any {
	return eris.ToJson(err, true)
}
