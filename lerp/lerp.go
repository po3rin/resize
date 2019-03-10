package lerp

// Point lerpで使う座標
type Point struct {
	X int
	Y int
}

// Points lerpで使う近傍4座標
type Points [4]Point

// PosDependFunc 座標依存関数
type PosDependFunc func(x, y int) float64

// Lerp calicurate relp
func Lerp(f PosDependFunc, a float64, b float64, ps Points) float64 {
	n := (1.0-b)*(1.0-a)*f(ps[0].X, ps[0].Y) +
		a*(1.0-b)*f(ps[1].X, ps[0].Y) +
		b*(1-a)*f(ps[0].X, ps[1].Y) +
		a*b*f(ps[1].X, ps[1].Y)
	return n
}
