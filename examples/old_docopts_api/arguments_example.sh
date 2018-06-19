#!/bin/bash
#
# Usage: arguments_example.sh [-vqrh] [FILE] ...
#           arguments_example.sh (--left | --right) CORRECTION FILE
#
# Process FILE and optionally apply correction to either left-hand side or
# right-hand side.
#
# Arguments:
#   FILE        optional input file
#   CORRECTION  correction angle, needs FILE, --left or --right to be present
#
# Options:
#   -h --help
#   -v       verbose mode
#   -q       quiet mode
#   -r       make report
#   --left   use left-hand side
#   --right  use right-hand side
#

# this example use no extra bash lib.
# only docopts

help=$(sed -n -e '/^# Usage:/,/^$/ s/^# \?//p' < $0)
# you must define PATH accordingly to find docopts
docopts -A ARGS -h "$help" : "$@"

# docopt_auto_parse use ARGS bash 4 globla assoc array
# main code
# on assoc array '!' before nane gike hash keys
for a in ${!ARGS[@]} ; do
    echo "$a = ${ARGS[$a]}"
done
