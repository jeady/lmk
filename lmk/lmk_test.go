package lmk

import (
  "testing"

  "launchpad.net/gocheck"
)

// Hooks GoCheck into go test for the entire package.
func Test(t *testing.T) { gocheck.TestingT(t) }
