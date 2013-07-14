package engine

import (
  "io/ioutil"
  go_log "log"
  "os"

  "github.com/jeady/go-logging"
  . "launchpad.net/gocheck"
)

type FileConfigTest struct{}

var _ = Suite(&FileConfigTest{})

func load_config(
  config string,
  rules string,
  c *C) (conf *FileConfig, cleanup func()) {
  var fconf, rconf *os.File
  var err error

  fconf, err = ioutil.TempFile("", "")
  if err != nil {
    panic(err)
  }
  rconf, err = ioutil.TempFile("", "")
  if err != nil {
    panic(err)
  }
  cleanup = func() {
    fconf.Close()
    rconf.Close()
  }

  // Trailing \n is required by go-config issue #3.
  fconf.WriteString("[Settings]\n" + config + "\n" + "rules=" + rconf.Name() + "\n")
  fconf.Sync()
  rconf.WriteString(rules + "\n")
  rconf.Sync()

  conf, err = NewFileConfig(fconf.Name())
  c.Assert(conf, Not(Equals), nil)
  c.Assert(err, Equals, nil)
  return
}

func (t *FileConfigTest) TestParseLogLevel(c *C) {

  // Wrapper function.
  test_log_level := func(in string, expected string) {
    conf, cleanup := load_config(in, "", c)
    defer cleanup()

    c.Check(conf.LogLevel(), Equals, expected)
  }

  // loglevel=debug
  test_log_level("loglevel=debug", "debug")
  test_log_level("loglevel=DEBUG", "DEBUG")
  test_log_level("loglevel=Info", "Info")
  test_log_level("garbage=trash", "Notice")
}

func (t *FileConfigTest) TestParseRuleFileConfig(c *C) {
  memlog := new(TestingLogger)
  logging.SetBackend(memlog)
  defer logging.SetBackend(
    logging.NewLogBackend(os.Stderr, "", go_log.LstdFlags))

  logging.SetLevel(logging.DEBUG, log.Module)
  var foo, bar, baz string

  test_parse_rule_config := func(rules string, check func(conf *FileConfig)) {
    foo = ""
    bar = ""
    baz = ""
    memlog.Reset()
    conf, cleanup := load_config("loglevel=info", rules, c)
    defer cleanup()

    check(conf)
  }

  // All present and account for.
  test_parse_rule_config(
    "[RuleA]\nfoo=abc\nbar=def\nbaz=123\nenabled=true",
    func(conf *FileConfig) {
      valid, enabled, opts := conf.parse_rule_config(
        "RuleA",
        map[string]*string{
          "foo": &foo,
          "bar": &bar,
          "baz": &baz,
        },
        []string{})
      c.Check(valid, Equals, true)
      c.Check(enabled, Equals, true)
      c.Check(foo, Equals, "abc")
      c.Check(bar, Equals, "def")
      c.Check(baz, Equals, "123")
      c.Check(len(opts), Equals, 0)
    })

  // Bunch of extra fields.
  test_parse_rule_config(
    "[RuleA]\nfoo=abc\nbar=def\nbaz=123\nenabled=true",
    func(conf *FileConfig) {
      valid, enabled, opts := conf.parse_rule_config(
        "RuleA",
        map[string]*string{},
        []string{})
      c.Check(valid, Equals, true)
      c.Check(enabled, Equals, true)
      c.Check(memlog.Logs(), Matches, `[\s\S]*foo[\s\S]*`)
      c.Check(memlog.Logs(), Matches, `[\s\S]*bar[\s\S]*`)
      c.Check(memlog.Logs(), Matches, `[\s\S]*baz[\s\S]*`)
      c.Check(len(opts), Equals, 0)
    })

  // Missing required field.
  test_parse_rule_config(
    "[RuleA]\nenabled=true",
    func(conf *FileConfig) {
      valid, enabled, opts := conf.parse_rule_config(
        "RuleA",
        map[string]*string{
          "foo": &bar,
        },
        []string{})
      c.Check(valid, Equals, false)
      c.Check(enabled, Equals, true)
      c.Check(foo, Equals, "")
      c.Check(len(opts), Equals, 0)
    })

  // Disabled due to missing enabled field.
  test_parse_rule_config(
    "[RuleA]\nfoo=abc\nbar=def\nbaz=123",
    func(conf *FileConfig) {
      valid, enabled, opts := conf.parse_rule_config(
        "RuleA",
        map[string]*string{
          "foo": &foo,
          "bar": &bar,
          "baz": &baz,
        },
        []string{})
      c.Check(valid, Equals, true)
      c.Check(enabled, Equals, false)
      c.Check(foo, Equals, "abc")
      c.Check(bar, Equals, "def")
      c.Check(baz, Equals, "123")
      c.Check(len(opts), Equals, 0)
    })

  // Disabled due to enabled=false.
  test_parse_rule_config(
    "[RuleA]\nfoo=abc\nbar=def\nbaz=123\nenabled=false",
    func(conf *FileConfig) {
      valid, enabled, opts := conf.parse_rule_config(
        "RuleA",
        map[string]*string{
          "foo": &foo,
          "bar": &bar,
          "baz": &baz,
        },
        []string{})
      c.Check(valid, Equals, true)
      c.Check(enabled, Equals, false)
      c.Check(foo, Equals, "abc")
      c.Check(bar, Equals, "def")
      c.Check(baz, Equals, "123")
      c.Check(len(opts), Equals, 0)
    })

  // Mising section.
  test_parse_rule_config(
    "[RuleB]\nfoo=abc\nbar=def\nbaz=123\nenabled=true",
    func(conf *FileConfig) {
      valid, enabled, opts := conf.parse_rule_config(
        "RuleA",
        map[string]*string{
          "foo": &foo,
          "bar": &bar,
          "baz": &baz,
        },
        []string{})
      c.Check(valid, Equals, false)
      c.Check(enabled, Equals, false)
      c.Check(foo, Equals, "")
      c.Check(bar, Equals, "")
      c.Check(baz, Equals, "")
      c.Check(len(opts), Equals, 0)
    })

  // Optional fields.
  test_parse_rule_config(
    "[RuleA]\nfoo=abc\nbar=def\nbaz=123\nenabled=true\nhello=goodbye",
    func(conf *FileConfig) {
      valid, enabled, opts := conf.parse_rule_config(
        "RuleA",
        map[string]*string{
          "foo": &foo,
          "bar": &bar,
          "baz": &baz,
        },
        []string{"hello", "world"})
      c.Check(valid, Equals, true)
      c.Check(enabled, Equals, true)
      c.Check(foo, Equals, "abc")
      c.Check(bar, Equals, "def")
      c.Check(baz, Equals, "123")
      c.Check(len(opts), Equals, 1)
      c.Check(opts["hello"], Equals, "goodbye")
    })

}

