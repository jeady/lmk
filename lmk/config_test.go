package lmk

import (
  "io/ioutil"
  go_log "log"
  "os"

  "github.com/op/go-logging"
  . "launchpad.net/gocheck"
)

type ConfigTest struct{}

var _ = Suite(&ConfigTest{})

func load_config(config_contents string, c *C) (conf *Config, f *os.File) {
  f, err := ioutil.TempFile("", "")
  if err != nil {
    panic(err)
  }

  // Trailing \n is required by go-config issue #3.
  f.WriteString(config_contents + "\n")
  f.Sync()

  conf, err = NewConfig(f.Name())
  c.Assert(conf, Not(Equals), nil)
  c.Assert(err, Equals, nil)
  return
}

func (t *ConfigTest) TestParseLogLevel(c *C) {

  // Wrapper function.
  test_log_level := func(in string, expected string) {
    conf, f := load_config("[Global]\n"+in, c)
    defer func(f *os.File) { f.Close() }(f)

    c.Check(conf.LogLevel(), Equals, expected)
  }

  // loglevel=debug
  test_log_level("loglevel=debug", "debug")
  test_log_level("loglevel=DEBUG", "DEBUG")
  test_log_level("loglevel=Info", "Info")
  test_log_level("garbage=trash", "Notice")
}

func (t *ConfigTest) TestParseRuleConfig(c *C) {
  memlog := new(TestingLogger)
  logging.SetBackend(memlog)
  defer logging.SetBackend(
    logging.NewLogBackend(os.Stderr, "", go_log.LstdFlags))

  logging.SetLevel(logging.DEBUG, log.Module)
  var foo, bar, baz string

  test_parse_rule_config := func(config string, check func(conf *Config)) {
    foo = ""
    bar = ""
    baz = ""
    memlog.Reset()
    conf, f := load_config("[Global]\nloglevel=info\n\n"+config, c)
    defer func(f *os.File) { f.Close() }(f)

    check(conf)
  }

  // All present and account for.
  test_parse_rule_config(
    "[RuleA]\nfoo=abc\nbar=def\nbaz=123\nenabled=true",
    func(conf *Config) {
      valid, enabled := conf.parse_rule_config("RuleA", map[string]*string{
        "foo": &foo,
        "bar": &bar,
        "baz": &baz,
      })
      c.Check(valid, Equals, true)
      c.Check(enabled, Equals, true)
      c.Check(foo, Equals, "abc")
      c.Check(bar, Equals, "def")
      c.Check(baz, Equals, "123")
    })

  // Bunch of extra fields.
  test_parse_rule_config(
    "[RuleA]\nfoo=abc\nbar=def\nbaz=123\nenabled=true",
    func(conf *Config) {
      valid, enabled := conf.parse_rule_config("RuleA", map[string]*string{})
      c.Check(valid, Equals, true)
      c.Check(enabled, Equals, true)
      c.Check(memlog.Logs(), Matches, `[\s\S]*foo[\s\S]*`)
      c.Check(memlog.Logs(), Matches, `[\s\S]*bar[\s\S]*`)
      c.Check(memlog.Logs(), Matches, `[\s\S]*baz[\s\S]*`)
    })

  // Missing required field.
  test_parse_rule_config(
    "[RuleA]\nenabled=true",
    func(conf *Config) {
      valid, enabled := conf.parse_rule_config("RuleA", map[string]*string{
        "foo": &bar,
      })
      c.Check(valid, Equals, false)
      c.Check(enabled, Equals, true)
      c.Check(foo, Equals, "")
    })

  // Disabled due to missing enabled field.
  test_parse_rule_config(
    "[RuleA]\nfoo=abc\nbar=def\nbaz=123",
    func(conf *Config) {
      valid, enabled := conf.parse_rule_config("RuleA", map[string]*string{
        "foo": &foo,
        "bar": &bar,
        "baz": &baz,
      })
      c.Check(valid, Equals, true)
      c.Check(enabled, Equals, false)
      c.Check(foo, Equals, "abc")
      c.Check(bar, Equals, "def")
      c.Check(baz, Equals, "123")
    })

  // Disabled due to enabled=false.
  test_parse_rule_config(
    "[RuleA]\nfoo=abc\nbar=def\nbaz=123\nenabled=false",
    func(conf *Config) {
      valid, enabled := conf.parse_rule_config("RuleA", map[string]*string{
        "foo": &foo,
        "bar": &bar,
        "baz": &baz,
      })
      c.Check(valid, Equals, true)
      c.Check(enabled, Equals, false)
      c.Check(foo, Equals, "abc")
      c.Check(bar, Equals, "def")
      c.Check(baz, Equals, "123")
    })

  // Mising section.
  test_parse_rule_config(
    "[RuleB]\nfoo=abc\nbar=def\nbaz=123\nenabled=true",
    func(conf *Config) {
      valid, enabled := conf.parse_rule_config("RuleA", map[string]*string{
        "foo": &foo,
        "bar": &bar,
        "baz": &baz,
      })
      c.Check(valid, Equals, false)
      c.Check(enabled, Equals, false)
      c.Check(foo, Equals, "")
      c.Check(bar, Equals, "")
      c.Check(baz, Equals, "")
    })
}

