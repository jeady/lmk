package main

import (
  "flag"
  "fmt"
  "strings"

  . "github.com/jeady/lmk/engine"
)

type TestCommand struct {
  flags *flag.FlagSet
}

func (cmd *TestCommand) Name() string {
  return "test"
}

func (cmd *TestCommand) Description() string {
  return "Tests if the given rules are sane and/or triggered"
}

func (cmd *TestCommand) PrintHelp() {
  fmt.Println("usage: lmk test $rule-name-1 ...")
  fmt.Println("")
  fmt.Println("test tests the specified rules and prints whether they are")
  fmt.Println("sane and/or triggered. test does NOT send out notifications.")
  fmt.Println("")
  fmt.Println("$rule-name may be any unambiguous case-insensitive subset of")
  fmt.Println("the rule's name.")
}

func (cmd *TestCommand) Init(f *flag.FlagSet) {
  cmd.flags = f
}

func (cmd *TestCommand) Main(e *Engine) int {
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
