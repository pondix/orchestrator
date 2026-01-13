package storage

import (
	"context"
	"sync"
	"time"
)

type MemoryStore struct {
	mu         sync.Mutex
	nextID     int64
	operations map[int64]OperationRecord
	logs       map[int64][]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		operations: make(map[int64]OperationRecord),
		logs:       make(map[int64][]string),
		nextID:     1,
	}
}

func (s *MemoryStore) CreateOperation(_ context.Context, record OperationRecord) (OperationRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	record.ID = s.nextID
	s.nextID++
	now := time.Now().Unix()
	record.StartedAtUnix = &now
	s.operations[record.ID] = record
	return record, nil
}

func (s *MemoryStore) UpdateOperation(_ context.Context, record OperationRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.operations[record.ID] = record
	return nil
}

func (s *MemoryStore) AppendOperationLog(_ context.Context, operationID int64, line string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logs[operationID] = append(s.logs[operationID], line)
	return nil
}
