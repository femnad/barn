package selection

func Purge(configFile string, buckets []string) error {
	cfg, err := getConfig(configFile)
	if err != nil {
		return err
	}

	return purgeBucket(cfg, buckets)
}
