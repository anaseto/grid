// Package grid provides a generic two-dimensional matrix slice type.
//
// The package provides utilities for iterating such grids, as well as
// manipulating positions and ranges.
//
// The API is inspired by the standard Go slice builtin functions and the image
// standard package.
package grid

import (
	"fmt"
)

// Point represents an (X,Y) position in a grid.
//
// It follows conventions similar to the ones used by the standard library
// image.Point.
type Point struct {
	X int
	Y int
}

// String returns a string representation of the form "(x,y)".
func (p Point) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

// Shift returns a new point with coordinates shifted by (x,y). It's a
// shorthand for p.Add(Point{x,y}).
func (p Point) Shift(x, y int) Point {
	return Point{X: p.X + x, Y: p.Y + y}
}

// Add returns vector p+q.
func (p Point) Add(q Point) Point {
	return Point{X: p.X + q.X, Y: p.Y + q.Y}
}

// Sub returns vector p-q.
func (p Point) Sub(q Point) Point {
	return Point{X: p.X - q.X, Y: p.Y - q.Y}
}

// In reports whether the position is within the given range.
func (p Point) In(rg Range) bool {
	return p.X >= rg.Min.X && p.X < rg.Max.X && p.Y >= rg.Min.Y && p.Y < rg.Max.Y
}

// Mul returns the vector p*k.
func (p Point) Mul(k int) Point {
	return Point{X: p.X * k, Y: p.Y * k}
}

// Div returns the vector p/k.
func (p Point) Div(k int) Point {
	return Point{X: p.X / k, Y: p.Y / k}
}

// Range represents a rectangle in a grid that contains all the positions P
// such that Min <= P < Max coordinate-wise. A range is well-formed if Min <=
// Max. When non-empty, Min represents the upper-left position in the range,
// and Max-(1,1) the lower-right one.
//
// It follows conventions similar to the ones used by the standard library
// image.Rectangle.
type Range struct {
	Min, Max Point
}

// NewRange returns a new Range with coordinates (x0, y0) for Min and (x1, y1)
// for Max. The returned range will have minimum and maximum coordinates
// swapped if necessary, so that the range is well-formed.
func NewRange(x0, y0, x1, y1 int) Range {
	if x1 < x0 {
		x0, x1 = x1, x0
	}
	if y1 < y0 {
		y0, y1 = y1, y0
	}
	return Range{Min: Point{X: x0, Y: y0}, Max: Point{X: x1, Y: y1}}
}

// String returns a string representation of the form "(x0,y0)-(x1,y1)".
func (rg Range) String() string {
	return fmt.Sprintf("%s-%s", rg.Min, rg.Max)
}

// Size returns the (width, height) of the range in cells.
func (rg Range) Size() Point {
	return rg.Max.Sub(rg.Min)
}

// Shift returns a new range with coordinates shifted by (x0,y0) and (x1,y1).
func (rg Range) Shift(x0, y0, x1, y1 int) Range {
	rg = Range{Min: rg.Min.Shift(x0, y0), Max: rg.Max.Shift(x1, y1)}
	if rg.Empty() {
		return Range{}
	}
	return rg
}

// Line reduces the range to relative line y, or an empty range if out of
// bounds.
func (rg Range) Line(y int) Range {
	if rg.Min.Shift(0, y).In(rg) {
		rg.Min.Y = rg.Min.Y + y
		rg.Max.Y = rg.Min.Y + 1
	} else {
		rg = Range{}
	}
	return rg
}

// Lines reduces the range to relative lines between y0 (included) and y1
// (excluded), or an empty range if out of bounds.
func (rg Range) Lines(y0, y1 int) Range {
	nrg := rg
	nrg.Min.Y = rg.Min.Y + y0
	nrg.Max.Y = rg.Min.Y + y1
	return rg.Intersect(nrg)
}

// Column reduces the range to relative column x, or an empty range if out of
// bounds.
func (rg Range) Column(x int) Range {
	if rg.Min.Shift(x, 0).In(rg) {
		rg.Min.X = rg.Min.X + x
		rg.Max.X = rg.Min.X + 1
	} else {
		rg = Range{}
	}
	return rg
}

