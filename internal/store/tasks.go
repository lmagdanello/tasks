package store

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/lmagdanello/tasks/internal/config"
	"go.etcd.io/bbolt"
)

var (
	ErrBucketNotFound = errors.New("bucket does not exist")
	ErrBucketExists   = errors.New("bucket already exists")
	ErrTaskNotFound   = errors.New("task not found")
)

const dbFileMode = 0o600
const dbTimeout = time.Second

// Priority level for a task.
type Priority string

const (
	PriorityHigh   Priority = "alta"
	PriorityMedium Priority = "media"
	PriorityLow    Priority = "baixa"
)

// Task status.
type Status string

const (
	StatusActive Status = "active"
	StatusDone   Status = "done"
)

// Task represents a task in the database.
type Task struct {
	Key      int
	Text     string
	Priority Priority
	Due      *time.Time
	Status   Status
}

// taskData is the JSON format stored in the database.
type taskData struct {
	Text     string     `json:"text"`
	Priority Priority   `json:"priority"`
	Due      *time.Time `json:"due,omitempty"`
	Status   Status     `json:"status"`
}

// Store wraps the database connection and task operations.
type Store struct {
	db     *bbolt.DB
	bucket string
	bKey   []byte
}

// Open opens the database for the given bucket.
func Open(bucket string) (*Store, error) {
	exists, err := config.BucketExists(bucket)
	if err != nil {
		return nil, fmt.Errorf("checking bucket: %w", err)
	}
	if !exists {
		return nil, ErrBucketNotFound
	}

	path, err := config.DBPath(bucket)
	if err != nil {
		return nil, fmt.Errorf("getting db path: %w", err)
	}

	db, err := bbolt.Open(path, dbFileMode, &bbolt.Options{Timeout: dbTimeout})
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	return &Store{db: db, bucket: bucket, bKey: []byte(bucket)}, nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// CreateBucket creates a new bucket (task list).
func CreateBucket(bucket string) error {
	exists, err := config.BucketExists(bucket)
	if err != nil {
		return fmt.Errorf("checking bucket: %w", err)
	}
	if exists {
		return ErrBucketExists
	}

	path, err := config.DBPath(bucket)
	if err != nil {
		return fmt.Errorf("getting db path: %w", err)
	}

	db, err := bbolt.Open(path, dbFileMode, &bbolt.Options{Timeout: dbTimeout})
	if err != nil {
		return fmt.Errorf("creating database: %w", err)
	}
	defer db.Close()

	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucket([]byte(bucket))
		return err
	})
}

func (s *Store) readTask(key int, data []byte) (Task, error) {
	var td taskData
	if err := json.Unmarshal(data, &td); err != nil {
		// Legacy: plain text
		return Task{
			Key:      key,
			Text:     string(data),
			Priority: PriorityMedium,
			Status:   StatusActive,
		}, nil
	}
	return Task{
		Key:      key,
		Text:     td.Text,
		Priority: td.Priority,
		Due:      td.Due,
		Status:   td.Status,
	}, nil
}

// CreateTask adds a new task with optional priority and due date.
func (s *Store) CreateTask(text string, opts ...CreateOption) (int, error) {
	opt := &createOpts{Priority: PriorityMedium}
	for _, o := range opts {
		o(opt)
	}

	td := taskData{
		Text:     text,
		Priority: opt.Priority,
		Due:      opt.Due,
		Status:   StatusActive,
	}
	data, err := json.Marshal(td)
	if err != nil {
		return -1, err
	}

	var id int
	err = s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(s.bKey)
		if b == nil {
			return fmt.Errorf("bucket %s not found", s.bucket)
		}
		id64, _ := b.NextSequence()
		id = int(id64)
		return b.Put(itob(id), data)
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

// CreateOption configures CreateTask.
type CreateOption func(*createOpts)

type createOpts struct {
	Priority Priority
	Due      *time.Time
}

func WithPriority(p Priority) CreateOption {
	return func(o *createOpts) { o.Priority = p }
}

func WithDue(t time.Time) CreateOption {
	return func(o *createOpts) { o.Due = &t }
}

// UpdateTask updates text, priority, or due date of an existing task.
func (s *Store) UpdateTask(key int, text string, opts ...UpdateOption) error {
	opt := &updateOpts{Text: text}
	for _, o := range opts {
		o(opt)
	}

	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(s.bKey)
		if b == nil {
			return fmt.Errorf("bucket %s not found", s.bucket)
		}
		v := b.Get(itob(key))
		if v == nil {
			return ErrTaskNotFound
		}
		task, err := s.readTask(key, v)
		if err != nil {
			return err
		}
		if opt.Text != "" {
			task.Text = opt.Text
		}
		if opt.Priority != nil {
			task.Priority = *opt.Priority
		}
		if opt.ClearDue {
			task.Due = nil
		} else if opt.Due != nil {
			task.Due = opt.Due
		}
		data, err := json.Marshal(taskData{
			Text:     task.Text,
			Priority: task.Priority,
			Due:      task.Due,
			Status:   task.Status,
		})
		if err != nil {
			return err
		}
		return b.Put(itob(key), data)
	})
}

