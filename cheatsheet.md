# Golang Cheatsheet

A communal collection of tips, tricks, and tools as we learn Go.

## Language features, idioms.

Concurrency

Closing channels


### Interfaces

An empty interface is like void\* in C++. Any type satisfies empty interface.
But in order to use it, need to cast.

An interface is represented as two words: pointer to itable (interface table,
method table like vtable) and pointer to data.
Further reading: 

* http://jordanorelli.com/post/32665860244/how-to-use-interfaces-in-go
* http://stackoverflow.com/questions/23148812/go-whats-the-meaning-of-interface 
* and testing: http://relistan.com/writing-testable-apps-in-go/

### Anonymous struct fields for struct composition.  Like mixins.

https://play.golang.org/p/AzYrp3HcDR

### Package initialization

[Package initialization](https://golang.org/ref/spec#Package_initialization)

    func init() {...}


### Empty structs

Empty struct as sentinel

* [playground example](https://play.golang.org/p/kL1OypyOZZ)
* [article](http://dave.cheney.net/2014/03/25/the-empty-struct)

### Empty interface ~ void pointer

An empty dict:

    d := make(map[interface{}]interface{})
    // or
    d := map[interface{}]interface{}{}

A set:

    map[string]struct{}

And other stuff.


### Defer

Example using defer with sleep: https://play.golang.org/p/7uOuUaFDhm


### Pass by value vs. pass by reference

If you're not used to languages that distinguish pointers/references from values,
this can trip you up.

See [this example](https://play.golang.org/p/vkvkj-Dpyd) of a bug where an object
is passed by value instead of reference, inadvertantly making a copy which breaks
the logic.


## Tools

### Testing

    go test -cover
    go test -coverprofile=cover.out
    go tool cover -html=cover.out

### Static analysis tools

    go vet
    go fmt
    go golint

### build constraints.

A special comment at the top of a file can contain build constraint labels, like:

    // +build foo

Then `go build -tags foo` will build that file, but `go build` will not.

See https://golang.org/pkg/go/build/#pkg-overview

This also works for tests, and could be a good way to run subsets of tests.

### Vim tools

[vim-go](http://blog.gopheracademy.com/vimgo-development-environment/):
The officially supported vim plugin.  Consolidates several earlier plugins.

For best autocomplete features, use it with YouCompleteMe. This also implies using
MacVim and not the standard vim on MacOs, because YCM wants a newer version of vim.

* Install MacVim following
  [this](https://github.com/macvim-dev/macvim/blob/master/README_mac.txt).
* Install [YouCompleteMe](https://github.com/Valloric/YouCompleteMe#mac-os-x-super-quick-installation)

Usage highlights:

    :GoInstallBinaries
    :GoUpdateBinaries # if things get stale
    :GoTest, :GoCoverage
    gd:  go to definition.
      <leader>gd -> godoc in vim
      <leader>gb -> godoc in browser
    <leader>s # go-implements, which interfaces does current symbol implement?
    :GoPlay or :<navigate>GoPlay # Copy to Go Playground, open in browser.
    snippets (haven’t figured it out yet)
    :GoCallers, :GoCalees

.vimrc:

    filetype plugin indent on
    " Whatever you want for <leader>
    let mapleader = " "
    
    " vim-go
    au FileType go nmap <leader>r <Plug>(go-run)
    au FileType go nmap <leader>b <Plug>(go-build)
    au FileType go nmap <leader>t <Plug>(go-test)
    au FileType go nmap <leader>c <Plug>(go-coverage)
    au FileType go nmap <Leader>s <Plug>(go-implements)
    au FileType go nmap <Leader>gb <Plug>(go-doc-browser)
    au FileType go nmap <Leader>i <Plug>(go-info)
    au FileType go nmap <Leader>e <Plug>(go-rename)

### Emacs tools

You're on your own. [FIXME!]


## Noteworthy packages


* https://github.com/codegangsta/cli

Really cool commandline tool helper, provides flags, help, bash complete.


