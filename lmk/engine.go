package lmk

import "github.com/op/go-logging"

type Engine struct {
  conf *Config
}

func NewEngine(config_file string) *Engine {
  var err error
  e := new(Engine)
  e.conf, err = NewConfig(config_file)

  if e.conf == nil {
    log.Error("Unable to load configuration file: " + err.Error())
    return nil
  }

  logging_loglevel, _ := logging.LogLevel(e.conf.LogLevel())
  logging.SetLevel(logging_loglevel, log.Module)
  log.Debug("Set logging level to %s", logging_loglevel.String())

  log.Info("Successfully initialized LetMeKnow engine.")

  return e
}
