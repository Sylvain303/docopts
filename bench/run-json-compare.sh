#!/bin/bash


pathadd() {
  local after=0
  if [[ "$1" == "after" ]] ; then
    after=1
    shift
  fi

  local p

  for p in $*
  do
    if [ -d "$p" ] && ! echo $PATH | grep -E -q "(^|:)$p($|:)" ; then
      if [[ $after -eq 1 ]]
      then
        PATH="$PATH:${p%/}"
      else
        PATH="${p%/}:$PATH"
      fi
    fi
  done
}

pathrm() {
	PATH="$(echo $PATH | \
		sed -e "s;\(^\|:\)${1%/}\(:\|\$\);\1\2;g" \
        -e 's;^:\|:$;;g' \
        -e 's;::;:;g')"
}

DOCOPTS_GOPATH=$(dirname $(readlink -f $0))/..
pathadd $DOCOPTS_GOPATH
./timing.sh ../examples/arguments_example.sh
./timing.sh ../examples/old_docopts_api/arguments_example.sh
pathrm $DOCOPTS_GOPATH
pathadd  $HOME/code/docopt/docopts
./timing.sh ../examples/old_docopts_api/arguments_example.sh
