Bash version of the [Pythonic command-line interface parser](http://docopt.org/) that will make you smile.

# WARNING: a new shell API is comming. 

* See Proposal lib API: [The Wiki](https://github.com/docopt/docopts/wiki)
* The [develop branch](https://github.com/docopt/docopts/tree/develop) is **abandoned**.

## docopts

See the Reference manual for the command line `docopts` (old
[README.rst](old_README.rst))

* `docopts` is the python wrapper that outputs valid Bash commands as the result. In this documention it is also refered as docopts.py
* `docopts.sh` is a Bash script that wraps docopts.py, provinding a full shell api, and also embed the python code in its source.

## Features

The current command line tools `docopts`, written in python, is maintained. A new
shell lib is added. 

The `docopts.sh` is a bash library you can source into your CLI script. It
automaticaly embed docopts.py and docopt.py, and is standalone. Just drop it
and source it.

Of course, it needs a python interperter in the $PATH.

You can still use `docopts` directly or both.

## Examples

More [examples/](examples/).

Those examples doesn't fully work for now. $verbose is not recognized.

### docopts.py

~~~bash
eval "$(docopts -V - -h - : "$@" <<EOF
Usage: rock [options] <argv>...

      --verbose  Generate verbose messages.
      --help     Show help options.
      --version  Print program version.
----
rock 0.1.0
Copyright (C) 200X Thomas Light
License RIT (Robot Institute of Technology)
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
EOF
)"

if $verbose ; then
    echo "Hello, world!"
fi
~~~

### docopts.sh

changes sourcing the code + `eval "$(docopt …)` <== without the **s**, its the shell function.

~~~bash
#!/bin/bash
source ../docopts.sh
eval "$(docopt -V - -h - : "$@" <<EOF
Usage: rock [options] <argv>...

      --verbose  Generate verbose messages.
      --help     Show help options.
      --version  Print program version.
----
rock 0.1.0
Copyright (C) 200X Thomas Light
License RIT (Robot Institute of Technology)
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
EOF
)"

if $verbose ; then
    echo "Hello, world!"
fi
~~~



## Shell helpers

~~~bash
# auto extract the Usage string from the top shell script comment
# ie: help=$(sed -n -e '/^# Usage:/,/^$/ s/^# \?//p' < $0)
help_string=$(docopt_get_help_string $0)

# if the option as multiples values, you can get it into an array
array_opt=( $(docopt_get_values args --multiple-time) )
~~~

## Sugar: --auto

Sugar taste good, it is just to make you smile, more.

The `--auto` feature:

`--auto` is for good-lazy sysadmin who codes small snippet of scripts but still want an option parser or doc inside the script header. It happens to me daily refactoring legacy scripts.

soucring `docopts.sh --auto "$@"` gives you an automatic parsing **and** eval of the docopt comment.

`$args` is evaled and pushed into global scope magically and ready for your script.

~~~bash
#!/bin/bash
# Usage: doit [cmd] FILE...
#
# do somethingdo something

source ../docopts.sh --auto "$@"

for a in ${!args[@]} ; do
    echo "$a = ${args[$a]}"
done
~~~


## Developpers

If you want to clone this repository and hack docopts.

Use `git clone --recursive`, to get submodules.

If you forgot, you can also run:

~~~bash
git submodule init
git submodule update
~~~

Folder structure:

~~~
.
├── API_proposal.md - doc See wiki - to be removed
├── build.sh        - build the embedded docopts.py into docopts.sh
├── docopt.py       - original copy of docopt.py
├── docopts         - current python wrapper - almost unmodified
├── docopts.py      - copy of docopts, See build.sh
├── docopts.sh      - bash lib - already embed both docopt.py + docopts.py
├── examples
│   ├── calculator_example.sh
│   ├── cat-n_wrapper_example.sh
│   ├── docopts_auto_examples.sh
│   └── quick_example.sh
├── language_agnostic_tester.py
├── setup.py
├── testcases.docopt
├── testee.sh
└── tests
    ├── bats/                - git submodules
    ├── bats.alias           - source it to have bats working
    ├── docopts-auto.bats    - unit test for --auto
    ├── docopts.bats         - unit test docopts.sh
    └── exit_handler.sh      - helper
~~~
