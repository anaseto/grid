// This example demontrates how to create a new grid and do some simple
// manipulations.
package grid_test

import (
	"fmt"

	"github.com/anaseto/grid"
)

func ExampleGrid() {
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
}

func ExampleGridIterator() {
	// Create a new 26x2 grid of runes.
	gd := grid.NewGrid[rune](26, 2)
	// Get an iterator.
	it := gd.Iterator()
	// Iterate on the grid and fill it with successive alphabetic
	// characters.
	r := 'a'
	max := gd.Size()
	for it.Next() {
		it.SetV(r)
		r++
		if it.P().X == max.X-1 {
			r = 'A'
		}
	}
	// Print the grid.
	it.Reset()
	for it.Next() {
		fmt.Printf("%c", it.V())
		if it.P().X == max.X-1 {
			fmt.Print("\n")
		}
	}
	// Output:
	// abcdefghijklmnopqrstuvwxyz
	// ABCDEFGHIJKLMNOPQRSTUVWXYZ
}
