#!/bin/bash
#
# Work in progress. Not working yet. Need RPC support inside docopts.
#
# Usage:
#   quick_example.sh tcp [<host>] [--force] [--timeout=<seconds>]
#   quick_example.sh serial <port> [--baud=<rate>] [--timeout=<seconds>]
#   quick_example.sh -h | --help | --version
#

echo "sorry Draft, not implemented yet!"
exit 0

# PATH modified outside see init_PATH.sh

output_json_config() {
  # Pretend that we load the following JSON file:
  cat << EOF
{"--force": true,
 "--timeout": "10",
 "--baud": "9600"}
EOF
}

output_ini_config() {
    # Pretend that we load the following INI file:
  cat << EOF
[default-arguments]
--force
--baud=19200
<host>=localhost
EOF
}

pprint_bash() {
  # docopts still uses DOCOPTS_JSON env var to fetch its parsed Json output
  for a in $(docopts get-keys); do
    # use -- to force stop option parsing as $a will be '--force' for example.
    echo "$a = $(docopts get $a)"
  done
}

# = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = =
# main code
# = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = =

# how does it works: 
# - docopts register itself as a deamon for this instance
# - auto-parse extract Usage from the source file
# - if docopt parsing fail, it send kill STOP to the parent process
# - if it succed options are stored in memory
# - subsequent call to docopts interact with in memory parsed data
# - so options is parsed once and stay active in memory
docopts auto-parse --auto-fail "$0" --version '0.1.1rc' -- "$@"

docopts merge ini - <<< "$(output_ini_config)"
docopts merge json - <<< "$(output_json_config)"

echo "JSON config:"
output_json_config
echo "INI config:"
output_ini_config

echo "Result:"
docopts dump json

echo "Result with bash loop:"
pprint_bash

# docopts also auto kill itself when the parent process ends
docopts kill 