type UpdateOption func(*updateOpts)

type updateOpts struct {
	Text     string
	Priority *Priority
	Due      *time.Time
	ClearDue bool
}

func UpdatePriority(p Priority) UpdateOption {
	return func(o *updateOpts) { o.Priority = &p }
}

func UpdateDue(t *time.Time) UpdateOption {
	return func(o *updateOpts) {
		if t == nil {
			o.ClearDue = true
		} else {
			o.Due = t
		}
	}
}

// ArchiveTask marks a task as done without deleting it.
func (s *Store) ArchiveTask(key int) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(s.bKey)
		if b == nil {
			return fmt.Errorf("bucket %s not found", s.bucket)
		}
		v := b.Get(itob(key))
		if v == nil {
			return ErrTaskNotFound
		}
		task, err := s.readTask(key, v)
		if err != nil {
			return err
		}
		task.Status = StatusDone
		data, err := json.Marshal(taskData{
			Text:     task.Text,
			Priority: task.Priority,
			Due:      task.Due,
			Status:   task.Status,
		})
		if err != nil {
			return err
		}
		return b.Put(itob(key), data)
	})
}

// DeleteTask permanently removes a task.
func (s *Store) DeleteTask(key int) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(s.bKey)
		if b == nil {
			return fmt.Errorf("bucket %s not found", s.bucket)
		}
		return b.Delete(itob(key))
	})
}

// AllTasks returns active tasks, optionally filtered.
func (s *Store) AllTasks(opts ...ListOption) ([]Task, error) {
	return s.listTasks(StatusActive, opts...)
}

// DoneTasks returns archived (status=done) tasks.
func (s *Store) DoneTasks(opts ...ListOption) ([]Task, error) {
	return s.listTasks(StatusDone, opts...)
}

type ListOption func(*listOpts)

type listOpts struct {
	Query    string
	Priority *Priority
}

func FilterQuery(q string) ListOption {
	return func(o *listOpts) { o.Query = q }
}

func FilterPriority(p Priority) ListOption {
	return func(o *listOpts) { o.Priority = &p }
}

func (s *Store) listTasks(status Status, opts ...ListOption) ([]Task, error) {
	opt := &listOpts{}
	for _, o := range opts {
		o(opt)
	}

	queryLower := strings.ToLower(opt.Query)
	var tasks []Task

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(s.bKey)
		if b == nil {
			return fmt.Errorf("bucket %s not found", s.bucket)
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			task, err := s.readTask(btoi(k), v)
			if err != nil {
				continue
			}
			if task.Status != status {
				continue
			}
			if opt.Query != "" && !strings.Contains(strings.ToLower(task.Text), queryLower) {
				continue
			}
			if opt.Priority != nil && task.Priority != *opt.Priority {
				continue
			}
			tasks = append(tasks, task)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sortTasks(tasks)
	return tasks, nil
}

// GetTask returns a task by its internal key.
func (s *Store) GetTask(key int) (Task, error) {
	var task Task
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(s.bKey)
		if b == nil {
			return fmt.Errorf("bucket %s not found", s.bucket)
		}
		v := b.Get(itob(key))
		if v == nil {
			return ErrTaskNotFound
		}
		var err error
		task, err = s.readTask(key, v)
		return err
	})
	return task, err
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}

func priorityOrder(p Priority) int {
	switch p {
	case PriorityHigh:
		return 0
	case PriorityMedium:
		return 1
	case PriorityLow:
		return 2
	}
	return 1
}

func sortTasks(tasks []Task) {
	sort.Slice(tasks, func(i, j int) bool {
		pi, pj := priorityOrder(tasks[i].Priority), priorityOrder(tasks[j].Priority)
		if pi != pj {
			return pi < pj
		}
		di, dj := tasks[i].Due, tasks[j].Due
		if di == nil && dj == nil {
			return tasks[i].Key < tasks[j].Key
		}
		if di == nil {
			return false
		}
		if dj == nil {
			return true
		}
		if !di.Equal(*dj) {
			return di.Before(*dj)
		}
		return tasks[i].Key < tasks[j].Key
	})
}
