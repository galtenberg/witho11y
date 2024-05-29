package witho11y

import (
  "context"

  "github.com/stretchr/testify/mock"
)

type MockBusinessLogic struct {
  mock.Mock
}

func (m *MockBusinessLogic) Execute(ctx context.Context, params ...any) ([]any, error) {
  args := m.Called(ctx, params)
  return args.Get(0).([]any), args.Error(1)
}
