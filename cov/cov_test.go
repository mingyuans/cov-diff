package cov

import (
	"github.com/panagiotisptr/cov-diff/interval"
	"testing"
)

func TestFilterBlocksBySearchingRange(t *testing.T) {
	rangeInterval := []interval.Interval{
		{
			Start: 1,
			End:   50,
		},
	}

	t.Run("base case", func(t *testing.T) {
		blocks := []CoverageBlock{
			{
				Block: interval.Interval{Start: 1, End: 10},
			},
			{
				Block: interval.Interval{Start: 11, End: 20},
			},
			{
				Block: interval.Interval{Start: 21, End: 30},
			},
			{
				Block: interval.Interval{Start: 50, End: 70},
			},
			{
				Block: interval.Interval{Start: 71, End: 80},
			},
		}

		filterBlocks := FilterBlocksBySearchingRange(rangeInterval, blocks)
		if len(filterBlocks) != 3 {
			t.Fatalf("expected 4 blocks, got %d", len(filterBlocks))
		}
	})

}
