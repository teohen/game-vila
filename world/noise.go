package world

import "math"

type Noise struct {
	perm [512]int
}

func NewNoise(seed int64) *Noise {
	n := &Noise{}
	p := make([]int, 256)
	for i := range p {
		p[i] = i
	}
	s := seed
	for i := 255; i > 0; i-- {
		s = s*6364136223846793005 + 1442695040888963407
		j := int(s & math.MaxInt64 % int64(i+1))
		if j < 0 {
			j = -j
		}
		p[i], p[j] = p[j], p[i]
	}
	for i := range n.perm {
		n.perm[i] = p[i%256]
	}
	return n
}

func fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

func grad(hash int, x, y float64) float64 {
	switch hash & 3 {
	case 0:
		return x + y
	case 1:
		return -x + y
	case 2:
		return x - y
	default:
		return -x - y
	}
}

func (n *Noise) Noise2D(x, y float64) float64 {
	ix := int(math.Floor(x)) & 255
	iy := int(math.Floor(y)) & 255
	fx := x - math.Floor(x)
	fy := y - math.Floor(y)

	u := fade(fx)
	v := fade(fy)

	a := n.perm[ix] + iy
	b := n.perm[ix+1] + iy

	aa := n.perm[a]
	ab := n.perm[a+1]
	ba := n.perm[b]
	bb := n.perm[b+1]

	x1 := lerp(grad(aa, fx, fy), grad(ba, fx-1, fy), u)
	x2 := lerp(grad(ab, fx, fy-1), grad(bb, fx-1, fy-1), u)

	return lerp(x1, x2, v)
}