func (t *ConfigTest) TestConfigFileDoesNotExist(c *C) {
  dir, err := ioutil.TempDir("", "")
  c.Check(err, Equals, nil)

  fname := dir + "/does/not/exist"
  _, err = os.Stat(fname)
  c.Assert(err, Not(Equals), nil)

  conf, err := NewConfig(fname)
  c.Check(conf, Equals, (*Config)(nil))
  c.Check(err, Not(Equals), nil)
}

func (t *ConfigTest) TestLoadsWebRules(c *C) {
  // Valid rule.
  conf, f := load_config(
    "[Global]\n"+
      "loglevel=debug\n"+
      "\n"+
      "[the rule]\n"+
      "url=http://google.com/\n"+
      "sanity=google\n"+
      "trigger=Mordor\n"+
      "enabled=true",
    c)
  defer func(f *os.File) { f.Close() }(f)

  c.Assert(len(conf.Rules()), Equals, 1)
  c.Check(conf.Rules()[0].Name(), Equals, "the rule")

  // Invalid rules.
  conf, f = load_config(
    "[Global]\n"+
      "loglevel=debug\n"+
      "\n"+
      "[the rule]\n"+
      "sanity=google\n"+
      "trigger=Mordor\n"+
      "enabled=true\n"+
      "\n"+
      "[the rule 2]\n"+
      "url=http://google.com/\n"+
      "trigger=Mordor\n"+
      "enabled=true\n"+
      "\n"+
      "[the rule 3]\n"+
      "url=http://google.com/\n"+
      "sanity=google\n"+
      "enabled=true",
    c)
  defer func(f *os.File) { f.Close() }(f)

  c.Assert(len(conf.Rules()), Equals, 0)
}

// Tests loading multiple rules, rules of different types.
func (t *ConfigTest) TestLoadsRules(c *C) {
  conf, f := load_config(
    "[Global]\n"+
      "loglevel=debug\n"+
      "\n"+
      "[the rule]\n"+
      "url=http://google.com/\n"+
      "sanity=google\n"+
      "trigger=Mordor\n"+
      "enabled=true\n"+
      "\n"+
      "[the rule 2]\n"+
      "url=http://google.com/\n"+
      "sanity=google\n"+
      "trigger=Mordor\n"+
      "enabled=true",
    c)
  defer func(f *os.File) { f.Close() }(f)

  c.Assert(len(conf.Rules()), Equals, 2)
  c.Check(conf.Rules()[0].Name(), Equals, "the rule")
  c.Check(conf.Rules()[1].Name(), Equals, "the rule 2")

  conf, f = load_config(
    "[Global]\n"+
      "loglevel=debug\n"+
      "\n"+
      "[the rule]\n"+
      "url=http://google.com/\n"+
      "sanity=google\n"+
      "trigger=Mordor\n"+
      "enabled=false\n"+
      "\n"+
      "[the rule 2]\n"+
      "url=http://google.com/\n"+
      "sanity=google\n"+
      "trigger=Mordor\n"+
      "enabled=true",
    c)
  defer func(f *os.File) { f.Close() }(f)

  c.Assert(len(conf.Rules()), Equals, 1)
  c.Check(conf.Rules()[0].Name(), Equals, "the rule 2")
}
