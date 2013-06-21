package engine

import (
  "time"

  . "launchpad.net/gocheck"
)

type BasicPollingRuleTest struct{}

var _ = Suite(&BasicPollingRuleTest{})

func (t *BasicPollingRuleTest) TestDontUpdate(c *C) {
  r := &BasicPollingRule{
    offset:    time.Now(),
    frequency: time.Hour,
  }

  c.Check(r.ShouldPoll(time.Now()), Equals, false)
}

func (t *BasicPollingRuleTest) TestUpdate(c *C) {
  r := &BasicPollingRule{
    offset:    time.Now().Add(-time.Hour),
    frequency: time.Second,
  }

  c.Check(r.ShouldPoll(time.Now()), Equals, false)
}
