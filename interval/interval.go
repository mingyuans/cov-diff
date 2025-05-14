package interval

import (
	"sort"
)

type Interval struct {
	Start int
	End   int
}

type FilesIntervals = map[string][]Interval

func joinSortedIntervals(a []Interval) []Interval {
	if len(a) == 0 {
		return a
	}

	result := []Interval{}
	last := a[0]
	for _, i := range a {
		if last.End >= i.Start && i.End >= last.End {
			last.End = i.End
			continue
		}
		result = append(result, last)
		last = i
	}

	return append(result, last)
}

func Sum(a []Interval) int {
	sum := 0
	for _, i := range a {
		sum += i.End - i.Start + 1
	}

	return sum
}

func JoinAndSortIntervals(a []Interval) []Interval {
	sort.Slice(a, func(i, j int) bool {
		return a[i].Start < a[j].Start
	})

	return joinSortedIntervals(a)
}

func Union(a []Interval, b []Interval) []Interval {
	if len(a) == 0 || len(b) == 0 {
		return []Interval{}
	}

	a = JoinAndSortIntervals(a)
	b = JoinAndSortIntervals(b)

	result := []Interval{}
	i := 0
	j := 0
	for i < len(a) && j < len(b) {
		if a[i].End < b[j].Start {
			i++
			continue
		}
		if a[i].Start > b[j].End {
			j++
			continue
		}
		start := a[i].Start
		if b[j].Start > start {
			start = b[j].Start
		}
		end := a[i].End
		if b[j].End < end {
			end = b[j].End
		}
		if a[i].End > b[j].End {
			j++
		} else {
			i++
		}

		result = append(result, Interval{
			Start: start,
			End:   end,
		})
	}

	return joinSortedIntervals(result)
}

func Diff(a []Interval, b []Interval) []Interval {
	if len(a) == 0 {
		return []Interval{}
	}
	if len(b) == 0 {
		return JoinAndSortIntervals(a)
	}

	a = JoinAndSortIntervals(a)
	b = JoinAndSortIntervals(b)

	result := []Interval{}
	i, j := 0, 0

	for i < len(a) {
		// If we've gone through all intervals in b or current a is before all remaining b
		if j >= len(b) || a[i].End < b[j].Start {
			// Add the entire interval from a
			result = append(result, a[i])
			i++
			continue
		}

		// If current a is after current b, move to next b
		if a[i].Start > b[j].End {
			j++
			continue
		}

		// Handle the part before overlap (if any)
		if a[i].Start < b[j].Start {
			result = append(result, Interval{
				Start: a[i].Start,
				End:   b[j].Start - 1,
			})
		}

		// Advance interval pointers based on which ends first
		if a[i].End > b[j].End {
			// Create new interval starting after current b
			a[i].Start = b[j].End + 1
			j++
		} else {
			i++
		}
	}

	return joinSortedIntervals(result)
}
