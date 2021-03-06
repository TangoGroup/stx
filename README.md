# stx

Go-based CLI for managing CloudFormation stacks written in Cue

## Installation

1. Clone the repo.
2. `cd` into the repo and run `go install ./...`
3. The `stx` binary will be placed in `$GOHOME/bin/stx` (e.g. `~/go/bin/stx`). Make sure this is in your `$PATH`

## Usage

`stx [global flags] <command> [command flags] [./... or specific cue files]`

If no args are present after <command>, stx will default to using `./...` as a way to find Cue files. This can be overriden with specific files: `stx print ./text.cue`

### Commands

- `add`        Writes scaffolding to template.cfn.cue
- `delete`     Deletes the stack along with .yml and .out.cue files
- `deploy`     Deploys a stack by creating a changeset, previews expected changes, and optionally executes.
- `diff`       DIFF against CloudFormation for the evaluted leaves.
- `events`     Shows the latest events from the evaluated stacks.
- `export`     Exports cue templates that implement the Stack pattern as yml files.
- `help`       Help about any command
- `import`     Imports an existing stack into Cue.
- `print`      Prints the Cue output as YAML
- `resources`  Lists the resources managed by the stack.
- `save`       Saves stack outputs as importable libraries to cue.mod
- `status`     Returns a stack status if it exists
- `notify`     Creates a light http server to listen for stack events from sns

### Roadmap

- Add color to yaml output of `print`
- Add sts, sdf, exe, and events commands
- Add config options to use ykman for automatic mfa retrieval
