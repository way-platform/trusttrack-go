# AGENTS.md

This SDK provides a client and CLI tool for the [TrustTrack FMS API](https://www.fmsdocumentation.com/apis/).

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

- Put all protobuf conversion functions in [./proto.go](./proto.go).

- Always handle enum conversion as separate functions in [./proto.go](./proto.go).

## Developing

- Run tests with `./tools/mage test`

- Lint with `./tools/mage lint`

- Re-generate code with `./tools/mage generate`

## OpenAPI schemas

- Use the generated types in [internal/oapi/ttoapi](./internal/oapi/ttoapi) to parse API responses.

- Always refer to the original OpenAPI spec in [internal/oapi/ttoapi/01-original.json](./internal/oapi/ttoapi/01-original.json) for API documentation and definitions.

## Protobuf schemas

- For enums:
  - always use `_UNSPECIFIED` as the first value. This value is always unused.
  - always use `_UNKNOWN` as the second value, used when the API value is not (yet) represented in the enum.
  - add a second `string unknown_<enum_name>` field, used when the API value is not (yet) represented in the enum.
  - if `UNKNOWN` is a valid value, represent this as a third `_NOT_AVAILABLE` value.
  - only put `ENUM_NAME_` prefix in front of `UNSPECIFIED`, `UNKNOWN` and `NOT_AVAILABLE` values.

- Always suffix fields with a physical unit, when applicable. E.g. `duration_s` and `odometer_km`.

## CLI Architecture

The CLI is split into two layers to keep credential storage pluggable:

```
cli/
├── cli.go       # Store interface, Credentials, FileStore, Options
└── command.go   # NewCommand() — full command tree
cmd/trusttrack/
└── main.go      # Thin wrapper: wires FileStore to XDG paths
```

- `cli.Store` — interface with `Read(any)`, `Write(any)`, `Clear()` methods
- `cli.NewCommand(...Option)` — builds the Cobra command tree; receives stores via functional options (`WithCredentialStore`)
- `cmd/trusttrack/main.go` — only wires `FileStore` instances and calls `cli.NewCommand()`

This separation lets consumers embed the CLI in a larger tool or swap the storage backend (e.g. use an in-memory store in tests, or a keychain-backed store) without forking.

### Embedding in a Parent CLI

The CLI can be embedded as a subcommand in a larger tool (e.g. a unified `way` CLI). Key design rules:

- **Never use `cmd.Root()`** — resolves to the parent CLI's root when embedded, breaking flag lookups. Use `cmd.Flags()` instead (works for both persistent and local flags).
- **`WithHTTPClient`** — the parent injects an `*http.Client` via `cli.WithHTTPClient()`. The SDK layers (auth, retry) stack on top of the injected client's transport.
- **`DebugTransport`** — exported in `debug.go` with a lazy `Enabled *bool` field. The parent owns the `--debug` flag and points `Enabled` at the flag variable. The transport checks the pointer at request time, solving the chicken-and-egg problem (transport constructed before flag parsing).

```go
var debug bool
cmd := cli.NewCommand(
    cli.WithCredentialStore(store),
    cli.WithHTTPClient(&http.Client{
        Transport: &trusttrack.DebugTransport{
            Enabled: &debug,
            Next:    http.DefaultTransport,
        },
    }),
)
cmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
```

### Module Structure

Three separate Go modules prevent Cobra/CLI dependencies from leaking into the SDK library:

```
go.mod              # SDK client library (no cobra, no CLI deps)
cli/go.mod          # CLI commands (depends on root SDK + cobra)
cmd/trusttrack/go.mod  # Standalone binary (depends on cli module)
```

Consumers who only need the Go client import the root module without pulling in CLI dependencies.

### Conventions

- Subcommands are organized by entity using `cobra.Group`
- Fully paginate results using the largest allowed page size
- Flat command structure: `objects`, `trips`, not `objects list`

## Code style

- Always represent identifiers and continuation tokens as strings.

- Use http.MethodGet constants instead of "GET", etc.

- Don't use pointers to `time.Time`. Use `.IsZero()` and the `omitzero` annotation.
