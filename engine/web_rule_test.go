package engine

import (
  "fmt"
  "net/http"
  "net/http/httptest"

  . "launchpad.net/gocheck"
)

type WebRuleTest struct{}

var _ = Suite(&WebRuleTest{})

func (t *WebRuleTest) TestName(c *C) {
  r := NewWebRule("a", "b", "c", "d", map[string]string{})

  c.Check(r.Name(), Equals, "a")
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

func (t *WebRuleTest) TestTriggerOnMatch(c *C) {
  r := NewWebRule("a", "b", "", "c", map[string]string{})
  c.Check(r.test_triggered("cat"), Equals, true)

  r = NewWebRule("a", "b", "", "c", map[string]string{
    "trigger-on-match": "true",
  })
  c.Check(r.test_triggered("cat"), Equals, true)

  r = NewWebRule("a", "b", "", "c", map[string]string{
    "trigger-on-match": "false",
  })
  c.Check(r.test_triggered("cat"), Equals, false)

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
