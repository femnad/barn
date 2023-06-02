package selection

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"path"

	bolt "go.etcd.io/bbolt"

	"github.com/femnad/barn/entity"
	"github.com/femnad/mare"
)

const (
	dbMode        = 0600
	defaultDbPath = "~/.config/barn/barn.boltdb"
)

type selectionMap map[string]entity.Entry
type bucketEntries map[string]selectionMap

func getDb(cfg entity.Config) (*bolt.DB, error) {
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

func decodeEntry(b []byte) (entity.Entry, error) {
	buffer := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buffer)

	e := entity.Entry{}
	err := dec.Decode(&e)
	if err != nil {
		return e, err
	}

	return e, err
}

func encodeEntry(entry entity.Entry) ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(entry)
	return buffer.Bytes(), err
}

func dbIncrementEntryCount(bucket, key string, tx *bolt.Tx) (entity.Entry, error) {
	var err error
	var bck *bolt.Bucket
	var value entity.Entry
	bucketName := []byte(bucket)
	keyName := []byte(key)

	bck = tx.Bucket(bucketName)
	if bck == nil {
		return value, fmt.Errorf("expected bucket %s to exists when incrementing entry for %s", bucket, key)
	}

	storedEntry := bck.Get(keyName)
	if storedEntry == nil {
		return value, fmt.Errorf("expected entry in bucket %s for %s to exists when incrementing", bucket, key)
	}
	value, err = decodeEntry(storedEntry)
	if err != nil {
		return value, err
	}
	value.Count = value.Count + 1

	encoded, err := encodeEntry(value)
	if err != nil {
		return value, err
	}

	return value, bck.Put(keyName, encoded)
}

func incrementEntryCount(cfg entity.Config, id, key string) (entity.Entry, error) {
	var value entity.Entry
	db, err := getDb(cfg)
	if err != nil {
		return value, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		value, err = dbIncrementEntryCount(id, key, tx)
		return err
	})
	return value, err
}

func dbEnsureEntry(bucket, key string, entry entity.Entry, tx *bolt.Tx) (entity.Entry, error) {
	var err error
	var bck *bolt.Bucket
	var value entity.Entry
	bucketName := []byte(bucket)
	keyName := []byte(key)

	bck = tx.Bucket(bucketName)
	if bck == nil {
		bck, err = tx.CreateBucket(bucketName)
		if err != nil {
			return value, err
		}
	}

	storedEntry := bck.Get(keyName)
	if storedEntry == nil {
		encoded, eErr := encodeEntry(entry)
		if eErr != nil {
			return value, eErr
		}

		eErr = bck.Put(keyName, encoded)
		if eErr != nil {
			return value, eErr
		}

		return entry, nil
	}

	return decodeEntry(storedEntry)
}

func getSelectionMap(cfg entity.Config, id string, entries []entity.Entry) (selectionMap, error) {
	countMap := make(selectionMap)

	db, err := getDb(cfg)
	if err != nil {
		return countMap, err
	}

	err = db.Batch(func(tx *bolt.Tx) error {
		for _, entry := range entries {
			value, eErr := dbEnsureEntry(id, entry.DisplayName, entry, tx)
			if eErr != nil {
				return eErr
			}
			countMap[value.DisplayName] = value
		}
		return nil
	})

	return countMap, err
}

func getStoredSelections(cfg entity.Config, id string) (bucketEntries, error) {
	bucketMap := make(bucketEntries)

	db, err := getDb(cfg)
	if err != nil {
		return bucketMap, err
	}

	err = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			bucketName := string(name)
			if id != "" && bucketName != id {
				return nil
			}

			bucketSelections := make(selectionMap)

			err = b.ForEach(func(k, v []byte) error {
				entry, dErr := decodeEntry(v)
				if dErr != nil {
					return dErr
				}
				bucketSelections[string(k)] = entry
				return nil
			})
			if err != nil {
				return err
			}

			bucketMap[bucketName] = bucketSelections
			return nil
		})
	})

	return bucketMap, err
}

func purgeBucket(cfg entity.Config, bucket string) error {
	db, err := getDb(cfg)
	if err != nil {
		return err
	}

	bucketName := []byte(bucket)

	return db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(bucketName)
	})
}

func getBuckets(cfg entity.Config) ([]string, error) {
	var buckets []string

	db, err := getDb(cfg)
	if err != nil {
		return buckets, err
	}

	err = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			buckets = append(buckets, string(name))
			return nil
		})
	})
	if err != nil {
		return buckets, err
	}

	return buckets, nil
}
