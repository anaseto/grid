package grid

import (
	"math/rand"
	"testing"
)

func randInt(n int) int {
	if n <= 0 {
		return 0
	}
	x := rand.Intn(n)
	return x
}

func TestPoint(t *testing.T) {
	p := Point{2, 3}
	if p.Mul(3).X != 6 {
		t.Errorf("bad mul: %v", p.Mul(3))
	}
	if p.Div(2).X != 1 {
		t.Errorf("bad mul: %v", p.Div(2))
	}
}

func TestPointString(t *testing.T) {
	p := Point{2, 3}
	if p.String() != "(2,3)" {
		t.Errorf("bad string representation: %s", p.String())
	}
}

func TestRange(t *testing.T) {
	rg := NewRange(2, 3, 20, 30)
	max := rg.Size()
	w, h := max.X, max.Y
	count := 0
	rg.Iter(func(p Point) {
		if !p.In(rg) {
			t.Errorf("bad position: %+v", p)
		}
		count++
	})
	if count != w*h {
		t.Errorf("bad count: %d", count)
	}
	rg = rg.Sub(rg.Min)
	max = rg.Size()
	nw, nh := max.X, max.Y
	if nw != w || nh != h {
		t.Errorf("bad size for range %+v", rg)
	}
	if rg.Min.X != 0 || rg.Min.Y != 0 {
		t.Errorf("bad min for range %+v", rg)
	}
	nrg := rg.Shift(1, 2, 3, 4)
	if rg.Min.Shift(1, 2) != nrg.Min {
		t.Errorf("bad min shift for range %+v", nrg)
	}
	if rg.Max.Shift(3, 4) != nrg.Max {
		t.Errorf("bad max shift for range %+v", nrg)
	}
}

func TestRangeString(t *testing.T) {
	p := Point{2, 3}
	q := Point{3, 4}
	rg := Range{p, q}
	if rg.String() != "(2,3)-(3,4)" {
		t.Errorf("bad string representation: %s", rg.String())
	}
}

func TestRangeIntersect(t *testing.T) {
	rg := NewRange(0, 1, 2, 3)
	rg2 := NewRange(4, 5, 6, 7)
	rg3 := rg.Intersect(rg2)
	var zero Range
	if rg3 != zero {
		t.Errorf("non zero Range")
	}
}

func TestRangeUnion(t *testing.T) {
	rg := NewRange(1, 2, 3, 4)
	org := NewRange(11, 12, 13, 14)
	union := NewRange(1, 2, 13, 14)
	if rg.Union(org) != union {
		t.Errorf("bad Union")
	}
	if !rg.In(union) {
		t.Errorf("bad In")
	}
}

func TestRangeUnion2(t *testing.T) {
	rg := NewRange(0, 1, 2, 3)
	rg2 := NewRange(4, 5, 6, 7)
	rg3 := rg2.Union(rg)
	if rg3.Min.X != 0 || rg3.Min.Y != 1 || rg3.Max.X != 6 || rg3.Max.Y != 7 {
		t.Errorf("bad range: %v", rg3)
	}
}

func TestRangeShift(t *testing.T) {
	rg := NewRange(1, 2, 3, 4)
	nrg := NewRange(2, 3, 4, 5)
	if rg.Shift(1, 1, 1, 1) != nrg {
		t.Errorf("bad shift: %v", rg.Shift(1, 1, 1, 1))
	}
	empty := Range{}
	if rg.Shift(0, 0, -5, 0) != empty {
		t.Errorf("bad shift: %v", rg.Shift(0, 0, -5, 0))
	}
	if rg.Shift(0, 0, 0, -5) != empty {
		t.Errorf("bad shift: %v", rg.Shift(0, 0, 0, -5))
	}
	if rg.Add(Point{1, 1}) != nrg {
		t.Errorf("bad add: %v", rg.Add(Point{1, 1}))
	}
}

