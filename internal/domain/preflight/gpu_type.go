package preflight

import (
	"fmt"
	"strings"
)

// GPUType represents a detected GPU
type GPUType struct {
	vendor GPUVendor
	model  string
	pciID  string
}

// NewGPUType creates a new GPU type value object
func NewGPUType(vendor GPUVendor, model string, pciID string) (GPUType, error) {
	if vendor == "" {
		return GPUType{}, ErrInvalidGPU
	}

	return GPUType{
		vendor: vendor,
		model:  strings.TrimSpace(model),
		pciID:  strings.TrimSpace(pciID),
	}, nil
}

// Vendor returns the GPU vendor
func (g GPUType) Vendor() GPUVendor {
	return g.vendor
}

// Model returns the GPU model name
func (g GPUType) Model() string {
	return g.model
}

// PCIID returns the PCI device ID
func (g GPUType) PCIID() string {
	return g.pciID
}

// IsNVIDIA returns true for NVIDIA GPUs
func (g GPUType) IsNVIDIA() bool {
	return g.vendor == GPUVendorNVIDIA
}

// IsAMD returns true for AMD GPUs
func (g GPUType) IsAMD() bool {
	return g.vendor == GPUVendorAMD
}

// IsIntel returns true for Intel GPUs
func (g GPUType) IsIntel() bool {
	return g.vendor == GPUVendorIntel
}

// RequiresProprietaryDriver returns true if proprietary drivers needed
func (g GPUType) RequiresProprietaryDriver() bool {
	return g.IsNVIDIA()
}

// RequiresSpecialConfiguration returns true if extra setup needed
func (g GPUType) RequiresSpecialConfiguration() bool {
	return g.IsNVIDIA()
}

// String returns the string representation
func (g GPUType) String() string {
	if g.model != "" {
		return fmt.Sprintf("%s %s", g.vendor, g.model)
	}
	return string(g.vendor)
}
