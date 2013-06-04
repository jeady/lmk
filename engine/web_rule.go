package engine

import (
  "io/ioutil"
  "net/http"
  "regexp"
)

type WebRule struct {
  name string

  // Required params.
  url           string
  sanity_check  string
  trigger_check string

  // Internal state.
  sanity_regex  *regexp.Regexp
  trigger_regex *regexp.Regexp
}

func NewWebRule(
  name string,
  url string,
  sanity_check string,
  trigger_check string) *WebRule {

  r := &WebRule{
    name:          name,
    url:           url,
    sanity_check:  sanity_check,
    trigger_check: trigger_check,
  }

  var err error
  r.sanity_regex, err = regexp.Compile(sanity_check)
  if err != nil {
    log.Warning("%s: Unable to compile sanity regex.", name)
    return nil
  }

  r.trigger_regex, err = regexp.Compile(trigger_check)
  if err != nil {
    log.Warning("%s: Unable to compile trigger regex.", name)
    return nil
  }
  return r
}

func (r *WebRule) Name() string {
  return r.name
}

func (r *WebRule) matches(content string, re *regexp.Regexp) bool {
  return re.Match([]byte(content))
}

func (r *WebRule) test_sane(page_content string) bool {
  return r.matches(page_content, r.sanity_regex)
}

func (r *WebRule) test_triggered(page_content string) bool {
  return r.matches(page_content, r.trigger_regex)
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