// Columns reduces the range to relative columns between x0 (included) and x1
// (excluded), or an empty range if out of bounds.
func (rg Range) Columns(x0, x1 int) Range {
	nrg := rg
	nrg.Min.X = rg.Min.X + x0
	nrg.Max.X = rg.Min.X + x1
	return rg.Intersect(nrg)
}

// Empty reports whether the range contains no positions.
func (rg Range) Empty() bool {
	return rg.Min.X >= rg.Max.X || rg.Min.Y >= rg.Max.Y
}

// Eq reports whether the two ranges contain the same set of points. All empty
// ranges are considered equal.
func (rg Range) Eq(r Range) bool {
	return rg == r || rg.Empty() && r.Empty()
}

// Sub returns a range of same size translated by -p.
func (rg Range) Sub(p Point) Range {
	rg.Max = rg.Max.Sub(p)
	rg.Min = rg.Min.Sub(p)
	return rg
}

// Add returns a range of same size translated by +p.
func (rg Range) Add(p Point) Range {
	rg.Max = rg.Max.Add(p)
	rg.Min = rg.Min.Add(p)
	return rg
}

// Intersect returns the largest range contained both by rg and r. If the two
// ranges do not overlap, the zero range will be returned.
func (rg Range) Intersect(r Range) Range {
	if rg.Max.X > r.Max.X {
		rg.Max.X = r.Max.X
	}
	if rg.Max.Y > r.Max.Y {
		rg.Max.Y = r.Max.Y
	}
	if rg.Min.X < r.Min.X {
		rg.Min.X = r.Min.X
	}
	if rg.Min.Y < r.Min.Y {
		rg.Min.Y = r.Min.Y
	}
	if rg.Min.X >= rg.Max.X || rg.Min.Y >= rg.Max.Y {
		return Range{}
	}
	return rg
}

// Union returns the smallest range containing both rg and r.
func (rg Range) Union(r Range) Range {
	if rg.Max.X < r.Max.X {
		rg.Max.X = r.Max.X
	}
	if rg.Max.Y < r.Max.Y {
		rg.Max.Y = r.Max.Y
	}
	if rg.Min.X > r.Min.X {
		rg.Min.X = r.Min.X
	}
	if rg.Min.Y > r.Min.Y {
		rg.Min.Y = r.Min.Y
	}
	return rg
}

// Overlaps reports whether the two ranges have a non-zero intersection.
func (rg Range) Overlaps(r Range) bool {
	return !rg.Intersect(r).Empty()
}

// In reports whether range rg is completely contained in range r.
func (rg Range) In(r Range) bool {
	return rg.Intersect(r) == rg
}

// Iter calls a given function for all the positions of the range.
func (rg Range) Iter(fn func(Point)) {
	for y := rg.Min.Y; y < rg.Max.Y; y++ {
		for x := rg.Min.X; x < rg.Max.X; x++ {
			p := Point{X: x, Y: y}
			fn(p)
		}
	}
}

// Grid represents a two-dimensional matrix of values of any type. It is a
// slice type, so it represents a rectangular range within an underlying
// original grid. Due to how it is represented internally, it is more efficient
// to iterate in row-major order, as in the following pattern:
//
//	max := gd.Size()
//	for y := 0; y < max.Y; y++ {
//		for x := 0; x < max.X; x++ {
//			p := Point{X: x, Y: y}
//			// do something with p and the grid gd
//		}
//	}
//
// Most iterations can be performed using the Slice, Fill, Copy, Map and Iter
// methods. An alternative choice is to use the Iterator method.
//
// Grid elements must be created with NewGrid.
type Grid[T any] struct {
	ug *grid[T] // underlying whole grid
	rg Range    // range within the whole grid
}

type grid[T any] struct {
	Cells []T
	Width int
}

