package main

import (
  "flag"
  "fmt"
)

const Version = "0.0"

type VersionCommand struct{}

func (cmd *VersionCommand) Name() string {
  return "version"
}

func (cmd *VersionCommand) Description() string {
  return "Print lmk! version"
}

func (cmd *VersionCommand) PrintHelp() {
  fmt.Println("usage: lmk version")
  fmt.Println("")
  fmt.Println("Version prints the lmk compile version.")
}

func (cmd *VersionCommand) Init(f *flag.FlagSet) {}

func (cmd *VersionCommand) Main() int {
  fmt.Printf("lmk! %s\n", Version)
  fmt.Printf("Copyright (c) 2013 James Eady\n")

  return 0
}
