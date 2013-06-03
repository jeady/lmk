package lmk

import (
  "io/ioutil"
  "net/http"
  "strconv"
  "strings"
)

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

func (r *WebRule) contains(content, substr string) bool {
  if !r.case_sensitive {
    content = strings.ToLower(content)
    substr = strings.ToLower(substr)
  }
  return strings.Contains(content, substr)
}

func (r *WebRule) test_sane(page_content string) bool {
  return r.contains(page_content, r.sanity_check)
}

func (r *WebRule) test_triggered(page_content string) bool {
  return r.contains(page_content, r.trigger_check)
}

func (r *WebRule) TestTriggered() (sane, triggered bool) {
  resp, err := http.Get(r.url)
  if err != nil {
    log.Info("WebRule '%s' not sane because of error GET'ing page: %s",
      r.Name(),
      err.Error())
    sane = false
    triggered = false
    return
  }
  defer resp.Body.Close()

  page_bytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Info("WebRule '%s' not sane because of error reading body: ",
      r.Name(),
      err)
    sane = false
    triggered = false
    return
  }

  page := string(page_bytes)

  sane = r.test_sane(page)
  triggered = r.test_triggered(page)
  return
}
