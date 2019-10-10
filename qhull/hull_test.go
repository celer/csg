package qhull

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/celer/csg/csg"
)

func init() {
	os.RemoveAll("output")
}

func SaveCSG(h *Hull, file string) {
	c := h.ToCSG()
	os.Mkdir("output", 0755)
	out, err := os.Create(filepath.Join("output", file))
	if err != nil {
		panic(err)
	}
	c.MarshalToASCIISTL(out)
	out.Close()
}

func TestHotDog(t *testing.T) {
	h := &Hull{}

	s1 := csg.NewSphere(&csg.SphereOptions{Center: &csg.Vector{0, 0, 0}})
	s2 := csg.NewSphere(&csg.SphereOptions{Center: &csg.Vector{4, 4, 4}})

	err := h.BuildFromCSG([]*csg.CSG{s1, s2})
	if err != nil {
		t.Fatal(err)
	}

	SaveCSG(h, "HotDog.stl")
}

func TestV(t *testing.T) {
	h := &Hull{}

	s1 := csg.NewSphere(&csg.SphereOptions{Center: &csg.Vector{0, 0, 0}})
	s2 := csg.NewSphere(&csg.SphereOptions{Center: &csg.Vector{4, 0, -4}})
	s3 := csg.NewSphere(&csg.SphereOptions{Center: &csg.Vector{4, 0, 4}})

	err := h.BuildFromCSG([]*csg.CSG{s1, s2, s3})
	if err != nil {
		t.Fatal(err)
	}

	SaveCSG(h, "V.stl")
}
func TestTri(t *testing.T) {
	h := &Hull{}

	s1 := csg.NewSphere(&csg.SphereOptions{Center: &csg.Vector{0, 0, 0}})
	s2 := csg.NewSphere(&csg.SphereOptions{Center: &csg.Vector{4, 0, -4}})
	s3 := csg.NewSphere(&csg.SphereOptions{Center: &csg.Vector{4, 0, 4}})
	s4 := csg.NewSphere(&csg.SphereOptions{Center: &csg.Vector{4, 4, 0}})

	err := h.BuildFromCSG([]*csg.CSG{s1, s2, s3, s4})
	if err != nil {
		t.Fatal(err)
	}

	SaveCSG(h, "Tri.stl")
}

func TestExtraLongHotDog(t *testing.T) {
	h := &Hull{}

	s1 := csg.NewSphere(&csg.SphereOptions{Center: &csg.Vector{0, 0, 0}, Slices: 50, Stacks: 50})
	s2 := csg.NewSphere(&csg.SphereOptions{Center: &csg.Vector{4, 4, 4}, Slices: 50, Stacks: 50})
	s3 := csg.NewSphere(&csg.SphereOptions{Center: &csg.Vector{-4, -4, -4}, Slices: 50, Stacks: 50})

	err := h.BuildFromCSG([]*csg.CSG{s1, s2, s3})
	if err != nil {
		t.Fatal(err)
	}

	SaveCSG(h, "ExtraLongHotDog.stl")
}

func TestBasicHull(t *testing.T) {
	vs := make([]*csg.Vector, 0)
	for i := 0; i < 100; i++ {
		v := &csg.Vector{}
		v.SetRandom(5, 100)
		vs = append(vs, v)
	}
	h := &Hull{}
	err := h.Build(vs, len(vs))
	if err != nil {
		t.Fatal(err)
	}

	SaveCSG(h, "BasicHull.stl")
}

func AssertVectorEq(t *testing.T, v1, v2 *csg.Vector) {
	if !v1.AlmostEquals(v2) {
		t.Fatalf("Expected %v to match %v", v1, v2)
	}
}

func TestSimpleHull(t *testing.T) {
	points := []*csg.Vector{
		&csg.Vector{0.0, 0.0, 0.0},
		&csg.Vector{1.0, 0.5, 0.0},
		&csg.Vector{2.0, 0.0, 0.0},
		&csg.Vector{0.5, 0.5, 0.5},
		&csg.Vector{0.0, 0.0, 2.0},
		&csg.Vector{0.1, 0.2, 0.3},
		&csg.Vector{0.0, 2.0, 0.0},
	}
	h := &Hull{}
	err := h.Build(points, len(points))
	if err != nil {
		t.Fatal(err)
	}

	vs := h.Vertices()
	AssertVectorEq(t, vs[0], &csg.Vector{0.0, 0.0, 0.0})
	AssertVectorEq(t, vs[1], &csg.Vector{2.0, 0.0, 0.0})
	AssertVectorEq(t, vs[2], &csg.Vector{0.0, 0.0, 2.0})
	AssertVectorEq(t, vs[3], &csg.Vector{0.0, 2.0, 0.0})

	faceIndices := h.Faces()

	expected := [][]int{
		[]int{1, 2, 0},
		[]int{3, 1, 0},
		[]int{3, 0, 2},
		[]int{3, 2, 1},
	}

	for i := 0; i < len(vs); i++ {
		if !reflect.DeepEqual(faceIndices[i], expected[i]) {
			t.Fatalf("Expected %#v to match %#v", faceIndices[i], expected[i])
		}
	}

	SaveCSG(h, "SimpleHull.stl")

}