// NewGrid returns a new grid with given width and height in cells. The width
// and height should be positive or null. The new grid contains all positions
// (X,Y) with 0 <= X < w and 0 <= Y < h. The grid is filled with the zero
// value for cells.
func NewGrid[T any](w, h int) Grid[T] {
	if w < 0 || h < 0 {
		panic(fmt.Sprintf("negative dimensions: NewGrid(%d,%d)", w, h))
	}
	gd := Grid[T]{}
	gd.ug = &grid[T]{}
	gd.rg.Max = Point{w, h}
	gd.ug.Width = w
	gd.ug.Cells = make([]T, w*h)
	return gd
}

// NewGridFromSlice builds a grid of width w with initial contents provided by
// slice s.  The slice's length should be a multiple of w. The slice's values
// are used in row-major order.
func NewGridFromSlice[T any](s []T, w int) Grid[T] {
	if w < 0 {
		panic(fmt.Sprintf("negative width: %d", w))
	}
	if w == 0 && len(s) > 0 {
		panic(fmt.Sprintf("zero width but len(s) > 0: %d", len(s)))
	}
	if w > 0 && len(s)%w != 0 {
		panic(fmt.Sprintf("bad length: %d (expected %d%%%d == 0, but got %d)", len(s), len(s), w, len(s)%w))
	}
	h := 0
	if w != 0 && len(s) > 0 {
		h += len(s) / w
	}
	gd := Grid[T]{}
	gd.ug = &grid[T]{}
	gd.rg.Max = Point{w, h}
	gd.ug.Width = w
	gd.ug.Cells = s
	return gd
}

// Bounds returns the range that is covered by this grid slice within the
// underlying original grid.
func (gd Grid[T]) Bounds() Range {
	return gd.rg
}

// Contents returns the grid's current underlying slice with the values of the
// whole underlying grid, in row-major order.
func (gd Grid[T]) Contents() []T {
	if gd.ug == nil {
		return nil
	}
	return gd.ug.Cells
}

// Range returns the range with Min set to (0,0) and Max set to gd.Size(). It
// may be convenient when using Slice with a range Shift.
func (gd Grid[T]) Range() Range {
	return gd.rg.Sub(gd.rg.Min)
}

// Cap returns the size (w,h) measuring the grid and the available space past
// it within the underlying whole grid. In other words,
// gd.Bounds().Min.Add(Point{x,y}) is the size of the underlying grid.
func (gd Grid[T]) Cap() Point {
	if gd.ug == nil {
		return Point{}
	}
	w := gd.ug.Width
	h := 0
	if w != 0 && len(gd.ug.Cells) > 0 {
		h += len(gd.ug.Cells) / w
	}
	return Point{w, h}.Sub(gd.rg.Min)
}

// Slice returns a rectangular slice of the grid given by a range relative to
// the grid. If the range is out of bounds of the parent grid, it will be
// reduced to fit to the available space. The returned grid shares memory with
// the parent.
func (gd Grid[T]) Slice(rg Range) Grid[T] {
	if rg.Min.X < 0 {
		rg.Min.X = 0
	}
	if rg.Min.Y < 0 {
		rg.Min.Y = 0
	}
	max := gd.rg.Size()
	if rg.Max.X > max.X {
		rg.Max.X = max.X
	}
	if rg.Max.Y > max.Y {
		rg.Max.Y = max.Y
	}
	min := gd.rg.Min
	rg.Min = rg.Min.Add(min)
	rg.Max = rg.Max.Add(min)
	return Grid[T]{ug: gd.ug, rg: rg}
}

// Size returns the grid (width, height) in cells, and is a shorthand for
// gd.Range().Size().
func (gd Grid[T]) Size() Point {
	return gd.rg.Size()
}

