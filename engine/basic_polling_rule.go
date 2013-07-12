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

func (r *BasicPollingRule) LastDeadline() time.Time {
  now := time.Now()
  deadline := r.offset

  for deadline.Before(now) {
    deadline = deadline.Add(r.frequency)
  }
  for deadline.After(now) {
    deadline = deadline.Add(-r.frequency)
  }

  return deadline
}

func (r *BasicPollingRule) NextDeadline() time.Time {
  deadline := r.LastDeadline()
  deadline = deadline.Add(r.frequency)

  return deadline
}

func (r *BasicPollingRule) ShouldPoll(last_update time.Time) bool {
  // Test to see if the last poll deadline was between the last update and now.
  // If so, it should be triggered ASAP.
  r.offset = r.LastDeadline()
  return r.offset.After(last_update) && r.offset.Before(time.Now())
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
