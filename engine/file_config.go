package engine

import (
  "errors"
  "path/filepath"

  "github.com/msbranco/goconfig"
)

// Convenience function to instantiate an Engine object using this config
// format.
func NewEngineFromFile(config_file string) *Engine {
  conf, err := NewFileConfig(config_file)

  if conf == nil {
    log.Error("Unable to load configuration file: " + err.Error())
    return nil
  }

  return NewEngine(conf)
}

type FileConfig struct {
  filename   string
  file       *goconfig.ConfigFile
  rules_file *goconfig.ConfigFile

  smtp_user string
  smtp_pass string
  smtp_host string

  recipient string
  loglevel  string
  rules     []Rule
}

func NewFileConfig(filename string) (*FileConfig, error) {
  var err error

  // Attempt to read the configuration file.
  c := new(FileConfig)
  c.filename = filename
  log.Debug("FileConfig filename: %s", c.filename)
  c.file, err = goconfig.ReadConfigFile(c.filename)
  if c.file == nil {
    return nil, err
  }

  // Make sure we can also read the rules file.
  rules_filename := c.get_config("rules", "rules.conf")
  if !filepath.IsAbs(rules_filename) {
    rules_filename = filepath.Join(filepath.Dir(c.filename), rules_filename)
  }
  log.Debug("Rules filename: %s", rules_filename)
  c.rules_file, err = goconfig.ReadConfigFile(rules_filename)
  if c.rules_file == nil {
    return nil, err
  }

  // Read the configuration.
  c.loglevel = c.get_config("loglevel", "Notice")
  c.smtp_user = c.get_config("smtp_user", "")
  c.smtp_pass = c.get_config("smtp_pass", "")
  c.smtp_host = c.get_config("smtp_host", "")
  c.recipient = c.get_config("recipient", "")

  // Parse the rules.
  log.Notice("Rules:")
  for _, name := range c.rules_file.GetSections() {
    if name == "default" {
      continue
    }

    var url, sanity, trigger string
    valid, enabled, opts := c.parse_rule_config(
      name,
      map[string]*string{
        "url":     &url,
        "sanity":  &sanity,
        "trigger": &trigger,
      },
      []string{"trigger-on-match"})

    if valid && enabled {
      r := NewWebRule(name, url, sanity, trigger)
      if r != nil {
        c.rules = append(c.rules, r)
        log.Notice("%s: Loaded WebRule", name)

        leftover := r.SetOptions(opts)
        if len(leftover) > 0 {
          log.Notice("%s: Unused options:", name)
          for k, v := range leftover {
            log.Notice("  %s=%s", k, v)
          }
        }
      } else {
        log.Error(name + ": Unable to load WebRule")
      }
    } else {
      log.Notice("%s: Ignoring rule.", name)
    }
  }
  log.Notice("")

  return c, nil
}

func (c *FileConfig) get_config(option, default_val string) string {
  result, err := c.file.GetString("settings", option)
  if err != nil {
    return default_val
  }
  return result
}

// Given a rule name and string => *string map, parses the config file and
// fills in the *strings with their entries in the configuration file. Returns
// enabled and valid. Enabled is true iff the config file contains an entry
// setting the rule to enabled. Valid is true iff all entries in the required
// map are present in the config file.
func (c *FileConfig) parse_rule_config(
  name string,
  required map[string]*string,
  optional []string,
) (valid bool, enabled bool, options map[string]string) {
  var err error
  valid = true
  options = make(map[string]string)

  if !c.rules_file.HasSection(name) {
    log.Error("Could not find section '" + name + "'")
    valid = false
    enabled = false
    return
  }

  log.Debug("Parsing rule '%s'", name)

  // Determine whether or not the rule is enabled.
  enabled, err = c.rules_file.GetBool(name, "enabled")
  if err != nil || !enabled {
    log.Notice("%s: not enabled.", name)
    enabled = false
  }

  // Parse the options for the rule.
  for opt, valp := range required {
    val, err := c.rules_file.GetString(name, opt)

    if err == nil {
      *valp = val
      log.Debug("  R %s=%s", opt, val)
    } else {
      log.Notice("%s: missing required option '%s'.", name, opt)
      valid = false
    }
  }

  for _, opt := range optional {
    val, err := c.rules_file.GetString(name, opt)

    if err == nil {
      options[opt] = val
      log.Debug("  O %s=%s", opt, val)
    }
  }

  // Output debugging notices about any unknown options
  opts, _ := c.rules_file.GetOptions(name)
  for _, opt := range opts {
    _, in_required := required[opt]
    _, in_optional := options[opt]
    if !in_required && !in_optional && opt != "enabled" {
      log.Notice("%s: ignoring unknown option '%s'.", name, opt)
    }
  }

  return
}

func (c *FileConfig) SmtpConfig() (user, pass, host string) {
  user = c.smtp_user
  pass = c.smtp_pass
  host = c.smtp_host
  return
}

func (c *FileConfig) DefaultNotifier() (Notifier, error) {
  if len(c.smtp_user) == 0 || len(c.smtp_host) == 0 {
    return nil, errors.New(
      "Could not create default notifier: SMTP credentials not specified.")
  }

  return NewSmtpNotifier(c.smtp_user, c.smtp_pass, c.smtp_host), nil
}

func (c *FileConfig) DefaultNotificationRecipient() string {
  return c.recipient
}

func (c *FileConfig) LogLevel() string {
  return c.loglevel
}

func (c *FileConfig) Rules() []Rule {
  return c.rules
}
