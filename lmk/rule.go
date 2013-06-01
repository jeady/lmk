package lmk

type Rule interface {
  Name() string
  TestTriggered() (sane, triggered bool)
}
