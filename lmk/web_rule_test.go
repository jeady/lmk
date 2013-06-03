package lmk

import (
  "fmt"
  go_log "log"
  "net/http"
  "net/http/httptest"
  "os"

  "github.com/op/go-logging"
  . "launchpad.net/gocheck"
)

type WebRuleTest struct{}

var _ = Suite(&WebRuleTest{})

func (t *WebRuleTest) TestName(c *C) {
  r := NewWebRule("a", "b", "c", "d", map[string]string{})

  c.Check(r.Name(), Equals, "a")
}

func (t *WebRuleTest) TestCaseSensitivity(c *C) {
  // Test that WebRules are case insensitive by default.
  r := NewWebRule("a", "b", "c", "d", map[string]string{})
  c.Check(r.contains("cat", "A"), Equals, true)

  // Test to ensure case-insensitive mode works.
  r = NewWebRule("a", "b", "c", "d", map[string]string{
    "case-sensitive": "true",
  })
  c.Check(r.contains("cat", "A"), Equals, false)
}

func (t *WebRuleTest) TestUnknownOption(c *C) {
  memlog := new(TestingLogger)
  logging.SetBackend(memlog)
  defer logging.SetBackend(
    logging.NewLogBackend(os.Stderr, "", go_log.LstdFlags))

  logging.SetLevel(logging.DEBUG, log.Module)

  NewWebRule("a", "b", "c", "d", map[string]string{
    "foobar": "true",
  })
  c.Check(memlog.Logs(), Matches, `[\s\S]*foobar[\s\S]*`)
  c.Check(memlog.Logs(), Matches, `[\s\S]*[uU]nknown[\s\S]*`)
}

func test_web_rule(
  page_content,
  sanity,
  trigger string) (sane, triggered bool) {
  s := httptest.NewServer(
    http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      fmt.Fprintln(w, page_content)
    }))
  defer s.Close()

  r := NewWebRule(
    "test rule",
    s.URL,
    sanity,
    trigger,
    map[string]string{})

  sane, triggered = r.TestTriggered()
  return
}

func (t *WebRuleTest) TestSanity(c *C) {
  sane, _ := test_web_rule("kittens", "cat", "")
  c.Check(sane, Equals, false)

  sane, _ = test_web_rule("kittens", "ten", "")
  c.Check(sane, Equals, true)
}

func (t *WebRuleTest) TestTriggers(c *C) {
  _, triggered := test_web_rule("kittens", "", "cat")
  c.Check(triggered, Equals, false)

  _, triggered = test_web_rule("kittens", "", "ten")
  c.Check(triggered, Equals, true)
}

func (t *WebRuleTest) TestHandleHttpFailures(c *C) {
  var s *httptest.Server
  s = httptest.NewServer(
    http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      s.CloseClientConnections()
    }))

  r := NewWebRule(
    "test rule",
    s.URL,
    "cats",
    "dogs",
    map[string]string{})

  sane, triggered := r.TestTriggered()
  c.Check(sane, Equals, false)
  c.Check(triggered, Equals, false)
}
