package csg

import "math"

//CylinderOptions specifies the options for the new cylinder
type CylinderOptions struct {
	// Start of the cylinder
	Start *Vector
	// End of the cylinder
	End *Vector
	// Radius of the cylinder
	Radius float64
	// Slices in the cylinder
	Slices int
}

//NewCylinder returns a new CSG cylinder
func NewCylinder(options *CylinderOptions) *CSG {
	s := &Vector{0, -1, 0}
	e := &Vector{0, 1, 0}
	radius := 1.0
	slices := 16

	if options != nil {
		if options.Start != nil {
			s = options.Start
		}
		if options.End != nil {
			e = options.End
		}
		if options.Radius != 0.0 {
			radius = options.Radius
		}
		if options.Slices != 0 {
			slices = options.Slices
		}
	}

	ray := e.Minus(s)
	axisZ := ray.Unit()

	isY := 0.0
	nisY := 1.0
	if math.Abs(axisZ.Y) > 0.5 {
		isY = 1.0
		nisY = 0.0
	}

	axisX := (&Vector{isY, nisY, 0}).Cross(axisZ).Unit()
	axisY := axisX.Cross(axisZ).Unit()

	start := NewVertexFromVectors(s, axisZ.Negated())
	end := NewVertexFromVectors(e, axisZ.Unit())

	point := func(stack float64, slice float64, normalBlend float64) *Vertex {
		angle := slice * math.Pi * 2.0
		out := axisX.Times(math.Cos(angle)).Plus(axisY.Times(math.Sin(angle)))
		pos := start.Position.Plus(ray.Times(stack)).Plus(out.Times(radius))
		normal := out.Times(1.0 - math.Abs(normalBlend)).Plus(axisZ.Times(normalBlend))
		return NewVertexFromVectors(pos, normal)
	}

	polygons := make([]*Polygon, 0)

	for i := 0.0; i < float64(slices); i++ {
		t0 := i / float64(slices)
		t1 := (i + 1) / float64(slices)
		polygons = append(polygons, NewPolygonFromVertices([]*Vertex{start, point(0.0, t0, -1), point(0.0, t1, -1)}))
		polygons = append(polygons, NewPolygonFromVertices([]*Vertex{point(0.0, t1, 0.0), point(0.0, t0, 0.0), point(1.0, t0, 0), point(1, t1, 0)}))
		polygons = append(polygons, NewPolygonFromVertices([]*Vertex{end, point(1.0, t1, 1.0), point(1.0, t0, 1.0)}))
	}

	return NewCSGFromPolygons(polygons)
}
