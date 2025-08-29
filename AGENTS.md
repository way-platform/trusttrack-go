# AGENTS.md

This SDK provides a client and CLI tool for the ABAX Open API.

## Docs

When developing this SDK, use the API docs and specs:

- [Getting Started](./docs/getting-started.md)
- [OpenAPI Spec](./internal/oapi/abaxoapi/01-original.json)

## Git

- Leave all version control and git to the user/developer.

- If you see a build error related to having a git diff, this is normal.

## Structure

- The project uses a [tools](./tools/) directory with a separate Go module containing tools for building, linting and generating code.

- The project uses [Mage](https://magefile.org) with build tasks declared in [magefile.go](./tools/magefile.go).

- The client API returns protobuf messages. See [proto](./proto) for schemas.

- One API operation per file. Naming convention is `client_<collection>_<operation>.go`.

## Developing

- Run tests with `./tools/mage test`

- Lint with `./tools/mage lint`

- Re-generate code with `./tools/mage generate`

## OpenAPI schemas

- Use the generated types in [internal/ttoapi](./internal/ttoapi) to parse API responses.

## Protobuf schemas

- For enums:

  - always use `_UNSPECIFIED` as the first value. This value is always unused.
  - always use `_UNKNOWN` as the second value, used when the API value is not (yet) represented in the enum.
  - add a second `string unknown_<enum_name>` field, used when the API value is not (yet) represented in the enum.
  - if `UNKNOWN` is a valid value, represent this as a third `_NOT_AVAILABLE` value.
  - only put `ENUM_NAME_` prefix in front of `UNSPECIFIED`, `UNKNOWN` and `NOT_AVAILABLE` values.

- Always suffix fields with a physical unit, when applicable. E.g. `duration_s` and `odometer_km`.

## CLI tool

- The CLI tool organizes subcommands by entity using `cobra.Group`.

- The CLI tool always fully paginates results and uses the largest allowed page size.

## Code style

- Always represent identifiers and continuation tokens as strings.

- Use http.MethodGet constants instead of "GET", etc.
