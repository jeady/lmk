package main

import (
  "flag"
  "fmt"

  "github.com/jeady/lmk/lmk"
)

type TestAllCommand struct {
  flags *flag.FlagSet
}

func (cmd *TestAllCommand) Name() string {
  return "test-all"
}

func (cmd *TestAllCommand) Description() string {
  return "List the notification rules"
}

func (cmd *TestAllCommand) PrintHelp() {
  fmt.Println("usage: lmk test-all")
  fmt.Println("")
  fmt.Println("test-all tests if all enabled rules are sane and/or triggered.")
}

func (cmd *TestAllCommand) Init(f *flag.FlagSet) {
  cmd.flags = f
}

func (cmd *TestAllCommand) Main(e *lmk.Engine) int {
  for _, rule := range e.Rules() {
    sane, triggered := rule.TestTriggered()
    if !sane {
      fmt.Println(rule.Name() + " is not sane.")
      continue
    }

    fmt.Println(rule.Name() + " is sane.")
    if !triggered {
      fmt.Println(rule.Name() + " has not been triggered.")
    } else {
      fmt.Println(rule.Name() + " has been triggered.")
    }
    fmt.Println("")
  }

  return 0
}
