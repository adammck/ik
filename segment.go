package math3d

import (
  "math"
)

type Segment struct {
  parent  *Segment
  Child   *Segment
  angle   *EulerAngles
  eaStart *EulerAngles
	eaEnd   *EulerAngles
  vec     *Vector3
}

func MakeSegment(parent *Segment, eaStart *EulerAngles, eaEnd *EulerAngles, vec *Vector3) *Segment {
  s := &Segment{parent, nil, eaStart, eaStart, eaEnd, vec}

  if parent != nil {
    parent.Child = s
  }

  return s
}

func MakeRootSegment(vec *Vector3) *Segment {
  return MakeSegment(nil, IdentityOrientation, IdentityOrientation, vec)
}

func (s *Segment) Range(step float64) []*EulerAngles {
  ea := make([]*EulerAngles, 0)

  minX := math.Min(s.eaStart.Heading, s.eaEnd.Heading)
  maxX := math.Max(s.eaStart.Heading, s.eaEnd.Heading)
  minY := math.Min(s.eaStart.Pitch,   s.eaEnd.Pitch)
  maxY := math.Max(s.eaStart.Pitch,   s.eaEnd.Pitch)
  minZ := math.Min(s.eaStart.Bank,    s.eaEnd.Bank)
  maxZ := math.Max(s.eaStart.Bank,    s.eaEnd.Bank)

  for x := minX; x <= maxX; x += step {
    for y := minY; y <= maxY; y += step {
      for z := minZ; z <= maxZ; z += step {
        ea = append(ea, MakeEulerAngles(x, y, z))
      }
    }
  }

  return ea
}

// TODO: GTFO, do this in range, don't mutate it from outside.
func (s *Segment) SetRotation(r *EulerAngles) {
  s.angle = r
}

// Start returns a vector3 with the coordinates of the start of this segment, in
// the world coordiante space.
func (s *Segment) Start() *Vector3{
  return s.Project(ZeroVector3)
}

// Start returns a vector3 with the coordinates of the end of this segment, in
// the world coordiante space.
func (s *Segment) End() *Vector3{
  return s.Project(s.vec)
}

// WorldMatrix returns a Matrix4 which can be applied to a vector in this
// segment's coordinate space to convert it to the world space.
func (s *Segment) WorldMatrix() *Matrix44 {

  // if this segment has a parent, our transformation will start at the zero of
  // that space, move by the vector (to the end of the segment), and rotate into
  // this coordinate space.
  if s.parent != nil {
    m := MakeMatrix44(s.parent.vec, s.angle)
    return MultiplyMatrices(m, s.parent.WorldMatrix())

  // no parent means that this is a root segment, so the origin is zero, and
  // transformations only need an angle.
  } else {
    return MakeMatrix44(ZeroVector3, s.angle)
  }
}

// Project transforms a vector in this segment's coordinate space into a vector3
// in the world space.
// (pointer to a) new vector in the world space.
func (s *Segment) Project(v *Vector3) *Vector3 {
  return v.MultiplyByMatrix44(s.WorldMatrix())
}
