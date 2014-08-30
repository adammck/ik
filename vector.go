package three

import "fmt"

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

func MakeVector3(x float64, y float64, z float64) *Vector3 {
	return &Vector3{x, y, z}
}

func (v *Vector3) String() string {
	return fmt.Sprintf("&Vec3{x=%g y=%g z=%g}", v.X, v.Y, v.Z)
}
