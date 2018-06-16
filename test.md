# testing docopts

## testee.sh

`testee.sh` is a bash wrapper which convert `docopts` bash output to JSON format.
Was required befor JSON version of docopts or for backward compatibily testing.

### testee.sh usage

recive usage on stdin and args on cmdline, should output JSON

```
echo "usage: prog (go <direction> --speed=<km/h>)..." | ./testee.sh go left --speed=5  go right --speed=9
```
