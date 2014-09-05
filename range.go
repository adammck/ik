package ik

// Given three EulerAngles (mid, start, end), returns two new (start, end)
// angles which are closer to mid.
// TODO: Encapsulate angle pairs in a Range struct
func NarrowRange(angle EulerAngles, start EulerAngles, end EulerAngles, factor float64) (EulerAngles, EulerAngles) {
  s := EulerAngles{
    Heading: mid(start.Heading, angle.Heading, factor),
    Pitch:   mid(start.Pitch, angle.Pitch, factor),
    Bank:    mid(start.Bank, angle.Bank, factor),
  }

  e := EulerAngles{
    Heading: mid(end.Heading, angle.Heading, factor),
    Pitch:   mid(end.Pitch, angle.Pitch, factor),
    Bank:    mid(end.Bank, angle.Bank, factor),
  }

  return s, e
}

func mid(a float64, b float64, f float64) float64 {
  return a + ((b - a) * f)
}
