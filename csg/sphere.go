package csg

import (
	"math"
)

// SphereOptions contains options for construction of this sphere
type SphereOptions struct {
	// Center the center of the sphere
	Center *Vector
	// Raidus the radius of the sphere
	Radius float64
	// Slice the number of slices for this sphere
	Slices int
	// Stacks the number of stacks for this sphere
	Stacks int
}

// NewSphere constructs a new sphere given the specified options, if no
// options are specified a default center of 0,0,0, radius of 1, and slice of 16 and stack of 8 is used.
func NewSphere(options *SphereOptions) *CSG {
	center := &Vector{X: 0.0, Y: 0.0, Z: 0.0}
	radius := 1.0
	slices := 16.0
	stacks := 8.0

	if options != nil {
		if options.Center != nil {
			center = options.Center
		}
		if options.Radius != 0.0 {
			radius = options.Radius
		}
		if options.Slices != 0 {
			slices = float64(options.Slices)
		}
		if options.Stacks != 0 {
			stacks = float64(options.Stacks)
		}
	}

	polygons := make([]*Polygon, 0, int(slices*stacks))
	vertices := make([]*Vertex, 0, 4)

	vertex := func(theta, phi float64) {
		theta *= math.Pi * 2.0
		phi *= math.Pi

		dir := &Vector{
			X: math.Cos(theta) * math.Sin(phi),
			Y: math.Cos(phi),
			Z: math.Sin(theta) * math.Sin(phi),
		}
		vertices = append(vertices, &Vertex{center.Plus(dir.Times(radius)), dir})
	}

	for i := 0.0; i < slices; i++ {
		for j := 0.0; j < stacks; j++ {
			vertices = make([]*Vertex, 0, 4)
			vertex(i/slices, j/stacks)
			if j > 0 {
				vertex((i+1)/slices, (j)/stacks)
			}
			if j < stacks-1 {
				vertex((i+1)/slices, (j+1)/stacks)
			}
			vertex(i/slices, (j+1)/stacks)
			polygons = append(polygons, NewPolygonFromVertices(vertices))
		}
	}
	return NewCSGFromPolygons(polygons)

}
