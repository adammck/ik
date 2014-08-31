package three

import (
  "testing"
)

func TestMakeMatrix44(t *testing.T) {
  ea := &EulerAngles{0.1, 0.2, 0.3}
  v3 := &Vector3{1, 2, 3}
  m := MakeMatrix44(v3, ea)

  exp := [4][4]float64{
    [4]float64{0.9362933635841992,   0.3129918257854679, -0.15934507930797787, 0},
    [4]float64{-0.28962947762551555, 0.9447024859948941,  0.1537919979889642,  0},
    [4]float64{0.19866933079506122, -0.09784339500725571, 0.9751703272018158,  0},
    [4]float64{1,                    2,                   3,                   1},
  }

  for r, row := range m.Elements() {
    for c, val := range row {
      if val != exp[r][c] {
        t.Errorf("m%d%d is %v, expected %v", (r+1), (c+1), val, exp[r][c])
      }
    }
  }
}

func TestMultiply(t *testing.T) {
}
