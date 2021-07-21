package relate

import (
	"github.com/spatial-go/geoos/algorithm/calc"
	"github.com/spatial-go/geoos/algorithm/matrix"
)

// Computer Computes the topological relationship between two Geometries.
type Computer struct {
	Arg            []matrix.Steric // the arg(s) of the operation
	IntersectBound bool
}

func (r *Computer) computeIM(arg []matrix.Steric, intersectBound bool) *matrix.IntersectionMatrix {
	r.Arg = arg
	r.IntersectBound = intersectBound
	im := matrix.IntersectionMatrixDefault()
	// since Geometries are finite and embedded in a 2-D space, the EE element must always be 2
	im.Set(calc.EXTERIOR, calc.EXTERIOR, 2)
	// if the Geometries don't overlap there is nothing to do
	if !intersectBound {
		r.computeDisjointIM(im)
		return im
	}
	if arg[0].Equals(arg[1]) {
		im.SetAtLeastString("T*F**FFF*")
	}
	switch r := arg[0].(type) {
	case matrix.Matrix:
		rp := &Point{r, arg[1]}
		return rp.IntersectionMatrix(im)
	case matrix.LineMatrix:
	case matrix.PolygonMatrix:
	}

	r.computeProperIntersectionIM(im)
	return im
}

// computeDisjointIM If the Geometries are disjoint, we need to enter their dimension and
// boundary dimension in the Ext rows in the IM
func (r *Computer) computeDisjointIM(im *matrix.IntersectionMatrix) {
	ga := r.Arg[0]
	if !ga.IsEmpty() {
		im.Set(calc.INTERIOR, calc.EXTERIOR, ga.Dimensions())
		im.Set(calc.BOUNDARY, calc.EXTERIOR, ga.BoundaryDimensions())
	}
	gb := r.Arg[1]
	if !gb.IsEmpty() {
		im.Set(calc.EXTERIOR, calc.INTERIOR, gb.Dimensions())
		im.Set(calc.EXTERIOR, calc.BOUNDARY, gb.BoundaryDimensions())
	}
}

func (r *Computer) computeProperIntersectionIM(im *matrix.IntersectionMatrix) {
	// If a proper intersection is found, we can set a lower bound on the IM.
	dimA := r.Arg[0].Dimensions()
	dimB := r.Arg[1].Dimensions()

	// For Geometry's of dim 0 there can never be proper intersections.

	/**
	 * If edge segments of Areas properly intersect, the areas must properly overlap.
	 */
	if dimA == 2 && dimB == 2 {
		im.SetAtLeastString("212101212")
	} else if dimA == 2 && dimB == 1 {
		//im.SetAtLeast("FFF0FFFF2")
		im.SetAtLeastString("1FFFFF1FF")
	} else if dimA == 1 && dimB == 2 {
		//im.SetAtLeast("F0FFFFFF2")
		im.SetAtLeastString("1F1FFFFFF")

	} else if dimA == 1 && dimB == 1 {
		im.SetAtLeastString("0FFFFFFFF")
	}
}

//   // updateIM update the IM with the sum of the IMs for each component
//    func (r *Computer)  updateIM(im *matrix.IntersectionMatrix) {
// 	 for (Iterator ei = isolatedEdges.iterator(); ei.hasNext(); ) {
// 	   Edge e = (Edge) ei.next();
// 	   e.updateIM(im);
// 	 }
// 	 for (Iterator ni = nodes.iterator(); ni.hasNext(); ) {
// 	   RelateNode node = (RelateNode) ni.next();
// 	   node.updateIM(im);
//
// 	   node.updateIMFromEdges(im);
// 	 }
//    }
//
