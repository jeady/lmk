package engine

import (
  "regexp"
  "strconv"
  "time"
)

type WebRule struct {
  BasicPollingRule

  name string

  // Required params.
  url           string
  sanity_check  string
  trigger_check string

  // Optional params.
  trigger_on_match bool

  // Internal state.
  sanity_regex  *regexp.Regexp
  trigger_regex *regexp.Regexp
  url_fetcher   UrlFetcher
}

// sanity_check and trigger_check should be valid regular expressions.
// Optional params:
//   * trigger-on-match: (default true)
//     If false, trigger when trigger_check does not match rather than when
//     trigger_check matches.
func NewWebRule(
  name string,
  url string,
  sanity_check string,
  trigger_check string) *WebRule {

  r := &WebRule{
    name:             name,
    url:              url,
    sanity_check:     sanity_check,
    trigger_check:    trigger_check,
    trigger_on_match: true,
    url_fetcher:      &NetHttpFetcher{},
    BasicPollingRule: *NewBasicPollingRule(
      time.Date(0, 0, 0, 24, 0, 0, 0, time.UTC),
      24*time.Hour),
  }

  var err error
  r.sanity_regex, err = regexp.Compile(sanity_check)
  if err != nil {
    log.Warning("%s: Unable to compile sanity regex: %s", name, err.Error())
    return nil
  }

  r.trigger_regex, err = regexp.Compile(trigger_check)
  if err != nil {
    log.Warning("%s: Unable to compile trigger regex: %s", name, err.Error())
    return nil
  }
  return r
}

func (r *WebRule) SetUrlFetcher(f UrlFetcher) UrlFetcher {
  old := r.url_fetcher
  r.url_fetcher = f
  return old
}

func (r *WebRule) Name() string {
  return r.name
}

func (r *WebRule) SetOptions(
  opts map[string]string) (unconsumed map[string]string) {
  opts = r.BasicPollingRule.SetOptions(opts)
  unconsumed = make(map[string]string)

  for name, value := range opts {
    switch name {
    case "trigger-on-match":
      b, err := strconv.ParseBool(value)
      if err == nil {
        r.trigger_on_match = b
      }
    default:
      unconsumed[name] = value
    }
  }
  return
}

func (r *WebRule) matches(content string, re *regexp.Regexp) bool {
  return re.Match([]byte(content))
}

func (r *WebRule) test_sane(page_content string) bool {
  return r.matches(page_content, r.sanity_regex)
}

func (r *WebRule) test_triggered(page_content string) bool {
  return r.matches(page_content, r.trigger_regex) == r.trigger_on_match
}

func (r *WebRule) TestTriggered() (sane, triggered bool) {
  page, err := r.url_fetcher.Get(r.url)
  if err != nil {
    log.Info("WebRule '%s' not sane because of error GET'ing page: %s",
      r.Name(),
      err.Error())
    sane = false
    triggered = false
    return
  }

  sane = r.test_sane(page)
  triggered = r.test_triggered(page)
  return
}
