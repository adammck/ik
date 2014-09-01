package math3d

import "fmt"

type Vector3 struct {
	X float64
	Y float64
	Z float64
}


// MakeVector3 returns a pointer to a new Vector3.
func MakeVector3(x float64, y float64, z float64) *Vector3 {
	return &Vector3{x, y, z}
}

func (v *Vector3) String() string {
	return fmt.Sprintf("&Vec3{x=%0.2f y=%0.2f z=%0.2f}", v.X, v.Y, v.Z)
}

// MultiplyByMatrix44 returns a new Vector3, by multiplying this vector my a 4x4
// matrix.
func (v *Vector3) MultiplyByMatrix44(m *Matrix44) *Vector3 {
  return &Vector3{
    (v.X * m.m11) + (v.Y * m.m21) + (v.Z * m.m31) + m.m41,
    (v.X * m.m12) + (v.Y * m.m22) + (v.Z * m.m32) + m.m42,
    (v.X * m.m13) + (v.Y * m.m23) + (v.Z * m.m33) + m.m43,
  }
}
