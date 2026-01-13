package operations

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"

	"github.com/github/orchestrator/backend/internal/events"
	"github.com/github/orchestrator/backend/internal/storage"
)

type Executor struct {
	store  storage.Store
	hub    *events.Hub
	logger *zap.Logger
}

type Runner func(ctx context.Context) (any, error)

func NewExecutor(store storage.Store, hub *events.Hub, logger *zap.Logger) *Executor {
	return &Executor{store: store, hub: hub, logger: logger}
}

func (e *Executor) Run(ctx context.Context, operationType string, commandText string, runner Runner) (storage.OperationRecord, error) {
	record := storage.OperationRecord{
		OperationType: operationType,
		Status:        storage.OperationRunning,
		CommandText:   commandText,
		Progress:      0,
	}

	created, err := e.store.CreateOperation(ctx, record)
	if err != nil {
		return storage.OperationRecord{}, err
	}

	e.hub.Broadcast(events.Event{
		Type:      "operation_started",
		Message:   operationType,
		Timestamp: time.Now().UTC(),
	})

	result, runErr := runner(ctx)
	completedAt := time.Now().Unix()
	created.CompletedAtUnix = &completedAt
	created.Progress = 100
	if runErr != nil {
		created.Status = storage.OperationFailed
		created.ErrorText = runErr.Error()
		e.logger.Warn("operation failed", zap.String("operation", operationType), zap.Error(runErr))
	} else {
		created.Status = storage.OperationCompleted
		payload, _ := json.Marshal(result)
		created.ResultJSON = payload
	}

	if err := e.store.UpdateOperation(ctx, created); err != nil {
		return created, err
	}

	e.hub.Broadcast(events.Event{
		Type:      "operation_completed",
		Message:   operationType,
		Timestamp: time.Now().UTC(),
	})

	return created, runErr
}
