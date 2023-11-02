# `sortigo`
![Sortigo Image](sortigo.png)

`sortigo` is an opinionated tool for sorting imports in Go source files.
It is different from `goimports` in that it does not permit different import blocks amongst stdlib, third-party, and local packages.

## Installation
```bash
go install github.com/michael-james-holloway/sortigo
```

## Usage
```bash
Usage:
  sortigo format [flags]

Flags:
  -c, --check                     Checks if any files are different, and exits with a non-zero exit code if so.
      --dont-consolidate-blocks   Don't consolidate existing separate blocks of the same group type (e.g. multiple third party blocks).
  -h, --help                      help for format
  -l, --local-prefixes strings    Local prefix(es) to consider first party imports (e.g. github.com/michael-james-holloway/sortigo).
  -v, --verbose                   Verbose output.
  -w, --write                     Write the formatted file back to the original file.
```
