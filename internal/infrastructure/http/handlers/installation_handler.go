package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rebelopsio/gohan/internal/application/installation/dto"
)

// StartInstallationUseCase defines the interface for starting an installation
type StartInstallationUseCase interface {
	Execute(ctx context.Context, request dto.InstallationRequest) (*dto.InstallationResponse, error)
}

// ExecuteInstallationUseCase defines the interface for executing an installation
type ExecuteInstallationUseCase interface {
	Execute(ctx context.Context, sessionID string) (*dto.InstallationProgressResponse, error)
}

// GetInstallationStatusUseCase defines the interface for getting installation status
type GetInstallationStatusUseCase interface {
	Execute(ctx context.Context, sessionID string) (*dto.InstallationProgressResponse, error)
}

// ListInstallationsUseCase defines the interface for listing all installations
type ListInstallationsUseCase interface {
	Execute(ctx context.Context) (*dto.ListInstallationsResponse, error)
}

// CancelInstallationUseCase defines the interface for cancelling an installation
type CancelInstallationUseCase interface {
	Execute(ctx context.Context, sessionID string) error
}

// InstallationHandler handles HTTP requests for installation operations
type InstallationHandler struct {
	startUseCase     StartInstallationUseCase
	executeUseCase   ExecuteInstallationUseCase
	getStatusUseCase GetInstallationStatusUseCase
	listUseCase      ListInstallationsUseCase
	cancelUseCase    CancelInstallationUseCase
}

// NewInstallationHandler creates a new installation handler
func NewInstallationHandler(
	startUseCase StartInstallationUseCase,
	executeUseCase ExecuteInstallationUseCase,
	getStatusUseCase GetInstallationStatusUseCase,
	listUseCase ListInstallationsUseCase,
	cancelUseCase CancelInstallationUseCase,
) *InstallationHandler {
	return &InstallationHandler{
		startUseCase:     startUseCase,
		executeUseCase:   executeUseCase,
		getStatusUseCase: getStatusUseCase,
		listUseCase:      listUseCase,
		cancelUseCase:    cancelUseCase,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// StartInstallation handles POST /api/installation/start
func (h *InstallationHandler) StartInstallation(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	var request dto.InstallationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Execute use case
	response, err := h.startUseCase.Execute(r.Context(), request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to start installation", err.Error())
		return
	}

	// Return successful response
	respondWithJSON(w, http.StatusCreated, response)
}

// ExecuteInstallation handles POST /api/installation/{sessionID}/execute
func (h *InstallationHandler) ExecuteInstallation(w http.ResponseWriter, r *http.Request) {
	// Get session ID from URL params
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		respondWithError(w, http.StatusBadRequest, "Session ID is required", "")
		return
	}

	// Execute use case
	response, err := h.executeUseCase.Execute(r.Context(), sessionID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to execute installation", err.Error())
		return
	}

	// Return successful response
	respondWithJSON(w, http.StatusOK, response)
}

// GetStatus handles GET /api/installation/{sessionID}/status
func (h *InstallationHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	// Get session ID from URL params
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		respondWithError(w, http.StatusBadRequest, "Session ID is required", "")
		return
	}

	// Execute use case
	response, err := h.getStatusUseCase.Execute(r.Context(), sessionID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Session not found", err.Error())
		return
	}

	// Return successful response
	respondWithJSON(w, http.StatusOK, response)
}

// ListInstallations handles GET /api/installation
func (h *InstallationHandler) ListInstallations(w http.ResponseWriter, r *http.Request) {
	// Execute use case
	response, err := h.listUseCase.Execute(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list installations", err.Error())
		return
	}

	// Return successful response
	respondWithJSON(w, http.StatusOK, response)
}

// CancelInstallation handles POST /api/installation/{sessionID}/cancel
func (h *InstallationHandler) CancelInstallation(w http.ResponseWriter, r *http.Request) {
	// Get session ID from URL params
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		respondWithError(w, http.StatusBadRequest, "Session ID is required", "")
		return
	}

	// Execute use case
	err := h.cancelUseCase.Execute(r.Context(), sessionID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to cancel installation", err.Error())
		return
	}

	// Return successful response
	respondWithJSON(w, http.StatusOK, map[string]string{
		"message":    "Installation cancelled successfully",
		"session_id": sessionID,
	})
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, status int, error string, message string) {
	response := ErrorResponse{
		Error:   error,
		Message: message,
	}
	respondWithJSON(w, status, response)
}
