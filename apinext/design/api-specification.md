# API Next Specification

This document describes the standards and conventions for specifying (and documenting) the Clarifai REST api.
It also describes related uses of the Swagger representation of the API specification. These include:
+ schema validation of http/REST request payloads
+ schema validation of http/REST response payloads
+ interactive discovery and exploration of the API (e.g. using the Swagger UI or our own version of that functionality)

## Specify the API using Swagger
We specify the REST endpoints using [Swagger](http://swagger.io/).

## Host the Swagger specification in a well-known place
e.g. api.clarifai.com/v2/swagger.  

# Related Uses of Swagger api spec

## Implement REST endpoint middleware to validate incoming and outgoing payloads
tbd

## Interactive discovery and exploration of API
Tools like the Swagger UI allow users to have an easy web-based
experience in which you can both learn what the API offers but
also to try it out interactively. Dan's current plan is to
write our version of the Swagger UI functionality, but we
want to offer the same types of capabilities.

# Related Issues

## Semantic Versioning
We've discussed implementing [semantic versioning](http://semver.org/). ~~We haven't really decided on whether or exactly
how we want to proceed.~~

We simply use /v2 in the api endpoint uri.

+ in the continuous delivery environment we are building, MAJOR.MINOR.PATCH seems too noisy (dan,jim).
+ MAJOR.MINOR can work, but does it add value over MAJOR i.e. /v2? (nope - see [discussion](https://github.com/Clarifai/go/pull/6/files#r44095347))


Issues:
+ we're addressing how the swagger spec is versioned along with the software, open is how *other* docs are versioned alongside

## Version Management
The swagger spec needs to be versioned along with the REST endpoint services. Changes to
the published spec *must* be accompanied by new software release that implements the API
as documented. Verifying this is a key element of acceptance testing for a new version.

Of course it's possible to implement REST endpoints that are *not documented*.

The official public API version (i.e. /v2 in the URL) is only bumped when we need to make breaking, non-backwards compatible changes.  Otherwise we are free to make backwards-compatible API changes within a single version.

We will reserve the right to make breaking, non-backwards compatible changes until a given feature is released to full Production.
The idea that once we release a feature to Production that we are committing to not change/break it will be challenging, but it
is important for customer/developer satisfaction.

Beyond that, the requirements for versioned releases are as follows:

* we can push freely without heavyweight communication / approval from external alpha/beta partners
* we have a predictable deprecation lifecycle, and communication mechanism, so that alpha clients can manage transitions when the API changes.
* documentation is versioned along with the releases
* we can push unannounced stuff to production and keep it secret, but let internal and alpha partners use it.
* internal developers can very easily find out exactly what code is running on a given server.


## Feature Flags
Feature flags are how we can push features to production but keep them secret generally while allowing selective access
to specific users and groups. There's a rough draft of a feature flags design [here](https://docs.google.com/document/d/1BAlix5klJ4EYIrwylnXNjaAfhmPfDHPI987zy3xHNP4/edit)

## Serving Swagger spec as a template parameterized by feature flags
It might be cool to be able to somehow annotate the swagger api spec with feature flags, and then be able to
serve a different version that includes all and only those endpoints / features that are enabled for the user based
on the feature flag settings.

The swagger spec describing features enabled for Production will be served by default. If an auth token is provided
(as it can be when the user is using Dev Hub Next and is logged in), then we can serve the version rendered for the
logged in user's enabled features.


# Parking Lot (Deferred)

## User and Application Profiles API Version
We've discussed adding an attribute to the application profile to set the version of the api to be used by the application.
For now, we are not pursuing this idea.
