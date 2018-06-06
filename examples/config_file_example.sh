#!/bin/bash
#
# Work in progress. Not working yet. Need Json support inside docopts.
#
# Usage:
#   quick_example.sh tcp [<host>] [--force] [--timeout=<seconds>]
#   quick_example.sh serial <port> [--baud=<rate>] [--timeout=<seconds>]
#   quick_example.sh -h | --help | --version
#

echo "sorry Draft, not implemented yet!"
exit 0

# if docopts is in PATH, not needed.
PATH=..:$PATH

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

# DOCOPTS_JSON is the default variable you need to set with docopts parsed output.
# You can change the variable's name See: --env
# In all case you need to export it once, so it will become available to docopts.
export DOCOPTS_JSON=$(docopts auto-parse --json "$0" --version '0.1.1rc' : "$@")
# parse failure or help display
[[ $? -ne 0 ]] && eval $(docopts fail)

# DOCOPTS_JSON is read and merged by docopts itself, use --env to change its name if needed
DOCOPTS_JSON=$(docopts merge ini - <<< "$(output_ini_config)")
DOCOPTS_JSON=$(docopts merge json - <<< "$(output_json_config)")

echo "JSON config:"
output_json_config
echo "INI config:"
output_ini_config
echo "Result:"
docopts dump-json

echo "Result with bash loop:"
pprint_bash

