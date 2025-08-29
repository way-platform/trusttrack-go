# AGENTS.md

This SDK provides a client and CLI tool for the ABAX Open API.

## Docs

When developing this SDK, use the API docs and specs:

- [Getting Started](./docs/getting-started.md)
- [OpenAPI Spec](./internal/oapi/abaxoapi/01-original.json)

## Structure

- The project uses a [tools](./tools/) directory with a separate Go module containing tools for building, linting and generating code.

- The project uses [Mage](https://magefile.org) with build tasks declared in [magefile.go](./tools/magefile.go).

- The client API returns protobuf messages. See [proto](./proto) for schemas.

## Developing

- Run tests with `./tools/mage test`

- Lint with `./tools/mage lint`

- Re-generate code with `./tools/mage generate`

- Leave all version control and git to the user/developer. If you see a build error related to having a git diff, this is normal.

## Project-specific guidance

- Always represent identifiers and continuation tokens as strings.

- One API operation per file. Naming convention is `client_<collection>_<operation>.go`.

- Use http.MethodGet constants instead of "GET", etc.

- The CLI tool organizes subcommands by entity using `cobra.Group`.
