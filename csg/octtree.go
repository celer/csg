package csg

/*
This is a defunct OctTree implementation, this was an attempt to optimize the polygon
splitting using and OctTree, essentially try to remove as many polygons from the
evaluation of polygons to BSP tree planes, but the results were dissappointing
because the primary overhead is in the number of recursive function calls made
to traverse the BSP tree, not the number of polygons to consider for splitting.

So an OctTree based optimization had such a small payoff that the additional complexity
was never worth it.


type OctNode struct {
	box             *Box
	nodes           []*OctNode
	fullPolygons    []*Polygon
	partialPolygons []*Polygon
}

func (o *OctNode) Build(polys []*Polygon, maxDepth int) {
	if o.box == nil {
		o.box = &Box{}
		for _, p := range polys {
			o.box.AddPolygon(p)
		}
	} else {
	}
	if maxDepth > 0 {
		boxes := o.box.Divide2x2x2()
		o.nodes = make([]*OctNode, 2*2*2)
		for i := 0; i < 2*2*2; i++ {
			o.nodes[i] = &OctNode{box: boxes[i]}
			o.nodes[i].Build(polys, maxDepth-1)
		}
	} else {
		o.fullPolygons = make([]*Polygon, 0)
		o.partialPolygons = make([]*Polygon, 0)
		for _, p := range polys {
			c := o.box.CountContainedPolygonVertices(p)
			if c > 0 {
				if p.onodes == nil {
					p.onodes = make([]*OctNode, 0)
				}
				//FIXME make it so the test for polygons either goes all the way or stops immediately
				//FIXME make it so that any node which would otherwise be empty, is discarded
				p.onodes = append(p.onodes, o)
				if c == len(p.Vertices) {
					o.fullPolygons = append(o.fullPolygons, p)
				} else {
					o.partialPolygons = append(o.partialPolygons, p)
				}
			}
		}
	}
}


This was an attempt at improving the performance of polygon splitting by first dividing the CSG into an
OctTree, then considering the orientation of the boxes in the octtree to avoid evaluating as many polygons
as possible. It did have a mild performance increase for very large meshes, but not enought to justfy the
additional complexity

type OctTreePolygonSplitter struct {
	Target IPolygonSplitter
}

func (ps *OctTreePolygonSplitter) SplitPolygons(plane *Plane, polygons []*Polygon, coplanarFront, coplanarBack, front, back *[]*Polygon) {

	if polygons[0].onodes == nil {
		o := OctNode{}
		o.Build(polygons, 2)
	}

	FRONT := 1
	BACK := 2

	oType := make(map[*OctNode]int)

	toProcess := make([]*Polygon, 0)

	for _, polygon := range polygons {
		ptype := 0
		for _, o := range polygon.onodes {
			if ot, ok := oType[o]; ok {
				ptype |= ot
			} else {
				ot := o.box.RelationToPlane(plane)
				oType[o] = ot
				ptype |= ot
			}
		}
		switch ptype {
		case FRONT:
			*front = append(*front, polygon)
		case BACK:
			*back = append(*back, polygon)
		default:
			toProcess = append(toProcess, polygon)
		}
	}

	ps.Target.SplitPolygons(plane, toProcess, coplanarFront, coplanarBack, front, back)
}
*/
