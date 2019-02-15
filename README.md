# gotty
A dotty getter for golang

Getting Started
===============

## Installing

To start using GJSON, install Go and run `go get`:

```sh
$ go get -u github.com/vingofe/gotty
```

This will retrieve the library.

## Get a value

```go
package main

import "github.com/vingofe/gotty"

type Object struct {
	Key interface{}
}

o := Object{
    Key: "foo",
}

o1 := Object{
    Key: o,
}

func main() {
	value := gotty.Get(1, "Key.Key")
	// "foo"
}
```

## License

gotty source code is available under the MIT [License](/LICENSE).