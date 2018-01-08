package main

import "github.com/boltdb/bolt"
import "fmt"

func createBucket(db *bolt.DB, name string) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return fmt.Errorf("Error creating bucket: %s", err)
		}
		return nil
	})
}

func updateDB(db *bolt.DB, msg []byte, date string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Readings"))
		err := b.Put([]byte(date), msg)
		return err
	})
}