// Resize is similar to Slice, but it only specifies new dimensions, and if the
// range goes beyond the underlying original grid range, it will grow the
// underlying grid. It preserves the content, and any new cells get the zero
// value.
func (gd Grid[T]) Resize(w, h int) Grid[T] {
	max := gd.Size()
	ow, oh := max.X, max.Y
	if ow == w && oh == h {
		return gd
	}
	if w <= 0 || h <= 0 {
		gd.rg.Max = gd.rg.Min
		return gd
	}
	if gd.ug == nil {
		gd.ug = &grid[T]{}
	}
	gd.rg.Max = gd.rg.Min.Shift(w, h)
	uh := 0
	if len(gd.ug.Cells) > 0 && gd.ug.Width > 0 {
		uh = len(gd.ug.Cells) / gd.ug.Width
	}
	nw := gd.ug.Width
	if w+gd.rg.Min.X > gd.ug.Width {
		nw = w + gd.rg.Min.X
	}
	nh := uh
	if h+gd.rg.Min.Y > uh {
		nh = h + gd.rg.Min.Y
	}
	if nw > gd.ug.Width || nh > uh {
		if nw == gd.ug.Width {
			var zero T
			for i := 0; i < nh-uh; i++ {
				for j := 0; j < gd.ug.Width; j++ {
					gd.ug.Cells = append(gd.ug.Cells, zero)
				}
			}
		} else {
			ngd := NewGrid[T](nw, nh)
			ngd.Copy(Grid[T]{ug: gd.ug, rg: NewRange(0, 0, gd.ug.Width, uh)})
			*gd.ug = *ngd.ug
		}
	}
	return gd
}

// Contains returns true if the given relative position is within the grid.
func (gd Grid[T]) Contains(p Point) bool {
	return p.Add(gd.rg.Min).In(gd.rg)
}

// Set draws a cell at a given position in the grid. If the position is out of
// range, the function does nothing.
func (gd Grid[T]) Set(p Point, c T) {
	q := p.Add(gd.rg.Min)
	if !q.In(gd.rg) {
		return
	}
	i := q.Y*gd.ug.Width + q.X
	gd.ug.Cells[i] = c
}

// At returns the cell at a given position. If the position is out of range, it
// returns the zero value.
func (gd Grid[T]) At(p Point) T {
	q := p.Add(gd.rg.Min)
	if !q.In(gd.rg) {
		var zero T
		return zero
	}
	i := q.Y*gd.ug.Width + q.X
	return gd.ug.Cells[i]
}

// AtU returns the cell at a given position without checking the grid slice
// bounds.  If the position is out of bounds, it returns a value corresponding
// to the position in the underlying grid, or the zero value if also out
// of the underlying grid's range.
//
// It may be somewhat faster than At in tight loops, but most of the time you
// can get the same performance using GridIterator or iteration functions,
// which are less error-prone.
func (gd Grid[T]) AtU(p Point) T {
	p = p.Add(gd.rg.Min)
	i := p.Y*gd.ug.Width + p.X
	if i < 0 || i >= len(gd.ug.Cells) {
		var zero T
		return zero
	}
	return gd.ug.Cells[i]
}

// Fill sets the given cell as content for all the grid positions.
func (gd Grid[T]) Fill(c T) {
	if gd.ug == nil {
		return
	}
	w := gd.rg.Max.X - gd.rg.Min.X
	switch {
	case w >= 8:
		// heuristic suited for T of small size
		gd.fillcp(c)
	case w == 1:
		gd.fillv(c)
	default:
		gd.fill(c)
	}
}

func (gd Grid[T]) fillcp(c T) {
	w := gd.ug.Width
	ymin := gd.rg.Min.Y * w
	gdw := gd.rg.Max.X - gd.rg.Min.X
	cells := gd.ug.Cells
	for xi := ymin + gd.rg.Min.X; xi < ymin+gd.rg.Max.X; xi++ {
		cells[xi] = c
	}
	idxmax := (gd.rg.Max.Y-1)*w + gd.rg.Max.X
	for idx := ymin + w + gd.rg.Min.X; idx < idxmax; idx += w {
		copy(cells[idx:idx+gdw], cells[ymin+gd.rg.Min.X:ymin+gd.rg.Max.X])
	}
}

func (gd Grid[T]) fill(c T) {
	w := gd.ug.Width
	cells := gd.ug.Cells
	yimax := gd.rg.Max.Y * w
	for yi := gd.rg.Min.Y * w; yi < yimax; yi += w {
		ximax := yi + gd.rg.Max.X
		for xi := yi + gd.rg.Min.X; xi < ximax; xi++ {
			cells[xi] = c
		}
	}
}

