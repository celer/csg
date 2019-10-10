package qhull

import (
	"github.com/celer/csg/csg"
)

// FaceState
type FaceState int

const (
	VISIBLE    FaceState = 1
	NON_CONVEX           = 2
	DELETED              = 3
)

// Face is a representation of a particular polygon, in a halfedge type structure
type Face struct {
	csg.Plane

	edge *HalfEdge

	area     float64
	index    int
	numVerts int

	centroid *csg.Vector

	next *Face

	mark FaceState

	outside *Vertex
}

// NewFaceFromTriangle constructs a new face from a triangular set of points
func NewFaceFromTriangle(v0, v1, v2 *Vertex) *Face {
	f := &Face{}

	he0 := &HalfEdge{Vertex: v0, Face: f}
	he1 := &HalfEdge{Vertex: v1, Face: f}
	he2 := &HalfEdge{Vertex: v2, Face: f}

	f.mark = VISIBLE

	he0.prev = he2
	he0.next = he1
	he1.prev = he0
	he1.next = he2
	he2.prev = he1
	he2.next = he0

	f.edge = he0

	f.SetFromPoints(v0.point, v1.point, v2.point)

	// compute the normal and offset
	f.update(0)

	return f
}

func (f *Face) update(m float64) {
	f.updateVertexCount()
	f.updateCentroid()
}

// ToPolygon converts the face to a polygon
func (f *Face) ToPolygon() *csg.Polygon {
	v := make([]*csg.Vertex, 0)
	he := f.edge
	for {
		v = append(v, &csg.Vertex{Position: he.Vertex.point, Normal: f.Normal})
		he = he.next
		if he == f.edge {
			break
		}
	}
	return csg.NewPolygonFromVertices(v)
}

func (f *Face) updateVertexCount() {

	he1 := f.edge.next
	he2 := he1.next

	f.numVerts = 2

	for he2 != f.edge {
		he1 = he2
		he2 = he2.next
		f.numVerts++
	}

}

func (f *Face) mergeAdjacentFace(hedgeAdj *HalfEdge, discarded []*Face) int {
	oppFace := hedgeAdj.OppositeFace()
	numDiscarded := 0

	discarded[numDiscarded] = oppFace
	numDiscarded++
	oppFace.mark = DELETED

	hedgeOpp := hedgeAdj.Opposite()

	hedgeAdjPrev := hedgeAdj.prev
	hedgeAdjNext := hedgeAdj.next
	hedgeOppPrev := hedgeOpp.prev
	hedgeOppNext := hedgeOpp.next

	for hedgeAdjPrev.OppositeFace() == oppFace {
		hedgeAdjPrev = hedgeAdjPrev.prev
		hedgeOppNext = hedgeOppNext.next
	}

	for hedgeAdjNext.OppositeFace() == oppFace {
		hedgeOppPrev = hedgeOppPrev.prev
		hedgeAdjNext = hedgeAdjNext.next
	}

	for hedge := hedgeOppNext; hedge != hedgeOppPrev.next; hedge = hedge.next {
		hedge.Face = f
	}

	if hedgeAdj == f.edge {
		f.edge = hedgeAdjNext
	}

	// handle the half edges at the head
	var discardedFace *Face

	discardedFace = f.connectHalfEdges(hedgeOppPrev, hedgeAdjNext)
	if discardedFace != nil {
		discarded[numDiscarded] = discardedFace
		numDiscarded++
	}

	// handle the half edges at the tail
	discardedFace = f.connectHalfEdges(hedgeAdjPrev, hedgeOppNext)
	if discardedFace != nil {
		discarded[numDiscarded] = discardedFace
		numDiscarded++
	}

	f.update(0)

	//TODO
	//f.checkConsistency()

	return numDiscarded
}

func (f *Face) connectHalfEdges(hedgePrev, hedge *HalfEdge) *Face {
	var discardedFace *Face

	if hedgePrev.OppositeFace() == hedge.OppositeFace() {
		// then there is a redundant edge that we can get rid off

		oppFace := hedge.OppositeFace()
		var hedgeOpp *HalfEdge

		if hedgePrev == f.edge {
			f.edge = hedge
		}
		if oppFace.numVerts == 3 {
			// then we can get rid of the opposite face altogether
			hedgeOpp = hedge.Opposite().prev.Opposite()

			oppFace.mark = DELETED
			discardedFace = oppFace
		} else {
			hedgeOpp = hedge.Opposite().next

			if oppFace.edge == hedgeOpp.prev {
				oppFace.edge = hedgeOpp
			}
			hedgeOpp.prev = hedgeOpp.prev.prev
			hedgeOpp.prev.next = hedgeOpp
		}
		hedge.prev = hedgePrev.prev
		hedge.prev.next = hedge

		hedge.opposite = hedgeOpp
		hedgeOpp.opposite = hedge

		// oppFace was modified, so need to recompute
		oppFace.update(0)
	} else {
		hedgePrev.next = hedge
		hedge.prev = hedgePrev
	}
	return discardedFace
}

// GetEdge returns a specific edge of the face
func (f *Face) GetEdge(i int) *HalfEdge {
	he := f.edge
	for i > 0 {
		he = he.next
		i--
	}
	for i < 0 {
		he = he.prev
		i++
	}
	return he
}

func (f *Face) updateCentroid() {
	v := &csg.Vector{}
	he := f.edge
	for {
		v.AddTo(he.Head().point)
		he = he.next
		if he == f.edge {
			break
		}
	}
	v.ScaleTo(1.0 / float64(f.numVerts))
	f.centroid = v
}

//FaceList is a list of faces
type FaceList struct {
	head *Face
	tail *Face
}

//Clear resets the list of faces
func (f *FaceList) Clear() {
	f.head = nil
	f.tail = nil
}

// Add adds a face to the end of our list
func (f *FaceList) Add(face *Face) {
	if f.head == nil {
		f.head = face
	} else {
		f.tail.next = face
	}
	face.next = nil
	f.tail = face
}

//First returns the first face in the list
func (f *FaceList) First() *Face {
	return f.head
}

//IsEmpty returns true if the list is empty
func (f *FaceList) IsEmpty() bool {
	return f.head == nil
}
