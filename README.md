# go-module

Parse a git remote URL into a Go module [`Path`](module.go), [`Name`](module.go), [`Identifier`](module.go), and [`EnvPrefix`](module.go) — pure string logic, no I/O.

## Install

```sh
go get github.com/gomatic/go-module
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/gomatic/go-module"
)

func main() {
	path, err := module.Parse(module.Remote("git@github.com:org/repo.git"))
	if err != nil {
		panic(err)
	}

	name := path.Repo()

	fmt.Println(path)              // github.com/org/repo
	fmt.Println(name)              // repo
	fmt.Println(name.Identifier()) // repo
	fmt.Println(name.EnvPrefix())  // REPO
}
```

`Parse` accepts the scp-like SSH form (`git@host:org/repo.git`) and URL forms (`https://host/org/repo.git`, `ssh://git@host/org/repo`), stripping the scheme, any userinfo, and a trailing `.git`. It returns [`ErrInvalidRemote`](errors.go) when the result is not a `host/org/repo` path.
