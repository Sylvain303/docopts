// vim: set ts=4 sw=4 sts=4 et:
//
// docopts.go is a command line wrapper for docopt.go to be used by bash scripts.
// docopt - command line arguments parser, that will make you smile.
//
package main

import (
    "fmt"
    "github.com/docopt/docopt-go"
    "regexp"
    "strings"
    "reflect"
    "os"
    "io"
    "io/ioutil"
    "sort"
    "encoding/json"
)

var Version string = `docopts 0.6.4
Copyleft (Ɔ) 2018 Sylvain Viart (golang version).
Copyright (C) 2013 Vladimir Keleshev, Lari Rasku.
License MIT <http://opensource.org/licenses/MIT>.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
`

var Usage string = `Shell interface for docopt, the CLI description language.

Usage:
  docopts [options] -h <msg> : [<argv>...]
  docopts [options] [-d] [--no-declare] -A <name>   -h <msg> : [<argv>...]
  docopts [options] -G <prefix>  -h <msg> : [<argv>...]
  docopts [options] --no-mangle  -h <msg> : [<argv>...]
  docopts [options] parse <msg> : [<argv>...]
  docopts [options] get <arg_name>
  docopts [options] get-keys
  docopts [options] auto-parse [--json|-A <name>|-G <prefix>] <bash_script> : [<argv>...]
  docopts [options] merge (ini|json)  <config_file>
  docopts [options] dump  (ini|json)
  docopts --howto

Actions:
  parse          Parse and outputs JSON result (shortcut for --json -h).
  auto-parse     Auto parse <bash_script> header and outputs parsed result.
  get            With DOCOPTS_JSON, get the value of <arg_name>
  fail           With DOCOPTS_JSON, set output the bash code with
                 error message, suitable for eval.
  merge          With DOCOPTS_JSON, output a new merge JSON reading
                 from ini file format or json file format.
  get-keys       With DOCOPTS_JSON, get all keys from JSON.

Options:
  -h <msg>, --help=<msg>        The help message in docopt format.
                                Without argument outputs this help.
                                If - is given, read the help message from
                                standard input.
                                If no argument is given, print docopts's own
                                help message and quit.
  -V <msg>, --version=<msg>     A version message.
                                If - is given, read the version message from
                                standard input.  If the help message is also
                                read from standard input, it is read first.
                                If no argument is given, print docopts's own
                                version message and quit.
  -s <str>, --separator=<str>   The string to use to separate the help message
                                from the version message when both are given
                                via standard input. [default: ----]
  -O, --options-first           Disallow interspersing options and positional
                                arguments: all arguments starting from the
                                first one that does not begin with a dash will
                                be treated as positional arguments.
  -H, --no-help                 Don't handle --help and --version specially.
  -A <name>                     Export the arguments as a Bash 4.x associative
                                array called <name>.
  -G <prefix>                   As without -A, but outputs Bash compatible
                                GLOBAL varibles assignment, uses the given
                                <prefix>_{option}={parsed_option}. Can be used
                                with numerical incompatible option as well.
                                See also: --no-mangle
  --no-mangle                   Output parsed option not suitable for bash eval.
                                As without -A but full option names are kept.
                                Rvalue is still shellquoted.
  --no-declare                  Don't output 'declare -A <name>', used only
                                with -A argument.
  --json                        Change output format to JSON. Activated
                                automaticaly for 'parse' action.
  --debug                       Output extra parsing information for debuging.
                                Output cannot be used in bash eval.
`

// testing trick, out can be mocked to catch stdout and validate
// https://stackoverflow.com/questions/34462355/how-to-deal-with-the-fmt-golang-library-package-for-cli-testing
var out io.Writer = os.Stdout

type DumpType uint

const (
    F_Json DumpType = 1 << iota
    F_Ini
    F_Bash_assoc
    F_Bash_global
)

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

// Store global behavior to avoid passing many optional arguments to methods.
type Docopts struct {
    Global_prefix string
    Mangle_key bool
    Output_declare bool
    Exit_function bool
    Format DumpType
    Json string
}

