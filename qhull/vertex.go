package qhull

import (
	"fmt"
	"github.com/celer/csg/csg"
)

type Vertex struct {
	point *csg.Vector
	index int
	prev  *Vertex
	next  *Vertex
	face  *Face
}

func NewVertex(v *csg.Vector, idx int) *Vertex {
	return &Vertex{point: v, index: idx}
}

func (v *Vertex) String() string {
	return fmt.Sprintf("Vertex[ %d: %v ]", v.index, v.point)
}

type VertexList struct {
	head *Vertex
	tail *Vertex
}

func (v *VertexList) Add(vtx *Vertex) {
	if v.head == nil {
		v.head = vtx
	} else {
		v.tail.next = vtx
	}
	vtx.prev = v.tail
	vtx.next = nil
	v.tail = vtx
}

func (v *VertexList) First() *Vertex {
	return v.head
}

func (v *VertexList) IsEmpty() bool {
	return v.head == nil
}

func (v *VertexList) AddAll(vtx *Vertex) {
	if v.head == nil {
		v.head = vtx
	} else {
		v.tail.next = vtx
	}
	vtx.prev = v.tail
	for vtx.next != nil {
		vtx = vtx.next
	}
	v.tail = vtx
}

func (v *VertexList) Delete(vtx *Vertex) {
	if vtx.prev == nil {
		v.head = vtx.next
	} else {
		vtx.prev.next = vtx.next
	}
	if vtx.next == nil {
		v.tail = vtx.prev
	} else {
		vtx.next.prev = vtx.prev
	}
}

func (v *VertexList) DeleteChain(vtx1, vtx2 *Vertex) {
	if vtx1.prev == nil {
		v.head = vtx2.next
	} else {
		vtx1.prev.next = vtx2.next
	}
	if vtx2.next == nil {
		v.tail = vtx1.prev
	} else {
		vtx2.next.prev = vtx1.prev
	}
}

func (v *VertexList) Clear() {
	v.head = nil
	v.tail = nil
}

func (v *VertexList) InsertBefore(vtx, next *Vertex) {
	vtx.prev = next.prev
	if next.prev == nil {
		v.head = vtx
	} else {
		next.prev.next = vtx
	}
	vtx.next = next
	next.prev = vtx
}
