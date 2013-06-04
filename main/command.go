package main

import (
  "flag"

  . "github.com/jeady/lmk/engine"
)

type Command interface {
  Name() string
  Description() string
  PrintHelp()

  Init(*flag.FlagSet)
  Main(*Engine) int
}
