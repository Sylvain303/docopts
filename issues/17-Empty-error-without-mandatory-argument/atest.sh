#!/usr/bin/env bash
# file: atest.sh

arg=''

eval "$(docopts -V - -h - : "$@" <<EOF
Usage: atest.sh <arg>
       atest.sh ( -h | --help )
       atest.sh ( -V | --version )

      arg				An argument.

Options:
      -h, --help        Show this help message and exits.
      -V, --version     Print version and copyright information.
----
atest 0.1.0
EOF
)"

echo "$arg"
