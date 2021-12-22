# Grid

[![pkg.go.dev](https://pkg.go.dev/badge/github.com/anaseto/grid.svg)](https://pkg.go.dev/github.com/anaseto/grid)
[![godocs.io](https://godocs.io/github.com/anaseto/grid?status.svg)](https://godocs.io/github.com/anaseto/grid)

The **grid** module and package provide a generic implementation of a
two-dimensional matrix slice type. The API is inspired both by the standard Go
slice builtin functions and the image standard package.

It makes no assumptions on the type of the cells, and provides optimized
generic methods for grid iteration and manipulation.

Here's a simple example of usage:

``` go
// create a 5x5 grid of int type (default 0)
gd := grid.NewGrid[int](5, 5)
// Fill all the cells with 2
gd.Fill(2)
```

See the documentation for more examples.
