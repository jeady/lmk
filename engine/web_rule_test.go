package engine

import (
  "fmt"
  "net/http"
  "net/http/httptest"

  "code.google.com/p/gomock/gomock"
  . "launchpad.net/gocheck"
)

type WebRuleTest struct{}

var _ = Suite(&WebRuleTest{})

func (t *WebRuleTest) TestName(c *C) {
  r := NewWebRule("a", "b", "c", "d")

  c.Check(r.Name(), Equals, "a")
}

func (t *WebRuleTest) TestBadRegexes(c *C) {
  // Unescaped backslash
  r := NewWebRule("a", "b", "\\", "d")
  c.Check(r, Equals, (*WebRule)(nil))

  r = NewWebRule("a", "b", "c", "\\")
  c.Check(r, Equals, (*WebRule)(nil))
}

func (t *WebRuleTest) TestUsesUrlFetcher(c *C) {
  m := gomock.NewController(c)
  defer m.Finish()

  u := NewMockUrlFetcher(m)
  u.EXPECT().Get("url").Return("", nil).Times(1)

  r := NewWebRule("test rule", "url", "a", "b")
  r.SetUrlFetcher(u)
  r.TestTriggered()
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
    trigger)

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
  r := NewWebRule("a", "b", "", "c")
  c.Check(r.test_triggered("cat"), Equals, true)

  r = NewWebRule("a", "b", "", "c")
  r.SetOptions(map[string]string{
    "trigger-on-match": "true",
  })
  c.Check(r.test_triggered("cat"), Equals, true)

  r = NewWebRule("a", "b", "", "c")
  r.SetOptions(map[string]string{
    "trigger-on-match": "false",
  })
  c.Check(r.test_triggered("cat"), Equals, false)

}

func (t *WebRuleTest) TestReturnsUnconsumed(c *C) {
  r := NewWebRule("a", "b", "", "c")
  u := r.SetOptions(map[string]string{
    "trigger-on-match": "false",
  })
  c.Check(u, HasLen, 0)

  r = NewWebRule("a", "b", "", "c")
  u = r.SetOptions(map[string]string{
    "trigger-on-match": "false",
    "extra-1":          "a",
  })
  c.Assert(len(u), Equals, 1)
  c.Check(u["extra-1"], Equals, "a")

  r = NewWebRule("a", "b", "", "c")
  u = r.SetOptions(map[string]string{
    "trigger-on-match": "false",
    "extra-1":          "a",
    "extra-2":          "b",
  })
  c.Assert(len(u), Equals, 2)
  c.Check(u["extra-1"], Equals, "a")
  c.Check(u["extra-2"], Equals, "b")
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
    "dogs")

  sane, triggered := r.TestTriggered()
  c.Check(sane, Equals, false)
  c.Check(triggered, Equals, false)
}
