package engine

import (
  "errors"
  "io/ioutil"
  "os"
  "time"

  "code.google.com/p/gomock/gomock"
  . "launchpad.net/gocheck"
)

type EngineTest struct {
  config_path string
  cleanup     func()
}

var _ = Suite(&EngineTest{})

func build_config(
  rules string,
  c *C) (conf string, cleanup func()) {
  var fconf, rconf *os.File
  var err error
  var conf_test Config

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
  fconf.WriteString(
    "[Settings]\nloglevel=debug\nrules=" + rconf.Name() + "\n")
  fconf.Sync()
  rconf.WriteString(rules + "\n")
  rconf.Sync()

  conf_test, err = NewFileConfig(fconf.Name())
  c.Assert(conf_test, Not(Equals), nil)
  c.Assert(err, Equals, nil)

  conf = fconf.Name()
  return
}

func (t *EngineTest) TestGoodConstruction(c *C) {
  conf, cleanup := build_config("", c)
  defer cleanup()

  e := NewEngineFromFile(conf)
  c.Check(e, Not(Equals), (*Engine)(nil))
}

func (t *EngineTest) TestBadConfigFile(c *C) {
  e := NewEngineFromFile("/dev/null/nothere")
  c.Check(e, Equals, (*Engine)(nil))
}

func (t *EngineTest) TestRunRequiresNotifier(c *C) {
  m := gomock.NewController(c)
  defer m.Finish()

  r := NewMockRule(m)
  r.EXPECT().TestTriggered().Times(0)
  r.EXPECT().Name().Return("Test Rule").AnyTimes()

  conf := NewMockConfig(m)
  conf.EXPECT().DefaultNotifier().Return(nil, errors.New("")).AnyTimes()
  conf.EXPECT().LogLevel().Return("DEBUG").AnyTimes()
  conf.EXPECT().Rules().Return([]Rule{r}).AnyTimes()

  e := NewEngine(conf)
  e.Run(r)
}

