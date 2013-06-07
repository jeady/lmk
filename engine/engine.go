package engine

import (
  "errors"
  "net/smtp"
  "strings"

  "github.com/jeady/go-logging"
)

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

func (e *Engine) Rules() []Rule {
  return e.conf.Rules()
}

// Attempts to send msg to the currently configured sink as a notification.
func (e *Engine) EmailNotification(recipient, subj, msg string) error {
  user, pass, host := e.conf.SmtpConfig()
  if len(user) == 0 || len(host) == 0 {
    return errors.New("SMTP is not configured.")
  }

  log.Info("Attempting to mail %s", recipient)
  err := smtp.SendMail(
    host,
    smtp.PlainAuth("", user, pass, strings.Split(host, ":")[0]),
    user,
    []string{recipient},
    []byte("Subject: "+subj+"\r\n\r\n"+msg))

  if err != nil {
    log.Error("Problem sending email: " + err.Error())
  }

  return err
}

// Tests the rule and sends out a notification if the rule is not sane or has
// been triggered.
func (e *Engine) Run(r Rule) {
  sane, triggered := r.TestTriggered()

  if !sane {
    log.Info("Rule '%s' is not sane, notifying.", r.Name())
    e.EmailNotification(
      e.conf.NotificationRecipient(),
      "[lmk] '"+r.Name()+"' is not sane <eom>",
      "")
  } else if triggered {
    log.Info("Rule '%s' has been triggered, notifying.", r.Name())
    e.EmailNotification(
      e.conf.NotificationRecipient(),
      "[lmk] '"+r.Name()+"' has been triggered <eom>",
      "")
  } else {
    log.Debug("Rule '%s' is dormant.", r.Name())
  }
}
