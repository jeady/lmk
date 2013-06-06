// The purpose of this interface is to make it possible to swap out the method
// lmk! uses to access the network. Specifically, Google App Engine doesn't
// allow you to simply use net/http. Instead, for GAE support, the UrlFetcher
// must be replaced with one that uses appengine/urlfetch.
package engine

type UrlFetcher interface {
  Get(url string) (page string, err error)
}

func SetUrlFetcher(f UrlFetcher) UrlFetcher {
  old := fetcher
  fetcher = f
  return old
}

func GetUrl(url string) (string, error) {
  return fetcher.Get(url)
}

var fetcher UrlFetcher

func init() {
  fetcher = &NetHttpFetcher{}
}
