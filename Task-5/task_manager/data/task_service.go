package data

import (
	"errors"
	"sync"
	"sync/atomic"
	"task_manager/models"
)

// TaskService provides an in-memory store and business logic for tasks.
type TaskService struct {
	mu    sync.RWMutex
	store map[int]models.Task
	idSeq int32
}

var svc *TaskService

func init() {
	svc = &TaskService{
		store: make(map[int]models.Task),
		idSeq: 0,
	}
}

// GetService returns the singleton task service
func GetService() *TaskService {
	return svc
}

// GetAll returns all tasks
func (s *TaskService) GetAll() []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]models.Task, 0, len(s.store))
	for _, t := range s.store {
		tasks = append(tasks, t)
	}
	return tasks
}

// GetByID returns a task by id or error if not found
func (s *TaskService) GetByID(id int) (models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.store[id]
	if !ok {
		return models.Task{}, errors.New("task not found")
	}
	return t, nil
}

// Create adds a new task and returns it with assigned ID
func (s *TaskService) Create(t models.Task) models.Task {
	next := int(atomic.AddInt32(&s.idSeq, 1))
	t.ID = next

	s.mu.Lock()
	s.store[next] = t
	s.mu.Unlock() 
	
	return t
}

// Update updates an existing task by id and returns updated task or error
func (s *TaskService) Update(id int, updated models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.store[id]
	if !ok {
		return models.Task{}, errors.New("task not found")
	}

	updated.ID = id
	s.store[id] = updated
	return updated, nil
}

// Delete removes a task by id
func (s *TaskService) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock() 

	if _, ok := s.store[id]; !ok {
		return errors.New("task not found")
	}
	delete(s.store, id)
	return nil
}
