package engine

import (
  "fmt"
  "net/http"
  "net/http/httptest"

  . "launchpad.net/gocheck"
)

type NetHttpFetcherTest struct{}

var _ = Suite(&NetHttpFetcherTest{})

func (t *NetHttpFetcherTest) TestGet(c *C) {
  s := httptest.NewServer(
    http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      fmt.Fprint(w, "foo bar baz")
    }))
  defer s.Close()

  f := NetHttpFetcher{}
  content, err := f.Get(s.URL)
  c.Check(content, Equals, "foo bar baz")
  c.Check(err, Equals, nil)
}