func TestRangeColumnsLines(t *testing.T) {
	rg := NewRange(1, 1, 30, 30)
	if rg.Columns(4, 10).Size().X != 6 {
		t.Errorf("bad number of columns for range %v", rg.Columns(4, 10).Size().Y)
	}
	if rg.Columns(4, 10).Min.X != 5 {
		t.Errorf("bad min.X for range %v", rg.Columns(4, 10).Min.X)
	}
	if rg.Lines(4, 10).Size().Y != 6 {
		t.Errorf("bad number of columns for range %v", rg.Columns(4, 10).Size().Y)
	}
	if rg.Lines(4, 10).Min.Y != 5 {
		t.Errorf("bad min.X for range %v", rg.Columns(4, 10).Min.Y)
	}
	if !rg.Column(200).Empty() {
		t.Errorf("not empty column")
	}
	if !rg.Line(200).Empty() {
		t.Errorf("not empty line")
	}
}

func TestRangeEq(t *testing.T) {
	rg := NewRange(1, 2, 3, 4)
	if !rg.Eq(rg) {
		t.Errorf("bad reflexive Eq for %v", rg)
	}
	if rg.Eq(rg.Shift(1, 0, 0, 0)) {
		t.Errorf("bad shift Eq for %v", rg)
	}
	erg := Range{Point{2, 3}, Point{-1, -4}}
	empty := Range{}
	if !erg.Eq(empty) {
		t.Errorf("bad empty range equality")
	}
}

func TestBounds(t *testing.T) {
	gd := NewGrid[int](10, 10)
	slice := gd.Slice(NewRange(2, 2, 4, 4))
	if slice.Bounds() != NewRange(2, 2, 4, 4) {
		t.Errorf("bad Bounds %v", slice.Bounds())
	}
}

func TestContents(t *testing.T) {
	gd := NewGrid[int](10, 10)
	values := gd.Contents()
	if len(values) != 100 {
		t.Errorf("bad length: %d", len(values))
	}
	var gd2 Grid[float64]
	if gd2.Contents() != nil {
		t.Errorf("non nil: %v", gd2.Contents())
	}
}

func TestNewGrid(t *testing.T) {
	gd := NewGrid[int](80, 24)
	max := gd.Size()
	if max.X != 80 && max.Y != 24 {
		t.Errorf("bad default size: (%d,%d)", max.X, max.Y)
	}
	gd = NewGrid[int](50, 50)
	max = gd.Size()
	w, h := max.X, max.Y
	if w != 50 && h != 50 {
		t.Errorf("grid size does not match configuration: (%d,%d)", w, h)
	}
	max = gd.Bounds().Size()
	rw, rh := max.X, max.Y
	if w != rw || rh != h {
		t.Errorf("incompatible sizes: grid (%d,%d) range (%d,%d)", w, h, rw, rh)
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := Point{x, y}
			c := gd.At(p)
			if c != 0 {
				t.Errorf("cell: bad content %c at %+v", c, p)
			}
		}
	}
}

func testPanic(t *testing.T, f func(), msg string) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic when " + msg)
		}
	}()
	f()
}

func TestNewGridPanic(t *testing.T) {
	testPanic(t, func() { NewGrid[int](-1, 2) }, "w < 0")
	testPanic(t, func() { NewGrid[int](2, -1) }, "h < 0")
}

func TestNewGridFromSlice(t *testing.T) {
	s := []int{}
	w := 5
	h := 10
	for i := 0; i < w*h; i++ {
		s = append(s, 1)
	}
	gd := NewGridFromSlice[int](s, w)
	max := gd.Size()
	gw, gh := max.X, max.Y
	if gw != w && gh != h {
		t.Errorf("incorrect grid size: (%d,%d)", gw, gh)
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := Point{x, y}
			c := gd.At(p)
			if c != 1 {
				t.Errorf("cell: bad content %c at %+v", c, p)
			}
		}
	}
}

