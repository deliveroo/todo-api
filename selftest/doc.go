/*
Package selftest implements an automated integration test strategy that launches
the application under as realistic conditions as possible:

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
*/
package selftest
