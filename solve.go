package ik

import (
  "math"
  "fmt"
)

type Result struct {
  Distance float64
  Segment  *Segment
}

func Solve(start *Segment, end *Segment, goal *Vector3, accuracy float64) *Result {
  best := &Result{
    math.Inf(1),
    nil,
  }

  segment := start
  step := 0.075
  n := 0

  for {
    n++
    fmt.Printf("Best%d=%v\n", n, best.Distance)

    prevBestDistance := best.Distance
    innerSolve(segment, segment, end, goal, step, best)

    improvement := math.Abs(best.Distance - prevBestDistance)
    if improvement < (accuracy * 0.1) {
      fmt.Printf("Stuck at: %0.4f after %d iterations\n", best.Distance, n)
      return best
    }

    //fmt.Printf("Best=%v\n", best)

    // TODO: Don't bother narrowing the range and increasing the resolution if
    //       the best solution is outside of the possible range.

    segment = best.Segment
    //step *= 0.4

    for s := segment; s.Child != nil; s = s.Child {
      s.Pair = s.Pair.Zoom(s.Angle, 0.25)
    }

    if n > 10 {
      fmt.Printf("Giving up at %0.4f after %d iterations\n", best.Distance, n)
      return best
    }

    if accuracy > best.Distance {
      fmt.Printf("Satisfied at %0.4f after %d iterations\n", best.Distance, n)
      return best
    }
  }

  panic("nope")
}

func innerSolve(root *Segment, s *Segment, end *Segment, goal *Vector3, step float64, best *Result) {
  for _, ea := range s.EulerAngles(step) {
    s.SetRotation(ea)

    // if this is the end segment, check its position against the best.
    // otherwise, keep recursing.
    //fmt.Printf("%p vs %p == %v\n", s, end, (s==end))
    if s.Name == end.Name {
    //s.Child == nil {
      d := goal.Distance(s.End())

      if best.Distance > d {
        best.Segment = root.Clone()
        best.Distance = d
        //fmt.Printf("New Best: %v\n", best)
      }

    } else if s.Child != nil {
      innerSolve(root, s.Child, end, goal, step, best)

    } else {
      fmt.Printf("no children!? %v\n", s)
    }
  }
}
