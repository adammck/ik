package ik

import (
  "math"
  "fmt"
)

type Result struct {
  distance float64
  angles   []EulerAngles
  segment  *Segment
}

func Solve(segment *Segment, goal *Vector3, f func(f Vector3)) (float64, *Segment) {
  best := &Result{
    math.Inf(1),
    []EulerAngles{},
    nil,
  }

  step := 18.0
  //min := 1.0

  //for best.distance > min {
    fmt.Printf("solving, step=%d\n", step)
    innerSolve(segment, segment, goal, (math.Pi/180) * step, best, f)
  //  step /= 2

  //  if min > best.distance {
      return best.distance, best.segment
  //  }
  //}

  //panic("nope")
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
        best.distance = d

        // count how deep we are
        n := 0
        for ss := s; ss.parent != nil; ss = ss.parent {
          n++
        }

        // copy the angle of each segment to the best result
        best.angles = make([]EulerAngles, n, n)
        for ss := s; ss.parent != nil; ss = ss.parent {
          best.angles[n-1] = ss.angle
          n--
        }

        best.segment = root.Clone()
      }

    } else {
      innerSolve(root, s.Child, goal, step, best, f)
    }
  }
}
