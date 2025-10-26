# CLAUDE.md - Go Application Development Guide

## Overview

This guide outlines the development approach for building Go applications using:

- **Behavior-Driven Development (BDD)** - Define features using Gherkin syntax
- **Acceptance Test Driven Development (ATDD)** - BDD scenarios inform acceptance tests
- **Test Driven Development (TDD)** - Write tests before implementation
- **Table-Driven Tests** - Go idiomatic testing patterns
- **Interface-Based Design** - Dependency injection and testability
- **Clean Architecture** - Domain-driven with clear boundaries

## Philosophy

**BDD → ATDD → TDD → Red → Green → Refactor**

1. Write Gherkin feature scenarios (BDD)
2. Convert scenarios to acceptance tests (ATDD)
3. Write failing unit/integration tests (TDD)
4. Write minimal code to pass tests
5. Refactor while keeping tests green
6. Repeat

## Expert Agents

This project uses specialized expert agents to ensure quality at each stage of development. These agents are integrated into the workflow and should be used proactively.

### Available Agents

#### 1. BDD Expert (`.claude/agents/bdd-expert.md`)

**Use for:**
- Reviewing and improving Gherkin feature files
- Ensuring scenarios are declarative (what, not how)
- Removing implementation details and technical coupling
- Making scenarios focus on user behavior
- Validating ubiquitous language usage

#### 2. DDD Expert (`.claude/agents/ddd-expert.md`)

**Use for:**
- Domain modeling and bounded context identification
- Designing entities, value objects, and aggregates
- Identifying ubiquitous language
- Ensuring domain logic stays in domain layer
- Preventing anemic domain models

#### 3. Test Quality Reviewer (`.claude/agents/test-quality-reviewer.md`)

**Use for:**
- Reviewing unit, integration, and E2E tests
- Ensuring tests focus on behavior, not implementation
- Validating meaningful test coverage
- Identifying brittle tests

#### 4. Go Expert (`.claude/agents/go-expert.md`)

**Use for:**
- Reviewing Go idioms and best practices
- Optimizing interface design and composition
- Ensuring proper error handling patterns
- Validating concurrency patterns and goroutine usage
- Identifying resource leaks (goroutines, files, connections)
- Reviewing context usage and cancellation
- Ensuring proper package structure and dependencies
- Validating test patterns (table-driven, subtests)

**When to use:**
- After writing or updating Go code
- When designing interfaces and abstractions
- When implementing concurrent operations
- During error handling implementation
- Before finalizing package APIs
- When experiencing race conditions or deadlocks

### Agent Integration in Development Workflow

```
┌─────────────────────────────────────────────────────────────┐
│ PHASE 1: BDD (Feature Definition)                          │
├─────────────────────────────────────────────────────────────┤
│ 1. Write Gherkin feature file                              │
│ 2. → RUN BDD EXPERT AGENT                                  │
│ 3. Apply recommendations                                    │
│ 4. Review with stakeholders                                │
│ 5. → RUN DDD EXPERT AGENT (for domain modeling)           │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│ PHASE 2: ATDD (Acceptance Tests)                           │
├─────────────────────────────────────────────────────────────┤
│ 1. Convert Gherkin to integration tests                    │
│ 2. Write failing E2E tests                                  │
│ 3. → RUN TEST QUALITY REVIEWER AGENT                       │
│ 4. Refactor tests based on feedback                        │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│ PHASE 3: TDD (Unit Tests)                                  │
├─────────────────────────────────────────────────────────────┤
│ 1. Write failing unit tests (table-driven)                 │
│ 2. Implement code (Red → Green)                            │
│ 3. → RUN GO EXPERT AGENT                                   │
│ 4. → RUN TEST QUALITY REVIEWER AGENT                       │
│ 5. Refactor tests and code                                 │
│ 6. Ensure all tests pass                                   │
│ 7. Run race detector: go test -race                        │
└─────────────────────────────────────────────────────────────┘
```

## BDD with Gherkin

All features are documented using Gherkin syntax in `docs/features/*.feature` files.

### Feature File Structure

```gherkin
Feature: Feature Name
  As a [role]
  I want [feature]
  So that [benefit]

  Background:
    Given [common setup for all scenarios]

  Scenario: Scenario Name
    Given [initial context]
    When [action occurs]
    Then [expected outcome]
    And [additional outcome]

  @not-implemented
  Scenario: Future Feature
    # Tagged scenarios are not yet implemented
```

### Example: User Authentication API

**Feature File** (`docs/features/authentication.feature`)

