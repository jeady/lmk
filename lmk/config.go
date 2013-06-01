package lmk

import (
  "github.com/msbranco/goconfig"
)

type Config struct {
  filename string
  file     *goconfig.ConfigFile

  loglevel string
  rules    []Rule
}

func NewConfig(filename string) (*Config, error) {
  var err error

  // Attempt to read the configuration file.
  c := new(Config)
  c.filename = filename
  c.file, err = goconfig.ReadConfigFile(c.filename)
  if c.file == nil {
    return nil, err
  }

  // Read the global section.
  c.loglevel, err = c.file.GetString("global", "loglevel")
  if err != nil {
    c.loglevel = "Notice"
  }

  // Parse the rules.
  log.Notice("Rules:")
  for _, name := range c.file.GetSections() {
    if name == "global" || name == "default" {
      continue
    }

    var url, sanity, trigger string
    valid, enabled := c.parse_rule_config(
      name,
      map[string]*string{
        "url":     &url,
        "sanity":  &sanity,
        "trigger": &trigger,
      })

    if valid && enabled {
      c.rules = append(c.rules, NewWebRule(name, url, sanity, trigger))
      log.Notice("%s: Loaded WebRule", name)
    } else {
      log.Notice("%s: Ignoring rule.", name)
    }
  }
  log.Notice("")

  return c, nil
}

// Given a rule name and string => *string map, parses the config file and
// fills in the *strings with their entries in the configuration file. Returns
// enabled and valid. Enabled is true iff the config file contains an entry
// setting the rule to enabled. Valid is true iff all entries in the required
// map are present in the config file.
func (c *Config) parse_rule_config(
  name string,
  required map[string]*string,
) (valid bool, enabled bool) {
  var err error
  valid = true

  if !c.file.HasSection(name) {
    log.Error("Could not find section '" + name + "'")
    valid = false
    enabled = false
    return
  }

  log.Debug("Parsing rule '%s'", name)

  // Determine whether or not the rule is enabled.
  enabled, err = c.file.GetBool(name, "enabled")
  if err != nil || !enabled {
    log.Notice("%s: not enabled.", name)
    enabled = false
  }

  // Parse the options for the rule.
  for opt, valp := range required {
    val, err := c.file.GetString(name, opt)

    if err == nil {
      *valp = val
    } else {
      log.Notice("%s: missing required option '%s'.", name, opt)
      valid = false
    }
  }

  // Output debugging notices about any unknown options
  opts, _ := c.file.GetOptions(name)
  for _, opt := range opts {
    if _, ok := required[opt]; !ok && opt != "enabled" {
      log.Notice("%s: ignoring unknown option '%s'.", name, opt)
    }
  }

  return
}

func (c *Config) LogLevel() string {
  return c.loglevel
}

func (c *Config) Rules() []Rule {
  return c.rules
}
