package three

import "fmt"

type Object3d struct {
	name     string
	parent   *Object3d
	children []*Object3d
}

func MakeObject3d(name string) *Object3d {
	return &Object3d{name: name}
}

func (obj *Object3d) String() string {
	return obj.name
}

// Parent returns the parent of an object.
func (obj *Object3d) Parent() *Object3d {
	return obj.parent
}

// Add appends a child object, and updates the child's parent.
//
// An error is returned if no change is necessary. We could allow this, but it
// seems preferable to return an error which can be caught and ignored if that
// is the correct course of action.
//
// An error is also returned if the object is an ancestor of the child. Allowing
// this would create a loop in the object hierarchy.
func (obj *Object3d) Add(child *Object3d) error {

	for _, c := range obj.children {
		if c == child {
			return fmt.Errorf("%v is already a child of %v", child, obj)
		}
	}

	for o := obj; o != nil; o = o.parent {
		if o == child {
			return fmt.Errorf("%v is an ancestor of %v", child, obj)
		}
	}

	if child.parent != nil {
		child.parent.Remove(child)
	}

	child.parent = obj
	obj.children = append(obj.children, child)

	return nil
}

// Remove removes a child from the object, and sets the child's parent to nil.
// Returns an error (and doesn't change the child's parent) if the child isn't
// actually a child of the object.
func (obj *Object3d) Remove(child *Object3d) error {

	for i, c := range obj.children {
		if c == child {
			copy(obj.children[i:], obj.children[i+1:])        // shift
			obj.children[len(obj.children)-1] = nil           // remove reference
			obj.children = obj.children[:len(obj.children)-1] // reslice
			child.parent = nil
			return nil
		}
	}

	return fmt.Errorf("%v is not a child of %v", child, obj)
}
