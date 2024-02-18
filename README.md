# got

A collection of Go Types and functions that I find myself creating in many projects.

In the spirit of "[A little copying is better than a little dependency][1]" - these packages have no
dependencies beyond the standard library, so they can be easily copied into a project directly.

## Packages

- `env` - Utilities for dealing with environment variables
- `ptr` - Creating pointers to literals and vice-versa.
- `fault` - Utilities for dealing with errors. Named so that it doesn't clash with the built-in `errors` package.
- `optional` - Implements an optional value type and some utility methods and functions to support it.
- `attempt` - Helper for calling functions with retry/backoff or timeout policy

[1]: https://www.youtube.com/watch?v=PAAkCSZUG1c&t=9m28s