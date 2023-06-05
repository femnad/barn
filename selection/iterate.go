package selection

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

func Iterate(configFile, id string, showZeroCounts bool) error {
	cfg, err := getConfig(configFile)
	if err != nil {
		return err
	}

	storedSelections, err := getStoredSelections(cfg, id)
	if err != nil {
		return err
	}

	var bucketIdx int
	numBuckets := len(storedSelections)

	for bucket, entryMap := range storedSelections {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Display Name", "Full Name", "Count"})

		var printedHeader bool
		sorted := sortEntries(entryMap)
		for _, entry := range sorted {
			if !showZeroCounts && entry.Count == 0 {
				continue
			}

			if !printedHeader {
				fmt.Printf("Bucket: %s\n", bucket)
				printedHeader = true
			}

			count := strconv.FormatInt(entry.Count, 10)
			table.Append([]string{entry.DisplayName, entry.FullName, count})
		}

		if printedHeader {
			table.Render()
			if bucketIdx < numBuckets-1 {
				fmt.Println()
			}
		}
		bucketIdx++
	}

	return nil
}
