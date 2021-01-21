package logger

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/rs/zerolog"
)

func testConsoleLogger(cw io.Writer) *zerolog.Logger {
	// c := config.GetConfigFromFile("test")
	return NewLogger(nil, cw, nil)
}

func TestGenerateMultipleConsoleLog(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("log :%d", i), func(t *testing.T) {
			out := &bytes.Buffer{}
			log := testConsoleLogger(out)
			log.Log().Msg(fmt.Sprintf("log :%d", i))
			if got := string(out.Bytes()); got == "" {
				t.Errorf("invalid log output:\n got:  %s", got)
			}
		})
	}
}
