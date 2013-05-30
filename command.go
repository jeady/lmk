package main

import "flag"

type Command interface {
  Name() string
  Description() string
  PrintHelp()

  Init(*flag.FlagSet)
  Main() int
}
