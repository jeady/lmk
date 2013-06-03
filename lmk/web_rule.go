package lmk

import "strconv"

type WebRule struct {
  name string

  url           string
  sanity_check  string
  trigger_check string

  case_sensitive bool
}

func NewWebRule(
  name string,
  url string,
  sanity_check string,
  trigger_check string,
  opts map[string]string) *WebRule {

  r := &WebRule{
    name:           name,
    url:            url,
    sanity_check:   sanity_check,
    trigger_check:  trigger_check,
    case_sensitive: false,
  }
  r.set_options(opts)
  return r
}

func (r *WebRule) Name() string {
  return r.name
}

// Sets optional parameters.
// option: default_value
// - case_sensitive: true
func (r *WebRule) set_options(opts map[string]string) {
  for k, v := range opts {
    switch k {
    case "case-sensitive":
      b, err := strconv.ParseBool(v)
      if err == nil {
        r.case_sensitive = b
      }
    default:
      log.Warning("WebRule: Unknown option '%s'", k)
    }
  }
}

func (r *WebRule) TestTriggered() (sane, triggered bool) {
  sane = false
  triggered = false
  return
}