```gherkin
Feature: User Authentication API
  As an API client
  I want to authenticate with the service
  So that I can access protected resources

  Scenario: Successful authentication
    Given a valid user account exists
    When I POST credentials to /api/auth/login
    Then I should receive a 200 status code
    And I should receive a JWT token
    And the token should be valid for 24 hours

  Scenario: Invalid credentials
    Given I have invalid credentials
    When I POST credentials to /api/auth/login
    Then I should receive a 401 status code
    And I should receive an error message
```

## Project Setup

### Project Structure

```
my-service/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── domain/               # Domain models and business logic
│   │   ├── user/
│   │   │   ├── user.go
│   │   │   ├── user_test.go
│   │   │   ├── repository.go # Repository interface
│   │   │   └── service.go    # Domain service
│   │   └── auth/
│   ├── infrastructure/       # External dependencies
│   │   ├── postgres/
│   │   │   └── user_repository.go
│   │   └── http/
│   │       ├── handlers/
│   │       └── middleware/
│   └── application/          # Use cases/application services
│       └── auth_service.go
├── pkg/                      # Public packages
│   ├── errors/
│   └── testutil/
├── tests/
│   ├── integration/
│   └── e2e/
├── docs/
│   └── features/
├── go.mod
└── go.sum
```

### Dependencies

```bash
# Initialize module
go mod init github.com/yourusername/my-service

# Testing
go get github.com/stretchr/testify
go get github.com/golang/mock/gomock

# HTTP
go get github.com/go-chi/chi/v5
go get github.com/go-chi/cors

# Database
go get github.com/jmoiron/sqlx
go get github.com/lib/pq

# Configuration
go get github.com/kelseyhightower/envconfig

# Logging
go get go.uber.org/zap
```

## Testing Pyramid

```
    /\
   /  \    E2E Tests (HTTP integration)
  /____\   - Full request/response cycle
  /      \  Integration Tests
 /________\ - Database, external services
/          \ Unit Tests (Table-driven)
/____________\ - Domain logic, pure functions
```

## Go Testing Patterns

### 1. Table-Driven Tests

```go
// internal/domain/user/validator_test.go
package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yourusername/my-service/internal/domain/user"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "valid email",
			email:   "user@example.com",
			wantErr: false,
		},
		{
			name:    "invalid email - missing @",
			email:   "userexample.com",
			wantErr: true,
		},
		{
			name:    "invalid email - missing domain",
			email:   "user@",
			wantErr: true,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user.ValidateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
```

### 2. Interface-Based Testing with Mocks

```go
// internal/domain/user/repository.go
package user

import (
	"context"
)

// Repository defines the interface for user persistence
type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

// internal/domain/user/service.go
package user

import (
	"context"
	"errors"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, email, password string) (*User, error) {
	// Validate input
	if err := ValidateEmail(email); err != nil {
		return nil, err
	}
	if err := ValidatePassword(password); err != nil {
		return nil, err
	}

	// Check if user exists
	existing, err := s.repo.FindByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, ErrDuplicateEmail
	}

	// Create user
	user := &User{
		Email:    email,
		Password: HashPassword(password),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// internal/domain/user/service_test.go
package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yourusername/my-service/internal/domain/user"
)

// MockRepository is a mock implementation of user.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_Register(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		password  string
		setupMock func(*MockRepository)
		wantErr   error
	}{
		{
			name:     "successful registration",
			email:    "user@example.com",
			password: "SecurePass123",
			setupMock: func(m *MockRepository) {
				m.On("FindByEmail", mock.Anything, "user@example.com").
					Return(nil, user.ErrUserNotFound)
				m.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:     "duplicate email",
			email:    "existing@example.com",
			password: "SecurePass123",
			setupMock: func(m *MockRepository) {
				m.On("FindByEmail", mock.Anything, "existing@example.com").
					Return(&user.User{Email: "existing@example.com"}, nil)
			},
			wantErr: user.ErrDuplicateEmail,
		},
		{
			name:     "invalid email",
			email:    "invalid-email",
			password: "SecurePass123",
			setupMock: func(m *MockRepository) {
				// No mock setup needed - validation happens before repo call
			},
			wantErr: user.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.setupMock(mockRepo)

			service := user.NewService(mockRepo)
			_, err := service.Register(context.Background(), tt.email, tt.password)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
```

### 3. Subtests and Test Helpers

