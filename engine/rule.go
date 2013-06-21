package engine

import "time"

type Rule interface {
  Name() string
  TestTriggered() (sane, triggered bool)

  SetOptions(opts map[string]string) (unconsumed map[string]string)
}

type UrlFetchingRule interface {
  SetUrlFetcher(f UrlFetcher) UrlFetcher
}

type PollingRule interface {
  ShouldPoll(last_update time.Time) bool
}
