---
name: mock-testing
description: Generate and use mocks for unit testing with the mockgen library
---

# Mock Testing with mockgen

## Purpose

Generate mock implementations of interfaces for isolation during unit testing without external dependencies.

## When to Use

- During REFACTOR phase of TDD (replace real deps with mocks)
- When testing UseCases in isolation
- When writing unit tests that depend on external APIs/DBs
- When testing error scenarios without real network calls

## Prerequisites

- Go 1.16+
- `mockgen` installed
- Test framework: `github.com/stretchr/testify`

## Installation

```bash
# Install mockgen
go install github.com/golang/mock/mockgen@latest

# Verify installation
mockgen -version
```

## Generating Mocks

### Generate Mocks for a File

Create a `mockgen.yaml` config in your project root or use flags.

```bash
# Option 1: Generate mocks for a specific interface
mockgen -source=internal/application/usecase/user_usecase.go -destination=mocks/user_usecase_mock.go -package=mocks

# Option 2: Generate mocks for all interfaces in a package
mockgen -source=./internal/infrastructure/persistence/... -destination=mocks/repository_mock.go -package=mocks

# Option 3: Using go:generate directive (Recommended)
# Add this comment to your interface file:
//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

# Then run:
go generate ./...
```

**Example `go:generate` usage**:

```go
// internal/application/usecase/user_usecase.go
package usecase

//go:generate mockgen -source=$GOFILE -destination=mocks/user_usecase_mock.go -package=mocks

type UserUseCase interface {
    CreateUser(ctx context.Context, email string, password string) (*entity.User, error)
    GetUser(ctx context.Context, id string) (*entity.User, error)
}
```

### Using Generated Mocks in Tests

```go
// Generated file (internal/application/usecase/mocks/user_usecase_mock.go)
package mocks

import (
    "context"
    "github.com/golang/mock/gomock"
)

// MockUserUseCase is a mock of UserUseCase interface
type MockUserUseCase struct {
    ctrl     *gomock.Controller
    recorder *MockUserUseCaseMockRecorder
}

// MockUserUseCaseMockRecorder is the mock recorder for MockUserUseCase
type MockUserUseCaseMockRecorder struct {
    mock *MockUserUseCase
}

// NewMockUserUseCase creates a new mock instance
func NewMockUserUseCase(ctrl *gomock.Controller) *MockUserUseCase {
    mock := &MockUserUseCase{ctrl: ctrl}
    mock.recorder = &MockUserUseCaseMockRecorder{mock}
    return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserUseCase) EXPECT() *MockUserUseCaseMockRecorder {
    return m.recorder
}

// CreateUser mocks base method
func (m *MockUserUseCase) CreateUser(ctx context.Context, email string, password string) (*entity.User, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "CreateUser", ctx, email, password)
    ret0, _ := ret[0].(*entity.User)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser
func (mr *MockUserUseCaseMockRecorder) CreateUser(ctx, email, password interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserUseCase)(nil).CreateUser), ctx, email, password)
}
```

**Using in Test**:

```go
func TestCreateUser(t *testing.T) {
    // Create mock controller
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // Create mock usecase from generated file
    mockUserUC := mocks.NewMockUserUseCase(ctrl)

    // Setup expectations
    user := &entity.User{ID: "123", Email: "test@test.com"}
    mockUserUC.EXPECT().
        CreateUser(gomock.Any(), "test@test.com", "password").
        Return(user, nil).
        Times(1)

    // Call the method
    result, err := mockUserUC.CreateUser(context.Background(), "test@test.com", "password")

    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, "test@test.com", result.Email)
}
```

## Common Mock Patterns

### 1. Basic Mocking

**Scenario**: Simple method call.

```go
// Setup mock controller
ctrl := gomock.NewController(t)
defer ctrl.Finish()

// Create mock repository
mockRepo := mocks.NewMockUserRepository(ctrl)

// ✅ Good - Verify expectations
mockRepo.EXPECT().FindById("123").Return(&user, nil)

// Use the mock
result, err := service.GetUser("123")
```

### 2. Partial Mocking

**Scenario**: Testing one method, ignoring others.

```go
// Mock only what you need
mockRepo.EXPECT().FindById("123").Return(&user, nil)

// Don't worry about other methods unless they're called
```

### 3. Returning Different Values

**Scenario**: Test success vs failure paths.

```go
// Test success path
mockRepo.EXPECT().FindByEmail("exists@example.com").Return(&user, nil)

// Test failure path
mockRepo.EXPECT().FindByEmail("notfound@example.com").Return(nil, repository.ErrNotFound)

// Usage in test
func TestGetUser_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockUserRepository(ctrl)

    mockRepo.EXPECT().FindByEmail("exists@example.com").Return(&user, nil)
    // ... test success ...
}

func TestGetUser_NotFound(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockUserRepository(ctrl)

    mockRepo.EXPECT().FindByEmail("notfound@example.com").Return(nil, repository.ErrNotFound)
    // ... test error handling ...
}
```

### 4. Mocking Interfaces in External Packages

**Scenario**: Mock a database driver or HTTP client.

```go
// ✅ Good - Use interface
type DB interface {
    Query(ctx context.Context, q string) (*sql.Rows, error)
}

// Create mock that satisfies interface
ctrl := gomock.NewController(t)
defer ctrl.Finish()
mockDB := mocks.NewMockDB(ctrl)

// Pass interface to function under test
processData(mockDB)
```

### 5. Mocking with Goroutines

**Scenario**: Testing concurrent operations.

