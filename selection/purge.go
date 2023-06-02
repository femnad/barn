package selection

func Purge(configFile, id string) error {
	cfg, err := getConfig(configFile)
	if err != nil {
		return err
	}

	return purgeBucket(cfg, id)
}
