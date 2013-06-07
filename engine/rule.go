package engine

type Rule interface {
  Name() string
  TestTriggered() (sane, triggered bool)
}

type UrlFetchingRule interface {
  SetUrlFetcher(f UrlFetcher) UrlFetcher
}
