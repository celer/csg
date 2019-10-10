package csg

import (
	"runtime"
	"sync"
)

//IPolygonSplitter is an interface for a specific implementation of a polygon splitter
type IPolygonSplitter interface {
	//SplitPolygons splits the polygons into various slices based upon their orientation to the specified plane
	SplitPolygons(plane *Plane, polygons []*Polygon, coplanarFront, coplanarBack, front, back *[]*Polygon)
}

// BasicPolygonSplitter is a basic implemenation of a polygon splitter
type BasicPolygonSplitter struct {
}

//SplitPolygons splits the polygons into various slices based upon their orientation to the specified plane
func (ps *BasicPolygonSplitter) SplitPolygons(plane *Plane, polygons []*Polygon, coplanarFront, coplanarBack, front, back *[]*Polygon) {

	types := make([]PlaneRelationship, 0, 20)
	for _, polygon := range polygons {

		var polygonType PlaneRelationship

		types = types[:0]
		for _, v := range polygon.Vertices {
			t := plane.Normal.Dot(v.Position) - plane.W
			var pType PlaneRelationship
			if t < (-EPSILON) {
				pType = BACK
			} else if t > EPSILON {
				pType = FRONT
			} else {
				pType = COPLANAR
			}
			polygonType |= pType
			types = append(types, pType)
		}

		switch polygonType {
		case COPLANAR:
			if plane.Normal.Dot(polygon.Plane.Normal) > 0 {
				*coplanarFront = append(*coplanarFront, polygon)
			} else {
				*coplanarBack = append(*coplanarBack, polygon)
			}
			break
		case FRONT:
			*front = append(*front, polygon)
			break
		case BACK:
			*back = append(*back, polygon)
			break
		case SPANNING:
			f := make([]*Vertex, 0)
			b := make([]*Vertex, 0)

			for i := range polygon.Vertices {
				j := (i + 1) % len(polygon.Vertices)
				ti := types[i]
				tj := types[j]

				vi := polygon.Vertices[i]
				vj := polygon.Vertices[j]

				if ti != BACK {
					f = append(f, vi)
				}
				if ti != FRONT {
					if ti != BACK {
						b = append(b, vi.Clone())
					} else {
						b = append(b, vi)
					}
				}
				if (ti | tj) == SPANNING {
					t := (plane.W - plane.Normal.Dot(vi.Position)) / plane.Normal.Dot(vj.Position.Minus(vi.Position))
					v := vi.Interpolate(vj, t)
					f = append(f, v)
					b = append(b, v.Clone())
				}
			}
			if len(f) >= 3 {
				*front = append(*front, NewPolygonFromVertices(f))
			}
			if len(b) >= 3 {
				*back = append(*back, NewPolygonFromVertices(b))
			}
			break
		}
	}
}

// MultiCorePolygonSplitter will utilize multiple goroutines to speed up the splitting of polygons
type MultiCorePolygonSplitter struct {
	// This is the target splitter to use - which should normally use the BasicPolygonSplitter
	Target IPolygonSplitter
}

//SplitPolygons splits the polygons into various slices based upon their orientation to the specified plane
func (ps *MultiCorePolygonSplitter) SplitPolygons(plane *Plane, polygons []*Polygon, coplanarFront, coplanarBack, front, back *[]*Polygon) {

	if len(polygons) > 1000 {
		var wg sync.WaitGroup
		var lock sync.Mutex

		cpus := runtime.NumCPU()

		batchSize := 500
		start := 0
		end := 0
		done := false

		for i := 0; i < cpus; i++ {
			wg.Add(1)
			go func() {
				cF := make([]*Polygon, 0)
				cB := make([]*Polygon, 0)
				f := make([]*Polygon, 0)
				b := make([]*Polygon, 0)

				for {
					lock.Lock()

					if done {
						wg.Done()
						lock.Unlock()
						return
					}

					end = start + batchSize
					if end > len(polygons) {
						end = len(polygons)
						done = true
					}

					p := polygons[start:end]
					start += batchSize

					lock.Unlock()

					cF = cF[:0]
					cB = cB[:0]
					f = f[:0]
					b = b[:0]

					ps.Target.SplitPolygons(plane, p, &cF, &cB, &f, &b)

					lock.Lock()
					*coplanarFront = append(*coplanarFront, cF...)
					*coplanarBack = append(*coplanarBack, cB...)
					*front = append(*front, f...)
					*back = append(*back, b...)
					lock.Unlock()

				}

			}()
		}
		wg.Wait()
	} else {
		ps.Target.SplitPolygons(plane, polygons, coplanarFront, coplanarBack, front, back)
	}
}
