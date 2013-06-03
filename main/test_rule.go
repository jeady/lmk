package main

import (
  "flag"
  "fmt"
  "strings"

  "github.com/jeady/lmk/lmk"
)

type TestRulesCommand struct {
  flags *flag.FlagSet
}

func (cmd *TestRulesCommand) Name() string {
  return "test-rule"
}

func (cmd *TestRulesCommand) Description() string {
  return "Tests if the given rules are sane and/or triggered"
}

func (cmd *TestRulesCommand) PrintHelp() {
  fmt.Println("usage: lmk test-rule $rule-name-1 ...")
  fmt.Println("")
  fmt.Println("test-rule tests the specified rules and prints whether they")
  fmt.Println("are sane and/or triggered.")
  fmt.Println("")
  fmt.Println("$rule-name may be any unambiguous case-insensitive subset of")
  fmt.Println("the rule's name.")
}

func (cmd *TestRulesCommand) Init(f *flag.FlagSet) {
  cmd.flags = f
}

func (cmd *TestRulesCommand) Main(e *lmk.Engine) int {
  rules := make([]lmk.Rule, 0)
  for _, name := range cmd.flags.Args() {
    lname := strings.ToLower(name)
    r := make([]lmk.Rule, 0)
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
