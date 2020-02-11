# todo-api

[![CircleCI](https://img.shields.io/circleci/build/github/deliveroo/todo-api)](https://circleci.com/gh/deliveroo/todo-api/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/deliveroo/todo-api)](https://goreportcard.com/report/github.com/deliveroo/todo-api)
[![GoDoc](https://godoc.org/net/http?status.svg)](https://godoc.org/github.com/deliveroo/todo-api)
[![go.dev](https://img.shields.io/badge/go.dev-pkg-007d9c.svg?style=flat)](https://pkg.go.dev/github.com/deliveroo/todo-api)

`todo-api` is a small API web server built with Go. It supports user account
creation, login, and basic CRUD functions for tasks via a simple RESTful API.

While application functionality itself is contrived to serve as an example, this
repository does demonstrate a few practical patterns for developing with Go:

- Use Docker Compose to manage service dependencies (Postgres, Redis) for local
  development
- Providing fast testing and linting feedback with `modd`, which is critical to
  developer productivity and confidence
- Using environment vars as the source for application configuration whether in
  deployed testing environments
- An integration testing strategy where the API is treated as a "black box,"
  where the only interface is via HTTP request/response

## Why?

Ultimately, every application has different needs and requirements—there is no
one-true-way™ to build applications.

The features demonstrated by this application were converged upon after shipping
dozens of small to medium APIs in Go over the past several years and constantly
re-evaluating tooling choices and patterns.

We think it might be helpful reference point when you need a starting point for
a new app, or as an answer to "how are others doing this" when you find yourself
questioning your own approach.

We believe the curation of tools and patterns around the development of an
application is important to the long term health and success of a project. They
can greatly help or hinder aspects like developer onboarding, speed of feature
development, and confidence around making code changes.

## Dependencies

### Docker Compose

Docker Compose is used to run ephemeral dependencies without having to install
them directly on the host machine.

### Rake

We love `make`, but this project uses a `Rakefile` and thus requires on Ruby
installed on the host machine.

We've found over time that a `Makefile` can get unwieldy when managing a complex
application (e.g. database seeding, translation sync, code generation,
image/file asset management, etc).

Using `rake` and having the full power of Ruby for development tasks seems to
grow alongside a long living and complex application a bit better.

### Postgres

Postgres is used for application persistence, storing user accounts and tasks in
a straight-forward relational data model.

### Redis

Redis is used to store and persist user login sessions.

In practice, it would be completely acceptable to store sessions in memory or
the database (depending on your application requirements), but we have added
Redis as dependency mostly for demonstration purposes.

## Packages and application structure

### /api

Package `api` is the API server backed by
[jsonrest-go](https://github.com/deliveroo/jsonrest-go). It contains routing,
handlers, and middleware.

#### /api/protocol

Package `protocol` translates `domain` objects to their API response format.
This requires you to be deliberate about exposing attributes of your domain
objects and allows domain objects to evolve without implicitly affecting their
response format.

### /conf

Package `conf` holds all of the configuration for the application: e.g. database
connection strings, port to listen on, external credentials, environment. By
locating all configuration in one place, it's easy to see all parameters at a
glance. This approach implies that no other packages should access environment
variables directly.

### /domain

Package `domain` contains the domain models of our application. It does not
contain any persistence logic; instead, that belongs in the repo package.

Our services are generally microservices, and therefore have a small enough
scope for all of our models to live in one package. If you're building something
with a larger scope (i.e. more monolith than microservice), it may make sense to
break this into subpackages. However, make this decision very thoughtfully, and
ensuring that your subpackages truly are independent, or else you may find
yourself running into circular dependencies (an indication your packages were
more coupled than you thought, and potentially shouldn't have been split up).

Expect this package to have few unit tests. For the most part, these are just
models, without much logic. When do you do have logic in here (e.g. custom JSON
marshaling, sort functions), consider writing unit tests.

### /cmd/todo-api/apicmd

The `apicmd` package hoists configuration and application startup together. This
allows the [application entrypoint](./cmd/todo-api/main.go) to be remain lean
and makes it easier to configure and run the application from `selftest` for
integration testing.

### /pkg

We treat the `pkg` directory as a "staging ground" for self-contained,
extractable packages. Once a package here is needed by multiple projects, it's a
signal that the package may be worth elevating to its own repo with a strong
commitment to semantic versioning.

### /repo

Package `repo` contains a database persistence strategy for our domain models.
All interactions with the database should be implemented in this package; no
higher-level concept (e.g. request handlers) should execute database queries
directly.

The `repo` package contains unit tests that are executed against an actual
database instance, provided by Docker Compose. This gives us confidence that our
repo code actually works, ensuring that our queries are valid, and our model
mapping is correct.

### /service

The `service` directory is where you should find self-contained, but
application-specific functionality, particularly chunks which can be tested in
isolation.

### /selftest

Package `selftest` implements an automated integration test strategy that
launches the application under as realistic conditions as possible:

- It provides a real Redis and database connection
- It starts the actual API server listening on a port in the same manner as the
  application's package main
- It provides configuration via environmental variables
- It may provide external resources (e.g. 3rd party APIs) as fake
  implementations (not demonstrated by this project)

Once the API server is running, it's primarily tested as a black box; requests
are made via HTTP, and responses are validated. The black box is only penetrated
for test setup (e.g. seeding) and teardown (e.g. wiping the database in between
tests).

This approach gives us extremely high confidence that our application will
behave properly in production, since it was tested with the full launch process
(including parsing env vars!) and real dependencies, no mocks. Because there's
absolutely no dependency on internal implementation details, you can refactor at
will without changing your tests so long as you maintain your API contracts.

## What's missing

A glaring omission from this project is both application tracing and logging
(!). We may layer it in in a future update.
