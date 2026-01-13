package storage

import "context"

type OperationStatus string

const (
	OperationPending   OperationStatus = "PENDING"
	OperationRunning   OperationStatus = "RUNNING"
	OperationCompleted OperationStatus = "COMPLETED"
	OperationFailed    OperationStatus = "FAILED"
	OperationCancelled OperationStatus = "CANCELLED"
)

type OperationRecord struct {
	ID              int64
	InstanceID      *int64
	OperationType   string
	Status          OperationStatus
	CommandText     string
	Progress        int
	ResultJSON      []byte
	ErrorText       string
	StartedAtUnix   *int64
	CompletedAtUnix *int64
}

type Store interface {
	CreateOperation(ctx context.Context, record OperationRecord) (OperationRecord, error)
	UpdateOperation(ctx context.Context, record OperationRecord) error
	AppendOperationLog(ctx context.Context, operationID int64, line string) error
}
