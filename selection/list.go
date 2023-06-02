package selection

import "fmt"

func ListBuckets(configFile string) error {
	cfg, err := getConfig(configFile)
	if err != nil {
		return err
	}

	buckets, err := getBuckets(cfg)
	if err != nil {
		return err
	}

	for _, bucket := range buckets {
		fmt.Println(bucket)
	}

	return nil
}
