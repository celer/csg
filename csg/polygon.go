package csg

import (
	"fmt"
	"io"
)

// NewTriangle creates a new polygon from 3 points
func NewTriangle(a *Vertex, b *Vertex, c *Vertex, plane *Plane) *Polygon {
	return &Polygon{Vertices: []*Vertex{a, b, c}, Plane: plane}
}

// Polygon is a 3 dimensional polygon with 3 or more vertices
type Polygon struct {
	Vertices []*Vertex
	Plane    *Plane
}

// NewPolygonFromVertices creates a new polygon from a set of vertices
func NewPolygonFromVertices(vertices []*Vertex) *Polygon {
	return &Polygon{Vertices: vertices, Plane: NewPlaneFromPoints(vertices[0].Position, vertices[1].Position, vertices[2].Position)}
}

//triangulate will (poorly) triangulate this polygon
func triangulate(vertices []*Vertex, plane *Plane) []*Polygon {
	t := make([]*Polygon, 0)
	l := len(vertices)
	if l == 3 {
		t = append(t, NewTriangle(vertices[0], vertices[1], vertices[2], plane))
	} else if l == 4 {
		t = append(t, NewTriangle(vertices[0], vertices[1], vertices[2], plane))
		t = append(t, NewTriangle(vertices[0], vertices[2], vertices[3], plane))
	} else if l == 5 {
		t = append(t, NewTriangle(vertices[0], vertices[1], vertices[2], plane))
		t = append(t, NewTriangle(vertices[0], vertices[2], vertices[4], plane))
		t = append(t, NewTriangle(vertices[2], vertices[3], vertices[4], plane))
	} else if l == 6 {
		t = append(t, NewTriangle(vertices[0], vertices[1], vertices[2], plane))
		t = append(t, NewTriangle(vertices[2], vertices[3], vertices[4], plane))
		t = append(t, NewTriangle(vertices[5], vertices[2], vertices[4], plane))
		t = append(t, NewTriangle(vertices[0], vertices[2], vertices[5], plane))
	} else if l > 6 {
		t = append(t, NewTriangle(vertices[0], vertices[1], vertices[2], plane))
		t = append(t, triangulate(vertices[2:], plane)...)
		t = append(t, NewTriangle(vertices[0], vertices[2], vertices[l-1], plane))
	}
	return t
}

// Triangles returns a triangulation of this polygon
func (p *Polygon) Triangles() []*Polygon {
	return triangulate(p.Vertices, p.Plane)
}

// IsTriangle returns true if this polygon is a triangle
func (p *Polygon) IsTriangle() bool {
	return len(p.Vertices) == 3
}

// Clone clones this polygon
func (p *Polygon) Clone() *Polygon {
	vs := make([]*Vertex, 0)
	for _, cp := range p.Vertices {
		vs = append(vs, cp.Clone())
	}
	return NewPolygonFromVertices(vs)
}

// Flip flips the normal of this polygon by reversing the ordering of points and flipping the normal on the associated plane
func (p *Polygon) Flip() {
	for i := len(p.Vertices)/2 - 1; i >= 0; i-- {
		opp := len(p.Vertices) - 1 - i
		p.Vertices[i], p.Vertices[opp] = p.Vertices[opp], p.Vertices[i]
	}
	p.Plane.Flip()
}

// MarshalToASCIISTL will write this polygon out as ASCII STL
func (p *Polygon) MarshalToASCIISTL(out io.Writer) {
	fmt.Fprintf(out, "facet Normal %f %f %f\n", p.Plane.Normal.X, p.Plane.Normal.Y, p.Plane.Normal.Z)
	fmt.Fprintf(out, "\touter loop\n")
	for _, v := range p.Vertices {
		v.Position.MarshalToASCIISTL(out)
	}
	fmt.Fprintf(out, "\tendloop\n")
	fmt.Fprintf(out, "endfacet\n")

}
