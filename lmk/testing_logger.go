package lmk

import "github.com/op/go-logging"

type TestingLogger struct {
  logs string
}

func (t *TestingLogger) Log(_ logging.Level, _ int, r *logging.Record) error {
  t.logs += r.Formatted() + "\n"
  return nil
}

func (t *TestingLogger) Logs() string {
  return t.logs
}

// Resets the stored logs so that they contain nothing.
func (t *TestingLogger) Reset() {
  t.logs = ""
}
