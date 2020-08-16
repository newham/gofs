package api

import "testing"

func TestMV(t *testing.T) {
	MV(map[string]string{"1": "tmp/test1.txt"}, "")
}
