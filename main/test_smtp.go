package main

import (
  "flag"
  "fmt"

  "github.com/jeady/lmk/lmk"
)

type TestSmtpCommand struct {
  flags *flag.FlagSet
}

func (cmd *TestSmtpCommand) Name() string {
  return "test-smtp"
}

func (cmd *TestSmtpCommand) Description() string {
  return "Sends a test message from the configured smtp user"
}

func (cmd *TestSmtpCommand) PrintHelp() {
  fmt.Println("usage: lmk test-smtp recipient@address")
  fmt.Println("")
  fmt.Println("test-smtp tests that smtp is correctly configured to send")
  fmt.Println("notifications by attempting to send a test message to the")
  fmt.Println("supplied email address.")
}

func (cmd *TestSmtpCommand) Init(f *flag.FlagSet) {
  cmd.flags = f
}

func (cmd *TestSmtpCommand) Main(e *lmk.Engine) int {
  if cmd.flags.NArg() < 1 {
    cmd.PrintHelp()
    return 1
  }

  err := e.EmailNotification(
    cmd.flags.Arg(0),
    "lmk! test notification",
    "lmk! is correctly configured to send email. Cheers!")

  if err != nil {
    fmt.Println("Problem sending mail.")
    return 1
  }
  fmt.Println("Successfully sent mail.")
  return 0
}
