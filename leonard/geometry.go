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
