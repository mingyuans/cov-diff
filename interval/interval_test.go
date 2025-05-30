package interval

import (
	"testing"
)

func TestUnion(t *testing.T) {
	type testcase struct {
		name string
		a    []Interval
		b    []Interval
		e    []Interval
	}

	testacses := []testcase{
		{
			name: "base case",
			a: []Interval{
				{Start: 0, End: 5},
				{Start: 7, End: 10},
				{Start: 12, End: 15},
			},
			b: []Interval{
				{Start: 0, End: 6},
				{Start: 6, End: 9},
				{Start: 11, End: 13},
			},
			e: []Interval{
				{Start: 0, End: 5},
				{Start: 7, End: 9},
				{Start: 12, End: 13},
			},
		},
		{
			name: "overlaps",
			a: []Interval{
				{Start: 0, End: 5},
				{Start: 5, End: 15},
				{Start: 12, End: 15},
				{Start: 18, End: 25},
			},
			b: []Interval{
				{Start: 0, End: 2},
				{Start: 0, End: 3},
				{Start: 4, End: 5},
				{Start: 12, End: 15},
				{Start: 17, End: 22},
				{Start: 23, End: 28},
			},
			e: []Interval{
				{Start: 0, End: 3},
				{Start: 4, End: 5},
				{Start: 12, End: 15},
				{Start: 18, End: 22},
				{Start: 23, End: 25},
			},
		},
		{
			name: "multi-splits",
			a: []Interval{
				{Start: 1, End: 2},
				{Start: 3, End: 4},
				{Start: 5, End: 6},
			},
			b: []Interval{
				{Start: 0, End: 1},
				{Start: 2, End: 3},
				{Start: 4, End: 5},
			},
			e: []Interval{
				{Start: 1, End: 1},
				{Start: 2, End: 2},
				{Start: 3, End: 3},
				{Start: 4, End: 4},
				{Start: 5, End: 5},
			},
		},
	}

	for _, tc := range testacses {
		t.Run(tc.name, func(t *testing.T) {
			r := Union(tc.a, tc.b)
			if len(r) != len(tc.e) {
				t.Fatalf("expected %d intervals for %d", len(tc.e), len(r))
			}
			for i := range r {
				if r[i].Start != tc.e[i].Start {
					t.Fatalf(
						"expected interval to start at %d starts at %d instead",
						tc.e[i].Start,
						r[i].Start,
					)
				}
				if r[i].End != tc.e[i].End {
					t.Fatalf(
						"expected interval to end at %d ends at %d instead",
						tc.e[i].End,
						r[i].End,
					)
				}
			}
		})
	}
}
