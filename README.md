# clarifai/go

Clarifai's monolithic (for now) repository of services and packages written in Go.

# Go at Clarifai

## Setting up

1. Install Go according to the [documentation](https://golang.org/doc/install).
1. Set up environment variables in your .bash\_profile:

        export GOPATH=$HOME/work/go
        export PATH=$PATH:$GOPATH/bin
        export GO15VENDOREXPERIMENT=1  # See below about vendoring.

__Make sure that $GOROOT is not set in your environment__. Some installers will set it to the wrong
location.  It shouldn't be needed for our setup.

## Checking out the clarifai/go repo

Ideally we could just run `go get github.com/clarifai/go` like you would with any package.
Unfortunately, the go tool doesn't natively support fetching from private repos with ssh.  There
is a workaround described [here](https://gist.github.com/shurcooL/6927554) but it currently
(2015.10.22) has some issues with updating later, so the best approach for now is to just manually
clone:

    mkdir -p $GOPATH/src/github.com/clarifai; cd $GOPATH/src/github.com/clarifai
    git clone git@github.com:clarifai/go.git


## Vendoring

_Vendoring_ is an approach to dependency management to overcome Go’s original shortcomings when it
came to dependency management.  A good summary of the issues is in this short article on 
[Go 1.5 Vendoring](http://engineeredweb.com/blog/2015/go-1.5-vendor-handling/) and the official
Go wiki page on [vendoring tools](https://github.com/golang/go/wiki/PackageManagementTools).

Quick summary of what you need to know:

* Add this to your .bash\_profile: `export GO15VENDOREXPERIMENT=1`
* We’re using [Glide](https://github.com/Masterminds/glide) to handle importing and updating
vendored packages.
  * Install on MacOS with `brew update; brew install glide`.
  * When adding a new dependency, instead of `go get <package>`, do `glide get <package>`
* After checking out or pulling new code from a Clarifai repo, before building you need to install
the dependencies with:
`glide install`

## Running Go programs in Docker containers

First read [https://blog.golang.org/docker](https://blog.golang.org/docker).
The standard approach is to build the Go program inside the container.  The docker image
needs to have the Go tools installed.  Note that if using vendoring tools like Glide, those
need to be installed as well.

See this example
[Dockerfile](https://github.com/Clarifai/clarifai-go/blob/master/apinext/Dockerfile).


## Recommended reading

### Style guides:
* [Effective Go](https://golang.org/doc/effective_go.html).
  Official collection of idioms and best practices.
* The semi-official [code review guide](https://github.com/golang/go/wiki/CodeReviewComments)
  (like a style guide but not).

### Tutorials
* Official tour: [tour.golang.org](https://tour.golang.org)
* Go for pythonistas:  [slides](https://talks.golang.org/2013/go4python.slide#1),
  [talk](https://www.youtube.com/watch?v=elu0VpLzJL8).

### Advanced topics, required reading.
* [Go in production](http://peter.bourgon.org/go-in-production/) by Peter Bourgon
* Concurrency topics
   * [LearnConcurrency](https://github.com/golang/go/wiki/LearnConcurrency) and the related 
[codewalk](https://golang.org/doc/codewalk/sharemem/).
  Introduction to the philosophy of “Do not communicate by sharing memory; instead, share memory by communicating.”.
   * [Concurrency pipelines](http://blog.golang.org/pipelines) -- detailed tutorial on how to put
concurrency primitives together for coordinating asynchronous work.

### Testimonials
*  The Facebook Parse team wrote a few articles on
[rewriting their API in Go from Ruby](http://blog.parse.com/learn/how-we-moved-our-api-from-ruby-to-go-and-saved-our-sanity/) and
[their open source contributions](http://blog.parse.com/learn/parse-loves-go/)
* Author of Glide on [why he likes Go](http://engineeredweb.com//blog/2013/why-go-excellent-programming-language/)

## Tips, tricks, and tools.

* Use the [Go Playground](https://play.golang.org/) for experimenting and sharing code snippets.
* [Build constraints](https://golang.org/pkg/go/build/#pkg-overview). Can also use for testing.
* Static analysis tools
    * go vet
    * go fmt
    * go golint
* Really cool [commandline tool helper](https://github.com/codegangsta/cli). Handles flags, help, bash complete.
