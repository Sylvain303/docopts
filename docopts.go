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
    "bufio"
)

var Docopts_Version string = `docopts 0.7.0 - with JSON support
Copyleft (Ɔ) 2018 Sylvain Viart (golang version).
Copyright (C) 2013 Vladimir Keleshev, Lari Rasku.
License MIT <http://opensource.org/licenses/MIT>.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
`

var Usage string = `bug auto-parse with extra options

Usage:
  docopts [options] -h <msg> : [<argv>...]
  docopts [options] [--no-declare] -A <name> -h <msg> : [<argv>...]
  docopts [options] [-G <prefix>] [--no-mangle] -h <msg> : [<argv>...]
  docopts [options] parse <msg> -- [<argv>...]
  docopts [options] get -- <arg_name>
  docopts [options] get-keys
  docopts [options] auto-parse [--json|-A <name>|-G <prefix>] <filename> -- [<argv>...]
  docopts [options] merge (ini|json)  <config_file>
  docopts [options] dump  (ini|json)
  docopts [options] fail
  docopts --howto

Actions:
  parse          Parse and outputs JSON result (shortcut for --json -h <msg>).
  auto-parse     Auto parse <filename> header and outputs parsed result
                 default output: json.
                 Warning: that <arg> separator is -- (ending option parsing).
  get            With DOCOPTS_JSON, get the value of <arg_name>.
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
var DEBUG bool

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
    Indent bool
    Assoc_name string
    Usage string
    Bash_Version string
    Separator string
}

func (d* Docopts) debug_print_Docopts() {
    enc := json.NewEncoder(os.Stdout)
    enc.SetIndent("", "    ")
    enc.SetEscapeHTML(false)

    fmt.Println("================ Docopts: =================")
	enc.Encode(d)
}

// output bash 4 compatible assoc array, suitable for eval.
func (d *Docopts) Print_bash_assoc(args docopt.Opts) {
    // Reuse python's fake nested Bash arrays for repeatable arguments with values.
    // The structure is:
    // bash_assoc[key,#]=length
    // bash_assoc[key,i]=value
    // 'i' is an integer from 0 to length-1
    // length can be 0, for empty array

    if d.Output_declare {
        fmt.Fprintf(out, "declare -A %s\n" , d.Assoc_name)
    }

    for key, value := range args {
        if IsArray(value) {
            // all array is outputed even 0 size
            val_arr := value.([]string)
            for index, v := range val_arr {
                fmt.Fprintf(out, "%s['%s,%d']=%s\n", d.Assoc_name, Shellquote(key), index, To_bash(v))
            }
            // size of the array
            fmt.Fprintf(out, "%s['%s,#']=%d\n", d.Assoc_name, Shellquote(key), len(val_arr))
        } else {
            // value is not an array
            fmt.Fprintf(out, "%s['%s']=%s\n", d.Assoc_name, Shellquote(key), To_bash(value))
        }
    }
}

func (d *Docopts) Print_value(key string) {
    parsed_args, err := d.Get_DOCOPTS_JSON()
    if err == nil {
        fmt.Fprintf(out, "%s\n", To_bash(parsed_args[key]))
    } else {
        docopts_error(fmt.Sprintf("get '%s': %%s", key), err)
    }
}

// fetch DOCOPTS_JSON to JSON data sturcture
func (d* Docopts) Get_DOCOPTS_JSON() (docopt.Opts, error) {
    var parsed_args docopt.Opts
    // JSON string is loaded from env DOCOPTS_JSON in main()
    if len(d.Json) > 0 {
        json.Unmarshal([]byte(d.Json), &parsed_args)
        return parsed_args, nil
    }
    return parsed_args, fmt.Errorf("Get_DOCOPTS_JSON '%s' is empty", "DOCOPTS_JSON")
}

func (d *Docopts) Print_keys() {
    parsed_args, err := d.Get_DOCOPTS_JSON()
    if err == nil {
        i := 0
        for k, _ := range parsed_args {
            if i > 0 {
                fmt.Fprintf(out, " ")
            }
            fmt.Fprintf(out, "%s", k)
            i++
        }
    } else {
        docopts_error("get-keys: %s", err)
    }
}

// Check if a value is an array
func IsArray(value interface{}) bool {
    // some golang tricks here using reflection to loop over the map[]
    rt := reflect.TypeOf(value)
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
        // this test is when called with DOCOPTS_JSON from Print_value()
        if IsArray(v) {
            arr := v.([]interface{})
            arr_out := make([]string, len(arr))
            for i, e := range arr {
                arr_out[i] = fmt.Sprintf("%v", e)
            }
            s = fmt.Sprintf("%s", strings.Join(arr_out[:],"\n"))
        } else {
            panic(fmt.Sprintf("To_bash():unsuported type: %v for '%v'", reflect.TypeOf(v), v ))
        }
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
        os.Exit(2)
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
                fmt.Println(strings.TrimSpace(Docopts_Version))
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
    var err error
    var b []byte

    switch(d.Format) {
    case F_Json:
        if d.Indent {
            b, err = json.MarshalIndent(arguments, "", "  ")
        } else {
            b, err = json.Marshal(arguments)
        }
        if err != nil {
            docopts_error("dump json: %v", err)
        }
        out.Write(b)
    case F_Ini:
        docopts_error("ini dump not supported yet", nil)
    }
}

// same as bash heler in docopts.sh docopt_get_help_string()
// extract help string from the file
func Get_help_string(filename string) []string {
    file, err := os.Open(filename)
    if err != nil {
        docopts_error("Get_help_string: %s", err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    found := false
    re, _ := regexp.Compile(`^(#\s?)[Uu]sage:`)
    prefix_l := 0
    var lines []string
    for scanner.Scan() {
        l := scanner.Text()
        if ! found {
            res := re.FindStringSubmatch(l)
            // check regexp match count
            if len(res) == 2 {
                found = true
                prefix_l = len(res[1])
            }
        }
        if found && l == "" {
            break
        }
        if found {
            if prefix_l < len(l) {
                lines = append(lines, l[prefix_l:])
            } else {
                lines = append(lines, "")
            }
        }
    }

    if err := scanner.Err(); err != nil {
        docopts_error("Get_help_string:scanner error: %s", err)
    }

    return lines
}

