package util

import "testing"

func TestSlicesEqual(t *testing.T) {
	for _, tt := range []struct {
		in0, in1 []int
		eq       bool
	}{
		{[]int{}, []int{}, true},
		{[]int{123}, []int{}, false},
		{[]int{}, []int{321}, false},
		{[]int{1}, []int{1}, true},
		{[]int{1, 6, 1, 7, 0}, []int{1, 6, 1, 7, 0}, true},
		{[]int{1, 6, 1, 7}, []int{1, 6, 1, 7, 0}, false},
		{[]int{1, 6, 1, 7, 0}, []int{1, 6, 1, 7, 1}, false},
	} {
		if got := SlicesEqual(tt.in0, tt.in1); got != tt.eq {
			t.Errorf("want %v, got %v", tt.eq, got)
		}
	}
}
