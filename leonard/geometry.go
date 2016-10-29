package leonard

type offset struct {
	X, Y int
}

func (o offset) apply(x, y int) (int, int) {
	return x + o.X, y + o.Y
}

func (o offset) applyN(x, y, n int) (int, int) {
	return x + n*o.X, y + n*o.Y
}

func (o offset) add(p offset) offset {
	return offset{o.X + p.X, o.Y + p.Y}
}

func (o offset) reverse() offset {
	return offset{-o.X, -o.Y}
}

var (
	north = offset{0, -1}
	east  = offset{1, 0}
	south = offset{0, 1}
	west  = offset{-1, 0}

	northwest = north.add(west)
	northeast = north.add(east)
	southwest = south.add(west)
	southeast = south.add(east)
)

var clockwiseOffsets = []offset{
	north,
	northeast,
	east,
	southeast,
	south,
	southwest,
	west,
	northwest,
}

type boolMatrix struct {
	cells         []bool
	height, width int
}

func newBoolMatrix(height, width int) *boolMatrix {
	return &boolMatrix{
		cells:  make([]bool, height*width),
		height: height,
		width:  width,
	}
}

func (m *boolMatrix) set(x, y int, value bool) {
	idx := y*m.width + x
	if idx < 0 || idx >= len(m.cells) {
		return
	}

	m.cells[idx] = value
}

func (m *boolMatrix) get(x, y int) bool {
	idx := y*m.width + x
	if idx < 0 || idx >= len(m.cells) {
		return false
	}
	return m.cells[idx]
}

func (m *boolMatrix) countNeighbours(x, y int) (n int) {
	for iy := y - 1; iy <= y+1; iy++ {
		for ix := x - 1; ix <= x+1; ix++ {
			if iy == y && ix == x {
				continue
			}

			if m.get(ix, iy) {
				n++
			}

		}
	}
	return
}
