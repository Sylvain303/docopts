#!/bin/bash 
# copied from docopt python quick_example.py
#
# Usage:
# quick_example.py tcp <host> <port> [--timeout=<seconds>]
# quick_example.py serial <port> [--baud=9600] [--timeout=<seconds>]
# quick_example.py -h | --help | --version
#

# auto fetch this comment above
version="0.1.1rc"
help=$(sed -n -e '/^# Usage:/,/^$/ s/^# \?//p' < $0)
libdir="$(readlink -f "$(dirname $0)/..")"

# python version do:
# docopt(doc, argv=None, help=True, version=None, options_first=False)
tmp=$(echo "$help" | python3 $libdir/docopts "$@")
echo $tmp
ls $tmp


#eval "$tmp"
