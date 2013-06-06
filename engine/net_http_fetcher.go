package engine

import (
  "io/ioutil"
  "net/http"
)

type NetHttpFetcher struct{}

func (*NetHttpFetcher) Get(url string) (page string, err error) {
  log.Debug("GET %s", url)

  page = ""

  var resp *http.Response
  resp, err = http.Get(url)
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