```go
// internal/domain/user/password_test.go
package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/my-service/internal/domain/user"
)

func TestPassword(t *testing.T) {
	t.Run("HashPassword", func(t *testing.T) {
		password := "MySecurePassword123"
		
		hash := user.HashPassword(password)
		
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
	})

	t.Run("VerifyPassword", func(t *testing.T) {
		t.Run("correct password", func(t *testing.T) {
			password := "MySecurePassword123"
			hash := user.HashPassword(password)

			assert.True(t, user.VerifyPassword(hash, password))
		})

		t.Run("incorrect password", func(t *testing.T) {
			hash := user.HashPassword("MySecurePassword123")

			assert.False(t, user.VerifyPassword(hash, "WrongPassword"))
		})
	})

	t.Run("ValidatePassword", func(t *testing.T) {
		tests := []struct {
			password string
			wantErr  bool
		}{
			{"SecurePass1", false},
			{"short", true},
			{"noupper1", true},
			{"NOLOWER1", true},
			{"NoNumber", true},
		}

		for _, tt := range tests {
			t.Run(tt.password, func(t *testing.T) {
				err := user.ValidatePassword(tt.password)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}
```

## Domain-Driven Design in Go

### 1. Entities and Value Objects

```go
// internal/domain/user/user.go
package user

import (
	"time"

	"github.com/google/uuid"
)

// User is an aggregate root
type User struct {
	ID        string    `db:"id"`
	Email     Email     `db:"email"`
	Password  string    `db:"password_hash"`
	Profile   Profile   `db:"-"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// NewUser creates a new user with validated data
