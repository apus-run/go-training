package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	testCases := []struct {
		name  string
		src   []int64
		index int

		wantSlice []int64
		wantVal   int64
	}{
		{
			name:      "empty slice",
			src:       nil,
			index:     0,
			wantSlice: nil,
			wantVal:   0,
		},
		{
			name:      "delete first element",
			src:       []int64{10, 20, 30, 40, 50},
			index:     0,
			wantSlice: []int64{20, 30, 40, 50},
			wantVal:   10,
		},
		{
			name:      "delete last element",
			src:       []int64{10, 20, 30, 40, 50},
			index:     4,
			wantSlice: []int64{10, 20, 30, 40},
			wantVal:   50,
		},
		{
			name:      "delete middle element",
			src:       []int64{10, 20, 30, 40, 50},
			index:     2,
			wantSlice: []int64{10, 20, 40, 50},
			wantVal:   30,
		},
		{
			name:      "delete out of range",
			src:       []int64{10, 20, 30, 40, 50},
			index:     5,
			wantSlice: nil,
			wantVal:   0,
		},
		{
			name:      "delete out of range",
			src:       []int64{10, 20, 30, 40, 50},
			index:     -1,
			wantSlice: nil,
			wantVal:   0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotSlice, gotVal := Delete[int64](tc.src, tc.index)
			assert.Equal(t, tc.wantSlice, gotSlice)
			assert.Equal(t, tc.wantVal, gotVal)
		})
	}
}

func TestDeleteV1(t *testing.T) {
	list := []int64{10, 20, 30, 40, 50}
	list = DeleteV1(list, 3)
	t.Logf("%v \n", list)
}

func TestDeleteV2(t *testing.T) {
	list := []int64{10, 20, 30, 40, 50}
	list = DeleteV2(list, 3)
	t.Logf("%v \n", list)
}
