package dot

import "testing"

func TestGenDateMixedNo(t *testing.T) {
	id := GenDateMixedNo(24, "20060102", "M", false)
	t.Log(id)
}
