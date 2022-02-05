package engine

import (
	"encoding/json"
	"igabir98/simpleTODO/models"
	"time"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	tasksBucket = "tasks"
)

type BoltDB struct {
	db *bolt.DB
}

func NewBoltDB() (*BoltDB, error) {
	result := BoltDB{db: &bolt.DB{}}
	db, err := bolt.Open("./my.db", 0600, &bolt.Options{})

	if err != nil {
		return nil, err
	}

	topBuckets := []string{tasksBucket}

	err = db.Update(func(tx *bolt.Tx) error {
		for _, bktName := range topBuckets {
			if _, e := tx.CreateBucketIfNotExists([]byte(bktName)); e != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	result.db = db

	return &result, nil
}

func (b *BoltDB) CreateTask(task *models.Task) (taskID string, err error) {
	task = b.prepareNewTask(task)
	taskBytes, err := json.Marshal(task)

	if err != nil {
		return "", err
	}

	err = b.db.Update(func(tx *bolt.Tx) (err error) {
		bkt := tx.Bucket([]byte(tasksBucket))

		err = bkt.Put([]byte(task.ID), taskBytes)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return task.ID, nil
}

func (s *BoltDB) prepareNewTask(task *models.Task) *models.Task {
	if task.ID == "" {
		task.ID = uuid.New().String()
	}

	if task.Timestamp.IsZero() {
		task.Timestamp = time.Now()
	}

	return task
}

func (b *BoltDB) Get(taskID string) (task models.Task, err error) {

	err = b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(tasksBucket))

		value := bkt.Get([]byte(taskID))
		if value == nil {
			return errors.Errorf("no value for %s", taskID)
		}

		if err := json.Unmarshal(value, &task); err != nil {
			return errors.Wrap(err, "failed to unmarshal")
		}
		return nil
	})
	return task, err
}

func (b *BoltDB) GetAll() (task *[]models.Task, err error) {
	var tasks []models.Task

	err = b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(tasksBucket))

		bkt.ForEach(func(k, v []byte) error {
			var task models.Task

			if err := json.Unmarshal(v, &task); err != nil {
				return errors.Wrap(err, "failed to unmarshal")
			}

			tasks = append(tasks, task)

			return nil
		})

		return nil

	})

	return &tasks, nil
}
