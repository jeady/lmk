package engine

import "time"

type Rule interface {
  Name() string
  TestTriggered() (sane, triggered bool)

  SetOptions(opts map[string]string) (unconsumed map[string]string)
}

type UrlFetchingRule interface {
  Rule
  SetUrlFetcher(f UrlFetcher) UrlFetcher
}

type PollingRule interface {
  Rule
  ShouldPoll(last_update time.Time) bool
  LastDeadline() time.Time
  NextDeadline() time.Time
}
