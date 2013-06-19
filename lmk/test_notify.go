package main

import (
  "flag"
  "fmt"

  . "github.com/jeady/lmk/engine"
)

type TestNotifyCommand struct {
  flags *flag.FlagSet
}

func (cmd *TestNotifyCommand) Name() string {
  return "test-notify"
}

func (cmd *TestNotifyCommand) Description() string {
  return "Sends a test message to the default notifier"
}

func (cmd *TestNotifyCommand) PrintHelp() {
  fmt.Println("usage: lmk test-notify recipient")
  fmt.Println("")
  fmt.Println("test-notify tests that the default notification mechanism is")
  fmt.Println("configured correctly my attempting to send a test message to")
  fmt.Println("the supplied recipient.")
  fmt.Println("")
  fmt.Println("For example, if smtp is configured as")
  fmt.Println("the default notification mechanism, run:")
  fmt.Println("`lmk test-notify jon@snow.com`")
}

func (cmd *TestNotifyCommand) Init(f *flag.FlagSet) {
  cmd.flags = f
}

func (cmd *TestNotifyCommand) Main(e *Engine) int {
  if cmd.flags.NArg() < 1 {
    cmd.PrintHelp()
    return 1
  }

  err := e.DefaultNotifier().Notify(
    cmd.flags.Arg(0),
    "test notification",
    "lmk! is correctly configured to send notifications. Cheers!")

  if err != nil {
    fmt.Println("Problem sending notifications.")
    return 1
  }
  fmt.Println("Successfully sent notifications.")
  return 0
}
