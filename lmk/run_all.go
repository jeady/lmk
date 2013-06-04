package main

import (
  "flag"
  "fmt"

  . "github.com/jeady/lmk/engine"
)

type RunAllCommand struct {
  flags *flag.FlagSet
}

func (cmd *RunAllCommand) Name() string {
  return "run-all"
}

func (cmd *RunAllCommand) Description() string {
  return "Sends out notifications for all rulese where appropriate"
}

func (cmd *RunAllCommand) PrintHelp() {
  fmt.Println("usage: lmk run-all")
  fmt.Println("")
  fmt.Println("run tests all rules and sends out notifications for rules that")
  fmt.Println("are not sane or have been triggered.")
}

func (cmd *RunAllCommand) Init(f *flag.FlagSet) {
  cmd.flags = f
}

func (cmd *RunAllCommand) Main(e *Engine) int {
  for _, rule := range e.Rules() {
    e.Run(rule)
  }

  return 0
}
