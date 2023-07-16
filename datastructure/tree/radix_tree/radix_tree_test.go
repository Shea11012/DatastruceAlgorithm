package radix_tree

import "testing"

func TestRadix_Insert(t *testing.T) {
	rdx := NewRadix()

	rdx.Insert("test")
	rdx.Insert("solo")
	rdx.Insert("tester")
}
