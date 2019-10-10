package qhull

// HalfEdge represents the edge of a particular polygon
type HalfEdge struct {
	Vertex *Vertex
	Face   *Face

	next     *HalfEdge
	prev     *HalfEdge
	opposite *HalfEdge
}

//SetOpposite the opposite half edge
func (h *HalfEdge) SetOpposite(e *HalfEdge) {
	h.opposite = e
	e.opposite = h
}

//Opposite returns the opposite edge of the halfedge
func (h *HalfEdge) Opposite() *HalfEdge {
	return h.opposite
}

//Head returns the first vertex in this halfedge
func (h *HalfEdge) Head() *Vertex {
	return h.Vertex
}

//Tail returns the tail vertex in this halfedge
func (h *HalfEdge) Tail() *Vertex {
	if h.prev != nil {
		return h.prev.Vertex
	}
	return nil
}

//OppositeFace returns the opposite face of this halfedge
func (h *HalfEdge) OppositeFace() *Face {
	if h.opposite != nil {
		return h.opposite.Face
	}
	return nil

}

//Length returns the length of this halfedge
func (h *HalfEdge) Length() float64 {
	if h.Tail() != nil {
		return h.Head().point.Distance(h.Tail().point)
	}
	return -1.0
}
