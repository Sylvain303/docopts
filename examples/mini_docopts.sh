#!/bin/bash 
#
# Usage: mini_docopts.sh <myparam>
#
#  hello world docopts myparam
#
# Option:
#   --help     display this message
#   --version  display the version


version="mini_docopts.sh 0.2"
# auto fetch this comment above
help=$(sed -n -e '/^# Usage:/,/^$/ s/^# \?//p' < $0)

libdir="$(readlink -f "$(dirname $0)/..")"

# -A output associative array named 'args'
# -V uses the version string
# -h uses the help string (fetched with sed above)
# "$@" nicely passes all script params as quoted strings (See. man bash)
output=$(echo "$help" | python3 $libdir/docopts "$@")
echo $output
#eval "$output"

#echo "hello world ${args[<myparam>]}"
