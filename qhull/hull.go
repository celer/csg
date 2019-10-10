package qhull

import (
	"fmt"
	"log"
	"math"

	"github.com/celer/csg/csg"
)

const AUTOMATIC_TOLERANCE = 0.0
const DOUBLE_PREC = 2.2204460492503131e-16

//Hull creates a hull between two 3d meshes
type Hull struct {
	findIndex          int
	charLength         float64
	Debug              bool
	points             []*Vertex
	vertexPointIndices []int

	maxVertex [3]*Vertex
	minVertex [3]*Vertex

	discardedFaces [3]*Face

	faces       []*Face
	horizon     []*HalfEdge
	claimed     *VertexList
	unclaimed   *VertexList
	newFaces    *FaceList
	numFaces    int
	numPoints   int
	numVertices int

	explicitTolerance float64
	tolerance         float64
}

func (q *Hull) markFaceVertices(face *Face, mark int) {
	he0 := face.edge
	he := he0
	for {
		he.Head().index = mark
		he = he.next
		if he == he0 {
			break
		}
	}
}

func (q *Hull) reindexFacesAndVertices() {

	for i := 0; i < q.numPoints; i++ {
		q.points[i].index = -1
	}
	// remove inactive faces and mark active vertices
	q.numFaces = 0
	for i := 0; i < len(q.faces); i++ {
		face := q.faces[i]
		if face.mark != VISIBLE {
			q.faces = append(q.faces[:i], q.faces[i+1:]...)
			i--
		} else {
			q.markFaceVertices(face, 0)
			q.numFaces++
		}
	}
	if q.Debug {
		log.Printf("Reindexing faces/verts - faces left after removing inactive faces %d", len(q.faces))
	}
	// reindex vertices
	q.numVertices = 0
	for i := 0; i < q.numPoints; i++ {
		vtx := q.points[i]
		if vtx.index == 0 {
			q.vertexPointIndices[q.numVertices] = i
			vtx.index = q.numVertices
			q.numVertices++
		}
	}
}

func (q *Hull) initBuffers(nump int) {
	q.vertexPointIndices = make([]int, nump)
	q.points = make([]*Vertex, nump)
	q.faces = make([]*Face, 0)
	q.horizon = make([]*HalfEdge, 0)
	q.claimed = &VertexList{}
	q.unclaimed = &VertexList{}
	q.newFaces = &FaceList{}
	q.numFaces = 0
	q.numPoints = nump
}

//Build a hull given a set of vectors (as points)
func (q *Hull) Build(points []*csg.Vector, nump int) error {

	if nump < 4 {
		return fmt.Errorf("Less than four input points specified")
	}
	if len(points) < nump {
		return fmt.Errorf("Point array too small for specified number of points")
	}

	q.initBuffers(nump)
	q.setPoints(points, nump)
	q.buildHull()

	return nil
}

func (q *Hull) setPoints(points []*csg.Vector, nump int) {
	for i := 0; i < nump; i++ {
		q.points[i] = NewVertex(points[i], i)
	}
}

func (q *Hull) computeMinAndMax() {
	var max *csg.Vector
	var min *csg.Vector

	max = q.points[0].point.Clone()
	min = q.points[0].point.Clone()

	for i := 0; i < 3; i++ {
		q.maxVertex[i] = q.points[0]
		q.minVertex[i] = q.points[0]
	}

	for _, p := range q.points {
		pnt := p.point

		if pnt.X > max.X {
			max.X = pnt.X
			q.maxVertex[0] = p
		} else if pnt.X < min.X {
			min.X = pnt.X
			q.minVertex[0] = p
		}

		if pnt.Y > max.Y {
			max.Y = pnt.Y
			q.maxVertex[1] = p
		} else if pnt.Y < min.Y {
			min.Y = pnt.Y
			q.minVertex[1] = p
		}

		if pnt.Z > max.Z {
			max.Z = pnt.Z
			q.maxVertex[2] = p
		} else if pnt.Z < min.Z {
			min.Z = pnt.Z
			q.minVertex[2] = p
		}

	}

	if q.Debug {
		log.Printf("Max %v", max)
		log.Printf("Min %v", min)
	}

	cl := math.Max(max.X-min.X, max.Y-min.Y)
	cl = math.Max(max.Z-min.Z, cl)
	if q.explicitTolerance == AUTOMATIC_TOLERANCE {
		q.tolerance = 3 * DOUBLE_PREC *
			(math.Max(math.Abs(max.X), math.Abs(min.X)) +
				math.Max(math.Abs(max.Y), math.Abs(min.Y)) +
				math.Max(math.Abs(max.Z), math.Abs(min.Z)))
	} else {
		q.tolerance = q.explicitTolerance
	}

	if q.Debug {
		log.Printf("Tolerance: %f", q.tolerance)

		for i := 0; i < 3; i++ {
			log.Printf("Max %d %v", i, q.maxVertex[i])
			log.Printf("Min %d %v", i, q.minVertex[i])
		}

	}

}

