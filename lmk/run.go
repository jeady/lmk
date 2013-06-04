package main

import (
  "flag"
  "fmt"
  "strings"

  . "github.com/jeady/lmk/engine"
)

type RunCommand struct {
  flags *flag.FlagSet
}

func (cmd *RunCommand) Name() string {
  return "run"
}

func (cmd *RunCommand) Description() string {
  return "Sends out notifications for the specified rule if appropriate"
}

func (cmd *RunCommand) PrintHelp() {
  fmt.Println("usage: lmk run $rule-name-1 ...")
  fmt.Println("")
  fmt.Println("run tests the specified rules and sends out notifications")
  fmt.Println("if the rule is not sane or has been triggered.")
  fmt.Println("")
  fmt.Println("$rule-name may be any unambiguous case-insensitive subset of")
  fmt.Println("the rule's name.")
}

func (cmd *RunCommand) Init(f *flag.FlagSet) {
  cmd.flags = f
}

func (cmd *RunCommand) Main(e *Engine) int {
  rules := make([]Rule, 0)
  for _, name := range cmd.flags.Args() {
    lname := strings.ToLower(name)
    r := make([]Rule, 0)
    for _, rule := range e.Rules() {
      if strings.Contains(strings.ToLower(rule.Name()), lname) {
        r = append(r, rule)
      }
    }

    if len(r) == 1 {
      rules = append(rules, r[0])
    } else if len(r) == 0 {
      fmt.Println("Could not locate a rule matching '" + name + "'")
      return 1
    } else {
      fmt.Println("Len: " + string(len(r)))
      fmt.Println("Could not disambiguate '" + name + "'")
      fmt.Println("")
      fmt.Println("Could be:")
      for _, rule := range r {
        fmt.Println("  * " + rule.Name())
      }
      return 1
    }
  }

  for _, rule := range rules {
    e.Run(rule)
  }

  return 0
}
