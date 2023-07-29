package slice

import "testing"

func TestRemove(t *testing.T) {
	list := []int64{10, 20, 30, 40, 50}
	list = RemoveV1(list, 10)
	t.Logf("%v \n", list)
}