func (q *Hull) createInitialSimplex() error {
	max := 0.0
	imax := 0

	for i := 0; i < 3; i++ {
		diff := q.maxVertex[i].point.Get(i) - q.minVertex[i].point.Get(i)
		if diff > max {
			max = diff
			imax = i
		}
	}

	if max <= q.tolerance {
		return fmt.Errorf("Input points appear to e coincident")
	}

	vtx := make([]*Vertex, 4)
	vtx[0] = q.maxVertex[imax]
	vtx[1] = q.minVertex[imax]

	var normal *csg.Vector

	n := vtx[1].point.Minus(vtx[0].point).Normalize()
	maxSqr := 0.0

	for _, p := range q.points {
		diff := p.point.Minus(vtx[0].point)
		xprod := n.Cross(diff)

		lenSqr := xprod.LengthSquared()

		if lenSqr > maxSqr && p != vtx[0] && p != vtx[1] {
			maxSqr = lenSqr
			vtx[2] = p
			normal = xprod.Clone()
		}
	}

	if math.Sqrt(maxSqr) < 100*q.tolerance {
		return fmt.Errorf("Input points appear to be colinear")
	}

	normal = normal.Normalize()

	maxDist := 0.0
	d0 := vtx[2].point.Dot(normal)
	for _, p := range q.points {
		dist := math.Abs(p.point.Dot(normal) - d0)
		if dist > maxDist && p != vtx[0] && p != vtx[1] && p != vtx[2] {
			maxDist = dist
			vtx[3] = p
		}
	}

	if math.Sqrt(maxDist) < 100*q.tolerance {
		return fmt.Errorf("Input points appear to be coplanar")
	}

	if q.Debug {
		log.Printf("Initial points")
		for _, v := range vtx {
			log.Printf("\t%v", v)
		}

	}

	tris := make([]*Face, 4)

	if vtx[3].point.Dot(normal)-d0 < 0 {
		tris[0] = NewFaceFromTriangle(vtx[0], vtx[1], vtx[2])
		tris[1] = NewFaceFromTriangle(vtx[3], vtx[1], vtx[0])
		tris[2] = NewFaceFromTriangle(vtx[3], vtx[2], vtx[1])
		tris[3] = NewFaceFromTriangle(vtx[3], vtx[0], vtx[2])

		for i := 0; i < 3; i++ {
			k := (i + 1) % 3
			tris[i+1].GetEdge(1).SetOpposite(tris[k+1].GetEdge(0))
			tris[i+1].GetEdge(2).SetOpposite(tris[0].GetEdge(k))
		}

	} else {
		tris[0] = NewFaceFromTriangle(vtx[0], vtx[2], vtx[1])
		tris[1] = NewFaceFromTriangle(vtx[3], vtx[0], vtx[1])
		tris[2] = NewFaceFromTriangle(vtx[3], vtx[1], vtx[2])
		tris[3] = NewFaceFromTriangle(vtx[3], vtx[2], vtx[0])

		for i := 0; i < 3; i++ {
			k := (i + 1) % 3
			tris[i+1].GetEdge(0).SetOpposite(tris[k+1].GetEdge(1))
			tris[i+1].GetEdge(2).SetOpposite(tris[0].GetEdge((3 - i) % 3))
		}

	}

	q.faces = append(q.faces, tris...)

	for _, v := range q.points {
		if v == vtx[0] || v == vtx[1] || v == vtx[2] || v == vtx[3] {
			continue
		}

		maxDist = q.tolerance
		var maxFace *Face
		for k := 0; k < 4; k++ {
			dist := tris[k].DistanceToPlane(v.point)
			if dist > maxDist {
				maxFace = tris[k]
				maxDist = dist
			}
		}
		if maxFace != nil {
			q.addPointToFace(v, maxFace)
		}
	}

	return nil

}