func (gd Grid[T]) fillv(c T) {
	w := gd.ug.Width
	cells := gd.ug.Cells
	ximax := gd.rg.Max.Y*w + gd.rg.Min.X
	for xi := gd.rg.Min.Y*w + gd.rg.Min.X; xi < ximax; xi += w {
		cells[xi] = c
	}
}

// FillFunc updates the content for all the grid positions, in row-major order,
// using the given function return value.
func (gd Grid[T]) FillFunc(fn func(Point) T) {
	if gd.ug == nil {
		return
	}
	w := gd.ug.Width
	yimax := gd.rg.Max.Y * w
	cells := gd.ug.Cells
	for y, yi := 0, gd.rg.Min.Y*w; yi < yimax; y, yi = y+1, yi+w {
		ximax := yi + gd.rg.Max.X
		for x, xi := 0, yi+gd.rg.Min.X; xi < ximax; x, xi = x+1, xi+1 {
			p := Point{X: x, Y: y}
			cells[xi] = fn(p)
		}
	}
}

// Iter iterates a function on all the grid positions and cells, in row-major
// order.
func (gd Grid[T]) Iter(fn func(Point, T)) {
	if gd.ug == nil {
		return
	}
	w := gd.ug.Width
	yimax := gd.rg.Max.Y * w
	cells := gd.ug.Cells
	for y, yi := 0, gd.rg.Min.Y*w; yi < yimax; y, yi = y+1, yi+w {
		ximax := yi + gd.rg.Max.X
		for x, xi := 0, yi+gd.rg.Min.X; xi < ximax; x, xi = x+1, xi+1 {
			c := cells[xi]
			p := Point{X: x, Y: y}
			fn(p, c)
		}
	}
}

// Map updates the grid content using the given mapping function. The iteration
// is done in row-major order.
func (gd Grid[T]) Map(fn func(Point, T) T) {
	if gd.ug == nil {
		return
	}
	w := gd.ug.Width
	cells := gd.ug.Cells
	yimax := gd.rg.Max.Y * w
	for y, yi := 0, gd.rg.Min.Y*w; yi < yimax; y, yi = y+1, yi+w {
		ximax := yi + gd.rg.Max.X
		for x, xi := 0, yi+gd.rg.Min.X; xi < ximax; x, xi = x+1, xi+1 {
			c := cells[xi]
			p := Point{X: x, Y: y}
			cells[xi] = fn(p, c)
		}
	}
}

// Copy copies elements from a source grid src into the destination grid gd,
// and returns the copied grid-slice size, which is the minimum of both grids
// for each dimension. The result is independent of whether the two grids
// referenced memory overlaps or not.
func (gd Grid[T]) Copy(src Grid[T]) Point {
	if gd.ug == nil || src.ug == nil {
		return Point{}
	}
	if gd.ug != src.ug {
		if src.rg.Max.X-src.rg.Min.X <= 4 {
			// heuristic suited for T of small size
			return gd.cpv(src)
		}
		return gd.cp(src)
	}
	if gd.rg == src.rg {
		return gd.rg.Size()
	}
	if !gd.rg.Overlaps(src.rg) || gd.rg.Min.Y <= src.rg.Min.Y {
		return gd.cp(src)
	}
	return gd.cprev(src)
}

func (gd Grid[T]) cp(src Grid[T]) Point {
	w := gd.ug.Width
	wsrc := src.ug.Width
	max := gd.Range().Intersect(src.Range()).Size()
	idxmin := gd.rg.Min.Y*w + gd.rg.Min.X
	idxsrcmin := src.rg.Min.Y*w + src.rg.Min.X
	idxmax := (gd.rg.Min.Y + max.Y) * w
	for idx, idxsrc := idxmin, idxsrcmin; idx < idxmax; idx, idxsrc = idx+w, idxsrc+wsrc {
		copy(gd.ug.Cells[idx:idx+max.X], src.ug.Cells[idxsrc:idxsrc+max.X])
	}
	return max
}

