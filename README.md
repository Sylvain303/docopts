# docopts

docopt for shell - make beautifull CLI with ease.

This branch is an **non-mergable branch**, it is for debuging purpose and sharing issue's code.

Main code is here: https://github.com/docopt/docopts

To use this branch, your code must reside in original golang source lib `$GOPATH/src/github.com/docopt/docopts`.
So we add a remote branch for this purpose:

```
go get https://github.com/docopt/docopts
cd $GOPATH/src/github.com/docopt/docopts
git remote add debug https://github.com/Sylvain303/docopts
git fetch debug debug-issues
git checkout debug/debug-issues
```
Ready!