func (q *Hull) addPointToFace(vtx *Vertex, face *Face) {
	vtx.face = face
	if face.outside == nil {
		q.claimed.Add(vtx)
	} else {
		q.claimed.InsertBefore(vtx, face.outside)
	}
	face.outside = vtx
}

func (q *Hull) nextPointToAdd() *Vertex {
	if !q.claimed.IsEmpty() {
		eyeFace := q.claimed.First().face
		var eyeVtx *Vertex
		maxDist := 0.0
		for vtx := eyeFace.outside; vtx != nil && vtx.face == eyeFace; vtx = vtx.next {
			dist := eyeFace.DistanceToPlane(vtx.point)
			if dist > maxDist {
				maxDist = dist
				eyeVtx = vtx
			}
		}
		return eyeVtx
	} else {
		return nil
	}
}

func (q *Hull) removePointFromFace(vtx *Vertex, face *Face) {
	if vtx == face.outside {
		if vtx.next != nil && vtx.next.face == face {
			face.outside = vtx.next
		} else {
			face.outside = nil
		}
	}
	q.claimed.Delete(vtx)
}

func (q *Hull) removeAllPointsFromFace(face *Face) *Vertex {
	if face.outside != nil {
		end := face.outside
		for end.next != nil && end.next.face == face {
			end = end.next
		}
		q.claimed.DeleteChain(face.outside, end)
		end.next = nil
		return face.outside
	}
	return nil
}

func (q *Hull) deleteFacePoints(face, absorbingFace *Face) {
	faceVtxs := q.removeAllPointsFromFace(face)
	if faceVtxs != nil {
		if absorbingFace == nil {
			q.unclaimed.AddAll(faceVtxs)
		} else {
			vtxNext := faceVtxs
			for vtx := vtxNext; vtx != nil; vtx = vtxNext {
				vtxNext = vtx.next
				dist := absorbingFace.DistanceToPlane(vtx.point)
				if dist > q.tolerance {
					q.addPointToFace(vtx, absorbingFace)
				} else {
					q.unclaimed.Add(vtx)
				}
			}
		}
	}
}

func (q *Hull) calculateHorizon(eyePnt *csg.Vector, edge0 *HalfEdge, face *Face, horizon []*HalfEdge) {
	q.deleteFacePoints(face, nil)
	face.mark = DELETED
	if q.Debug {
		log.Printf("  visiting face %v", face)
	}
	var edge *HalfEdge
	if edge0 == nil {
		edge0 = face.GetEdge(0)
		edge = edge0

	} else {
		edge = edge0.next
	}
	for {
		oppFace := edge.OppositeFace()
		if oppFace.mark == VISIBLE {
			if oppFace.DistanceToPlane(eyePnt) > q.tolerance {
				q.calculateHorizon(eyePnt, edge.Opposite(), oppFace, horizon)
			} else {
				q.horizon = append(q.horizon, edge)
				if q.Debug {
					log.Printf("  adding horizon edge %v", edge)
				}
			}
		}
		edge = edge.next
		if edge == edge0 {
			break
		}
	}
}

func (q *Hull) addAdjoiningFace(eyeVtx *Vertex, he *HalfEdge) *HalfEdge {
	face := NewFaceFromTriangle(eyeVtx, he.Tail(), he.Head())
	q.faces = append(q.faces, face)
	face.GetEdge(-1).SetOpposite(he.Opposite())
	return face.GetEdge(0)
}

