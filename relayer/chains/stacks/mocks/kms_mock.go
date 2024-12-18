package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockKMS struct {
	mock.Mock
}

func (m *MockKMS) Init(ctx context.Context) (*string, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		s := args.String(0)
		return &s, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockKMS) Encrypt(ctx context.Context, data []byte) ([]byte, error) {
	args := m.Called(ctx, data)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockKMS) Decrypt(ctx context.Context, data []byte) ([]byte, error) {
	args := m.Called(ctx, data)
	return args.Get(0).([]byte), args.Error(1)
}
