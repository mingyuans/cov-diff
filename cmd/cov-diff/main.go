package main

import (
	"flag"
	"fmt"
	"github.com/panagiotisptr/cov-diff/interval"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/actions-go/toolkit/core"
	"github.com/panagiotisptr/cov-diff/cov"
	"github.com/panagiotisptr/cov-diff/diff"
	"github.com/panagiotisptr/cov-diff/files"
)

var path = flag.String("path", "", "path to the git repository")
var coverageFile = flag.String("coverprofile", "", "location of the coverage file")
var diffFile = flag.String("diff", "", "location of the diff file")
var moduleName = flag.String("module", "", "the name of module")
var ignoreMain = flag.String("ignore-main", "", "ignore main package")

func emptyValAndActionInputSet(val string, input string) bool {
	return val == "" && os.Getenv(
		fmt.Sprintf("INPUT_%s", strings.ToUpper(input)),
	) != ""
}

func getActionInput(input string) string {
	return os.Getenv(
		fmt.Sprintf("INPUT_%s", strings.ToUpper(input)),
	)
}

func populateFlagsFromActionEnvs() {
	if emptyValAndActionInputSet(*path, "path") {
		*path = getActionInput("path")
	}
	if emptyValAndActionInputSet(*coverageFile, "coverprofile") {
		*coverageFile = getActionInput("coverprofile")
	}
	if emptyValAndActionInputSet(*diffFile, "diff") {
		*diffFile = getActionInput("diff")
	}
	if emptyValAndActionInputSet(*moduleName, "module") {
		*moduleName = getActionInput("module")
	}
	if emptyValAndActionInputSet(*ignoreMain, "ignore-main") {
		*ignoreMain = getActionInput("ignore-main")
	}
}

func main() {
	flag.Parse()
	populateFlagsFromActionEnvs()

	if *coverageFile == "" {
		log.Fatal("missing coverage file")
	}

	diffIntervals, err := diff.GetFilesIntervalsFromDiffFile(*diffFile)
	if err != nil {
		log.Fatal(err)
	}

	coverageBlocks, err := cov.GetFilesIntervalsFromCoverageFile(*coverageFile)
	if err != nil {
		log.Fatal(err)
	}

	var totalCovBlocks, coveredBlocks int
	for filename, di := range diffIntervals {
		fmt.Printf("Processing file: %s\n", filename)
		// intervals which functions are in the file
		fi, getFuncIntervalsErr := files.GetFuncIntervalsFromFilePath(filepath.Join(*path, filename), *ignoreMain == "true")
		if getFuncIntervalsErr != nil {
			log.Fatal(getFuncIntervalsErr)
		}

		fullFilename := filepath.Join(*moduleName, filename)

		// intervals that changed and are parts of the code we care about
		measuredIntervals := interval.Union(di, fi)
		si, ok := coverageBlocks[fullFilename]
		if !ok {
			continue
		}

		covBlocks := cov.FilterBlocksBySearchingRange(measuredIntervals, si)
		fmt.Printf("Total coverage blocks: %d\n", len(covBlocks))
		for _, cb := range covBlocks {
			fmt.Printf("Block: %s, Start: %d, End: %d, Count: %d\n", cb.FileName, cb.Block.Start, cb.Block.End, cb.ExecutionCount)
		}

		totalCovBlocks += len(covBlocks)

		for _, cb := range covBlocks {
			if cb.ExecutionCount > 0 {
				coveredBlocks++
			}
		}
	}
	var percentCoverage = 100
	if totalCovBlocks > 0 {
		percentCoverage = coveredBlocks * 100 / totalCovBlocks
	}

	fmt.Printf("Coverage on new lines: %d%%\n", percentCoverage)
	if getActionInput("coverprofile") != "" {
		core.SetOutput("covdiff", fmt.Sprintf("%d", percentCoverage))
	}
}
