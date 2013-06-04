package main

import (
  "flag"
  "fmt"

  . "github.com/jeady/lmk/engine"
)

type ListRulesCommand struct{}

func (cmd *ListRulesCommand) Name() string {
  return "list-rules"
}

func (cmd *ListRulesCommand) Description() string {
  return "List the notification rules"
}

func (cmd *ListRulesCommand) PrintHelp() {
  fmt.Println("usage: lmk list-rules")
  fmt.Println("")
  fmt.Println("list-rules prints the valid and enabled notification rules.")
}

func (cmd *ListRulesCommand) Init(f *flag.FlagSet) {}

func (cmd *ListRulesCommand) Main(e *Engine) int {
  fmt.Println("Enabled rules:")
  for _, rule := range e.Rules() {
    fmt.Println("  ", rule.Name())
  }

  return 0
}
