package logger

import (
	"bytes"
	"errors"
	"log"
	"strings"
	"testing"
)

func TestErrorJSONAndOptions(t *testing.T) {
	var out bytes.Buffer
	old := log.Writer()
	defer log.SetOutput(old)
	log.SetOutput(&out)
	New().Error("failed", errors.New("boom"), "worker")
	line := out.String()
	for _, part := range []string{`"type":"ERROR"`, `"message":"failed"`, `"error":"boom"`, `"component":"worker"`} {
		if !strings.Contains(line, part) {
			t.Errorf("log missing %s: %s", part, line)
		}
	}
}
