package selection

func Truncate(configFile, id string, keys []string) error {
	cfg, err := getConfig(configFile)
	if err != nil {
		return err
	}

	return truncateKeys(cfg, id, keys)
}