func TestNewGridFromSlicePanic(t *testing.T) {
	s := []int{1, 2}
	testPanic(t, func() {
		NewGridFromSlice[int](s, -2)
	}, "w < 0")
	testPanic(t, func() {
		NewGridFromSlice[int](s, 0)
	}, "w == 0 && len(s) > 0")
	testPanic(t, func() {
		NewGridFromSlice[int](s, 3)
	}, "len(s) %%w != 0")
}

func TestSetV(t *testing.T) {
	gd := NewGrid[int](80, 24)
	max := gd.Size()
	w, h := max.X, max.Y
	gd.Fill(1)
	for i := 0; i < w*h; i++ {
		p := Point{X: randInt(2*w) - w/2, Y: randInt(2*h) - h/2}
		if gd.Contains(p) {
			c := gd.At(p)
			if c != 1 && c != 2 {
				t.Errorf("Bad fill or setcell %+v at p %+v", c, p)
			}
		}
		gd.Set(p, 2)
		c := gd.At(p)
		if gd.Contains(p) {
			if c != 2 {
				t.Errorf("Bad content %+v at %+v", c, p)
			}
		} else if c != 0 {
			t.Errorf("Bad out of range content: %+v at %+v", c, p)
		}
	}
}

func TestGridCap(t *testing.T) {
	gd := NewGrid[int](10, 10)
	gd.Fill(1)
	slice := gd.Slice(NewRange(2, 3, 7, 7))
	c := slice.Cap()
	if c.X != 8 || c.Y != 7 {
		t.Errorf("Bad capacity: %+v", c)
	}
	var gd2 Grid[int]
	c = gd2.Cap()
	if c.X != 0 || c.Y != 0 {
		t.Errorf("non zero capacity: %+v", c)
	}
}

func TestGridAt(t *testing.T) {
	gd := NewGrid[int](10, 10)
	gd.FillFunc(func(p Point) int {
		return 100*p.X + p.Y
	})
	gd.Iter(func(p Point, c int) {
		if c != 100*p.X+p.Y {
			t.Errorf("bad value %d at %v", c, p)
		}
		if c != gd.AtU(p) {
			t.Errorf("got %d instead of %d at %v", gd.AtU(p), c, p)
		}
	})
	v := gd.AtU(Point{20, 20})
	if v != 0 {
		t.Errorf("non zero: %v", v)
	}
}

func TestGridNil(t *testing.T) {
	var gd Grid[int]
	gd.Fill(3)
	gd.FillFunc(func(p Point) int { return 1 })
	gd.Iter(func(p Point, v int) {})
	gd.Map(func(p Point, v int) int { return 1 })
	gd.Copy(gd)
	gd.Iterator()
	// does nothing, no panic
}

func TestGridSlice(t *testing.T) {
	gd := NewGrid[int](80, 24)
	max := gd.Size()
	w, h := max.X, max.Y
	gd.Fill(1)
	slice := gd.Slice(NewRange(5, 5, 10, 10))
	slice.Fill(3)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := Point{x, y}
			c := gd.At(p)
			if p.In(slice.Bounds()) {
				if c != 3 {
					t.Errorf("bad slice cell: %c at %+v", c, p)
				}
			} else if c != 1 {
				t.Errorf("bad grid non-slice cell: %c at %+v", c, p)
			}
		}
	}
}

func TestGridSlice2(t *testing.T) {
	gd := NewGrid[int](80, 24)
	gd.Fill(1)
	slice := gd.Slice(NewRange(0, 0, 0, 0))
	if !slice.Bounds().Empty() {
		t.Errorf("non empty range %v", slice.Bounds())
	}
	slice = gd.Slice(NewRange(0, 0, -5, -5))
	if !slice.Bounds().Empty() {
		t.Errorf("non empty negative range %v", slice.Bounds())
	}
	slice = gd.Slice(NewRange(5, 5, 0, 0))
	rg := slice.Bounds()
	if rg.Max.X != 5 || rg.Max.Y != 5 || rg.Min.X != 0 || rg.Min.Y != 0 {
		t.Errorf("bad inversed range %+v", slice.Bounds())
	}
}

