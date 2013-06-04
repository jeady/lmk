package main

import (
  "flag"
  "fmt"
  "os"

  . "github.com/jeady/lmk/engine"
)

func main() {
  flag.Parse()

  cmd_name := "help"

  if flag.NArg() >= 1 {
    cmd_name = flag.Arg(0)
  }

  commands := []Command{
    new(ListRulesCommand),
    new(TestCommand),
    new(TestAllCommand),
    new(TestSmtpCommand),
    new(RunCommand),
    new(RunAllCommand),
    new(VersionCommand),
  }

  // Try to run the command the user requested.
  for _, cmd := range commands {
    if cmd.Name() == cmd_name {
      f := flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
      cmd.Init(f)

      var config_filename string
      f.StringVar(
        &config_filename,
        "config",
        "lmk.cfg",
        "The configuration file to use")

      f.Parse(flag.Args()[1:])
      e := NewEngine(config_filename)

      os.Exit(cmd.Main(e))
    }
  }

  // We'll only ever get this far if either the command was "help" or
  // unrecognized.
  if flag.NArg() >= 2 {
    for _, cmd := range commands {
      if cmd.Name() == flag.Arg(1) {
        f := flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
        cmd.Init(f)
        cmd.PrintHelp()
        os.Exit(0)
      }
    }
  }

  // If all else fails, show the most general help screen we've got.
  fmt.Println("lmk! is a service to notify you of things.")
  fmt.Println("")
  fmt.Println("Usage:")
  fmt.Println("")
  fmt.Println("    lmk command [arguments] [--config configfile]")
  fmt.Println("")
  fmt.Println("The commands are:")
  fmt.Println("")
  for _, cmd := range commands {
    fmt.Printf("    %-12s    %s\n", cmd.Name(), cmd.Description())
  }
  fmt.Println("")
  fmt.Println("Every command accepts the -c option to specify the config")
  fmt.Println("file to be used.")
  fmt.Println("")
  fmt.Println("Use `lmk help [command]` for more information about a command.")
  fmt.Println("")

  os.Exit(0)
}
