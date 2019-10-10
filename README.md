# Introduction

This is a Constructive Solid Geometry and Quick Hull implemention in golang. 

The CSG implementation is based upon https://github.com/jscad/csg.js/
The QuickHull implementation is based upon https://www.cs.ubc.ca/~lloyd/java/quickhull3d.html

You should also investigate SDFs as a better alternative to traditional CSG methods, see https://github.com/deadsy/sdfx

# Basic usage

## CSG

```golang
  	s1 := NewCube(&CubeOptions{Size: &Vector{2, 2, 2}})
	s2 := NewSphere(&SphereOptions{Center: &Vector{1, 1, 1}, Radius: 1.2, Slices: 15, Stacks: 15})
    
    c:=s1.Subtract(s2)

    out, err := os.Create("v.stl")
	if err != nil {
		panic(err)
	}
	c.MarshalToASCIISTL(out)
    out.Close()

```

will result in:

![Image of resulting subtraction](/images/subtract.png)


## Hulling shapes:

```golang
    h := &Hull{}

	s1 := csg.NewSphere(&csg.SphereOptions{Center: &csg.Vector{0, 0, 0}})
	s2 := csg.NewSphere(&csg.SphereOptions{Center: &csg.Vector{4, 0, -4}})
	s3 := csg.NewSphere(&csg.SphereOptions{Center: &csg.Vector{4, 0, 4}})

	err := h.BuildFromCSG([]*csg.CSG{s1, s2, s3})
	if err != nil {
		panic(err)
	}

    c := h.ToCSG()
	out, err := os.Create("v.stl")
	if err != nil {
		panic(err)
	}
	c.MarshalToASCIISTL(out)
    out.Close()
```

will result in:





# Why?

I do a lot of OpenSCAD programming and I was hoping to find a way to build a more preformant and automatable tool-chain for 3D programming, plus I wanted to get my 3D math and understanding to a better place. Anyways while building this library I discovered SDFs (Signed Distance Functions), which is where I'll re-focus my efforts because it they seem like ultimately a more flexible approach to CSG rather then a mesh or polygon approach, I recommend seeing https://github.com/deadsy/sdfx as it's where I'm going to refocus my efforts. 
