# docopts with JSON support

This branch is implementing JSON support in ’docopts’
See Proposal lib API: [on `docopts` Wiki](https://github.com/docopt/docopts/wiki)

See: current go version is [100% compatible with python's docopts](https://github.com/Sylvain303/docopts/tree/docopts-go).

Please report any non working code with [issue](https://github.com/docopt/docopts/issues) and examples.

## docopts

Status: Draft - work in progress

This is the golang version  of `docopts` : the command line wrapper for bash.

## Install

```
go get github.com/docopt/docopt-go
go get github.com/Sylvain303/docopts
cd src/github.com/Sylvain303/docopts
go build docopts.go
```

or via Makefile (generate 64bits, 32bits, arm and OSX-64bits version of docopts)

```
cd src/github.com/Sylvain303/docopts
make all
```

Tested built on: `go version go1.10.2 linux/amd64`

The var name could be explicitly set to any user need (instead of default `DOCOPTS_JSON`):

### pre-built binary

pre-built binary are attached to [releases](https://github.com/Sylvain303/docopts/releases)

download and rename it as `docopts` and put in your `PATH`

```
mv docopts-32bit docopts
cp docopts docopts.sh /usr/local/bin
```

You are strongly encouraged to build your own binary. Find a local golang developper in whom you trust and ask her, for a beer or two, if she could build it for you. ;)

## Usage

See [Examples](examples/)


See: [`config_file_example.sh`](examples/config_file_example.sh) for a detailed prototype code.

We propose to embed JSON support into `docopts` (single binary, no need to extra `jq` for handling
JSON in bash.

The idea is to store arguments parsed result into a shell env variable, and to reuse it by
`docopts` sub-call with action expecting this variable to be filled with JSON output.

### Usage examples
With a go workspace.

```bash
DOCOPTS_JSON=$(docopts --json --h "Usage: mystuff [--code] INFILE [--out=OUTFILE]" : "$@")

# automaticly use $DOCOPTS_JSON
if [[ $(docopts get --code) == checkit ]]
then
  action
fi
```


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

See branch docpot-go
