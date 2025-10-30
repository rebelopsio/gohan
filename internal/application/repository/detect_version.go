package repository

import (
	"context"

	"github.com/rebelopsio/gohan/internal/domain/repository"
)

// DetectVersionRequest contains parameters for detecting Debian version
type DetectVersionRequest struct {
	// Empty for now - detection happens automatically
}

// DetectVersionResponse contains the detected version information
type DetectVersionResponse struct {
	Codename       string // e.g., "sid", "trixie", "bookworm"
	Version        string // e.g., "unstable", "testing", "12"
	IsSupported    bool   // Whether this version is supported
	IsSid          bool   // True if Debian Sid
	IsTrixie       bool   // True if Debian Trixie
	IsBookworm     bool   // True if Debian Bookworm
	IsUbuntu       bool   // True if Ubuntu
	SupportMessage string // Warning or error message if not fully supported
	DisplayString  string // Human-readable version string
}

// VersionDetector is the interface for detecting system version
type VersionDetector interface {
	DetectVersion() (*repository.DebianVersion, error)
}

// DetectVersionUseCase handles detecting and validating the Debian version
type DetectVersionUseCase struct {
	detector VersionDetector
}

// NewDetectVersionUseCase creates a new use case instance
func NewDetectVersionUseCase(detector VersionDetector) *DetectVersionUseCase {
	return &DetectVersionUseCase{
		detector: detector,
	}
}

// Execute detects the current Debian version
func (uc *DetectVersionUseCase) Execute(ctx context.Context, req DetectVersionRequest) (*DetectVersionResponse, error) {
	// Detect version from system
	dv, err := uc.detector.DetectVersion()
	if err != nil {
		return nil, err
	}

	// Build response
	return &DetectVersionResponse{
		Codename:       dv.Codename(),
		Version:        dv.Version(),
		IsSupported:    dv.IsSupported(),
		IsSid:          dv.IsSid(),
		IsTrixie:       dv.IsTrixie(),
		IsBookworm:     dv.IsBookworm(),
		IsUbuntu:       dv.IsUbuntu(),
		SupportMessage: dv.SupportMessage(),
		DisplayString:  dv.String(),
	}, nil
}