func NewUser(email, password string) (*User, error) {
	e, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	if err := ValidatePassword(password); err != nil {
		return nil, err
	}

	return &User{
		ID:        uuid.New().String(),
		Email:     e,
		Password:  HashPassword(password),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// Email is a value object
type Email struct {
	value string
}

func NewEmail(email string) (Email, error) {
	if err := ValidateEmail(email); err != nil {
		return Email{}, err
	}
	return Email{value: email}, nil
}

func (e Email) String() string {
	return e.value
}

// Profile is a value object
type Profile struct {
	FirstName string
	LastName  string
	Bio       string
}

func (p Profile) FullName() string {
	return p.FirstName + " " + p.LastName
}
```

### 2. Repository Pattern

```go
// internal/infrastructure/postgres/user_repository.go
package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/yourusername/my-service/internal/domain/user"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		u.ID,
		u.Email.String(),
		u.Password,
		u.CreatedAt,
		u.UpdatedAt,
	)

	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = $1`

	var u user.User
	err := r.db.GetContext(ctx, &u, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}
```

## HTTP Handler Testing

```go
// internal/infrastructure/http/handlers/auth_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yourusername/my-service/internal/domain/user"
)

type AuthHandler struct {
	userService *user.Service
}

func NewAuthHandler(userService *user.Service) *AuthHandler {
	return &AuthHandler{userService: userService}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	u, err := h.userService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case user.ErrDuplicateEmail:
			http.Error(w, err.Error(), http.StatusConflict)
		case user.ErrInvalidEmail, user.ErrInvalidPassword:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	resp := RegisterResponse{
		ID:    u.ID,
		Email: u.Email.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// internal/infrastructure/http/handlers/auth_handler_test.go
package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yourusername/my-service/internal/domain/user"
	"github.com/yourusername/my-service/internal/infrastructure/http/handlers"
)

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]string
		setupMock      func(*MockUserService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful registration",
			requestBody: map[string]string{
				"email":    "user@example.com",
				"password": "SecurePass123",
			},
			setupMock: func(m *MockUserService) {
				m.On("Register", mock.Anything, "user@example.com", "SecurePass123").
					Return(&user.User{ID: "123", Email: "user@example.com"}, nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp handlers.RegisterResponse
				err := json.NewDecoder(rec.Body).Decode(&resp)
				assert.NoError(t, err)
				assert.Equal(t, "123", resp.ID)
				assert.Equal(t, "user@example.com", resp.Email)
			},
		},
		{
			name: "duplicate email",
			requestBody: map[string]string{
				"email":    "existing@example.com",
				"password": "SecurePass123",
			},
			setupMock: func(m *MockUserService) {
				m.On("Register", mock.Anything, "existing@example.com", "SecurePass123").
					Return(nil, user.ErrDuplicateEmail)
			},
			expectedStatus: http.StatusConflict,
			checkResponse:  func(t *testing.T, rec *httptest.ResponseRecorder) {},
		},
		{
			name: "invalid email",
			requestBody: map[string]string{
				"email":    "invalid-email",
				"password": "SecurePass123",
			},
			setupMock: func(m *MockUserService) {
				m.On("Register", mock.Anything, "invalid-email", "SecurePass123").
					Return(nil, user.ErrInvalidEmail)
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  func(t *testing.T, rec *httptest.ResponseRecorder) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			tt.setupMock(mockService)

			handler := handlers.NewAuthHandler(mockService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			handler.Register(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			tt.checkResponse(t, rec)
			mockService.AssertExpectations(t)
		})
	}
}
```

## Integration Testing

```go
// tests/integration/user_test.go
//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/my-service/internal/domain/user"
	"github.com/yourusername/my-service/internal/infrastructure/postgres"
)

func TestUserRepository_Integration(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	t.Run("Create and Find user", func(t *testing.T) {
		// Create user
		u, err := user.NewUser("test@example.com", "SecurePass123")
		require.NoError(t, err)

		err = repo.Create(ctx, u)
		require.NoError(t, err)

		// Find user
		found, err := repo.FindByEmail(ctx, "test@example.com")
		require.NoError(t, err)
		assert.Equal(t, u.ID, found.ID)
		assert.Equal(t, u.Email.String(), found.Email.String())
	})

	t.Run("FindByEmail returns error for non-existent user", func(t *testing.T) {
		_, err := repo.FindByEmail(ctx, "nonexistent@example.com")
		assert.ErrorIs(t, err, user.ErrUserNotFound)
	})
}
```

## Concurrency Patterns

### 1. Context and Cancellation

```go
// internal/application/auth_service.go
package application

import (
	"context"
	"time"
)

func (s *AuthService) ProcessWithTimeout(ctx context.Context, data string) error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Use errgroup for concurrent operations
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return s.validateData(ctx, data)
	})

	g.Go(func() error {
		return s.enrichData(ctx, data)
	})

	return g.Wait()
}
```

### 2. Worker Pools

```go
// pkg/worker/pool.go
package worker

import (
	"context"
	"sync"
)

type Task func(context.Context) error

type Pool struct {
	workers int
	tasks   chan Task
}

func NewPool(workers int) *Pool {
	return &Pool{
		workers: workers,
		tasks:   make(chan Task, workers*2),
	}
}

func (p *Pool) Start(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case task, ok := <-p.tasks:
					if !ok {
						return
					}
					_ = task(ctx)
				}
			}
		}()
	}

	wg.Wait()
}

func (p *Pool) Submit(task Task) {
	p.tasks <- task
}

func (p *Pool) Close() {
	close(p.tasks)
}
```

## Error Handling

```go
// pkg/errors/errors.go
package errors

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(code, message string) *Error {
	return &Error{Code: code, Message: message}
}

func Wrap(err error, code, message string) *Error {
	return &Error{Code: code, Message: message, Err: err}
}

// Common errors
var (
	ErrNotFound     = New("NOT_FOUND", "resource not found")
	ErrUnauthorized = New("UNAUTHORIZED", "unauthorized")
	ErrForbidden    = New("FORBIDDEN", "forbidden")
	ErrBadRequest   = New("BAD_REQUEST", "bad request")
)
```

## Development Workflow

1. **Write BDD Feature File** (`docs/features/*.feature`)
2. **Write Integration Tests** (HTTP handlers, database)
3. **Write Unit Tests** (Domain logic, table-driven)
4. **Implement** (Make tests pass)
5. **Run Race Detector** (`go test -race ./...`)
6. **Refactor** (Improve code quality)
7. **Run All Tests** (`go test -v -cover ./...`)
8. **Repeat**

## Running Tests

```bash
# Run all tests
go test ./...

# Run with race detector
go test -race ./...

# Run with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package
go test ./internal/domain/user/...

# Run integration tests only
go test -tags=integration ./tests/integration/...

# Run benchmarks
go test -bench=. ./...
```

## Continuous Integration

```yaml
# .github/workflows/test.yml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: testdb
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
          cache: true
      
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
      
      - name: Run integration tests
        run: go test -v -tags=integration ./tests/integration/...
        env:
          DATABASE_URL: postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
```

## Best Practices

1. **Accept interfaces, return structs** - Flexibility at boundaries
2. **Use table-driven tests** - Go idiomatic testing
3. **Context everywhere** - Cancellation and timeouts
4. **Error wrapping** - Preserve context with `fmt.Errorf` or custom errors
5. **Defer cleanup** - Resource management with defer
6. **Small interfaces** - Single-method interfaces are powerful
7. **Package organization** - Clear boundaries between layers
8. **Race detection** - Always run with `-race` flag

## Summary

This approach ensures:

- ✅ Clean architecture with clear boundaries
- ✅ Interface-based design for testability
- ✅ Table-driven tests following Go idioms
- ✅ Proper error handling and context usage
- ✅ Concurrency safety with race detection
- ✅ Integration testing with real dependencies

Remember: **Test First, Interface-Driven, Concurrent-Safe**
