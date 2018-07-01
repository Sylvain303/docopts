package main

import (
  "github.com/docopt/docopt-go"
  "fmt"
  "sort"
)

var Usage string = `
Usage: atest.sh <arg>
       atest.sh ( -h | --help )
       atest.sh ( -V | --version )

      arg				An argument.

Options:
      -h, --help        Show this help message and exits.
      -V, --version     Print version and copyright information.
`

// debug helper
func print_args(args docopt.Opts, message string) {
    // sort keys
    mk := make([]string, len(args))
    i := 0
    for k, _ := range args {
        mk[i] = k
        i++
    }
    sort.Strings(mk)
    fmt.Printf("################## %s ##################\n", message)
    for _, key := range mk {
        fmt.Printf("%20s : %v\n", key, args[key])
    }
}

func main() {
    golang_parser := &docopt.Parser{}
    arguments, err := golang_parser.ParseArgs(Usage, nil, "")

    if err != nil {
        msg := fmt.Sprintf("mypanic: %v\n", err)
        panic(msg)
    }

    print_args(arguments, "debug")
}
