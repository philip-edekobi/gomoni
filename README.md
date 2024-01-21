# gomoni

gomoni is a tool which mimics the behavior of nodemon. It takes an optional dir argument and keeps track of all files which are members of packages used in the project.

It excludes test files and any file that is never used in the project.

## Installation

```bash
$ go install github.com/philip-edekobi/gomoni@latest
```

## Usage

```bash
$ gomoni [folder]
```

`folder` is optional if the current directory is the target folder
