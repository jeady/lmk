#!/bin/sh

if [ $# -eq 0 ]
then
  echo "Usage: ./format.sh file1.go file2.go ..."
fi

set -x

for file in "$@"
do
  gofmt -w -tabs=false -tabwidth=2 "$file"
done
