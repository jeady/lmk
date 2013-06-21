package engine

import (
  "time"

  "github.com/jeady/go-logging"
)

type Engine struct {
  conf             *Config
  default_notifier Notifier
}

func NewEngine(config_file string) *Engine {
  var err error
  e := new(Engine)
  e.conf, err = NewConfig(config_file)

  if e.conf == nil {
    log.Error("Unable to load configuration file: " + err.Error())
    return nil
  }

  e.default_notifier, err = e.conf.DefaultNotifier()
  if e.default_notifier == nil {
    log.Warning("Unable to create default notifier: " + err.Error())
  }

  logging_loglevel, _ := logging.LogLevel(e.conf.LogLevel())
  logging.SetLevel(logging_loglevel, log.Module)
  log.Debug("Set logging level to %s", logging_loglevel.String())

  log.Info("Successfully initialized LetMeKnow engine.")

  return e
}

func (e *Engine) Rules() []Rule {
  return e.conf.Rules()
}

func (e *Engine) RulesToPoll(last_update time.Time) []Rule {
  rules := []Rule{}

  for _, r := range e.conf.Rules() {
    if pr, ok := r.(PollingRule); ok && pr.ShouldPoll(last_update) {
      rules = append(rules, r)
    }
  }

  return rules
}

// Tests the rule and sends out a notification if the rule is not sane or has
// been triggered.
func (e *Engine) Run(r Rule) {
  if e.default_notifier == nil {
    log.Error("Cannot run rule '" + r.Name() + "': No default notifier.")
    return
  }

  sane, triggered := r.TestTriggered()

  if !sane {
    log.Info("Rule '%s' is not sane, notifying.", r.Name())
    e.default_notifier.Notify(
      e.conf.DefaultNotificationRecipient(),
      r.Name(),
      "'"+r.Name()+"' is not sane")
  } else if triggered {
    log.Info("Rule '%s' has been triggered, notifying.", r.Name())
    e.default_notifier.Notify(
      e.conf.DefaultNotificationRecipient(),
      r.Name(),
      "'"+r.Name()+"' has been triggered")
  } else {
    log.Debug("Rule '%s' is dormant.", r.Name())
  }
}

func (e *Engine) Config() *Config {
  return e.conf
}

func (e *Engine) DefaultNotifier() Notifier {
  return e.default_notifier
}

func (e *Engine) SetDefaultNotifier(n Notifier) Notifier {
  old := e.default_notifier
  e.default_notifier = n
  return old
}