func TestGridSlice3(t *testing.T) {
	gd := NewGrid[int](80, 24)
	max := gd.Size()
	w, h := max.X, max.Y
	gd.Fill(1)
	slice := gd.Slice(gd.Range().Line(1))
	slice.Fill(11)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := Point{x, y}
			c := gd.At(p)
			if y == 1 {
				if c != 11 {
					t.Errorf("bad line slice: %c", c)
				}
			} else if c != 1 {
				t.Errorf("bad outside line slice: %c", c)
			}
		}
	}
	slice = gd.Slice(gd.Range().Column(2))
	gd.Fill(1)
	slice.Fill(12)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := Point{x, y}
			c := gd.At(p)
			if x == 2 {
				if c != 12 {
					t.Errorf("bad column slice: %c", c)
				}
			} else if c != 1 {
				t.Errorf("bad outside column slice: %c", c)
			}
		}
	}
}

func TestGridSlice4(t *testing.T) {
	gd := NewGrid[int](10, 10)
	if gd.Slice(NewRange(-5, -5, 20, 20)).Range() != gd.Range() {
		t.Errorf("bad oversized slice")
	}
}

func TestGridSlice5(t *testing.T) {
	gd := NewGrid[int](80, 24)
	max := gd.Size()
	w, h := max.X, max.Y
	gd.Fill(1)
	slice := gd.Slice(NewRange(5, 5, 15, 15))
	slice.Fill(3)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := Point{x, y}
			c := gd.At(p)
			if p.In(slice.Bounds()) {
				if c != 3 {
					t.Errorf("bad slice cell: %c at %+v", c, p)
				}
			} else if c != 1 {
				t.Errorf("bad grid non-slice cell: %c at %+v", c, p)
			}
		}
	}
}

func TestIterMap(t *testing.T) {
	gd := NewGrid[int](10, 10)
	gd.Map(func(p Point, c int) int { return 4 })
	gd.Iter(func(p Point, c int) {
		if c != 4 {
			t.Errorf("bad cell %c at %v", c, p)
		}
	})
}

func TestCopy(t *testing.T) {
	gd := NewGrid[int](80, 30)
	max := gd.Size()
	w, h := max.X, max.Y
	gd.Fill(1)
	gd2 := NewGrid[int](10, 10)
	gd2.Fill(4)
	gd.Copy(gd2)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := Point{x, y}
			c := gd.At(p)
			if p.In(gd2.Range()) {
				if c != 4 {
					t.Errorf("bad copy at cell: %c at %+v", c, p)
				}
			} else if c != 1 {
				t.Errorf("bad grid non-slice cell: %c at %+v", c, p)
			}
		}
	}
}

func TestCopy2(t *testing.T) {
	gd := NewGrid[int](80, 10)
	max := gd.Size()
	w, h := max.X, max.Y
	gd.Fill(1)
	rg := gd.Bounds()
	slice := gd.Slice(rg.Lines(1, 3))
	slice2 := gd.Slice(rg.Line(2))
	slice3 := gd.Slice(rg.Lines(2, 4))
	slice.Fill(11)  // line 1
	slice3.Fill(13) // line 3
	slice2.Fill(12) // line 2
	slice.Copy(slice3)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := Point{x, y}
			c := gd.At(p)
			switch {
			case p.In(rg.Line(1)):
				if c != 12 {
					t.Errorf("bad line 1: %c at %+v", c, p)
				}
			case p.In(rg.Line(2)):
				if c != 13 {
					t.Errorf("bad line 2: %c at %+v", c, p)
				}
			case p.In(rg.Line(3)):
				if c != 13 {
					t.Errorf("bad line 3: %c at %+v", c, p)
				}
			default:
				if c != 1 {
					t.Errorf("bad line 3: %c at %+v", c, p)
				}
			}
		}
	}
}

