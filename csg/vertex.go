package csg

//Vertex is a point with an associated normal
type Vertex struct {
	//Position of this vertex
	Position *Vector
	// Normal of this vertex
	Normal *Vector
}

//NewVertexFromVectors constructs a new vertex from vectors and points
func NewVertexFromVectors(p *Vector, n *Vector) *Vertex {
	return &Vertex{Position: p, Normal: n}
}

//NewVertex constructs a new vertex
func NewVertex() *Vertex {
	return &Vertex{Position: &Vector{}, Normal: &Vector{}}
}

//Clone copies this vertex
func (v *Vertex) Clone() *Vertex {
	return &Vertex{Position: v.Position.Clone(), Normal: v.Normal.Clone()}
}

//Flip flips the normal of this vertex
func (v *Vertex) Flip() {
	v.Normal = v.Normal.Negated()
}

//Interpolate lerp's this vertex and it's normal
func (v *Vertex) Interpolate(other *Vertex, t float64) *Vertex {
	return &Vertex{
		Position: v.Position.Lerp(other.Position, t),
		Normal:   v.Normal.Lerp(other.Normal, t),
	}
}
