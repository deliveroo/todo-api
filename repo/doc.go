/*
Package repo contains a database persistence strategy for our domain models. All
interactions with the database should be implemented in this package; no
higher-level concept (e.g. request handlers) should execute database queries
directly.

The repo package contains unit tests that are executed against an actual
database instance, provided by Docker Compose. This gives us confidence that our
repo code actually works, ensuring that our queries are valid, and our model
mapping is correct.
*/
package repo