```go
// Mock expects to be called multiple times
mockRepo.EXPECT().Save(gomock.Any()).Return(nil).Times(3)

// Use WaitGroup
var wg sync.WaitGroup
for i := 0; i < 3; i++ {
    wg.Add(1)
    go func() {
        repo.Save(data)
        defer wg.Done()
    }()
}
wg.Wait()
```

### 6. Verifying Call Counts

**Scenario**: Ensure methods are called correct number of times.

```go
// mockgen automatically verifies expectations when ctrl.Finish() is called
// Use .Times(n) to specify exact count
mockRepo.EXPECT().Save(gomock.Any()).Return(nil).Times(1)
mockRepo.EXPECT().Update(gomock.Any()).Return(nil).Times(0)
```

### 7. Using gomock.Matchers

**Scenario**: Match complex arguments.

```go
// Match any value
mockRepo.EXPECT().FindById(gomock.Any()).Return(&user, nil)

// Match exact value
mockRepo.EXPECT().FindById(gomock.Eq("123")).Return(&user, nil)

// Match using custom function
mockRepo.EXPECT().FindByEmail(gomock.HasPrefix("admin@")).Return(&adminUser, nil)

// Match length
mockRepo.EXPECT().GetUsers(gomock.Len(3)).Return(users, nil)
```

## Best Practices

### DO Mock

- **Interfaces**, not concrete structs.
- **Public methods** (services, repositories), not private internals.
- **Database queries** (slow, I/O).
- **External APIs** (network, third-party).

### DON'T Mock

- **Value objects** (DTOs, entities that just hold data).
- **Data transformation** (formatting, calculating).
- **Business logic** (rules, validations).

### Write Good Tests

1. **Setup expectations before calling** code.
2. **One expectation per test** (clear assertions).
3. **Use `gomock.Any()`** for arguments you don't care about.
4. **Use `gomock.Eq()`** for specific arguments.
5. **Always call `defer ctrl.Finish()`** to verify expectations.
6. **Use table-driven tests** for multiple scenarios.

## Integrating with go:generate

The best way to manage mocks is using `go:generate` directives:

```go
// internal/application/usecase/user_usecase.go
package usecase

//go:generate mockgen -source=$GOFILE -destination=mocks/user_usecase_mock.go -package=mocks

type UserUseCase interface {
    CreateUser(ctx context.Context, email string, password string) (*entity.User, error)
    GetUser(ctx context.Context, id string) (*entity.User, error)
}
```

Then generate all mocks at once:

```bash
go generate ./...
```

This keeps your mocks in sync with your interfaces.

## Troubleshooting

### "No matching method found"

Ensure the method name in `EXPECT().Method()` matches the interface method exactly.

### "Expected call doesn't match"

Check types. mockgen generates types based on the interface definition.

### "Uninitialized Mock"

Always use `gomock.NewController(t)` and `defer ctrl.Finish()`.

### "mockgen: no import path for package"

When using `-source` flag, ensure the package has a proper import path or use `-package` and `-imports` flags.

## Example: Testing a UseCase with mockgen

```go
// internal/application/usecase/user_usecase_test.go
package usecase_test

import (
    "context"
    "testing"
    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"

    "github.com/aiagent/internal/application/usecase"
    "github.com/aiagent/internal/application/usecase/mocks"
    "github.com/aiagent/internal/domain/entity"
)

func TestUserUseCase_CreateUser(t *testing.T) {
    // 1. Create mock controller
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // 2. Mock the dependencies (Repository)
    mockRepo := mocks.NewMockUserRepository(ctrl)

    // 3. Setup expectations
    user := &entity.User{ID: "123", Email: "test@example.com"}
    mockRepo.EXPECT().
        Save(gomock.Any(), gomock.Any()).
        Return(nil).
        Times(1)

    // 4. Create UseCase with mocked repo
    useCase := usecase.NewUserUseCase(mockRepo)

    // 5. Execute
    ctx := context.Background()
    err := useCase.CreateUser(ctx, "test@example.com", "password")

    // 6. Assertions
    assert.NoError(t, err)
}
```

## Table-Driven Tests with Mocks

```go
func TestUserUseCase_CreateUser(t *testing.T) {
    tests := []struct {
        name        string
        email       string
        password    string
        setupMock   func(*mocks.MockUserRepository)
        wantErr     bool
        errContains string
    }{
        {
            name:     "success",
            email:    "test@example.com",
            password: "password123",
            setupMock: func(m *mocks.MockUserRepository) {
                m.EXPECT().
                    Save(gomock.Any(), gomock.Any()).
                    Return(nil).
                    Times(1)
            },
            wantErr: false,
        },
        {
            name:     "invalid email",
            email:    "invalid-email",
            password: "password123",
            setupMock: func(m *mocks.MockUserRepository) {
                m.EXPECT().
                    Save(gomock.Any(), gomock.Any()).
                    Return(errors.New("invalid email")).
                    Times(0)
            },
            wantErr:     true,
            errContains: "invalid email",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockRepo := mocks.NewMockUserRepository(ctrl)
            tt.setupMock(mockRepo)

            useCase := usecase.NewUserUseCase(mockRepo)

            ctx := context.Background()
            err := useCase.CreateUser(ctx, tt.email, tt.password)

            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errContains)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## Checklist

- [ ] mockgen installed in project
- [ ] `//go:generate` directives added to interface files
- [ ] Mocks generated (`go generate ./...`)
- [ ] Tests use generated mocks instead of real dependencies
- [ ] Expectations are set before calling code
- [ ] `defer ctrl.Finish()` is called in every test
- [ ] Mocks are regenerated when interfaces change
