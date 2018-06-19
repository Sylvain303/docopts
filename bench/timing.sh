#!/bin/bash

# float computation US number format
export LANG=C

tmpf=/dev/shm/timing_$$.out
START=$(date +%s.%N)
# do something #######################
echo "running: $@"

eval "$@" &> $tmpf

#######################################
END=$(date +%s.%N)
DIFF=$( echo "scale=3; (${END} - ${START})*1000/1" | bc )
echo "${DIFF} ms output: $(wc -l $tmpf)"
