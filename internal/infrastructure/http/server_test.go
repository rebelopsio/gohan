package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/application/installation/dto"
	httpinfra "github.com/rebelopsio/gohan/internal/infrastructure/http"
	"github.com/rebelopsio/gohan/internal/infrastructure/http/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockStartInstallationUseCase is a mock for testing
type MockStartInstallationUseCase struct {
	mock.Mock
}

func (m *MockStartInstallationUseCase) Execute(ctx context.Context, request dto.InstallationRequest) (*dto.InstallationResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.InstallationResponse), args.Error(1)
}

// MockExecuteInstallationUseCase is a mock for testing
type MockExecuteInstallationUseCase struct {
	mock.Mock
}

func (m *MockExecuteInstallationUseCase) Execute(ctx context.Context, sessionID string) (*dto.InstallationProgressResponse, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.InstallationProgressResponse), args.Error(1)
}

// MockGetInstallationStatusUseCase is a mock for testing
type MockGetInstallationStatusUseCase struct {
	mock.Mock
}

func (m *MockGetInstallationStatusUseCase) Execute(ctx context.Context, sessionID string) (*dto.InstallationProgressResponse, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.InstallationProgressResponse), args.Error(1)
}

// MockListInstallationsUseCase is a mock for testing
type MockListInstallationsUseCase struct {
	mock.Mock
}

func (m *MockListInstallationsUseCase) Execute(ctx context.Context) (*dto.ListInstallationsResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ListInstallationsResponse), args.Error(1)
}

// MockCancelInstallationUseCase is a mock for testing
type MockCancelInstallationUseCase struct {
	mock.Mock
}

func (m *MockCancelInstallationUseCase) Execute(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func TestServer_Routes(t *testing.T) {
	mockStartUseCase := new(MockStartInstallationUseCase)
	mockExecuteUseCase := new(MockExecuteInstallationUseCase)
	mockGetStatusUseCase := new(MockGetInstallationStatusUseCase)
	mockListUseCase := new(MockListInstallationsUseCase)
	mockCancelUseCase := new(MockCancelInstallationUseCase)
	installationHandler := handlers.NewInstallationHandler(
		mockStartUseCase,
		mockExecuteUseCase,
		mockGetStatusUseCase,
		mockListUseCase,
		mockCancelUseCase,
	)

	config := httpinfra.Config{
		Host:         "localhost",
		Port:         8080,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server := httpinfra.NewServer(config, installationHandler, false)
	router := server.Router()

	t.Run("health check endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "ok")
	})

	t.Run("installation start endpoint", func(t *testing.T) {
		requestBody := dto.InstallationRequest{
			Components: []dto.ComponentRequest{
				{Name: "hyprland", Version: "0.35.0"},
			},
			AvailableSpace: 100 * 1024 * 1024 * 1024,
			RequiredSpace:  10 * 1024 * 1024 * 1024,
		}

		expectedResponse := &dto.InstallationResponse{
			SessionID:      "session-123",
			Status:         "pending",
			Message:        "Installation session created successfully",
			ComponentCount: 1,
		}

		mockStartUseCase.On("Execute", mock.Anything, requestBody).
			Return(expectedResponse, nil)

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/api/installation/start", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response dto.InstallationResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "session-123", response.SessionID)

		mockStartUseCase.AssertExpectations(t)
	})

	t.Run("installation execute endpoint", func(t *testing.T) {
		sessionID := "session-123"

		expectedResponse := &dto.InstallationProgressResponse{
			SessionID:           sessionID,
			Status:              "completed",
			CurrentPhase:        "completed",
			PercentComplete:     100,
			Message:             "Installation completed successfully",
			ComponentsInstalled: 1,
			ComponentsTotal:     1,
		}

		mockExecuteUseCase.On("Execute", mock.Anything, sessionID).
			Return(expectedResponse, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/installation/"+sessionID+"/execute", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.InstallationProgressResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, sessionID, response.SessionID)
		assert.Equal(t, 100, response.PercentComplete)

		mockExecuteUseCase.AssertExpectations(t)
	})

	t.Run("404 for unknown routes", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/unknown", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("CORS headers are set", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/api/installation/start", nil)
		req.Header.Set("Origin", "http://example.com")
		req.Header.Set("Access-Control-Request-Method", "POST")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), "POST")
	})
}
