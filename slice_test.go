package dot

import "testing"

func TestSliceUnique(t *testing.T) {
	input := &[]string{"aa", "bb", "cc", "aa", "dd", "dd"}
	SliceUnique[string](input)
	t.Log(*input)
}