func (d* Docopts) Check_output_format(arguments docopt.Opts) {
    if arguments["parse"].(bool) || arguments["--json"].(bool) {
        d.Format = F_Json
    } else if name, err := arguments.String("-A"); err == nil {
        if ! IsBashIdentifier(name) {
            docopts_error("-A switch:%v", fmt.Errorf("not a valid Bash identifier: '%s'", name))
        }
        d.Format = F_Bash_assoc
        d.Output_declare = ! arguments["--no-declare"].(bool)
        d.Assoc_name = name
    } else if global_prefix, err := arguments.String("-G"); err == nil {
        d.Global_prefix = global_prefix
        d.Format = F_Bash_global
        d.Mangle_key = ! arguments["--no-mangle"].(bool)
    } else {
        // force default as JSON
        d.Format = F_Json
    }
}

func (d* Docopts) Dumper(bash_args docopt.Opts) {
    switch d.Format {
    case F_Json:
        d.Dump(bash_args)
    case F_Bash_assoc:
        d.Print_bash_assoc(bash_args)
    case F_Bash_global:
        d.Print_bash_global(bash_args)
    }
}

func (d* Docopts) Read_Doc_Version_from_stdin() {
    // read from stdin
    if d.Usage == "-" && d.Bash_Version == "-" {
        bytes, _ := ioutil.ReadAll(os.Stdin)
        arr := strings.Split(string(bytes), d.Separator)
        if len(arr) == 2 {
            d.Usage, d.Bash_Version = arr[0], arr[1]
        } else {
            msg := "error: help + version stdin, not found"
            if DEBUG {
                msg += fmt.Sprintf("\nseparator is: '%s'\n", d.Separator)
                msg += fmt.Sprintf("spliting has given %d blocs, exactly 2 are expected\n", len(arr))
            }
            panic(msg)
        }
    } else if d.Usage == "-" {
        bytes, _ := ioutil.ReadAll(os.Stdin)
        d.Usage = string(bytes)
    } else if d.Bash_Version == "-" {
        bytes, _ := ioutil.ReadAll(os.Stdin)
        d.Bash_Version = string(bytes)
    }
    d.Usage = strings.TrimSpace(d.Usage)
    d.Bash_Version = strings.TrimSpace(d.Bash_Version)
}

