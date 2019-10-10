package csg

import (
	"fmt"
)

// Box is a bounding box representation
type Box struct {
	Min Vector
	Max Vector
}

// Center returns the center of the bounding box
func (b *Box) Center() *Vector {
	return b.Max.Minus(&b.Min).DividedBy(2.0).Plus(&b.Min)
}

//AddVector increases the size of the boundig box to include the vector (as a point)
func (b *Box) AddVector(v *Vector) {
	b.Min.Min(v)
	b.Max.Max(v)
}

//String returns a string representation of the bounding box
func (b *Box) String() string {
	return fmt.Sprintf("box [ %v %v ]", b.Min, b.Max)
}

//AddVertex increases the size of the bounding box to include the vertex
func (b *Box) AddVertex(v *Vertex) {
	b.Min.Min(v.Position)
	b.Max.Max(v.Position)
}

//AddPolygon increases the size of the bounting box to include all vertices from the polygon
func (b *Box) AddPolygon(p *Polygon) {
	for _, v := range p.Vertices {
		b.Min.Min(v.Position)
		b.Max.Max(v.Position)
	}
}

//Divide2x2x2 divides the bouding box into a 2x2x2 set of bounding boxes for use in an OctTree
func (b *Box) Divide2x2x2() []*Box {

	dSize := b.Size().DividedBy(2)

	bs := make([]*Box, 8)
	for x := 0; x < 2; x++ {
		for y := 0; y < 2; y++ {
			for z := 0; z < 2; z++ {
				bc := &Box{}

				min := b.Min.Clone()
				min.X += dSize.X * float64(x)
				min.Y += dSize.Y * float64(y)
				min.Z += dSize.Z * float64(z)

				bc.Min.CopyFrom(min)

				max := min.Plus(dSize)

				bc.Max.CopyFrom(max)
				i := x*4 + y*2 + z
				bs[i] = bc
			}
		}
	}
	return bs
}

//CountContainedPolygonVertices returns a count of the number of vertices in the polygon which are contained within the bounding box
func (b *Box) CountContainedPolygonVertices(p *Polygon) int {
	i := 0
	for _, v := range p.Vertices {
		if b.Contains(v.Position) {
			i++
		}
	}
	return i
}

//Contains determines if the vector (treated as a point) is included in the bounding box
func (b *Box) Contains(v *Vector) bool {
	if v.X > b.Min.X && v.X < b.Max.X {
		if v.Y > b.Min.Y && v.Y < b.Max.Y {
			if v.Z > b.Min.Z && v.Z < b.Max.Z {
				return true
			}
		}
	}
	return false
}

//Corners returns an slice of the vectors (as points) of the corners of the bounding box
func (b *Box) Corners() []*Vector {
	size := b.Size()
	return []*Vector{
		b.Min.Clone(),
		b.Min.Plus(&Vector{X: size.X}),
		b.Min.Plus(&Vector{Y: size.Y}),
		b.Min.Plus(&Vector{Z: size.Z}),
		b.Max.Clone(),
		b.Max.Minus(&Vector{X: size.X}),
		b.Max.Minus(&Vector{Y: size.Y}),
		b.Max.Minus(&Vector{Z: size.Z}),
	}
}

//RelationToPlane returns a PlaneRelationship for the bounding box
func (b *Box) RelationToPlane(p *Plane) PlaneRelationship {

	var boxType PlaneRelationship
	corners := b.Corners()
	for _, corner := range corners {
		t := p.Normal.Dot(corner) - p.W
		var pType PlaneRelationship
		if t < (-EPSILON) {
			pType = BACK
		} else if t > EPSILON {
			pType = FRONT
		} else {
			pType = COPLANAR
		}
		boxType |= pType
	}
	return boxType
}

//Size returns the dimensions of the bounding box as a vector
func (b *Box) Size() *Vector {
	return &Vector{
		X: b.Max.X - b.Min.X,
		Y: b.Max.Y - b.Min.Y,
		Z: b.Max.Z - b.Min.Z,
	}
}
