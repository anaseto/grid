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
// Create a new 10x10 grid of runes.
gd := grid.NewGrid[rune](10, 10)
// Fill the whole grid with dots.
gd.Fill('.')
// Define a range (3,3)-(7,7).
rg := grid.NewRange(3, 3, 7, 7)
// Define a slice of the grid using the range.
rectangle := gd.Slice(rg)
// Fill the rectangle with #.
rectangle.Fill('#')
// Print the grid.
it := gd.Iterator()
max := gd.Size()
for it.Next() {
	fmt.Printf("%c", it.V())
	if it.P().X == max.X-1 {
		fmt.Print("\n")
	}
}
// Output:
// ..........
// ..........
// ..........
// ...####...
// ...####...
// ...####...
// ...####...
// ..........
// ..........
// ..........
```

See the documentation for more examples.
