# begit

A beginner-friendly snapshot tool that follows the philosophy of Git.

`begit` lets you snapshot your working directory, browse and reset to past snapshots, and clean up files across all snapshots — without branches, staging areas, or commit hashes to memorize.

## Build

Requires Go 1.16+.

```sh
make build    # builds ./bin/begit
# or
make install  # installs to $GOPATH/bin
```

## Quickstart

```sh
cd my-project
begit init               # create .__begit__/
begit send "first try"   # snapshot the current directory

# modify your files ...
begit send "second try"  # another snapshot
begit select             # browse snapshots: diff / reset / remove
```

## Commands

| Command | Description |
|---|---|
| `init` | Initialize a begit repository in the current directory |
| `send "message"` | Create a snapshot |
| `select` | Interactively select and operate on snapshots |
| `remove-all <path>` | Remove a relative path from every snapshot |
| `md` | Export the directory as a Markdown report |
| `typ` | Export the directory as a Typst report |
| `help [command]` | Show help |

## How it works

- Snapshots live under `.__begit__/snapshots/<id>/`, where `<id>` is a UTC timestamp.
- Snapshot metadata is stored in `.__begit__.json`.
