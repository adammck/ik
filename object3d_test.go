package three

import (
	"testing"
)

func TestMakeObject3d(t *testing.T) {
	name := "AAA"
	pos := &Vector3{1, 2, 3}
	rot := &EulerAngles{90, 180, 270}
	o := MakeObject3d("AAA", pos, rot)

	if o.String() != name {
		t.Errorf("o.String() returned %v, expected %v", o.String(), name)
	}

	if o.Position() != pos {
		t.Errorf("o.Position() returned %v, expected %v", o.Position(), pos)
	}

	if o.Rotation() != rot {
		t.Errorf("o.Rotation() returned %v, expected %v", o.Rotation(), rot)
	}
}

func TestAdd(t *testing.T) {
	a := &Object3d{name: "a"}
	b := &Object3d{name: "b"}
	c := &Object3d{name: "c"}

	// a <- b

	err := a.Add(b)
	if err != nil {
		t.Fatalf("a.Add(b) returned error: %s", err)
	}

	if len(a.children) != 1 || a.children[0] != b {
		t.Fatalf("a.children is %v, expected [%v]", a.children, b)
	}

	if b.parent != a {
		t.Fatalf("b.parent is %v, expected %v", b.parent, a)
	}

	// a <- b
	//   <- c

	err = a.Add(c)
	if err != nil {
		t.Fatalf("a.Add(c) returned error: %s", err)
	}

	if len(a.children) != 2 {
		t.Fatalf("a has %d children, expected 2", len(a.children))
	}

	// a <- b <- a
	//   <- c

	err = b.Add(a)
	if err == nil {
		t.Fatalf("expected b.Add(a) to return error")
	}

	if len(b.children) != 0 {
		t.Fatalf("b has %d children, expected none", len(a.children))
	}

	// a <- b <- c

	err = b.Add(c)
	if err != nil {
		t.Fatalf("b.Add(c) returned error: %s", err)
	}

	if len(a.children) != 1 || a.children[0] != b {
		t.Fatalf("a.children is %v, expected [%v]", a.children, b)
	}

	if len(b.children) != 1 || b.children[0] != c {
		t.Fatalf("b.children is %v, expected [%v]", b.children, c)
	}
}

func TestRemove(t *testing.T) {
	a := &Object3d{name: "a"}
	b := &Object3d{name: "b"}
	a.Add(b)

	err := a.Remove(b)
	if err != nil {
		t.Fatalf("a.Remove(b) returned error: %s", err)
	}

	if b.parent != nil {
		t.Fatalf("b.parent is %v, expected nil", b.parent)
	}

	if len(a.children) != 0 {
		t.Fatalf("a has %d children, expected none", len(a.children))
	}
}