func TestCopy3(t *testing.T) {
	gd := NewGrid[int](80, 10)
	max := gd.Size()
	w, h := max.X, max.Y
	gd.Fill(1)
	rg := gd.Bounds()
	slice := gd.Slice(rg.Lines(1, 3))
	slice2 := gd.Slice(rg.Line(2))
	slice3 := gd.Slice(rg.Lines(2, 4))
	slice.Fill(11)  // line 1
	slice3.Fill(13) // line 3
	slice2.Fill(12) // line 2
	slice.Copy(slice2)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := Point{x, y}
			c := gd.At(p)
			switch {
			case p.In(rg.Line(1)):
				if c != 12 {
					t.Errorf("bad line 1: %c at %+v", c, p)
				}
			case p.In(rg.Line(2)):
				if c != 12 {
					t.Errorf("bad line 2: %c at %+v", c, p)
				}
			case p.In(rg.Line(3)):
				if c != 13 {
					t.Errorf("bad line 3: %c at %+v", c, p)
				}
			default:
				if c != 1 {
					t.Errorf("bad line 3: %c at %+v", c, p)
				}
			}
		}
	}
}

func TestCopy4(t *testing.T) {
	gd := NewGrid[int](80, 10)
	max := gd.Size()
	w, h := max.X, max.Y
	gd.Fill(1)
	rg := gd.Bounds()
	slice := gd.Slice(rg.Lines(1, 3))
	slice2 := gd.Slice(rg.Line(2))
	slice3 := gd.Slice(rg.Lines(2, 4))
	slice.Fill(11)  // line 1
	slice3.Fill(13) // line 3
	slice2.Fill(12) // line 2
	slice3.Copy(slice2)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := Point{x, y}
			c := gd.At(p)
			switch {
			case p.In(rg.Line(1)):
				if c != 11 {
					t.Errorf("bad line 1: %c at %+v", c, p)
				}
			case p.In(rg.Line(2)):
				if c != 12 {
					t.Errorf("bad line 2: %c at %+v", c, p)
				}
			case p.In(rg.Line(3)):
				if c != 13 {
					t.Errorf("bad line 3: %c at %+v", c, p)
				}
			default:
				if c != 1 {
					t.Errorf("bad line 3: %c at %+v", c, p)
				}
			}
		}
	}
}

func TestCopy5(t *testing.T) {
	gd := NewGrid[int](80, 10)
	max := gd.Size()
	w, h := max.X, max.Y
	gd.Fill(1)
	rg := gd.Bounds()
	slice := gd.Slice(rg.Lines(1, 3))
	slice2 := gd.Slice(rg.Line(2))
	slice3 := gd.Slice(rg.Lines(2, 4))
	slice.Fill(11)  // line 1
	slice3.Fill(13) // line 3
	slice2.Fill(12) // line 2
	slice2.Copy(slice)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := Point{x, y}
			c := gd.At(p)
			switch {
			case p.In(rg.Line(1)):
				if c != 11 {
					t.Errorf("bad line 1: %c at %+v", c, p)
				}
			case p.In(rg.Line(2)):
				if c != 11 {
					t.Errorf("bad line 2: %c at %+v", c, p)
				}
			case p.In(rg.Line(3)):
				if c != 13 {
					t.Errorf("bad line 3: %c at %+v", c, p)
				}
			default:
				if c != 1 {
					t.Errorf("bad line 3: %c at %+v", c, p)
				}
			}
		}
	}
}

