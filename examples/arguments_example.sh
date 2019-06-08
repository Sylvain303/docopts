#!/usr/bin/env bash
#
# Usage:
#   arguments_example.sh [-vqrh] [FILE] ...
#   arguments_example.sh (--left | --right) CORRECTION FILE
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

# if docopts is in PATH, not needed.
PATH=..:$PATH

# auto parse the header above, See: docopt_get_help_string
# docopt_auto_parse use DOCOPTS_JSON exported global
export DOCOPTS_JSON
DOCOPTS_JSON=$(docopts auto-parse "$0" -- "$@")
[[ $? -ne 0 ]] && eval "$DOCOPTS_JSON"

# main code
for a in $(docopts get-keys) ; do
    echo "$a = $(docopts get -- "$a")"
done
