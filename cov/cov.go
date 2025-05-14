package cov

import (
	"bytes"
	"os"

	"github.com/panagiotisptr/cov-diff/files"
	"github.com/panagiotisptr/cov-diff/interval"
	"golang.org/x/tools/cover"
)

func GetFilesIntervalsFromCoverageFile(coverageFilePath string) (map[string][]CoverageBlock, error) {
	covFileBytes, err := os.ReadFile(coverageFilePath)
	if err != nil {
		return nil, err
	}

	coverageBlocks, err := getFilesIntervalsFromCoverage(covFileBytes)
	if err != nil {
		return nil, err
	}
	// de-allocate covFileBytes
	covFileBytes = nil
	return coverageBlocks, nil
}

// GetFilesIntervalsFromCoverage processes coverage data and extracts coverage blocks for each file.
func getFilesIntervalsFromCoverage(covBytes []byte) (map[string][]CoverageBlock, error) {
	blocks := map[string][]CoverageBlock{}
	cps, err := cover.ParseProfilesFromReader(bytes.NewReader(covBytes))
	if err != nil {
		return blocks, err
	}

	for _, cp := range cps {
		if files.ShouldSkipFile(cp.FileName) {
			continue
		}

		if _, ok := blocks[cp.FileName]; !ok {
			blocks[cp.FileName] = []CoverageBlock{}
		}

		for _, b := range cp.Blocks {
			block := CoverageBlock{
				FileName:       cp.FileName,
				Block:          interval.Interval{Start: b.StartLine, End: b.EndLine},
				StatementCount: b.NumStmt,
				ExecutionCount: b.Count,
			}
			blocks[cp.FileName] = append(blocks[cp.FileName], block)
		}
	}

	return blocks, nil
}

type CoverageBlock struct {
	FileName       string
	Block          interval.Interval
	StatementCount int
	ExecutionCount int
}

func FilterBlocksBySearchingRange(rangeIntervals []interval.Interval, blocks []CoverageBlock) []CoverageBlock {
	var filteredBlocks []CoverageBlock
	for _, block := range blocks {
		for _, r := range rangeIntervals {
			if block.Block.Start >= r.Start && block.Block.End <= r.End {
				filteredBlocks = append(filteredBlocks, block)
			}
		}
	}
	return filteredBlocks
}
