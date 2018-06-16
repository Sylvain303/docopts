# docopts (docopt for bash) TODO list

## functional testing for all option

## return or kill for function instead of exit

Add a parameter to handle return or kill instead of exit so it can be launched inside a function.

## embeded JSON

See [API_proposal.md](API_proposal.md)

## build and publish binary

Makefile OK
publish it as a new release too.
https://docs.travis-ci.com/user/deployment/releases/
https://blog.questionable.services/article/build-go-binaries-travis-ci-github/

## generate bash completion from usage

```
docopts -h "$help" --generate-completion
```

## embed test routine (validation)?

may we cat interract with the caller to eval some validation…
It is needed? Is it our goal?

```
# with tests
# pass value to parent: JSON or some_thing_else
eval $(docopts --eval --json --help="Usage: mystuff [--code] INFILE [--out=OUTFILE]" -- "$@")
if docopts test -- file_exists:--code !file_exists:--out

eval $(docopts --eval --json --help="Usage: prog [--count=NUM] INFILE..."  -- "$@")
if docopts test -- num:gt:1:--count file_exists:INFILE
```

## config file parse config to option format

À la nslcd… ?
