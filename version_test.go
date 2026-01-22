package dot

import "testing"

func TestShowVCSInfo(t *testing.T) {
	ShowVCSInfo()
}

func TestGetVCSInfo(t *testing.T) {
	t.Log(GetVCSInfo().String())
}
