package eeris

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rotisserie/eris"
)

type PrinterError interface {
	Error(...string)
}

type defaultPrinter struct{}

func (d defaultPrinter) Error(content ...string) {
	c := make([]any, len(content))
	for i, v := range content {
		c[i] = v
	}
	log.Fatalf(content[0], c[1:]...)
}

func NewError(cls, msg string) error {
	return eris.New(fmt.Sprintf("[%s] %s", cls, msg))
}

func NewWrapError(err error, cls, msg string) error {
	return eris.Wrap(err, fmt.Sprintf("[%s] %s", cls, msg))
}

func TraceAndExit(err error, printer PrinterError) {
	format := eris.NewDefaultStringFormat(eris.FormatOptions{
		InvertOutput: true,
		WithTrace:    true,
		InvertTrace:  true,
	})

	strerr := eris.ToCustomString(err, format)

	p := printer
	if p == nil {
		p = defaultPrinter{}
	}
	p.Error(strings.Split(strerr, "\n")...)
	os.Exit(3)
}

func ToJsonString(err error, indent bool) string {
	value := eris.ToJSON(err, true)

	dat, _ := json.Marshal(value)
	return string(dat)
}

func ToJsonMap(err error, indent bool) map[string]any {
	return eris.ToJSON(err, true)
}
