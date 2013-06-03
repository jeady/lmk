package main

import (
  "flag"

  "github.com/jeady/lmk/lmk"
)

type Command interface {
  Name() string
  Description() string
  PrintHelp()

  Init(*flag.FlagSet)
  Main(*lmk.Engine) int
}
