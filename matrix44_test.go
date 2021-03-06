package ik

import (
	"testing"
)

func TestMakeMatrix44(t *testing.T) {
	ea := &EulerAngles{0.1, 0.2, 0.3}
	v3 := &Vector3{1, 2, 3}
	m := MakeMatrix44(v3, ea)

	exp := [4][4]float64{
		[4]float64{0.9362933635841992, 0.3129918257854679, -0.15934507930797787, 0},
		[4]float64{-0.28962947762551555, 0.9447024859948941, 0.1537919979889642, 0},
		[4]float64{0.19866933079506122, -0.09784339500725571, 0.9751703272018158, 0},
		[4]float64{1, 2, 3, 1},
	}

	for r, row := range m.Elements() {
		for c, val := range row {
			if val != exp[r][c] {
				t.Errorf("m%d%d is %v, expected %v", (r + 1), (c + 1), val, exp[r][c])
			}
		}
	}
}

// http://www.wolframalpha.com/input/?i=%7B%7B1%2C2%2C3%2C4%7D%2C%7B4%2C3%2C2%2C1%7D%2C%7B1%2C3%2C2%2C4%7D%2C%7B4%2C2%2C3%2C1%7D%7D+*+%7B%7B4%2C5%2C6%2C7%7D%2C%7B7%2C6%2C5%2C4%7D%2C%7B4%2C6%2C5%2C7%7D%2C%7B7%2C5%2C6%2C4%7D%7D
func TestMultiply(t *testing.T) {
	a := &Matrix44{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	b := &Matrix44{17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	m := MultiplyMatrices(a, b)

	exp := [4][4]float64{
		[4]float64{250, 260, 270, 280},
		[4]float64{618, 644, 670, 696},
		[4]float64{986, 1028, 1070, 1112},
		[4]float64{1354, 1412, 1470, 1528},
	}

	for r, row := range m.Elements() {
		for c, val := range row {
			if val != exp[r][c] {
				t.Errorf("m%d%d is %v, expected %v", (r + 1), (c + 1), val, exp[r][c])
			}
		}
	}
}
