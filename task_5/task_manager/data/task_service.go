package data

import (
	"sync"
	"errors"
	"task_manager/models"
)

var (
	ErrNotFound = errors.New("task not found")
)

type TaskService struct {
	mu     sync.Mutex
	tasks  map[int]models.Task
	nextID int
}

func NewTaskService() *TaskService {
	return &TaskService{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

func (s *TaskService) GetAll() []models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]models.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		out = append(out, t)
	}
	return out
}

func (s *TaskService) GetByID(id int) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return models.Task{}, ErrNotFound
	}
	return t, nil
}

func (s *TaskService) Create(t models.Task) models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	t.ID = s.nextID
	s.nextID++
	s.tasks[t.ID] = t
	return t
}

func (s *TaskService) Update(id int, t models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.tasks[id]
	if !ok {
		return models.Task{}, ErrNotFound
	}
	t.ID = id
	s.tasks[id] = t
	return t, nil
}

func (s *TaskService) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.tasks[id]
	if !ok {
		return ErrNotFound
	}
	delete(s.tasks, id)
	return nil
}

