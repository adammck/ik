package three

import (
	"testing"
)

func TestMakeVector3(t *testing.T) {
	v := MakeVector3(1, 2, 3)
	if v.X != 1 {
		t.Errorf("v.X returned %v, expected %v", v.X, 1)
	}
	if v.Y != 2 {
		t.Errorf("v.Y returned %v, expected %v", v.Y, 2)
	}
	if v.Z != 3 {
		t.Errorf("v.Z returned %v, expected %v", v.Z, 3)
	}
}