// split argv on colon ':' for old -h <msg> option format and OptionsFirst = false
// before_colon contains ':'
func Preprocess_agrv(argv []string) (before_colon []string, bash_argv []string) {
    // empty not nil array
    bash_argv = []string{}
    found_colon := false
    for _, v := range argv {
        if found_colon {
            bash_argv = append(bash_argv, v)
            continue
        }

        if v == ":" {
            found_colon = true
        }
        before_colon = append(before_colon, v)
    }
    return
}

func main() {
    // optimize get no docopt parsing
    if os.Args[1] == "get" {
        d := &Docopts{
            Json:  os.Getenv("DOCOPTS_JSON"),
        }
        d.Print_value(os.Args[2])
        os.Exit(0)
    }

    golang_parser := &docopt.Parser{
      OptionsFirst: false,
      SkipHelpFlags: true,
      HelpHandler: HelpHandler_golang,
    }

    argv, bash_argv := Preprocess_agrv(os.Args[1:])

    arguments, err_parse := golang_parser.ParseArgs(Usage, argv, Docopts_Version)

    if err_parse != nil {
        msg := fmt.Sprintf("mypanic: %v\n", err_parse)
        panic(msg)
    }

    DEBUG := arguments["--debug"].(bool)
    if DEBUG {
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
        Indent: false,
        Assoc_name: "",
    }

    // parse docopts's own arguments
    options_first := arguments["--options-first"].(bool)
    no_help :=  arguments["--no-help"].(bool)
    d.Separator = arguments["--separator"].(string)
    var err error

    // d.Bash_Version will be empty if error, so we dont care about error
    d.Bash_Version, _ = arguments.String("--version")
    d.Usage, err = arguments.String("--help")
    if err == nil && d.Usage == "-" {
        d.Read_Doc_Version_from_stdin()
    }

    // mode parse: read from another arg
    if d.Usage == "" && arguments["parse"].(bool) {
        d.Usage, err = arguments.String("<msg>")
        if err != nil {
            panic(err)
        } else if d.Usage == "-" {
            d.Read_Doc_Version_from_stdin()
        }
    }

    if ! arguments[":"].(bool) {
        // for non ':' colon new action arguments are in <argv>
        bash_argv = arguments["<argv>"].([]string)
    }

    if DEBUG {
        fmt.Printf("%20s : %v\n", "doc", d.Usage)
        fmt.Printf("%20s : %v\n", "bash_version", d.Bash_Version)
        fmt.Printf("%20s : %v\n", "argv", argv)
        fmt.Printf("%20s : %v\n", "bash_argv", bash_argv)
        d.debug_print_Docopts()
    }

    // ================================================================================
    // docopts JSON action modes
    // ================================================================================

    if arguments["fail"].(bool) {
        fmt.Println(d.Json)
        os.Exit(0)
    }

    if arguments["get"].(bool) {
        d.Print_value(arguments["<arg_name>"].(string))
        os.Exit(0)
    }

    if arguments["get-keys"].(bool) {
        // From DOCOPTS_JSON
        d.Print_keys()
        os.Exit(0)
    }

    if arguments["dump"].(bool) {
        if arguments["json"].(bool) {
            d.Format = F_Json
            d.Indent = true
        } else if arguments["ini"].(bool) {
            d.Format = F_Ini
        }
        if parsed_args, err := d.Get_DOCOPTS_JSON(); err == nil {
            d.Dump(parsed_args)
        }
        os.Exit(0)
    }

    if arguments["auto-parse"].(bool) {
        help_string := Get_help_string(arguments["<filename>"].(string))
        d.Usage = strings.Join(help_string[:], "\n")
        d.Format = F_Json
    }

	// ================================================================================
	// need to parse arguments after this point
	// ================================================================================

    // now parses bash program's arguments
    parser := &docopt.Parser{
      HelpHandler: d.HelpHandler_for_bash_eval,
      OptionsFirst: options_first,
      SkipHelpFlags: no_help,
    }

    bash_args, err_parse_bash := parser.ParseArgs(d.Usage, bash_argv, d.Bash_Version)
    if err_parse_bash != nil {
        panic(err_parse_bash)
    }

    // ========================================
    // arguments parsing modes
    // ========================================

    if DEBUG {
        print_args(bash_args, "bash")
        fmt.Println("----------------------------------------")
    }

    // outputer
    d.Check_output_format(arguments)
    d.Dumper(bash_args)

    os.Exit(0)
}
