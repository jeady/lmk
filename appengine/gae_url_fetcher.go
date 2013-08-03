// +build appengine

package appengine

import (
  "io/ioutil"
  "net/http"

  "appengine"
  "appengine/urlfetch"
)

// Complies with lmk.UrlFetcher
type GaeUrlFetcher struct {
  client *http.Client
}

func (f *GaeUrlFetcher) Get(url string) (page string, err error) {
  page = ""

  var resp *http.Response
  resp, err = f.client.Get(url)
  if err != nil {
    return
  }
  defer resp.Body.Close()

  var page_bytes []byte
  page_bytes, err = ioutil.ReadAll(resp.Body)
  if err == nil {
    page = string(page_bytes)
  }
  return
}

func NewGaeUrlFetcher(c appengine.Context) *GaeUrlFetcher {
  return &GaeUrlFetcher{client: urlfetch.Client(c)}
}