func (gd Grid[T]) cpv(src Grid[T]) Point {
	w := gd.ug.Width
	wsrc := src.ug.Width
	max := gd.Range().Intersect(src.Range()).Size()
	yimax := (gd.rg.Min.Y + max.Y) * w
	cells := gd.ug.Cells
	srccells := src.ug.Cells
	for yi, yisrc := gd.rg.Min.Y*w, src.rg.Min.Y*wsrc; yi < yimax; yi, yisrc = yi+w, yisrc+wsrc {
		ximax := yi + max.X
		for xi, xisrc := yi+gd.rg.Min.X, yisrc+src.rg.Min.X; xi < ximax; xi, xisrc = xi+1, xisrc+1 {
			cells[xi] = srccells[xisrc]
		}
	}
	return max
}

func (gd Grid[T]) cprev(src Grid[T]) Point {
	w := gd.ug.Width
	wsrc := src.ug.Width
	max := gd.Range().Intersect(src.Range()).Size()
	idxmax := (gd.rg.Min.Y+max.Y-1)*w + gd.rg.Min.X
	idxsrcmax := (src.rg.Min.Y+max.Y-1)*w + src.rg.Min.X
	idxmin := gd.rg.Min.Y * w
	for idx, idxsrc := idxmax, idxsrcmax; idx >= idxmin; idx, idxsrc = idx-w, idxsrc-wsrc {
		copy(gd.ug.Cells[idx:idx+max.X], src.ug.Cells[idxsrc:idxsrc+max.X])
	}
	return max
}

// GridIterator represents a stateful iterator for a grid. They are created
// with the Iterator method.
type GridIterator[T any] struct {
	cells  []T   // grid cells
	p      Point // iterator's current position
	max    Point // last position
	i      int   // current position's index
	w      int   // underlying grid's width
	nlstep int   // newline step
	rg     Range // grid range
}

// Iterator returns an iterator that can be used to iterate on the grid. It may
// be convenient when more flexibility than the provided by the other iteration
// functions is needed. It is used as follows:
//
// 	it := gd.Iterator()
// 	for it.Next() {
// 		// call it.P() or it.Cell() or it.SetCell() as appropriate
// 	}
func (gd Grid[T]) Iterator() GridIterator[T] {
	if gd.ug == nil {
		return GridIterator[T]{}
	}
	w := gd.ug.Width
	it := GridIterator[T]{
		w:      w,
		cells:  gd.ug.Cells,
		max:    gd.Size().Shift(-1, -1),
		rg:     gd.rg,
		nlstep: gd.rg.Min.X + (w - gd.rg.Max.X + 1),
	}
	it.Reset()
	return it
}

// Reset resets the iterator's state so that it can be used again.
func (it *GridIterator[T]) Reset() {
	it.p = Point{-1, 0}
	it.i = it.rg.Min.Y*it.w + it.rg.Min.X - 1
}

// Next advances the iterator the next position in the grid, using row-major
// ordering.
func (it *GridIterator[T]) Next() bool {
	if it.p.X < it.max.X {
		it.p.X++
		it.i++
		return true
	}
	if it.p.Y < it.max.Y {
		it.p.Y++
		it.p.X = 0
		it.i += it.nlstep
		return true
	}
	return false
}

// P returns the iterator's current position.
func (it *GridIterator[T]) P() Point {
	return it.p
}

// SetP sets the iterator's current position.
func (it *GridIterator[T]) SetP(p Point) {
	q := p.Add(it.rg.Min)
	if !q.In(it.rg) {
		return
	}
	it.p = p
	it.i = q.Y*it.w + q.X
}

// V returns the cell value at the iterator's current position.
func (it *GridIterator[T]) V() T {
	return it.cells[it.i]
}

// SetV updates cell value at the iterator's current position. It's faster than
// calling Set on the grid.
func (it *GridIterator[T]) SetV(c T) {
	it.cells[it.i] = c
}
