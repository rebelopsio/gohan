package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/rebelopsio/gohan/internal/application/installation/dto"
	"github.com/rebelopsio/gohan/internal/infrastructure/http/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockStartInstallationUseCase is a mock for the StartInstallationUseCase
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

// MockExecuteInstallationUseCase is a mock for the ExecuteInstallationUseCase
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

func TestInstallationHandler_StartInstallation(t *testing.T) {
	t.Run("successfully starts installation", func(t *testing.T) {
		mockUseCase := new(MockStartInstallationUseCase)
		handler := handlers.NewInstallationHandler(mockUseCase, nil, nil, nil, nil)

		requestBody := dto.InstallationRequest{
			Components: []dto.ComponentRequest{
				{
					Name:    "hyprland",
					Version: "0.35.0",
				},
			},
			AvailableSpace: 100 * 1024 * 1024 * 1024, // 100 GB
			RequiredSpace:  10 * 1024 * 1024 * 1024,  // 10 GB
		}

		expectedResponse := &dto.InstallationResponse{
			SessionID:      "session-123",
			Status:         "pending",
			Message:        "Installation session created successfully",
			StartedAt:      "2024-01-01T00:00:00Z",
			ComponentCount: 1,
		}

		mockUseCase.On("Execute", mock.Anything, requestBody).
			Return(expectedResponse, nil)

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/api/installation/start", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.StartInstallation(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response dto.InstallationResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "session-123", response.SessionID)
		assert.Equal(t, "pending", response.Status)
		assert.Equal(t, 1, response.ComponentCount)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("returns bad request for invalid JSON", func(t *testing.T) {
		mockUseCase := new(MockStartInstallationUseCase)
		handler := handlers.NewInstallationHandler(mockUseCase, nil, nil, nil, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/installation/start", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.StartInstallation(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns bad request for empty components", func(t *testing.T) {
		mockUseCase := new(MockStartInstallationUseCase)
		handler := handlers.NewInstallationHandler(mockUseCase, nil, nil, nil, nil)

		requestBody := dto.InstallationRequest{
			Components:     []dto.ComponentRequest{},
			AvailableSpace: 100 * 1024 * 1024 * 1024,
			RequiredSpace:  10 * 1024 * 1024 * 1024,
		}

		mockUseCase.On("Execute", mock.Anything, requestBody).
			Return(nil, assert.AnError)

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/api/installation/start", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.StartInstallation(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("returns service unavailable for internal errors", func(t *testing.T) {
		mockUseCase := new(MockStartInstallationUseCase)
		handler := handlers.NewInstallationHandler(mockUseCase, nil, nil, nil, nil)

		requestBody := dto.InstallationRequest{
			Components: []dto.ComponentRequest{
				{Name: "hyprland", Version: "0.35.0"},
			},
			AvailableSpace: 100 * 1024 * 1024 * 1024,
			RequiredSpace:  10 * 1024 * 1024 * 1024,
		}

		mockUseCase.On("Execute", mock.Anything, requestBody).
			Return(nil, assert.AnError)

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/api/installation/start", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.StartInstallation(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestInstallationHandler_ExecuteInstallation(t *testing.T) {
	t.Run("successfully executes installation", func(t *testing.T) {
		mockUseCase := new(MockExecuteInstallationUseCase)
		handler := handlers.NewInstallationHandler(nil, mockUseCase, nil, nil, nil)

		sessionID := "session-123"

		expectedResponse := &dto.InstallationProgressResponse{
			SessionID:           sessionID,
			Status:              "completed",
			CurrentPhase:        "completed",
			PercentComplete:     100,
			Message:             "Installation completed successfully",
			EstimatedRemaining:  "0s",
			ComponentsInstalled: 1,
			ComponentsTotal:     1,
		}

		mockUseCase.On("Execute", mock.Anything, sessionID).
			Return(expectedResponse, nil)

		// Create request with chi URL params
		req := httptest.NewRequest(http.MethodPost, "/api/installation/"+sessionID+"/execute", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("sessionID", sessionID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler.ExecuteInstallation(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.InstallationProgressResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, sessionID, response.SessionID)
		assert.Equal(t, "completed", response.Status)
		assert.Equal(t, 100, response.PercentComplete)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("returns not found for non-existent session", func(t *testing.T) {
		mockUseCase := new(MockExecuteInstallationUseCase)
		handler := handlers.NewInstallationHandler(nil, mockUseCase, nil, nil, nil)

		sessionID := "non-existent"

		mockUseCase.On("Execute", mock.Anything, sessionID).
			Return(nil, assert.AnError)

		req := httptest.NewRequest(http.MethodPost, "/api/installation/"+sessionID+"/execute", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("sessionID", sessionID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler.ExecuteInstallation(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("returns bad request for empty session ID", func(t *testing.T) {
		mockUseCase := new(MockExecuteInstallationUseCase)
		handler := handlers.NewInstallationHandler(nil, mockUseCase, nil, nil, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/installation//execute", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("sessionID", "")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler.ExecuteInstallation(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestInstallationHandler_ContentTypeValidation(t *testing.T) {
	t.Run("accepts application/json content type", func(t *testing.T) {
		mockUseCase := new(MockStartInstallationUseCase)
		handler := handlers.NewInstallationHandler(mockUseCase, nil, nil, nil, nil)

		requestBody := dto.InstallationRequest{
			Components: []dto.ComponentRequest{
				{Name: "hyprland", Version: "0.35.0"},
			},
			AvailableSpace: 100 * 1024 * 1024 * 1024,
			RequiredSpace:  10 * 1024 * 1024 * 1024,
		}

		expectedResponse := &dto.InstallationResponse{
			SessionID: "session-123",
			Status:    "pending",
		}

		mockUseCase.On("Execute", mock.Anything, requestBody).
			Return(expectedResponse, nil)

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/api/installation/start", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.StartInstallation(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		mockUseCase.AssertExpectations(t)
	})
}
