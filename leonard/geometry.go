package leonard

type offset struct {
	X, Y int
}

func (o offset) apply(x, y int) (int, int) {
	return x + o.X, y + o.Y
}

func (o offset) add(p offset) offset {
	return offset{o.X + p.X, o.Y + p.Y}
}

var (
	north = offset{0, -1}
	east  = offset{1, 0}
	south = offset{0, 1}
	west  = offset{-1, 0}
)

var neighboursOffsets = []offset{
	north,
	north.add(east),
	east,
	south.add(east),
	south,
	south.add(west),
	west,
	north.add(west),
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

func (m *boolMatrix) count(b bool) int {
	n := 0
	for _, c := range m.cells {
		if c == b {
			n++
		}
	}
	return n
}