func TestCopy6(t *testing.T) {
	gd := NewGrid[int](80, 10)
	max := gd.Size()
	w, h := max.X, max.Y
	gd.Fill(1)
	rg := gd.Bounds()
	slice := gd.Slice(rg.Lines(1, 3))
	slice2 := gd.Slice(rg.Line(2))
	slice3 := gd.Slice(rg.Lines(2, 4))
	slice.Fill(11)  // line 1
	slice3.Fill(13) // line 3
	slice2.Fill(12) // line 2
	slice2.Copy(slice3)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := Point{x, y}
			c := gd.At(p)
			switch {
			case p.In(rg.Line(1)):
				if c != 11 {
					t.Errorf("bad line 1: %c at %+v", c, p)
				}
			case p.In(rg.Line(2)):
				if c != 12 {
					t.Errorf("bad line 2: %c at %+v", c, p)
				}
			case p.In(rg.Line(3)):
				if c != 13 {
					t.Errorf("bad line 3: %c at %+v", c, p)
				}
			default:
				if c != 1 {
					t.Errorf("bad line 3: %c at %+v", c, p)
				}
			}
		}
	}
}

func TestCopy7(t *testing.T) {
	gd := NewGrid[int](80, 10)
	max := gd.Size()
	w, h := max.X, max.Y
	gd.Fill(1)
	rg := gd.Bounds()
	slice := gd.Slice(rg.Lines(1, 3))
	slice2 := gd.Slice(rg.Line(2))
	slice3 := gd.Slice(rg.Lines(2, 4))
	slice.Fill(11)  // line 1
	slice3.Fill(13) // line 3
	slice2.Fill(12) // line 2
	slice3.Copy(slice)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := Point{x, y}
			c := gd.At(p)
			switch {
			case p.In(rg.Line(1)):
				if c != 11 {
					t.Errorf("bad line 1: %c at %+v", c, p)
				}
			case p.In(rg.Line(2)):
				if c != 11 {
					t.Errorf("bad line 2: %c at %+v", c, p)
				}
			case p.In(rg.Line(3)):
				if c != 12 {
					t.Errorf("bad line 3: %c at %+v", c, p)
				}
			default:
				if c != 1 {
					t.Errorf("bad line 3: %c at %+v", c, p)
				}
			}
		}
	}
}

func TestCopy8(t *testing.T) {
	gd := NewGrid[int](3, 10)
	gd.Fill(1)
	gd2 := NewGrid[int](3, 10)
	gd2.Fill(4)
	gd.Copy(gd2)
	gd.Iter(func(p Point, c int) {
		if c != 4 {
			t.Errorf("bad num %c at %v", c, p)
		}
	})
}

func TestCopySelf(t *testing.T) {
	gd := NewGrid[int](80, 10)
	if gd.Copy(gd) != gd.Range().Size() {
		t.Errorf("bad same range copy")
	}
}

func TestCopyShiftX(t *testing.T) {
	gd := NewGrid[int](80, 10)
	gd.Fill(5)
	ngd := NewGrid[int](80, 10)
	ngd.Fill(6)
	slice := gd.Slice(NewRange(20, 0, 80, 10))
	slice.Copy(ngd)
	gd.Iter(func(p Point, c int) {
		if p.X >= 20 && c != 6 {
			t.Errorf("not 6 in slice")
		}
		if p.X < 20 && c != 5 {
			t.Errorf("not 5 outside slice")
		}
	})
}

func TestResize(t *testing.T) {
	gd := NewGrid[int](20, 10)
	gd.Fill(1)
	rg := gd.Range()
	gd = gd.Resize(30, 20)
	if gd.Size().X != 30 || gd.Size().Y != 20 {
		t.Errorf("bad size: %v", gd.Size())
	}
	gd.Iter(func(p Point, c int) {
		if p.In(rg) {
			if c != 1 {
				t.Error("bad preservation of content")
			}
		} else if c != 0 {
			t.Error("bad new content")
		}
	})
}

