package ik

import (
  "math"
)

type Result struct {
  distance float64
  segment  *Segment
}

func Solve(segment *Segment, goal *Vector3, f func(f Vector3)) (float64, *Segment) {
  best := &Result{
    math.Inf(1),
    nil,
  }

  step := 36.0
  min := 0.1
  n := 0

  for {
    innerSolve(segment, segment, goal, (math.Pi/180) * step, best, f)

    segment = best.segment
    step *= 0.5

    for s := segment; s.Child != nil; s = s.Child {
      s.eaStart, s.eaEnd = NarrowRange(s.angle, s.eaStart, s.eaEnd, 0.5)
    }

    n++
    if min > best.distance || n > 10 {
      return best.distance, best.segment
    }
  }

  panic("nope")
}

func innerSolve(root *Segment, s *Segment, goal *Vector3, step float64, best *Result, f func(f Vector3)) {
  for _, ea := range s.Range(step) {
    s.SetRotation(ea)

    // if this is the last segment (i.e. it has no children), check its position
    // against the best. otherwise, keep recursing.
    if s.Child == nil {
      d := goal.Distance(s.End())
      f(s.End())

      if best.distance > d {
        best.segment = root.Clone()
        best.distance = d
      }

    } else {
      innerSolve(root, s.Child, goal, step, best, f)
    }
  }
}
