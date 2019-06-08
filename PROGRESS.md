# current work in PROGRRESS json-api

## Dumper() handle error

On master merge, now handle err returned by `Print_bash_assoc`

Remove `Print_bash_args` now handled by `Print_bash_assoc`

Commit and go to master and code release helper.

## argument api change compatibility

`auto-parse --json`

add to functionnal testing :

added `Preprocess_agrv` => before + bash
with + without ':' (test true)

with empty args => print help ??

OK
```
echo "usage: prog" | ./docopts parse - --debug --
```

KO
```
echo "usage: prog" | ./docopts -h - --debug :
```

## code docopts fail API

## release and binary

publish release and pre-build binaries

## CI
integration with automated tests

* travis deploy? https://docs.travis-ci.com/user/deployment/releases/

## provide test on old environment

docker?
32bits
bash 3

