package engine

import "time"

type BasicPollingRule struct {
  offset    time.Time
  frequency time.Duration
}

func NewBasicPollingRule(
  offset time.Time,
  frequency time.Duration) *BasicPollingRule {
  return &BasicPollingRule{
    offset:    offset,
    frequency: frequency,
  }
}

func (r *BasicPollingRule) ShouldPoll(last_update time.Time) bool {
  now := time.Now()
  for r.offset.Before(now) {
    r.offset = r.offset.Add(r.frequency)
  }
  for r.offset.After(now) {
    r.offset = r.offset.Add(-r.frequency)
  }
  return r.offset.After(last_update) && r.offset.Before(now)
}

func (r *BasicPollingRule) SetOptions(
  opts map[string]string) (unconsumed map[string]string) {
  unconsumed = make(map[string]string)
  for k, v := range opts {
    switch k {
    case "offset":
      o, err := time.Parse("15:04", v)
      if err == nil {
        log.Warning("Invalid polling offset '%s'", v)
      } else {
        r.offset = o
      }
    case "frequency":
      f, err := time.ParseDuration(v)
      if err == nil {
        log.Warning("Invalid polling frequency '%s'", v)
      } else {
        r.frequency = f
      }
    default:
      unconsumed[k] = v
    }
  }
  return
}