func TestResize2(t *testing.T) {
	gd := NewGrid[int](20, 10)
	gd2 := gd.Resize(20, 10)
	if gd != gd2 {
		t.Error("same dimensions but different")
	}
	gd3 := gd.Resize(-20, 10)
	if !gd3.Range().Empty() {
		t.Errorf("non empty range: %v", gd3.Range())
	}
	var gd4 Grid[int]
	gd5 := gd4.Resize(20, 10)
	rg := gd5.Range()
	if rg.Max.X != 20 || rg.Max.Y != 10 {
		t.Errorf("bad range: %v", rg)
	}
	gd5.Fill(2)
	gd6 := gd5.Resize(20, 30)
	rg = gd6.Range()
	if rg.Max.X != 20 || rg.Max.Y != 30 {
		t.Errorf("bad range: %v", rg)
	}
	if gd6.AtU(Point{10, 5}) != 2 {
		t.Errorf("expected 2")
	}
	if gd6.AtU(Point{10, 25}) != 0 {
		t.Errorf("non zero")
	}
}

func TestIterator(t *testing.T) {
	gd := NewGrid[int](10, 10)
	slice := gd.Slice(NewRange(2, 2, 5, 5))
	it := slice.Iterator()
	for it.Next() {
		if it.V() != 0 {
			t.Errorf("not zero: %c", it.V())
		}
		it.SetV(2)
		if it.V() != 2 {
			t.Errorf("not x 2: %c", it.V())
		}
		if slice.At(it.P()) != 2 {
			t.Errorf("not x 2 at %v: %c", it.P(), slice.At(it.P()))
		}
	}
	gd.Iter(func(p Point, c int) {
		if p.In(slice.Bounds()) {
			if c != 2 {
				t.Errorf("bad num at %v: %c", p, c)
			}
		} else if c != 0 {
			t.Errorf("not zero at %v: %c", p, c)
		}

	})
	it.SetP(Point{1, 1})
	if it.P().X != 1 || it.P().Y != 1 {
		t.Errorf("bad SetP: %v", it.P())
	}
	it.SetP(Point{20, 20}) // does nothing, no panic
	it.SetV(7)
	if slice.At(Point{1, 1}) != 7 {
		t.Errorf("not 7: %c", slice.At(Point{1, 1}))
	}
	it.Reset()
	for it.Next() {
		it.SetV(8)
	}
}

func BenchmarkGridIter(b *testing.B) {
	gd := NewGrid[int](80, 24)
	gd.Fill(1)
	for i := 0; i < b.N; i++ {
		n := 0
		gd.Iter(func(p Point, c int) {
			if c == 1 {
				n++
			}
		})
	}
}

func BenchmarkGridIterator(b *testing.B) {
	gd := NewGrid[int](80, 24)
	gd.Fill(1)
	for i := 0; i < b.N; i++ {
		n := 0
		it := gd.Iterator()
		for it.Next() {
			if it.V() == 1 {
				n++
			}
		}
	}
}

func BenchmarkGridLoopAt(b *testing.B) {
	gd := NewGrid[int](80, 24)
	gd.Fill(1)
	for i := 0; i < b.N; i++ {
		n := 0
		max := gd.Size()
		for y := 0; y < max.Y; y++ {
			for x := 0; x < max.X; x++ {
				p := Point{x, y}
				n += gd.At(p)
			}
		}
	}
}

func BenchmarkGridLoopAtU(b *testing.B) {
	gd := NewGrid[int](80, 24)
	gd.Fill(1)
	for i := 0; i < b.N; i++ {
		n := 0
		max := gd.Size()
		for y := 0; y < max.Y; y++ {
			for x := 0; x < max.X; x++ {
				p := Point{x, y}
				n += gd.AtU(p)
			}
		}
	}
}

func BenchmarkGridIterSet(b *testing.B) {
	gd := NewGrid[int](80, 24)
	gd.Fill(1)
	for i := 0; i < b.N; i++ {
		gd.Iter(func(p Point, c int) {
			gd.Set(p, 2)
		})
	}
}

