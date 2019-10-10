package csg

import (
	"fmt"
	"io"
	"math"
	"math/rand"
)

// F64Epsilon is the epsilon utilized for AlmostEqual
var F64Epsilon float64

func init() {
	// Calculate the epsilon
	F64Epsilon = math.Nextafter(1, 2) - 1
}

// Vector representation of a vector point in 3 dimensional space
type Vector struct {
	X float64
	Y float64
	Z float64
}

// Get the i'th component of the vector, used for iteration purposes
func (v *Vector) Get(i int) float64 {
	switch i {
	case 0:
		return v.X
	case 1:
		return v.Y
	case 2:
		return v.Z
	}
	return 0.0
}

// LengthSquared returns the length of this vector squared
func (v *Vector) LengthSquared() float64 {
	return v.Dot(v)
}

// AlmostEquals returns true if this vector almost equals another vector (assuming a max delta of F64Epsilon)
func (v *Vector) AlmostEquals(e *Vector) bool {
	return math.Abs(e.X-v.X) < F64Epsilon && math.Abs(e.Y-v.Y) < F64Epsilon && math.Abs(e.Z-v.Z) < F64Epsilon
}

// Equals returns true if the vectors are equal
func (v *Vector) Equals(e *Vector) bool {
	return e.X == v.X && e.Y == v.Y && e.Z == v.Z
}

// Normalize returns a new vector which represents the normalized version of this vector
func (v *Vector) Normalize() *Vector {
	l := v.Length()
	return &Vector{X: v.X / l, Y: v.Y / l, Z: v.Z / l}
}

// Clone returns a clone of this vector
func (v *Vector) Clone() *Vector {
	return &Vector{X: v.X, Y: v.Y, Z: v.Z}
}

// Negated returns a new vector which is negated
func (v *Vector) Negated() *Vector {
	return &Vector{X: -v.X, Y: -v.Y, Z: -v.Z}
}

// Plus returns a new vector which is the resulting addition of these two vectors
func (v *Vector) Plus(a *Vector) *Vector {
	return &Vector{X: v.X + a.X, Y: v.Y + a.Y, Z: v.Z + a.Z}
}

// Minus returns a new vector which is the resulting subtraction of these two vectors
func (v *Vector) Minus(a *Vector) *Vector {
	return &Vector{X: v.X - a.X, Y: v.Y - a.Y, Z: v.Z - a.Z}
}

// Times returns a new vector which is the resulting multiplication of this vector and a scalar
func (v *Vector) Times(a float64) *Vector {
	return &Vector{X: v.X * a, Y: v.Y * a, Z: v.Z * a}
}

// DividedBy returns a new vector which is the resulting division of this vector with a scalar
func (v *Vector) DividedBy(a float64) *Vector {
	return &Vector{X: v.X / a, Y: v.Y / a, Z: v.Z / a}
}

// Dot returns the dot product of this vector and another
func (v *Vector) Dot(a *Vector) float64 {
	return v.X*a.X + v.Y*a.Y + v.Z*a.Z
}

// Lerp returns the linear interpolation of this vector at t as a new vector
func (v *Vector) Lerp(a *Vector, t float64) *Vector {
	return v.Plus(a.Minus(v).Times(t))
}

// Length returns the length of this vector
func (v *Vector) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

// Unit returns a new vector with the unit vector of this vector
func (v *Vector) Unit() *Vector {
	return v.DividedBy(v.Length())
}

// Cross returns the cross project of these two vectors as a new vector
func (v *Vector) Cross(a *Vector) *Vector {
	return &Vector{
		X: v.Y*a.Z - v.Z*a.Y,
		Y: v.Z*a.X - v.X*a.Z,
		Z: v.X*a.Y - v.Y*a.X,
	}
}

// Distance returns the distance of this vector (as a point) and another vector (as a point)
func (v *Vector) Distance(e *Vector) float64 {
	dX := v.X - e.X
	dY := v.Y - e.Y
	dZ := v.Z - e.Z
	return math.Sqrt(dX*dX + dY*dY + dZ*dZ)
}

// Max modifies this vector to have the max X, Y, Z values from the supplied vector
func (v *Vector) Max(m *Vector) {
	if m.X > v.X {
		v.X = m.X
	}
	if m.Y > v.Y {
		v.Y = m.Y
	}
	if m.Z > v.Z {
		v.Z = m.Z
	}
}

// Min modifies this vector to have the min X, Y, Z values from the supplied vector
func (v *Vector) Min(m *Vector) {
	if m.X < v.X {
		v.X = m.X
	}
	if m.Y < v.Y {
		v.Y = m.Y
	}
	if m.Z < v.Z {
		v.Z = m.Z
	}
}

// CopyFrom copies values from the specified vector
func (v *Vector) CopyFrom(e *Vector) {
	v.X = e.X
	v.Y = e.Y
	v.Z = e.Z
}

// Set explicitly sets this vector
func (v *Vector) Set(x, y, z float64) {
	v.X = x
	v.Y = y
	v.Z = z
}

// SetRandom sets this vector randomly using an upper and lower bounds
func (v *Vector) SetRandom(lower, upper float64) {
	r := upper - lower
	v.X = rand.Float64()*r + lower
	v.Y = rand.Float64()*r + lower
	v.Z = rand.Float64()*r + lower
}

// AddTo modifies this vector, adding the supplied vector to it
func (v *Vector) AddTo(e *Vector) {
	v.X += e.X
	v.Y += e.Y
	v.Z += e.Z
}

//SetZero modifies this vector setting it to zero
func (v *Vector) SetZero() {
	v.X = 0
	v.Y = 0
	v.Z = 0
}

//ScaleTo modifies this vector scaling it to the specified s value
func (v *Vector) ScaleTo(s float64) {
	v.X = s * v.X
	v.Y = s * v.Y
	v.Z = s * v.Z
}

// MarshalToASCIISTL marshals this vector to ASCII STL
func (v *Vector) MarshalToASCIISTL(out io.Writer) {
	fmt.Fprintf(out, "\t\tvertex %f %f %f\n", v.X, v.Y, v.Z)
}

//String returns a string representation of this vector
func (v *Vector) String() string {
	return fmt.Sprintf("[%f %f %f]", v.X, v.Y, v.Z)
}
