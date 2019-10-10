package csg

// Node is a node from a BSP tree
type Node struct {
	plane    *Plane
	front    *Node
	back     *Node
	polygons []*Polygon
}

// NewNodeFromPolygons constructs a node from a slice of polygons
func NewNodeFromPolygons(p []*Polygon) *Node {
	n := &Node{}
	n.Build(p)
	return n
}

// Clone will clone this node and it's children
func (n *Node) Clone() *Node {
	r := &Node{}
	if n.plane != nil {
		r.plane = n.plane.Clone()
	}
	if n.front != nil {
		r.front = n.front.Clone()
	}
	if n.back != nil {
		r.back = n.back.Clone()
	}
	r.polygons = make([]*Polygon, 0)
	for _, p := range n.polygons {
		r.polygons = append(r.polygons, p)
	}
	return r
}

// Invert flips all the normals of the polygons in this node and it's children
func (n *Node) Invert() {
	for _, p := range n.polygons {
		p.Flip()
	}
	n.plane.Flip()
	if n.front != nil {
		n.front.Invert()
	}
	if n.back != nil {
		n.back.Invert()
	}
	temp := n.front
	n.front = n.back
	n.back = temp
}

func (n *Node) clipPolygons(splitter IPolygonSplitter, polygons []*Polygon) []*Polygon {
	if n.plane == nil {
		p := make([]*Polygon, 0)
		p = append(p, polygons...)
		return p
	}

	front := make([]*Polygon, 0, len(polygons)/5)
	back := make([]*Polygon, 0, len(polygons)/5)

	splitter.SplitPolygons(n.plane, polygons, &front, &back, &front, &back)

	if n.front != nil {
		front = n.front.clipPolygons(splitter, front)
	}
	if n.back != nil {
		back = n.back.clipPolygons(splitter, back)
		return append(front, back...)
	}
	return front

}

// ClipPolygons will clip the slice of polygons to this node and it's children
func (n *Node) ClipPolygons(polygons []*Polygon) []*Polygon {
	return n.clipPolygons(n.getPolygonSplitter(len(polygons)), polygons)
}

// ClipTo will clip the node to this node and vice versa
func (n *Node) ClipTo(bsp *Node) {
	n.polygons = bsp.ClipPolygons(n.polygons)
	if n.front != nil {
		n.front.ClipTo(bsp)
	}
	if n.back != nil {
		n.back.ClipTo(bsp)
	}
}

// AllPolygons will return all the polygons associated with this node and it's children
func (n *Node) AllPolygons() []*Polygon {
	polygons := make([]*Polygon, 0)

	polygons = append(polygons, n.polygons...)
	if n.front != nil {
		polygons = append(polygons, n.front.AllPolygons()...)
	}
	if n.back != nil {
		polygons = append(polygons, n.back.AllPolygons()...)
	}
	return polygons
}

func (n *Node) build(splitter IPolygonSplitter, polygons []*Polygon) {
	if len(polygons) == 0 {
		return
	}
	if n.plane == nil {
		n.plane = polygons[0].Plane.Clone()
	}

	front := make([]*Polygon, 0)
	back := make([]*Polygon, 0)

	splitter.SplitPolygons(n.plane, polygons, &n.polygons, &n.polygons, &front, &back)

	if len(front) > 0 {
		if n.front == nil {
			n.front = &Node{}
		}
		n.front.Build(front)
	}
	if len(back) > 0 {
		if n.back == nil {
			n.back = &Node{}
		}
		n.back.Build(back)
	}
}

// getPolygonSplitter will use different polygon splitter implementations depending upon the
// number of polygons to split, so the trade off here is managing multiple goroutines vs
// a single go routine for a small number of polygons
func (n *Node) getPolygonSplitter(numPolys int) IPolygonSplitter {
	var splitter IPolygonSplitter
	if numPolys > 1000 {
		splitter = &MultiCorePolygonSplitter{Target: &BasicPolygonSplitter{}}
	} else {
		splitter = &BasicPolygonSplitter{}
	}
	return splitter
}

//Build constructs a BSP tree for the given slice of polygons
func (n *Node) Build(polygons []*Polygon) {
	n.build(n.getPolygonSplitter(len(polygons)), polygons)
}