func (t *FileConfigTest) TestFileConfigFileDoesNotExist(c *C) {
  dir, err := ioutil.TempDir("", "")
  c.Check(err, Equals, nil)

  fname := dir + "/does/not/exist"
  _, err = os.Stat(fname)
  c.Assert(err, Not(Equals), nil)

  conf, err := NewFileConfig(fname)
  c.Check(conf, Equals, (*FileConfig)(nil))
  c.Check(err, Not(Equals), nil)
}

func (t *FileConfigTest) TestLoadsWebRules(c *C) {
  // Valid rule.
  conf, cleanup := load_config(
    "loglevel=debug",
    "[the rule]\n"+
      "url=http://google.com/\n"+
      "sanity=google\n"+
      "trigger=Mordor\n"+
      "enabled=true",
    c)
  defer func(c func()) { c() }(cleanup)

  c.Assert(len(conf.Rules()), Equals, 1)
  c.Check(conf.Rules()[0].Name(), Equals, "the rule")

  // Invalid rules.
  conf, cleanup = load_config(
    "loglevel=debug",
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
  defer func(c func()) { c() }(cleanup)

  c.Assert(len(conf.Rules()), Equals, 0)
}

// Tests loading multiple rules, rules of different types.
func (t *FileConfigTest) TestLoadsRules(c *C) {
  conf, cleanup := load_config(
    "loglevel=debug",
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
  defer func(c func()) { c() }(cleanup)

  c.Assert(len(conf.Rules()), Equals, 2)
  c.Check(conf.Rules()[0].Name(), Equals, "the rule")
  c.Check(conf.Rules()[1].Name(), Equals, "the rule 2")

  conf, cleanup = load_config(
    "loglevel=debug",
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
  defer func(c func()) { c() }(cleanup)

  c.Assert(len(conf.Rules()), Equals, 1)
  c.Check(conf.Rules()[0].Name(), Equals, "the rule 2")
}
