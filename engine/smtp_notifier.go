package engine

import (
  "net/smtp"
  "strings"
)

type SmtpNotifier struct {
  user string
  pass string
  host string
}

func NewSmtpNotifier(user, pass, host string) *SmtpNotifier {
  return &SmtpNotifier{
    user: user,
    pass: pass,
    host: host,
  }
}

func (n *SmtpNotifier) Notify(who, rule_name, msg string) error {
  log.Info("Attempting to mail %s", who)
  err := smtp.SendMail(
    n.host,
    smtp.PlainAuth("", n.user, n.pass, strings.Split(n.host, ":")[0]),
    n.user,
    []string{who},
    []byte("Subject: lmk! update for "+rule_name+"\r\n\r\n"+msg))

  if err != nil {
    log.Error("Problem sending email: " + err.Error())
  }

  return err
}
