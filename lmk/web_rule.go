package lmk

type WebRule struct {
  name string

  url           string
  sanity_check  string
  trigger_check string
}

func NewWebRule(
  name string,
  url string,
  sanity_check string,
  trigger_check string) *WebRule {

  return &WebRule{
    name:          name,
    url:           url,
    sanity_check:  sanity_check,
    trigger_check: trigger_check,
  }
}

func (r *WebRule) Name() string {
  return r.name
}

func (r *WebRule) TestTriggered() (sane, triggered bool) {
  sane = false
  triggered = false
  return
}
