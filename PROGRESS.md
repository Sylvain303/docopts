# current work in PROGRRESS json-api

argument api change compatibility

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

# next

code docopts fail API

## CI
integration with automated tests

## provide test on old environment

docker?
32bits
bash 3