func BenchmarkGridRangeIterSet(b *testing.B) {
	gd := NewGrid[int](80, 24)
	gd.Fill(1)
	for i := 0; i < b.N; i++ {
		gd.Range().Iter(func(p Point) {
			gd.Set(p, 2)
		})
	}
}

func BenchmarkGridLoopSet(b *testing.B) {
	gd := NewGrid[int](80, 24)
	gd.Fill(1)
	for i := 0; i < b.N; i++ {
		max := gd.Size()
		for y := 0; y < max.Y; y++ {
			for x := 0; x < max.X; x++ {
				p := Point{x, y}
				gd.Set(p, 2)
			}
		}
	}
}

func BenchmarkGridIteratorSet(b *testing.B) {
	gd := NewGrid[int](80, 24)
	gd.Fill(1)
	for i := 0; i < b.N; i++ {
		it := gd.Iterator()
		for it.Next() {
			it.SetV(2)
		}
	}
}

func BenchmarkGridIteratorNew(b *testing.B) {
	gd := NewGrid[int](80, 24)
	for i := 0; i < b.N; i++ {
		it := gd.Iterator()
		it.Next()
	}
}

func BenchmarkGridFill(b *testing.B) {
	gd := NewGrid[int](80, 24)
	for i := 0; i < b.N; i++ {
		gd.Fill(1)
	}
}

func BenchmarkGridMap(b *testing.B) {
	gd := NewGrid[int](80, 24)
	for i := 0; i < b.N; i++ {
		gd.Map(func(p Point, c int) int { return 1 })
	}
}

func BenchmarkGridFillFunc(b *testing.B) {
	gd := NewGrid[int](80, 24)
	for i := 0; i < b.N; i++ {
		gd.FillFunc(func(p Point) int { return 1 })
	}
}

func BenchmarkGridCopy(b *testing.B) {
	gd := NewGrid[int](80, 24)
	gd2 := NewGrid[int](80, 24)
	for i := 0; i < b.N; i++ {
		gd.Copy(gd2)
	}
}

func BenchmarkGridCopyVertical2(b *testing.B) {
	gd := NewGrid[int](2, 40*24)
	gd2 := NewGrid[int](2, 40*24)
	for i := 0; i < b.N; i++ {
		gd.Copy(gd2)
	}
}

func BenchmarkGridCopyVertical4(b *testing.B) {
	gd := NewGrid[int](4, 20*24)
	gd2 := NewGrid[int](4, 20*24)
	for i := 0; i < b.N; i++ {
		gd.Copy(gd2)
	}
}

func BenchmarkGridCopyVertical8(b *testing.B) {
	gd := NewGrid[int](8, 10*24)
	gd2 := NewGrid[int](8, 10*24)
	for i := 0; i < b.N; i++ {
		gd.Copy(gd2)
	}
}

func BenchmarkGridFillVertical(b *testing.B) {
	gd := NewGrid[int](1, 24*80)
	for i := 0; i < b.N; i++ {
		gd.Fill(1)
	}
}

func BenchmarkGridFillVertical2(b *testing.B) {
	gd := NewGrid[int](2, 24*40)
	for i := 0; i < b.N; i++ {
		gd.Fill(1)
	}
}

func BenchmarkGridFillVertical4(b *testing.B) {
	gd := NewGrid[int](4, 24*20)
	for i := 0; i < b.N; i++ {
		gd.Fill(1)
	}
}

func BenchmarkGridFillVertical8(b *testing.B) {
	gd := NewGrid[int](8, 24*10)
	for i := 0; i < b.N; i++ {
		gd.Fill(1)
	}
}

func BenchmarkGridFillVertical16(b *testing.B) {
	gd := NewGrid[int](16, 12*10)
	for i := 0; i < b.N; i++ {
		gd.Fill(1)
	}
}
