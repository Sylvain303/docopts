# New API proposal for docopts - docopt on shell (bash)

## Update: 2015-03-13

See [issue #7](https://github.com/docopt/docopts/issues/7) for the discussion.

This file content is moved to the wiki, update the following wiki page:

Published in the [docopts Wiki](https://github.com/docopt/docopts/wiki)

Submit pull request for changes.

## docopts go, with JSON support
Date: 2018-06-01

See: [`config_file_example.sh`](examples/config_file_example.sh) for a detailed prototype code.

We propose to embed JSON support into `docopts` (single binary, no need to extra `jq` for handling
JSON in bash.

The idea is to store arguments parsed result into a shell env variable, and to reuse it by
`docopts` sub-call with action expecting this variable to be filled with JSON output.

### Usage examples

```bash
# require to separate export and assignment to get $? after $()
export DOCOPTS_JSON
DOCOPTS_JSON=$(docopts parse "Usage: mystuff [--code] INFILE [--out=OUTFILE]" -- "$@")

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
export DOCOPTS_JSON
DOCOPTS_JSON=$(docopts auto-parse "$0" --version '0.1.1rc' -- "$@")
# docopts fail : display error stored in DOCOPTS_JSON and output exit code for
# caller
[[ $? -ne 0 ]] && eval $(docopts fail)
```

### RPC support and session

See: [examples/rpc_suppots_api_example.sh](examples/rpc_suppots_api_example.sh)

`docopts` becomes an daemon program which interract with itself and the API looks like if it was part of bash language.


```bash
# create daemon RPC process
docopts auto-parse --auto-fail "$0" --version '0.1.1rc' -- "$@"

# all calls are using first call init
docopts merge ini - <<< "$(output_ini_config)"
docopts merge json - <<< "$(output_json_config)"

docopts dump json

for a in $(docopts get-keys); do
  echo "$a = $(docopts get $a)"
done

docopts kill
```
