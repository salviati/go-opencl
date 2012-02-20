package main

import (
	"math"
)

type matrix []float32

// Rotate.
func R(angle float32) matrix {
	s64, c64 := math.Sincos(float64(angle))
	s := float32(s64)
	c := float32(c64)

	return []float32{c, -s, s, c}
}

// Scale.
func S(sx, sy float32) matrix {
	return []float32{sx, 0, 0, sy}
}

// Shear.
func H(hx, hy float32) matrix {
	return []float32{1, hx, hy, 1}
}

// Inverts m.
func (m matrix) inv() {
	det := m[0]*m[3] - m[1]*m[2]
	m[1], m[2] = -m[1]/det, -m[2]/det
	m[0], m[3] = m[3]/det, m[0]/det
}

// Returns a*b.
func mul(a, b matrix) matrix {
	return []float32{
		a[0]*b[0] + a[1]*b[2],
		a[0]*b[1] + a[1]*b[3],
		a[2]*b[0] + a[3]*b[2],
		a[2]*b[1] + a[3]*b[3],
	}
}
