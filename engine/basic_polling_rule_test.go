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

func (t *BasicPollingRuleTest) TestLastDeadline(c *C) {
  r := &BasicPollingRule{
    offset:    time.Date(2000, 1, 1, 5, 30, 0, 0, time.UTC),
    frequency: time.Hour,
  }
  d := r.LastDeadlineBefore(time.Date(2005, 2, 3, 4, 5, 6, 7, time.UTC))
  c.Check(d, Equals, time.Date(2005, 2, 3, 3, 30, 0, 0, time.UTC))

  r = &BasicPollingRule{
    offset:    time.Date(2000, 1, 1, 5, 30, 0, 0, time.UTC),
    frequency: time.Hour,
  }
  d = r.LastDeadlineBefore(time.Date(2005, 2, 3, 0, 5, 6, 7, time.UTC))
  c.Check(d, Equals, time.Date(2005, 2, 2, 23, 30, 0, 0, time.UTC))

  r = &BasicPollingRule{
    offset:    time.Date(2000, 1, 1, 5, 30, 0, 0, time.UTC),
    frequency: time.Hour,
  }
  d = r.LastDeadlineBefore(time.Date(2005, 2, 3, 23, 35, 6, 7, time.UTC))
  c.Check(d, Equals, time.Date(2005, 2, 3, 23, 30, 0, 0, time.UTC))
}

func (t *BasicPollingRuleTest) TestNextDeadline(c *C) {
  r := &BasicPollingRule{
    offset:    time.Date(2000, 1, 1, 5, 30, 0, 0, time.UTC),
    frequency: time.Hour,
  }
  d := r.NextDeadlineAfter(time.Date(2005, 2, 3, 4, 5, 6, 7, time.UTC))
  c.Check(d, Equals, time.Date(2005, 2, 3, 4, 30, 0, 0, time.UTC))

  r = &BasicPollingRule{
    offset:    time.Date(2000, 1, 1, 5, 30, 0, 0, time.UTC),
    frequency: time.Hour,
  }
  d = r.NextDeadlineAfter(time.Date(2005, 2, 3, 0, 5, 6, 7, time.UTC))
  c.Check(d, Equals, time.Date(2005, 2, 3, 0, 30, 0, 0, time.UTC))

  r = &BasicPollingRule{
    offset:    time.Date(2000, 1, 1, 5, 30, 0, 0, time.UTC),
    frequency: time.Hour,
  }
  d = r.NextDeadlineAfter(time.Date(2005, 2, 3, 23, 35, 6, 7, time.UTC))
  c.Check(d, Equals, time.Date(2005, 2, 4, 0, 30, 0, 0, time.UTC))
}

func (t *BasicPollingRuleTest) TestParseGoodOptions(c *C) {
  r := &BasicPollingRule{
    offset:    time.Now(),
    frequency: time.Second,
  }
  u := r.SetOptions(map[string]string{
    "offset": "01:02",
  })
  c.Check(u, HasLen, 0)
  c.Check(r.offset.Hour(), Equals, 1)
  c.Check(r.offset.Minute(), Equals, 2)

  r = &BasicPollingRule{
    offset:    time.Now(),
    frequency: time.Second,
  }
  u = r.SetOptions(map[string]string{
    "offset": "03:04",
  })
  c.Check(u, HasLen, 0)
  c.Check(r.offset.Hour(), Equals, 3)
  c.Check(r.offset.Minute(), Equals, 4)

  r = &BasicPollingRule{
    offset:    time.Now(),
    frequency: time.Second,
  }
  u = r.SetOptions(map[string]string{
    "frequency": "1h",
  })
  c.Check(u, HasLen, 0)
  c.Check(r.frequency, Equals, time.Hour)

  r = &BasicPollingRule{
    offset:    time.Now(),
    frequency: time.Second,
  }
  u = r.SetOptions(map[string]string{
    "frequency": "5m",
  })
  c.Check(u, HasLen, 0)
  c.Check(r.frequency, Equals, 5*time.Minute)
}

func (t *BasicPollingRuleTest) TestParseBadOptions(c *C) {
  start_time := time.Now()
  r := &BasicPollingRule{
    offset:    start_time,
    frequency: time.Second,
  }
  u := r.SetOptions(map[string]string{
    "frequency": "-1h",
  })
  c.Check(u, HasLen, 0)
  c.Check(r.frequency, Equals, time.Second)

  r = &BasicPollingRule{
    offset:    start_time,
    frequency: time.Second,
  }
  u = r.SetOptions(map[string]string{
    "frequency": "h",
  })
  c.Check(u, HasLen, 0)
  c.Check(r.frequency, Equals, time.Second)

  r = &BasicPollingRule{
    offset:    start_time,
    frequency: time.Second,
  }
  u = r.SetOptions(map[string]string{
    "offset": "1:2:3",
  })
  c.Check(u, HasLen, 0)
  c.Check(r.offset, Equals, start_time)
}
