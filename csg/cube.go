package csg

// CubeOptions contains the basic options used to construct this cube.
type CubeOptions struct {
	// Center of the cube to be constructed
	Center *Vector
	// Size of the cube to be constructed
	Size *Vector
}

// NewCube creates a new CSG cube with the specified options, if no options are
// specified then a default center of 0,0,0 and size of 1,1,1 will be used
func NewCube(options *CubeOptions) *CSG {
	center := &Vector{X: 0.0, Y: 0.0, Z: 0.0}
	size := &Vector{X: 1.0, Y: 1.0, Z: 1.0}

	if options != nil {
		if options.Center != nil {
			center = options.Center
		}
		if options.Size != nil {
			size = options.Size
		}
	}

	order := [][]int{
		[]int{0, 4, 6, 2},
		[]int{1, 3, 7, 5},
		[]int{0, 1, 5, 4},
		[]int{2, 6, 7, 3},
		[]int{0, 2, 3, 1},
		[]int{4, 5, 7, 6},
	}

	normals := []*Vector{
		&Vector{-1, 0, 0},
		&Vector{+1, 0, 0},
		&Vector{0, -1, 0},
		&Vector{0, +1, 0},
		&Vector{0, 0, -1},
		&Vector{0, 0, +1},
	}

	polygons := make([]*Polygon, 0, 6)

	for j, o := range order {
		vx := make([]*Vertex, 0)
		for _, i := range o {
			v := &Vector{
				X: center.X + size.X/2.0*(2.0*float64(i&1)-1.0),
				Y: center.Y + size.Y/2.0*(2.0*float64(i&2>>1)-1.0),
				Z: center.Z + size.Z/2.0*(2.0*float64(i&4>>2)-1.0),
			}
			vx = append(vx, &Vertex{Position: v, Normal: normals[j]})
		}
		polygons = append(polygons, NewPolygonFromVertices(vx))
	}

	return NewCSGFromPolygons(polygons)
}