// output bash 4 compatible assoc array, suitable for eval.
func (d *Docopts) Print_bash_assoc(bash_assoc string, args docopt.Opts) {
    // Reuse python's fake nested Bash arrays for repeatable arguments with values.
    // The structure is:
    // bash_assoc[key,#]=length
    // bash_assoc[key,i]=value
    // 'i' is an integer from 0 to length-1
    // length can be 0, for empty array

    if d.Output_declare {
        fmt.Fprintf(out, "declare -A %s\n" ,bash_assoc)
    }

    for key, value := range args {
        // some golang tricks here using reflection to loop over the map[]
        rt := reflect.TypeOf(value)
        if IsArray(rt) {
            // all array is outputed even 0 size
            val_arr := value.([]string)
            for index, v := range val_arr {
                fmt.Fprintf(out, "%s['%s,%d']=%s\n", bash_assoc, Shellquote(key), index, To_bash(v))
            }
            // size of the array
            fmt.Fprintf(out, "%s['%s,#']=%d\n", bash_assoc, Shellquote(key), len(val_arr))
        } else {
            // value is not an array
            fmt.Fprintf(out, "%s['%s']=%s\n", bash_assoc, Shellquote(key), To_bash(value))
        }
    }
}

// Check if a value is an array
func IsArray(rt reflect.Type) bool {
    if rt == nil {
        return false
    }
    switch rt.Kind() {
    case reflect.Slice:
        return true
    case reflect.Array:
        return true
    default:
        return false
    }
}

func Shellquote(s string) string {
    return strings.Replace(s, "'", `'\''`, -1)
}

func IsBashIdentifier(s string) bool {
    identifier := regexp.MustCompile(`^([A-Za-z]|[A-Za-z_][0-9A-Za-z_]+)$`)
    return identifier.MatchString(s)
}

// Convert a parsed type to a text string suitable for bash eval
// as a right-hand side of an assignment.
// Handles quoting for string, no quote for number or bool.
func To_bash(v interface{}) string {
    var s string
    switch v.(type) {
    case bool:
        s = fmt.Sprintf("%v", v.(bool))
    case int:
        s = fmt.Sprintf("%d", v.(int))
    case string:
        s = fmt.Sprintf("'%s'", Shellquote(v.(string)))
    case []string:
        // escape all strings
        arr := v.([]string)
        arr_out := make([]string, len(arr))
        for i, e := range arr {
            arr_out[i] = Shellquote(e)
        }
        s = fmt.Sprintf("('%s')", strings.Join(arr_out[:],"' '"))
    case nil:
        s = ""
    default:
        panic(fmt.Sprintf("To_bash():unsuported type: %v for '%v'", reflect.TypeOf(v), v ))
    }

    return s
}

// Performs output for bash Globals (not bash 4 assoc) Names are mangled to became
// suitable for bash eval.
// If Docopts.Mangle_key: false simply print left-hand side assignment verbatim.
// used for --no-mangle
func (d *Docopts) Print_bash_global(args docopt.Opts) {
    var new_name string
    var err error
    var out_buf string

    // value is an interface{}
    for key, value := range args {
        if d.Mangle_key {
            new_name, err = d.Name_mangle(key)
            if err != nil {
                docopts_error("%v", err)
            }
        } else {
            new_name = key
        }

        out_buf += fmt.Sprintf("%s=%s\n", new_name, To_bash(value))
    }

    // final output
    fmt.Fprintf(out, "%s", out_buf)
}

