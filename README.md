plugs
=====

`plugs` is a simple WIP package to ease the process of building Go plugins.

A simple usage looks like the following:

```go
p, err := plugs.Open("github.com/sbinet/plugs/testdata/p1")
if err != nil {
	panic(err)
}
v, err := p.Lookup("V")
f, err := p.Lookup("F")
*v.(*int) = 7
f.(func() string)() // returns "Hello, number 7"
```

`plugs` will invoke the Go compiler toolchain (via `go build -buildmode=plugin`) before trying to load the `*plugin.Plugin` via the `plugin` module.
