package mock

import (
	"context"
	"errors"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

type MockTxManager struct {
	ShouldFail bool
	FailReason string
}

// NewMockTxManager creates a new MockTxManager that succeeds
func NewMockTxManager() *MockTxManager {
	return &MockTxManager{
		ShouldFail: false,
		FailReason: "",
	}
}

// NewMockTxManagerWithError creates a new MockTxManager that fails with the given reason
func NewMockTxManagerWithError(reason string) *MockTxManager {
	return &MockTxManager{
		ShouldFail: true,
		FailReason: reason,
	}
}

func (m *MockTxManager) Do(ctx context.Context, fn func(context.Context) error) error {
	if m.ShouldFail {
		return errors.New(m.FailReason)
	}
	return fn(ctx)
}

func (m *MockTxManager) DoWithSettings(ctx context.Context, settings trm.Settings, fn func(context.Context) error) error {
	if m.ShouldFail {
		return errors.New(m.FailReason)
	}
	return fn(ctx)
}
