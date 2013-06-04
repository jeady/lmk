#!/bin/sh

set -x

gocov test github.com/jeady/lmk/engine > /tmp/lmk.cov.json
gocov-html /tmp/lmk.cov.json > /tmp/lmk.cov.html
open /tmp/lmk.cov.html
