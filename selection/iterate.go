package selection

import (
	"fmt"
)

func showEntry(p pair) {
	k, v := p.key, p.value
	fmt.Printf("Key: %s => Display Name: %s, Full Name: %s, Count: %d\n", k, v.DisplayName, v.FullName, v.Count)
}

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
		var printedHeader bool
		sorted := sortEntries(entryMap, true)
		for _, entryPair := range sorted {
			if !showZeroCounts && entryPair.value.Count == 0 {
				continue
			}

			if !printedHeader {
				fmt.Printf("Bucket: %s\n", bucket)
				printedHeader = true
			}

			showEntry(entryPair)
		}

		if printedHeader && bucketIdx < numBuckets-1 {
			fmt.Println()
		}
		bucketIdx++
	}

	return nil
}
