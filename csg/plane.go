package csg

// EPSILON used for determining the relationship of a point to the plane
const EPSILON = 1e-5

// PlaneRelationship signifies the relationship of a point/polygon/box/etc to a plane
type PlaneRelationship int

const (
	// COPLANAR the point/polygon is on the plane (given the EPSILON value used)
	COPLANAR PlaneRelationship = 0
	// FRONT the point/polygon is on the front of the plane (considering the normal of the plane & EPSILON)
	FRONT = 1
	// BACK the point/polygon is on the back of the plane (considering the normal of the plane & EPSILON)
	BACK = 2
	// SPANNING the polygon or box or mesh spans the plane (has points both on the front and back)
	SPANNING = 3
)

// Plane is a 3 dimensional plane
type Plane struct {
	// Normal of the plane
	Normal *Vector
	// W of the plane
	W float64
}

// NewPlaneFromPoints construct a new plane from 3 vectors (as points)
func NewPlaneFromPoints(a *Vector, b *Vector, c *Vector) *Plane {
	p := &Plane{}
	p.SetFromPoints(a, b, c)
	return p
}

// Clone the plane
func (p *Plane) Clone() *Plane {
	return &Plane{Normal: p.Normal.Clone(), W: p.W}
}

// SetFromPoints set the Normal and W for this plane from the specified vectors (as points)
func (p *Plane) SetFromPoints(a *Vector, b *Vector, c *Vector) {
	n := b.Minus(a).Cross(c.Minus(a)).Unit()
	p.Normal = n
	p.W = n.Dot(a)
}

// Flip flips the normal of this plane
func (p *Plane) Flip() {
	p.Normal = p.Normal.Negated()
	p.W = -p.W
}

//DistancesToPlane calculates the distance of each specified vector (as a point) from the plane
func (p *Plane) DistancesToPlane(vs []*Vector) []float64 {
	d := make([]float64, len(vs))
	for i, v := range vs {
		d[i] = p.Normal.Dot(v) - p.W
	}
	return d
}

// DistanceToPlane returns the distance of the vector to this plane
func (p *Plane) DistanceToPlane(v *Vector) float64 {
	return p.Normal.Dot(v) - p.W
}