func (q *Hull) addNewFaces(newFaces *FaceList, eyeVtx *Vertex, horizon []*HalfEdge) {
	newFaces.Clear()

	var hedgeSidePrev *HalfEdge
	var hedgeSideBegin *HalfEdge

	for _, horizonHe := range horizon {
		{
			hedgeSide := q.addAdjoiningFace(eyeVtx, horizonHe)
			if q.Debug {
				log.Printf("new face: %v", hedgeSide.Face)
			}
			if hedgeSidePrev != nil {
				hedgeSide.next.SetOpposite(hedgeSidePrev)
			} else {
				hedgeSideBegin = hedgeSide
			}
			newFaces.Add(hedgeSide.Face)
			hedgeSidePrev = hedgeSide
		}
		hedgeSideBegin.next.SetOpposite(hedgeSidePrev)
	}
}

const NONCONVEX_WRT_LARGER_FACE = 1
const NONCONVEX = 2

func (q *Hull) oppFaceDistance(he *HalfEdge) float64 {
	return he.Face.DistanceToPlane(he.opposite.Face.centroid)
}

func (q *Hull) doAdjacentMerge(face *Face, mergeType int) bool {
	hedge := face.edge

	convex := true
	for {
		oppFace := hedge.OppositeFace()
		merge := false
		var dist1 float64

		if mergeType == NONCONVEX {
			// then merge faces if they are definitively non-convex
			if q.oppFaceDistance(hedge) > -q.tolerance ||
				q.oppFaceDistance(hedge.opposite) > -q.tolerance {
				merge = true
			}
		} else {
			// merge faces if they are parallel or non-convex
			// wrt to the larger face; otherwise, just mark
			// the face non-convex for the second pass.
			if face.area > oppFace.area {
				dist1 = q.oppFaceDistance(hedge)
				if dist1 > -q.tolerance {
					merge = true
				} else if q.oppFaceDistance(hedge.opposite) > -q.tolerance {
					convex = false
				}
			} else {
				if q.oppFaceDistance(hedge.opposite) > -q.tolerance {
					merge = true
				} else if q.oppFaceDistance(hedge) > -q.tolerance {
					convex = false
				}
			}
		}

		if merge {
			if q.Debug {
				log.Printf("  merging %v and %v", face, oppFace)
			}

			numd := face.mergeAdjacentFace(hedge, q.discardedFaces[0:])
			for i := 0; i < numd; i++ {
				q.deleteFacePoints(q.discardedFaces[i], face)
			}
			if q.Debug {
				log.Printf("  result %v", face)
			}
			return true
		}
		hedge = hedge.next
		if hedge == face.edge {
			break
		}
	}

	if !convex {
		face.mark = NON_CONVEX
	}
	return false
}

func (q *Hull) addPointToHull(eyeVtx *Vertex) {
	q.horizon = make([]*HalfEdge, 0)
	q.unclaimed.Clear()

	if q.Debug {
		log.Printf("Adding point: %v", eyeVtx)
		log.Printf("which is %f above face %v", eyeVtx.face.DistanceToPlane(eyeVtx.point), eyeVtx.face)
	}

	q.removePointFromFace(eyeVtx, eyeVtx.face)
	q.calculateHorizon(eyeVtx.point, nil, eyeVtx.face, q.horizon)
	q.newFaces.Clear()
	q.addNewFaces(q.newFaces, eyeVtx, q.horizon)

	// first merge pass ... merge faces which are non-convex
	// as determined by the larger face
	if q.Debug {
		log.Printf("First merge")
	}
	for face := q.newFaces.First(); face != nil; face = face.next {
		if face.mark == VISIBLE {
			for q.doAdjacentMerge(face, NONCONVEX_WRT_LARGER_FACE) {
			}
		}
	}

	// second merge pass ... merge faces which are non-convex
	// wrt either face
	if q.Debug {
		log.Printf("Second merge")
	}
	for face := q.newFaces.First(); face != nil; face = face.next {
		if face.mark == NON_CONVEX {
			face.mark = VISIBLE
			for q.doAdjacentMerge(face, NONCONVEX) {
			}
		}
	}
	q.resolveUnclaimedPoints(q.newFaces)
}

