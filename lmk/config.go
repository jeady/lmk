package lmk

import (
  "github.com/msbranco/goconfig"
)

type Config struct {
  filename string
  file     *goconfig.ConfigFile

  loglevel string
}

func new_config(filename string) (*Config, error) {
  var err error
  goconfig.DefaultSection = "global"

  c := new(Config)
  c.filename = filename
  c.file, err = goconfig.ReadConfigFile(c.filename)

  if c.file == nil {
    return nil, err
  }

  c.loglevel, _ = c.file.GetString("global", "loglevel")

  return c, nil
}

func (c *Config) LogLevel() string {
  return c.loglevel
}
