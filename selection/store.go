package selection

import (
	"fmt"
	"path"
	"strconv"

	bolt "go.etcd.io/bbolt"

	"github.com/femnad/barn/config"
	"github.com/femnad/mare"
)

const (
	countBase     = 10
	countBit      = 64
	dbMode        = 0600
	defaultDbPath = "~/.config/barn/barn.boltdb"
)

type selectionMap map[string]int64

func getDb(cfg config.Config) (*bolt.DB, error) {
	dbPath := cfg.Options.DatabasePath
	if dbPath == "" {
		dbPath = defaultDbPath
	}
	dbPath = mare.ExpandUser(dbPath)
	dbDir, _ := path.Split(dbPath)
	err := mare.EnsureDir(dbDir)
	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(dbPath, dbMode, nil)
	if err != nil {
		return nil, fmt.Errorf("error opening database at path %s: %v", dbPath, err)
	}

	return db, nil
}

func storeSelection(cfg config.Config, id, selection string) error {
	db, err := getDb(cfg)
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		var uErr error
		var count int64
		var bucket *bolt.Bucket
		bucketName := []byte(id)
		key := []byte(selection)

		bucket = tx.Bucket(bucketName)
		if bucket == nil {
			bucket, err = tx.CreateBucket(bucketName)
			if err != nil {
				return uErr
			}
		}

		countByte := bucket.Get(key)
		if countByte == nil {
			count = 0
		} else {
			count, err = strconv.ParseInt(string(countByte), countBase, countBit)
			if err != nil {
				return err
			}
		}

		countNew := strconv.FormatInt(count+1, countBase)
		return bucket.Put(key, []byte(countNew))
	})
	return err
}

func getSelectionMap(cfg config.Config, id string) (selectionMap, error) {
	countMap := make(selectionMap)
	bucketName := []byte(id)

	db, err := getDb(cfg)
	if err != nil {
		return countMap, err
	}

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return nil
		}

		err = bucket.ForEach(func(k, v []byte) error {
			key := string(k)
			value, cErr := strconv.ParseInt(string(v), countBase, countBit)
			if cErr != nil {
				return cErr
			}
			countMap[key] = value
			return nil
		})
		return nil
	})

	return countMap, err
}