func (q *Hull) resolveUnclaimedPoints(newFaces *FaceList) {
	vtxNext := q.unclaimed.First()
	for vtx := vtxNext; vtx != nil; vtx = vtxNext {
		vtxNext = vtx.next

		maxDist := q.tolerance
		var maxFace *Face
		for newFace := newFaces.First(); newFace != nil; newFace = newFace.next {
			if newFace.mark == VISIBLE {
				dist := newFace.DistanceToPlane(vtx.point)
				if dist > maxDist {
					maxDist = dist
					maxFace = newFace
				}
				if maxDist > 1000*q.tolerance {
					break
				}
			}
		}
		if maxFace != nil {
			q.addPointToFace(vtx, maxFace)
			if q.Debug && vtx.index == q.findIndex {
				log.Printf("%d CLAIMED BY %v", q.findIndex, maxFace)
			}
		} else {
			if q.Debug && vtx.index == q.findIndex {
				log.Printf("%d DISCARDED", q.findIndex)
			}
		}
	}
}

//Vertices returns the vertices used to construct this hull
func (q *Hull) Vertices() []*csg.Vector {
	ret := make([]*csg.Vector, q.numVertices)
	for i := 0; i < q.numVertices; i++ {
		ret[i] = q.points[q.vertexPointIndices[i]].point
	}
	return ret
}

//BuildFromCSG builds a hull from some number of CSGs
func (q *Hull) BuildFromCSG(csgs []*csg.CSG) error {
	points := make([]*csg.Vector, 0)
	for _, c := range csgs {
		for _, p := range c.ToPolygons() {
			for _, v := range p.Vertices {
				points = append(points, v.Position)
			}
		}
	}
	return q.Build(points, len(points))
}

// ToCSG converts this hull into a CSG object
func (q *Hull) ToCSG() *csg.CSG {
	polys := make([]*csg.Polygon, len(q.faces))
	for i, face := range q.faces {
		polys[i] = face.ToPolygon()
	}
	return csg.NewCSGFromPolygons(polys)
}

//Faces returns the faces which constitude this hull
func (q *Hull) Faces() [][]int {
	indexFlags := 0
	allFaces := make([][]int, len(q.faces))
	k := 0
	for _, face := range q.faces {
		allFaces[k] = q.getFaceIndices(face, indexFlags)
		k++
	}
	return allFaces

}

const CLOCKWISE = 0x1
const INDEXED_FROM_ONE = 0x2
const INDEXED_FROM_ZERO = 0x4
const POINT_RELATIVE = 0x8

func (q *Hull) getFaceIndices(face *Face, flags int) []int {
	ccw := ((flags & CLOCKWISE) == 0)
	indexedFromOne := ((flags & INDEXED_FROM_ONE) != 0)
	pointRelative := ((flags & POINT_RELATIVE) != 0)

	indices := make([]int, face.numVerts)

	hedge := face.edge
	k := 0
	for {
		idx := hedge.Head().index
		if pointRelative {
			idx = q.vertexPointIndices[idx]
		}
		if indexedFromOne {
			idx++
		}
		indices[k] = idx
		k++
		if ccw {
			hedge = hedge.next
		} else {
			hedge = hedge.prev
		}
		if hedge == face.edge {
			break
		}
	}
	return indices
}

func (q *Hull) buildHull() error {
	cnt := 0
	eyeVtx := &Vertex{}

	q.computeMinAndMax()
	err := q.createInitialSimplex()
	if err != nil {
		return err
	}

	for {
		eyeVtx = q.nextPointToAdd()
		if eyeVtx == nil {
			break
		}
		q.addPointToHull(eyeVtx)
		cnt++
		if q.Debug {
			log.Printf("iteration %d done", cnt)
		}
	}
	q.reindexFacesAndVertices()
	if q.Debug {
		log.Printf("hull done")
	}
	return nil

}
