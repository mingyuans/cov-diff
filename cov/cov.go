package cov

import (
	"bytes"

	"github.com/panagiotisptr/cov-diff/files"
	"github.com/panagiotisptr/cov-diff/interval"
	"golang.org/x/tools/cover"
)

// GetFilesIntervalsFromCoverage parses coverage data and extracts covered and uncovered intervals for files.
//
// Parameters:
//   - covBytes: A byte slice containing the coverage data in a format compatible with `cover.ParseProfilesFromReader`.
//
// Returns:
//   - coveredIntervals: A map where the keys are file names and the values are slices of intervals representing covered lines.
//   - allStatementIntervals: A map where the keys are file names and the values are slices of intervals representing all statement lines.
//   - error: An error if parsing the coverage data fails.
func GetFilesIntervalsFromCoverage(
	covBytes []byte,
) (interval.FilesIntervals, interval.FilesIntervals, error) {
	coveredIntervals := interval.FilesIntervals{}
	allStatementIntervals := interval.FilesIntervals{}

	cps, err := cover.ParseProfilesFromReader(bytes.NewReader(covBytes))
	if err != nil {
		return coveredIntervals, allStatementIntervals, err
	}

	for _, cp := range cps {
		if files.ShouldSkipFile(cp.FileName) {
			continue
		}

		if _, ok := coveredIntervals[cp.FileName]; !ok {
			coveredIntervals[cp.FileName] = []interval.Interval{}
		}

		if _, ok := allStatementIntervals[cp.FileName]; !ok {
			allStatementIntervals[cp.FileName] = []interval.Interval{}
		}

		for _, b := range cp.Blocks {
			allStatementIntervals[cp.FileName] = append(allStatementIntervals[cp.FileName], interval.Interval{
				Start: b.StartLine,
				End:   b.EndLine,
			})

			if b.Count == 0 {
				continue
			}

			coveredIntervals[cp.FileName] = append(coveredIntervals[cp.FileName], interval.Interval{
				Start: b.StartLine,
				End:   b.EndLine,
			})
		}
	}

	return coveredIntervals, allStatementIntervals, nil
}
