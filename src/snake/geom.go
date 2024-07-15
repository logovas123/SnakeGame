package snake

type Point struct {
	X, Y float64
}

const (
	Top    = iota
	Right  = iota
	Bottom = iota
	Left   = iota
)

type Dir int

func (d Dir) Exec(p Point) Point {
	switch d {
	case Top:
		return Point{X: p.X, Y: p.Y + 1}
	case Right:
		return Point{X: p.X + 1, Y: p.Y}
	case Bottom:
		return Point{X: p.X, Y: p.Y - 1}
	case Left:
		return Point{X: p.X - 1, Y: p.Y}
	}
	return Point{X: -1, Y: -1}
}

func (d Dir) CheckParallel(d2 Dir) bool {
	switch d {
	case Top:
		return d2 == Bottom
	case Right:
		return d2 == Left
	case Bottom:
		return d2 == Top
	case Left:
		return d2 == Right
	default:
		return false
	}

}
