package db

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

type Task struct {
	Key   int
	Value string
}

func Create(dbPath string, bucketName []byte) error {
	var err error
	db, err = Init(dbPath)
	if err != nil {
		log.Println("Bucket failed: ", string(bucketName))
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
}

func Init(dbPath string) (*bolt.DB, error) {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Println("Error opening the bucket: \n", err)
	}

	return db, err
}

func CreateTask(task string, bucketName []byte) (int, error) {
	var id int
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		id64, _ := b.NextSequence()
		id = int(id64)
		key := itob(id)
		return b.Put(key, []byte(task))
	})

	if err != nil {
		return -1, err
	}

	return id, nil
}

func AllTasks(bucketName []byte) ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, Task{
				Key:   btoi(k),
				Value: string(v),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func DeleteTask(key int, bucketName []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		return b.Delete(itob(key))
	})
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
