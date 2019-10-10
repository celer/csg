package csg

import (
	"fmt"
	"io"
)

// CSG is a mesh which represents some constructive solid geometry, it's made up of polygons and
// can be unioned, subtracted or intersected with other CSG meshs. It's important to note that
// polygons in the CSG mesh are not constrained to being triangles, but must be coplanar.
//
// For a more comprehensive discussion of the algorithm see: https://github.com/evanw/csg.js/blob/master/csg.js:w
type CSG struct {
	polygons []*Polygon
}

// NewCSGFromPolygons constructs a new CSG from a slice of polygons
func NewCSGFromPolygons(polygons []*Polygon) *CSG {
	csg := &CSG{}
	csg.polygons = polygons
	return csg
}

// BoundingBox returns the bounding box for the CSG
func (c *CSG) BoundingBox() *Box {
	b := &Box{}
	for _, p := range c.polygons {
		b.AddPolygon(p)
	}
	return b
}

// Clone copies this CSG into a new CSG
func (c *CSG) Clone() *CSG {
	n := &CSG{}
	n.polygons = make([]*Polygon, 0, len(c.polygons))
	n.polygons = append(n.polygons, c.polygons...)
	return n
}

// ToPolygons returns the list of polygons constituting this CSG
func (c *CSG) ToPolygons() []*Polygon {
	return c.polygons
}

// MarshalToASCIISTL writes out this CSG object to an ASCII STL representation
func (c *CSG) MarshalToASCIISTL(out io.Writer) {
	fmt.Fprintf(out, "solid %s\n", "name")

	for _, p := range c.polygons {
		for _, t := range p.Triangles() {
			t.MarshalToASCIISTL(out)
		}
	}
	fmt.Fprintf(out, "endsolid %s\n", "name")
}

// Union combines this CSG object with another CSG object and returns the newly combined mesh.
func (c *CSG) Union(csg *CSG) *CSG {
	a := NewNodeFromPolygons(c.Clone().polygons)
	b := NewNodeFromPolygons(csg.Clone().polygons)

	a.ClipTo(b)
	b.ClipTo(a)
	b.Invert()
	b.ClipTo(a)
	b.Invert()
	a.Build(b.AllPolygons())

	return NewCSGFromPolygons(a.AllPolygons())
}

// Subtract subtracts another CSG object from this object returning the resulting mesh.
func (c *CSG) Subtract(csg *CSG) *CSG {
	a := NewNodeFromPolygons(c.Clone().polygons)
	b := NewNodeFromPolygons(csg.Clone().polygons)

	a.Invert()
	a.ClipTo(b)
	b.ClipTo(a)
	b.Invert()
	b.ClipTo(a)
	b.Invert()
	a.Build(b.AllPolygons())
	a.Invert()

	return NewCSGFromPolygons(a.AllPolygons())
}

// Intersect returns the intersection of two CSGs
func (c *CSG) Intersect(csg *CSG) *CSG {
	a := NewNodeFromPolygons(c.Clone().polygons)
	b := NewNodeFromPolygons(csg.Clone().polygons)

	a.Invert()
	b.ClipTo(a)
	b.Invert()
	a.ClipTo(b)
	b.ClipTo(a)
	a.Build(b.AllPolygons())
	a.Invert()

	return NewCSGFromPolygons(a.AllPolygons())
}

// Inverse clones this CSG and returns a CSG with the normals flipped on all the polygons
func (c *CSG) Inverse() *CSG {
	csg := c.Clone()
	for _, p := range csg.polygons {
		p.Flip()
	}
	return csg
}
