package ik

import (
  "math"
)

type Result struct {
  Distance float64
  Segment  *Segment
}

func Solve(segment *Segment, goal *Vector3, f func(Vector3, float64)) *Result {
  best := &Result{
    math.Inf(1),
    nil,
  }

  step := 36.0
  min := 0.1
  n := 0

  for {
    innerSolve(segment, segment, goal, (math.Pi/180) * step, best, f)

    segment = best.Segment
    step *= 0.5

    for s := segment; s.Child != nil; s = s.Child {
      s.eaStart, s.eaEnd = NarrowRange(s.angle, s.eaStart, s.eaEnd, 0.5)
    }

    n++
    if min > best.Distance || n > 10 {
      return best
    }
  }

  panic("nope")
}

func innerSolve(root *Segment, s *Segment, goal *Vector3, step float64, best *Result, f func(Vector3, float64)) {
  for _, ea := range s.Range(step) {
    s.SetRotation(ea)

    // if this is the last segment (i.e. it has no children), check its position
    // against the best. otherwise, keep recursing.
    if s.Child == nil {
      e := s.End()
      d := goal.Distance(e)
      f(e, d)

      if best.Distance > d {
        best.Segment = root.Clone()
        best.Distance = d
      }

    } else {
      innerSolve(root, s.Child, goal, step, best, f)
    }
  }
}
