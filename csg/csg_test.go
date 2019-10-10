package csg

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func AssertVectorEq(t *testing.T, v *Vector, X, Y, Z float64) {
	if v.X == X && v.Y == Y && v.Z == Z {

	} else {
		t.Fatal(fmt.Sprintf("Expected vector %v to equal %f %f %f", v, X, Y, Z))
	}
}

func SaveCSG(c *CSG, file string) {
	os.Mkdir("output", 0755)
	out, err := os.Create(filepath.Join("output", file))
	if err != nil {
		panic(err)
	}
	c.MarshalToASCIISTL(out)
	out.Close()
}

func BenchmarkSubtraction(b *testing.B) {
	s1 := NewSphere(&SphereOptions{Center: &Vector{0, 0, 0}, Radius: 5.0, Slices: 5, Stacks: 5})
	s2 := NewSphere(&SphereOptions{Center: &Vector{1, 1, 1}, Radius: 5.0, Slices: 5, Stacks: 5})

	for i := 0; i < b.N; i++ {
		s1.Subtract(s2)
	}

}

func TestSubtraction(b *testing.T) {
	s1 := NewCube(&CubeOptions{Size: &Vector{2, 2, 2}})
	s2 := NewSphere(&SphereOptions{Center: &Vector{1, 1, 1}, Radius: 1.2, Slices: 15, Stacks: 15})

	c := s1.Subtract(s2)

	SaveCSG(c, "basic_sub.stl")

}

func TestSphereGrid(b *testing.T) {
	var last *CSG
	now := time.Now()
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			s := NewSphere(&SphereOptions{Center: &Vector{float64(i), float64(j), 0}, Radius: 2.0, Slices: 6, Stacks: 6})
			if last != nil {
				last = last.Union(s)
			} else {
				last = s
			}
		}
	}

	s := NewSphere(&SphereOptions{Center: &Vector{0, 0, 0}, Radius: 5.0, Slices: 6, Stacks: 6})

	last = last.Subtract(s)

	fmt.Printf("Total %v\n", time.Since(now))
	SaveCSG(last, "SphereGrid.stl")
}

func TestBasicCube(t *testing.T) {

	c := NewCube(&CubeOptions{})

	bb := c.BoundingBox()

	AssertVectorEq(t, &bb.Min, -0.5, -0.5, -0.5)
	AssertVectorEq(t, &bb.Max, 0.5, 0.5, 0.5)
	AssertVectorEq(t, bb.Size(), 1, 1, 1)
	AssertVectorEq(t, bb.Center(), 0, 0, 0)

	c = NewCube(&CubeOptions{Size: &Vector{2, 2, 2}})

	bb = c.BoundingBox()

	AssertVectorEq(t, &bb.Min, -1, -1, -1)
	AssertVectorEq(t, &bb.Max, 1, 1, 1)
	AssertVectorEq(t, bb.Center(), 0, 0, 0)

	c = NewCube(&CubeOptions{Size: &Vector{2, 4, 2}})

	bb = c.BoundingBox()

	AssertVectorEq(t, &bb.Min, -1, -2, -1)
	AssertVectorEq(t, &bb.Max, 1, 2, 1)
	AssertVectorEq(t, bb.Size(), 2, 4, 2)
	AssertVectorEq(t, bb.Center(), 0, 0, 0)

	c = NewCube(&CubeOptions{Size: &Vector{2, 4, 2}, Center: &Vector{1, 1, 1}})

	bb = c.BoundingBox()
	AssertVectorEq(t, bb.Size(), 2, 4, 2)
	AssertVectorEq(t, bb.Center(), 1, 1, 1)

}
