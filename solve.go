package ik

import (
  "math"
)

type Result struct {
  distance float64
  angles   []EulerAngles
}

func Solve(segment *Segment, goal *Vector3) (float64, []EulerAngles) {
  best := &Result{
    math.Inf(1),
    []EulerAngles{},
  }

  step := (math.Pi/180) * 4.5
  innerSolve(segment, goal, step, best)
  return best.distance, best.angles
}

func innerSolve(s *Segment, goal *Vector3, step float64, best *Result) {
  for _, ea := range s.Range(step) {
    s.SetRotation(ea)

    // if this is the last segment (i.e. it has no children), check its position
    // against the best. otherwise, keep recursing.
    if s.Child == nil {
      d := goal.Distance(s.End())
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
          best.angles[n-1] = *ss.angle
          n--
        }
      }

    } else {
      innerSolve(s.Child, goal, step, best)
    }
  }
}

