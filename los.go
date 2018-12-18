// Package los calculates line-of-sight for roguelike maps
package los

import "math"

// Map is a map of a roguelike level. Make sure all functions are
// linear time
type Map interface {
	OOB(int, int) bool            // Is the point x, y out of bounds?
	Activate(int, int)            // Activate x, y
	Discover(int, int)            // Mark x, y discovered
	Lit(int, int)                 // Light x, y
	UnLit(int, int)               // Darken x, y
	CantSeeThrough(int, int) bool // Can you see through x, y?
}

// CalcVisibility updates the visibility for the given map, about the
// given point (px, py) and with the radius given.
func CalcVisibility(m Map, px, py, radius int) {
	clearlight(m, px, py, radius)
	fov(m, px, py, radius)
}

func clearlight(m Map, px, py, light int) {
	for x := px - light - 1; x < px+light+1; x++ {
		for y := py - light - 1; y < py+light+1; y++ {
			m.UnLit(x, y)
		}
	}
}

func fov(m Map, x, y int, radius int) {
	for i := -radius; i <= radius; i++ { //iterate out of map bounds as well (radius^1)
		for j := -radius; j <= radius; j++ { //(radius^2)
			if i*i+j*j < radius*radius {
				los(m, x, y, x+i, y+j)
			}
		}
	}
}

/* Los calculation http://www.roguebasin.com/index.php?title=LOS_using_strict_definition */
func los(m Map, x0, y0, x1, y1 int) {
	// By taking source by reference, litting can be done outside of this function which would be better made generic.
	var sx int
	var sy int
	var dx int
	var dy int
	var dist float64

	dx = x1 - x0
	dy = y1 - y0

	//determine which quadrant to we're calculating: we climb in these two directions
	if x0 < x1 { //sx = (x0 < x1) ? 1 : -1;
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 { //sy = (y0 < y1) ? 1 : -1;
		sy = 1
	} else {
		sy = -1
	}

	xnext := x0
	ynext := y0

	//calculate length of line to cast (distance from start to final tile)
	dist = sqrt(dx*dx + dy*dy)

	for xnext != x1 || ynext != y1 { //essentially casting a ray of length radius: (radius^3)
		if m.OOB(xnext, ynext) {
			return
		}
		if m.CantSeeThrough(xnext, ynext) {
			m.Discover(xnext, ynext)
			return
		}

		// Line-to-point distance formula < 0.5
		if abs(dy*(xnext-x0+sx)-dx*(ynext-y0))/dist < 0.5 {
			xnext += sx
		} else if abs(dy*(xnext-x0)-dx*(ynext-y0+sy))/dist < 0.5 {
			ynext += sy
		} else {
			xnext += sx
			ynext += sy
		}
	}
	m.Lit(x1, y1)
	m.Discover(x1, y1)
	if !m.OOB(x1, y1) {
		m.Activate(x1, y1)
	}
}

func sqrt(x int) float64 {
	return math.Sqrt(float64(x))
}

func abs(x int) float64 {
	return math.Abs(float64(x))
}
