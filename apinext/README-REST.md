# Building REST in Go: Design Doc

This doc outlines our preferred approach for building an HTTP REST service in Go.

There is significant overlap between the design of a REST service and an RPC service; they share many features common in a general service-oriented-architecture.  But those issues are discussed elsewhere.  Hence:

Not in scope for this doc:
* Generic build, test, and deployment issues, unless they are specific to HTTP/REST services.
* Generic service-oriented-architecture issues (RPC framework, logging, alerting, service discovery), unless they are specific to HTTP/REST services.

## HTTP Routing
There are many HTTP router/multiplexer (mux) packages in Go.

Requirements:
* standard http.Handler interface, so it plays well with middleware, testing, go-kit and other packages.
* flexible routing:  regex, variables, subrouters.  easy to write handlers
* fast
* swagger-friendly?

Current favorites:
*  [goji](https://github.com/zenazn/goji)
*  [bone](https://github.com/go-zoo/bone)
*  [httprouter](https://github.com/julienschmidt/httprouter)

Others: gorilla mux, Gin/Martini, Negroni.

A simple example with three different routers implementing the same endpoints is at
[//apinext/main.go](https://github.com/Clarifai/go/blob/4ccb64920c603735f1ae81126eedab4d7ea6063e/apinext/main.go).

Notes:

[Goji](https://goji.io/) (2792 stars, 183 forks, 16 contribs, 9 issues)
* minimal, Flask-like philosophy.
* parameterized URLs, regex
* Composable.  Extensible middleware.  compatible with stdlib net/http handlers.
* list of [contributed middleware](https://github.com/zenazn/goji/wiki/Third-Party-Libraries) including csrf, gzip, cors, httpauth, sessions graceful shutdown
* Fast, does well in [http routing benchmarks](https://github.com/julienschmidt/go-http-routing-benchmark).
* Not a standalone router, might not be go-kit friendly?
* goji/web web.C context object.  contributed middleware to bridge that to x/net/context (used by go-kit Endpoint).
web.Mux satisfies net/http.Handler.

[httprouter](https://github.com/julienschmidt/httprouter) (2313 stars, 175 forks, 16 contribs, 13 issues)
* dominates benchmarks
* not fully compatible w/ http.Handler; not as middleware friendly. → Deal killer.
* docs not quite as nice as Goji

[bone](https://github.com/go-zoo/bone)  (807 stars, 42 forks, 11 contribs, 0 issues)
* URL params, regex, subrouters
* implements http.Handler
* good benchmark performance (limited set of benchmarks), on par with httprouter on some, but [HN review](https://news.ycombinator.com/item?id=8737574) shows it lagging httprouter by 10x on some.
* See [my favorite go multiplexer](http://www.peterbe.com/plog/my-favorite-go-multiplexer)

[Gorilla mux](http://www.gorillatoolkit.org/) (1709 stars, 284 forks, 26 contribs, 12 issues)
* URL reverse
* nested subrouters
* Router is an http.Handler.
* significantly worse benchmark performance compared to httprouter and bone.
* More features, more complicated API.  perhaps bloated.

Also-rans:
* Martini.  batteries included:  authentication, static routing.
* [Negroni](https://github.com/codegangsta/negroni). web middle framework, net/http comptible, router-agnostic. same author as martini. 3rd-party middleware includes:  restgate http header auth, oauth2, cors, logging, gzip, render.  (most of these are Goji-compatible too?)
* Revel. bigger web framework.  Peter Bourgon considers harmful.

Related, but not an HTTP router:
[cmux](https://github.com/soheilhy/cmux) - A connection multiplexer, to serve multiple protocols on the same port, e.g. gRPC, HTTP.

__2015.10.29__:  Go with Goji for first experiments.  Standard http.Handler, good enough benchmarks, clean interface, and actively developed.
Schema validation
Swagger support?

[gorilla/schema](http://www.gorillatoolkit.org/pkg/schema)

## Proxying to v1
Until APIv1 is shut down (which may be a long time…) we need a strategy for serving v1 traffic alongside v2.  Furthermore, during the migration/development period, some v2 functionality will be implemented by putting a translation layer in front of existing v1 endpoints and new endpoints in the v1 binary.

Options for serving v1 traffic include:

1. Point all v1 traffic to a new v2 binary that performs transparent passthrough routing to the v1 servers.  Use DNS for discovery (initially api.clarifai.com, moving to api-1.clarifai.com when the new binary is ready to assume api.clarifai.com traffic).
    1. pro: Can wrap the proxied requests in new middleware, e.g. logging.
    1. pro: Can do forking/mirroring
    1. pro: Doesn’t change the existing v1 deployment/serving story.
    1. con: ?  Extra traffic through the new binary, with some undesirable legacy characteristics (long synchronous requests).

![Alt text](http://g.gravizo.com/g?
@startuml;
interface "v1 traffic" as v1in;
interface "v2 traffic" as v2in;
cloud ELB;
node "API-Next\\nREST" as APIv2;
node APIv1 %23lightblue;
node "Other services" as Other;
v1in --> ELB;
v2in --> ELB;
ELB --> APIv2;
APIv2 --> APIv1 : v1 passthrough;
APIv2 --> APIv1 : v2 translated;
APIv2 --> Other : v2 new backends;
@enduml;
)


1. Put both v1 and v2 servers behind a new proxy (nginx or haproxy) that routes v1/v2 to their respective servers.
    1. con: Introduces a new proxy component that needs to be managed and configured.
    1. con: Introduces new service discovery requirements for the proxy to discover the REST backends; needs to support blue/green traffic ramping.
    1. pro: extensible, could add more REST servers for new endpoints over time.
    1. pro: haproxy and nginx are battle-tested and well supported.

![Alt text](http://g.gravizo.com/g?
@startuml;
interface "v1 traffic" as v1in;
interface "v2 traffic" as v2in;
cloud ELB;
node Proxy;
node "API-Next\\nREST" as APIv2;
node APIv1 %23lightblue;
node "Other services" as Other;
v1in --> ELB;
v2in --> ELB;
ELB --> Proxy;
Proxy --> APIv2;
Proxy --> APIv1 : v1 traffic;
APIv2 --> APIv1 : v2 translated;
APIv2 --> Other : v2 new backends;
@enduml;
)

__2015.10.29__:  Let’s try (1) first and see how it feels.

Other proxy-related issues:
* Authentication: Auth forwarding, or use an internal key for proxied requests, or authentication-free internal service?
    * -> In passthrough mode, auth tokens are simply forwarded.
    * -> In translation mode, we can authenticate tokens to the new user/auth service once it’s available, and meanwhile forward them until that’s ready.
* When proxying, copy all other http headers, accept-encoding, X-Forwarded-For, Host.  To preserve nginx rate limiting and server_name, django allowed hosts, etc.
* Avoid double logging and throttle/usage counting?  Once we start logging in the v2 binary, we could set should_log_requests=false.  However, it might make sense to double-log because the v2 log format will be different, and having both copies may help ensure we’re logging the right stuff and the translation layer is working right.
* Lame-ducking: long-running requests sent to v1 may not complete within the v2 service's
  lame-duck window.  We can either make the window long enough to avoid this, or just fail those
  requests until the vision/predict service is refactored around an async work queue.

## API discovery and self-documenting APIs

[TODO: investigate Swagger support in Go]

## Middleware

On one hand, writing http middleware is straightforward in Go -- you can simple write an http.Handler that wraps another Handler and chain them.  See [here](https://justinas.org/writing-http-middleware-in-go/) for a discussion with some simple examples.  The author of that article has written several open-source middleware packages (some listed below) and a simple chaining wrapper called [Alice](https://github.com/justinas/alice).

Some example open-source middleware projects:
* Authentication:  OAuth2, jwt?
* CSRF anti-spoofing:  [nosurf](https://github.com/justinas/nosurf)
* Logging



## Other interesting things

Configurable proxy server: [goproxy](https://github.com/elazarl/goproxy)

