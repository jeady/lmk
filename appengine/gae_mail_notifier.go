// +build appengine

package appengine

import (
  "fmt"

  "appengine"
  "appengine/mail"
)

type GaeMailNotifier struct {
  from    string
  context appengine.Context
}

func NewGaeMailNotifier(f string, c appengine.Context) *GaeMailNotifier {
  return &GaeMailNotifier{
    from:    f,
    context: c,
  }
}

// TODO(jmeady): Make lmk/engine.log externally available and use here for
// logging.
func (n *GaeMailNotifier) Notify(who, rule_name, msg string) error {
  fmt.Println("Attempting to mail %s", who)
  fmt.Println("From: " + n.from)
  fmt.Println("To: " + who)
  m := &mail.Message{
    Sender:  n.from,
    To:      []string{who},
    Subject: "lmk! update for " + rule_name,
    Body:    msg,
  }

  if err := mail.Send(n.context, m); err != nil {
    fmt.Println("Problem sending email: " + err.Error())
  }

  return nil
}