// Transform a parsed option or place-holder name into a bash identifier if possible.
// It Docopts.Global_prefix is prepended if given, wrong prefix may produce invalid
// bash identifier and this method will fail.
func (d *Docopts) Name_mangle(elem string) (string, error) {
    var v string

    if elem == "-" || elem == "--" {
        return "", fmt.Errorf("not supported")
    }

    if Match(`^<.*>$`, elem) {
        v = elem[1:len(elem)-1]
    } else if Match(`^-[^-]$`, elem) {
        v = fmt.Sprintf("%c", elem[1])
    } else if Match(`^--.+$`, elem) {
        v = elem[2:]
    } else {
        v = elem
    }

    // alter output if we have a prefix
    key_fmt := "%s"
    if d.Global_prefix != "" {
        key_fmt = fmt.Sprintf("%s_%%s", d.Global_prefix)
    }

    v = fmt.Sprintf(key_fmt, strings.Replace(v, "-", "_", -1))

    if ! IsBashIdentifier(v) {
        return "", fmt.Errorf("cannot transform into a bash identifier: '%s' => '%s'", elem, v)
    }

    return v, nil
}

// helper for lazy typing
func Match(regex string, source string) bool {
    matched, _ := regexp.MatchString(regex, source)
    return matched
}

// Experimental: Change bash exit source code based on '--function' parameter
func (d *Docopts) Get_exit_code(exit_code int) (str_code string) {
    if d.Exit_function {
        str_code = fmt.Sprintf("return %d", exit_code)
    } else {
        str_code = fmt.Sprintf("exit %d", exit_code)
    }
    return
}

// Our HelpHandler which outputs bash source code to be evaled as error and stop or
// display program's help or version.
func (d *Docopts) HelpHandler_for_bash_eval (err error, usage string) {
    if err != nil {
        fmt.Printf("echo 'error: %s\n%s' >&2\n%s\n",
            Shellquote(err.Error()),
            Shellquote(usage),
            d.Get_exit_code(64),
        )
        os.Exit(1)
    } else {
        // --help or --version found and --no-help was not given
        fmt.Printf("echo '%s'\n%s\n", Shellquote(usage), d.Get_exit_code(0))
        os.Exit(0)
    }
}

// HelpHandler for go parser which parses docopts options. See: HelpHandler_for_bash_eval for parsing
// bash options. This handler is called when docopts itself detects a parse error on docopts usage.
// If docopts parsing is OK, then HelpHandler_for_bash_eval will be called by a second parser based on the
// help string given with -h <msg> or --help=<msg>. This behavior is a legacy behavior from docopts python
// previous version. This introduce strange hack in option parsing and may be changed after initial docopts go
// version release.
func HelpHandler_golang(err error, usage string) {
    if err != nil {
        err_str := err.Error()
        // we hack for our polymorphic argument -h or -V
        // it was the same hack in python version
        if len(err_str) >= 9 {
            if err_str[0:2] == "-h" || err_str[0:6] == "--help" {
                // print full usage message (global var)
                fmt.Println(strings.TrimSpace(Usage))
                os.Exit(0)
            }
            if err_str[0:2] == "-V" || err_str[0:9] == "--version" {
                fmt.Println(strings.TrimSpace(Version))
                os.Exit(0)
            }
        }

        // When we have an error with err_str empty, this is a special case:
        // we received an usage string which MUST receive an argument and no argument has been
        // given by the user. So this is a valid, from golang point of view but not for bash.
        if len(err_str) == 0 {
            // no arg at all, display small usage, also exits 1
            d := &Docopts{Exit_function: false}
            d.HelpHandler_for_bash_eval(fmt.Errorf("no argument"), usage)
        }

        // real error
        fmt.Fprintf(os.Stderr, "my error: %v, %v\n", err, usage)
        os.Exit(1)
    } else {
        // no error, never reached?
        fmt.Println(usage)
        os.Exit(0)
    }
}

func docopts_error(msg string, err error) {
    if err != nil {
        msg = fmt.Sprintf(msg, err)
    }
    fmt.Fprintf(os.Stderr, "docopts:error: %s\n", msg)
    os.Exit(1)
}

func (d* Docopts) Dump(arguments docopt.Opts) {
    switch(d.Format) {
    case F_Json:
        b, err := json.Marshal(arguments)
        if err != nil {
            docopts_error("dump json: %v", err)
        }
        out.Write(b)
    case F_Ini:
        fmt.Fprintf(out, "ini dump not supported yet")
    }
}

