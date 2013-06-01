# lmk!
Let Me Know! is a user-facing notification service. Want to be the first to
know when the lineout for ACL goes live? Woot.com is selling a macbook air?
George R. R. Martin has blogged about the next ASOIAF book? lmk! is the answer.
lmk! will notify you by email when an event you've set occurs.

## Installation
``` bash
go get -u github.com/jeady/lmk/lmk
cp $GOPATH/pkg/github.com/jeady/lmk/config/lmk.conf.sample lmk.conf
cp $GOPATH/pkg/github.com/jeady/lmk/config/rules.conf.sample rules.conf
vim -p lmk.conf rules.conf
lmk test-all
```

It is highly recommended that you also run `lmk test-smtp` to test that
notifications will be delivered.

It is also recommended that you add `lmk run-all` to a cron job so that it will
run periodically. In the future, lmk! will provide a daemon mode that may be
used instead.

A package for Arch Linux will be provided in the future.

For more options, run `lmk help`.

## Configuring
Global program configuration is done by modifying lmk.conf (or whatever file
is specified using -c on the command line). Rules are entered into rules.conf
(or whatever rule file is specified in lmk.conf). For examples and
available options, read config/lmk.conf.sample and config/rules.conf.sample.

## Contributing
I am always looking for ideas to add as triggers (or anything else). Please do
feel free to open an issue tagged as 'question' or 'enhancement' if you'd like
to comment or suggest.
