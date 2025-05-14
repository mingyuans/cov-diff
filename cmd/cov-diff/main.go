package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/actions-go/toolkit/core"
	"github.com/panagiotisptr/cov-diff/cov"
	"github.com/panagiotisptr/cov-diff/diff"
	"github.com/panagiotisptr/cov-diff/files"
	"github.com/panagiotisptr/cov-diff/interval"
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

	diffBytes, err := os.ReadFile(*diffFile)
	if err != nil {
		log.Fatal("failed to read diff file: ", err)
	}

	diffIntervals, err := diff.GetFilesIntervalsFromDiff(diffBytes)
	if err != nil {
		log.Fatal(err)
	}
	// de-allocate diffBytes
	diffBytes = nil

	covFileBytes, err := os.ReadFile(*coverageFile)
	if err != nil {
		log.Fatal(err)
	}

	coverIntervals, allStatementIntervals, err := cov.GetFilesIntervalsFromCoverage(covFileBytes)
	if err != nil {
		log.Fatal(err)
	}
	// de-allocate covFileBytes
	covFileBytes = nil

	total := 0
	covered := 0
	for filename, di := range diffIntervals {
		fmt.Printf("Processing file: %s\n", filename)
		fileBytes, err := os.ReadFile(filepath.Join(*path, filename))
		if err != nil {
			log.Fatal(err)
		}

		// intervals which functions are in the file
		fi, err := files.GetIntervalsFromFile(fileBytes, *ignoreMain == "true")
		if err != nil {
			log.Fatal(err)
		}

		fullFilename := filepath.Join(*moduleName, filename)

		// intervals that changed and are parts of the code we care about
		measuredIntervals := interval.Union(di, fi)
		si, ok := allStatementIntervals[fullFilename]
		if !ok {
			continue
		}
		measuredIntervals = interval.Union(measuredIntervals, si)
		total += interval.Sum(measuredIntervals)
		fmt.Printf("measured new lines - sum: %d\n", interval.Sum(measuredIntervals))
		for _, i := range measuredIntervals {
			fmt.Printf("new lines - interval: %d-%d, sum:%d\n", i.Start, i.End, interval.Sum([]interval.Interval{i}))
		}

		ci, ok := coverIntervals[fullFilename]
		if !ok {
			continue
		}

		coveredMeasuredIntervals := interval.Union(measuredIntervals, ci)
		covered += interval.Sum(coveredMeasuredIntervals)
		fmt.Printf("covered new lines - sum: %d\n", interval.Sum(coveredMeasuredIntervals))
		for _, i := range coveredMeasuredIntervals {
			fmt.Printf("new lines - interval: %d-%d\n", i.Start, i.End)
		}

		diffMeasuredIntervals := interval.Diff(measuredIntervals, coveredMeasuredIntervals)
		fmt.Printf("diff new lines - sum: %d\n", interval.Sum(diffMeasuredIntervals))
		for _, i := range diffMeasuredIntervals {
			fmt.Printf("new lines - interval: %d-%d\n", i.Start, i.End)
		}
	}

	percentCoverage := 100
	if total != 0 {
		percentCoverage = (100 * covered) / total
	}

	fmt.Printf("Coverage on new lines: %d%%\n", percentCoverage)
	if getActionInput("coverprofile") != "" {
		core.SetOutput("covdiff", fmt.Sprintf("%d", percentCoverage))
	}
}