func main() {
    golang_parser := &docopt.Parser{
      OptionsFirst: true,
      SkipHelpFlags: true,
      HelpHandler: HelpHandler_golang,
    }

    arguments, err_parse := golang_parser.ParseArgs(Usage, nil, Version)

    if err_parse != nil {
        msg := fmt.Sprintf("mypanic: %v\n", err_parse)
        panic(msg)
    }

    debug := arguments["--debug"].(bool)
    if debug {
        print_args(arguments, "golang")
    }

    // create our Docopts struct
    d := &Docopts{
        Global_prefix: "",
        Mangle_key: true,
        Output_declare: true,
        // Exit_function is experimental
        Exit_function: false,
        Format: F_Json,
        Json:  os.Getenv("DOCOPTS_JSON"),
    }

    // parse docopts's own arguments
    argv := arguments["<argv>"].([]string)
    options_first := arguments["--options-first"].(bool)
    no_help :=  arguments["--no-help"].(bool)
    separator := arguments["--separator"].(string)

    var doc string
    var err error
    // bash_version will be empty if error, so we dont care
    bash_version, _ := arguments.String("--version")

    if doc, err = arguments.String("--help"); err != nil {
        // read from stdin
        if doc == "-" && bash_version == "-" {
            bytes, _ := ioutil.ReadAll(os.Stdin)
            arr := strings.Split(string(bytes), separator)
            if len(arr) == 2 {
                doc, bash_version = arr[0], arr[1]
            } else {
                msg := "error: help + version stdin, not found"
                if debug {
                    msg += fmt.Sprintf("\nseparator is: '%s'\n", separator)
                    msg += fmt.Sprintf("spliting has given %d blocs, exactly 2 are expected\n", len(arr))
                }
                panic(msg)
            }
        } else if doc == "-" {
            bytes, _ := ioutil.ReadAll(os.Stdin)
            doc = string(bytes)
        } else if bash_version == "-" {
            bytes, _ := ioutil.ReadAll(os.Stdin)
            bash_version = string(bytes)
        }
        doc = strings.TrimSpace(doc)
        bash_version = strings.TrimSpace(bash_version)
    }

    // mode parse read from another arg
    if doc == "" {
        doc, err = arguments.String("<msg>")
        if err != nil {
            panic(err)
        }
    }

    if debug {
        fmt.Printf("%20s : %v\n", "doc", doc)
        fmt.Printf("%20s : %v\n", "bash_version", bash_version)
    }

    // now parses bash program's arguments
    parser := &docopt.Parser{
      HelpHandler: d.HelpHandler_for_bash_eval,
      OptionsFirst: options_first,
      SkipHelpFlags: no_help,
    }

    bash_args, err_parse_bash := parser.ParseArgs(doc, argv, bash_version)
    if err_parse_bash != nil {
        panic(err_parse_bash)
    }

    // ========================================
    // main action code
    // ========================================

    if debug {
        print_args(bash_args, "bash")
        fmt.Println("----------------------------------------")
    }

    if arguments["dump"].(bool) {
        if arguments["json"].(bool) {
            d.Format = F_Json
        } else if arguments["ini"].(bool) {
            d.Format = F_Ini
        }
        d.Dump(bash_args)
        os.Exit(0)
    }

    // outputer

    if name, err := arguments.String("-A"); err == nil {
        if ! IsBashIdentifier(name) {
            docopts_error("-A switch:%v", fmt.Errorf("not a valid Bash identifier: '%s'", name))
        }
        d.Format = F_Bash_assoc
        d.Output_declare = ! arguments["--no-declare"].(bool)
        d.Print_bash_assoc(name, bash_args)
    }

    if arguments["parse"].(bool) {
        // force option
        d.Format = F_Json
        d.Dump(bash_args)
    } else {
        d.Format = F_Bash_global
        d.Mangle_key = ! arguments["--no-mangle"].(bool)
        if global_prefix, err := arguments.String("-G"); err == nil {
            d.Global_prefix = global_prefix
        }
        d.Print_bash_global(bash_args)
    }

    os.Exit(0)
}