func (t *EngineTest) TestNotifiesOnInsane(c *C) {
  m := gomock.NewController(c)
  defer m.Finish()

  r := NewMockRule(m)
  r.EXPECT().TestTriggered().Return(false, false).Times(1)
  r.EXPECT().Name().Return("Test Rule").AnyTimes()

  n := NewMockNotifier(m)
  n.EXPECT().Notify(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

  conf := NewMockConfig(m)
  conf.EXPECT().DefaultNotifier().Return(n, nil).AnyTimes()
  conf.EXPECT().DefaultNotificationRecipient().Return("").AnyTimes()
  conf.EXPECT().LogLevel().Return("DEBUG").AnyTimes()
  conf.EXPECT().Rules().Return([]Rule{r}).AnyTimes()

  e := NewEngine(conf)
  e.Run(r)
}

func (t *EngineTest) TestNotifiesOnTrigger(c *C) {
  m := gomock.NewController(c)
  defer m.Finish()

  r := NewMockRule(m)
  r.EXPECT().TestTriggered().Return(true, true).Times(1)
  r.EXPECT().Name().Return("Test Rule").AnyTimes()

  n := NewMockNotifier(m)
  n.EXPECT().Notify(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

  conf := NewMockConfig(m)
  conf.EXPECT().DefaultNotifier().Return(n, nil).AnyTimes()
  conf.EXPECT().DefaultNotificationRecipient().Return("").AnyTimes()
  conf.EXPECT().LogLevel().Return("DEBUG").AnyTimes()
  conf.EXPECT().Rules().Return([]Rule{r}).AnyTimes()

  e := NewEngine(conf)
  e.Run(r)
}

func (t *EngineTest) TestNoNotifyOnDormant(c *C) {
  m := gomock.NewController(c)
  defer m.Finish()

  r := NewMockRule(m)
  r.EXPECT().TestTriggered().Return(true, false).Times(1)
  r.EXPECT().Name().Return("Test Rule").AnyTimes()

  n := NewMockNotifier(m)
  n.EXPECT().Notify(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

  conf := NewMockConfig(m)
  conf.EXPECT().DefaultNotifier().Return(n, nil).AnyTimes()
  conf.EXPECT().DefaultNotificationRecipient().Return("").AnyTimes()
  conf.EXPECT().LogLevel().Return("DEBUG").AnyTimes()
  conf.EXPECT().Rules().Return([]Rule{r}).AnyTimes()

  e := NewEngine(conf)
  e.Run(r)
}

func (t *EngineTest) TestChoosesRulesToPoll(c *C) {
  m := gomock.NewController(c)
  defer m.Finish()

  r1 := NewMockRule(m)
  r1.EXPECT().Name().Return("Non-pollable 1").AnyTimes()
  r2 := NewMockRule(m)
  r2.EXPECT().Name().Return("Non-pollable 2").AnyTimes()
  r3 := NewMockPollingRule(m)
  r3.EXPECT().Name().Return("Pollable 1 on").AnyTimes()
  r3.EXPECT().ShouldPoll(gomock.Any()).Return(true).Times(1)
  r3.EXPECT().NextDeadline().Return(time.Now()).AnyTimes()
  r4 := NewMockPollingRule(m)
  r4.EXPECT().Name().Return("Pollable 2 on").AnyTimes()
  r4.EXPECT().ShouldPoll(gomock.Any()).Return(true).Times(1)
  r4.EXPECT().NextDeadline().Return(time.Now()).AnyTimes()
  r5 := NewMockPollingRule(m)
  r5.EXPECT().Name().Return("Pollable 3 off").AnyTimes()
  r5.EXPECT().ShouldPoll(gomock.Any()).Return(false).Times(1)
  r5.EXPECT().NextDeadline().Return(time.Now()).AnyTimes()
  r6 := NewMockPollingRule(m)
  r6.EXPECT().Name().Return("Pollable 4 off").AnyTimes()
  r6.EXPECT().ShouldPoll(gomock.Any()).Return(false).Times(1)
  r6.EXPECT().NextDeadline().Return(time.Now()).AnyTimes()
  rules := []Rule{r1, r2, r3, r4, r5, r6}

  n := NewMockNotifier(m)
  n.EXPECT().Notify(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

  conf := NewMockConfig(m)
  conf.EXPECT().DefaultNotifier().Return(n, nil).AnyTimes()
  conf.EXPECT().DefaultNotificationRecipient().Return("").AnyTimes()
  conf.EXPECT().LogLevel().Return("DEBUG").AnyTimes()
  conf.EXPECT().Rules().Return(rules).AnyTimes()

  e := NewEngine(conf)
  r := e.RulesToPoll(time.Now())
  c.Assert(r, HasLen, 2)
  c.Assert(r[0], Equals, r3)
  c.Assert(r[1], Equals, r4)
}

func (t *EngineTest) TestSetsUrlFetcher(c *C) {
  m := gomock.NewController(c)
  defer m.Finish()

  u := NewMockUrlFetcher(m)

  r1 := NewMockUrlFetchingRule(m)
  r1.EXPECT().Name().Return("Non-pollable 1").AnyTimes()
  r1.EXPECT().SetUrlFetcher(u).Times(1)
  r2 := NewMockUrlFetchingRule(m)
  r2.EXPECT().Name().Return("Non-pollable 2").AnyTimes()
  r2.EXPECT().SetUrlFetcher(u).Times(1)
  rules := []Rule{r1, r2}

  n := NewMockNotifier(m)
  n.EXPECT().Notify(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

  conf := NewMockConfig(m)
  conf.EXPECT().DefaultNotifier().Return(n, nil).AnyTimes()
  conf.EXPECT().DefaultNotificationRecipient().Return("").AnyTimes()
  conf.EXPECT().LogLevel().Return("DEBUG").AnyTimes()
  conf.EXPECT().Rules().Return(rules).AnyTimes()

  e := NewEngine(conf)
  e.SetUrlFetcher(u)
}
