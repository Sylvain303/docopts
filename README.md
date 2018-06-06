# docopts with JSON support

This branch is implementing JSON support in ’docopts’
See Proposal lib API: [on `docopts` Wiki](https://github.com/docopt/docopts/wiki)

See: current go version is [100% compatible with python's docopts](https://github.com/Sylvain303/docopts/tree/docopts-go).

Please report any non working code with [issue](https://github.com/docopt/docopts/issues) and examples.

## docopts

Status: Draft - work in progress

This is the golang version  of `docopts` : the command line wrapper for bash.

## Usage

See [Examples](examples/)

Date: 2018-06-01

See: [`config_file_example.sh`](examples/config_file_example.sh) for a detailed prototype code.

We propose to embed JSON support into `docopts` (single binary, no need to extra `jq` for handling
JSON in bash.

The idea is to store arguments parsed result into a shell env variable, and to reuse it by
`docopts` sub-call with action expecting this variable to be filled with JSON output.

### Usage examples

```bash
DOCOPTS_JSON=$(docopts --json --h "Usage: mystuff [--code] INFILE [--out=OUTFILE]" : "$@")

# automaticly use $DOCOPTS_JSON
if [[ $(docopts get --code) == checkit ]]
then
  action
fi
```

The var name could be explicitly set to any user need (instead of default `DOCOPTS_JSON`):

```bash
docopts --env SOME_DOCOPTS_JSON get --code

# or using an env var naming the env var...
DOCOPTS_JSON_VAR=SOME_DOCOPTS_JSON
docopts get --code
```

### parse failure detection and reaction inside bash code:

* On parse success, the JSON is filled as expected => exit 0 (`DOCOPTS_JSON` is filled with parsed options)
* if -h is given, exit code is 42 and questions need to be answered (`DOCOPTS_JSON` is filled with help message)
* else exit 1 => some argument parsing error (`DOCOPTS_JSON` is filled with error message)

`DOCOPTS_JSON` also contains exit code for the caller if necessary.

```bash
DOCOPTS_JSON=$(docopts --json --auto-parse "$0" --version '0.1.1rc' : "$@")
# docopts fail : display error stored in DOCOPTS_JSON and output exit code for
# caller
[[ $? -ne 0 ]] && eval $(docopts fail)
```

## Developpers

If you want to clone this repository and hack docopts:

Use `git clone --recursive`, to get submodules only required for testing with `bats`.

Fetch the extra golang version of `docopt` (required for building `docopts`)

```bash
go get github.com/docopt/docopt-go
```

If you forgot `--recursive`, you can also run afterward:

~~~bash
git submodule init
git submodule update
~~~

## Tests

Some tests are coded along this code base.

- bats bash unit tests and functionnal testing
- `language_agnostic_tester.py` (old python wrapper, full docopt compatibily tests)
- See Also: docopt.go own tests in golang
- `docopts_test.go` go unit test for `docopts.go`

### Runing tests

#### bats
```
cd ./tests
. bats.alias
bats .
```

#### `language_agnostic_tester`

```
python language_agnostic_tester.py ./testee.sh
```

#### golang docopt.go (golang parser lib)

```
cd PATH/to/go/src/github.com/docopt/docopt-go/
go test -v .
```

#### golang docopts (our bash wrapper)

```
go test -v
```
