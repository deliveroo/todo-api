/*
Package domain contains the domain models of our application. It does not
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
*/
package domain
