#!/bin/bash
#
# require jq (JSON command line parser to be installed)

PATH=..:$PATH

@test "parse" {
  run docopts parse "usage: pipo [-v] [--count=N] FILE..." -- pipoland deux "troisieme argument" ici "et la"
  echo "output=$output"
  [[ -n $output ]]
  [[ $status -eq 0 ]]
  run jq '."-v"' <<< "$output"
  [[ $output == false ]]
}

@test "get" {
  export DOCOPTS_JSON
  DOCOPTS_JSON=$(docopts parse "usage: pipo [-v] [--count=N] FILE..." -- pipoland deux "troisieme argument" ici "et la")
  echo $DOCOPTS_JSON
  run docopts get -- FILE
  [[ $status -eq 0 ]]
  [[ ${lines[0]} == "pipoland" ]]
  [[ ${#lines[@]} -eq 5 ]]
}

@test "get-keys" {
  export DOCOPTS_JSON=$(docopts parse "usage: pipo [-v] [--count=N] FILE..." -- pipoland deux "troisieme argument" ici "et la")
  echo $DOCOPTS_JSON
  run docopts get-keys
  [[ $status -eq 0 ]]
  keys=( $output )
  [[ ${#keys[@]} -eq 3 ]]
}

@test "dump" {
  export DOCOPTS_JSON=$(docopts parse "usage: pipo [-v] [--count=N] FILE..." -- pipoland deux "troisieme argument" ici "et la")
  echo $DOCOPTS_JSON
  run docopts dump json
  [[ $status -eq 0 ]]
  run jq -r '.FILE[0]' <<< "$output"
  [[ $output == "pipoland" ]]
}
