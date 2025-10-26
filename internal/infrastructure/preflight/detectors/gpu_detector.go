package detectors

import (
	"context"
	"os/exec"
	"regexp"
	"strings"

	"github.com/rebelopsio/gohan/internal/domain/preflight"
)

// SystemGPUDetector implements preflight.GPUDetector using lspci
type SystemGPUDetector struct{}

// NewSystemGPUDetector creates a new GPU detector
func NewSystemGPUDetector() *SystemGPUDetector {
	return &SystemGPUDetector{}
}

// DetectGPUs returns all detected GPUs
func (d *SystemGPUDetector) DetectGPUs(ctx context.Context) ([]preflight.GPUType, error) {
	cmd := exec.CommandContext(ctx, "lspci", "-nn")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	gpus := make([]preflight.GPUType, 0)
	lines := strings.Split(string(output), "\n")

	// Regex to match VGA/3D controller lines
	// Example: 01:00.0 VGA compatible controller [0300]: NVIDIA Corporation [10de:2684] (rev a1)
	vgaRegex := regexp.MustCompile(`(?i)(VGA|3D|Display).*controller`)
	pciIDRegex := regexp.MustCompile(`\[([0-9a-fA-F]{4}):([0-9a-fA-F]{4})\]`)

	for _, line := range lines {
		if !vgaRegex.MatchString(line) {
			continue
		}

		gpu, err := d.parseGPULine(line, pciIDRegex)
		if err == nil {
			gpus = append(gpus, gpu)
		}
	}

	if len(gpus) == 0 {
		return nil, preflight.ErrInvalidGPU
	}

	return gpus, nil
}

// PrimaryGPU returns the primary GPU (first detected)
func (d *SystemGPUDetector) PrimaryGPU(ctx context.Context) (preflight.GPUType, error) {
	gpus, err := d.DetectGPUs(ctx)
	if err != nil {
		return preflight.GPUType{}, err
	}

	if len(gpus) == 0 {
		return preflight.GPUType{}, preflight.ErrInvalidGPU
	}

	return gpus[0], nil
}

func (d *SystemGPUDetector) parseGPULine(line string, pciIDRegex *regexp.Regexp) (preflight.GPUType, error) {
	// Extract PCI ID
	pciMatches := pciIDRegex.FindStringSubmatch(line)
	var pciID string
	if len(pciMatches) >= 3 {
		pciID = pciMatches[1] + ":" + pciMatches[2]
	}

	// Determine vendor and extract model
	vendor := d.detectVendor(line)
	model := d.extractModel(line, vendor)

	return preflight.NewGPUType(vendor, model, pciID)
}

func (d *SystemGPUDetector) detectVendor(line string) preflight.GPUVendor {
	lineLower := strings.ToLower(line)

	if strings.Contains(lineLower, "nvidia") {
		return preflight.GPUVendorNVIDIA
	}
	if strings.Contains(lineLower, "amd") || strings.Contains(lineLower, "radeon") {
		return preflight.GPUVendorAMD
	}
	if strings.Contains(lineLower, "intel") {
		return preflight.GPUVendorIntel
	}

	return preflight.GPUVendorUnknown
}

func (d *SystemGPUDetector) extractModel(line string, vendor preflight.GPUVendor) string {
	// Split on ": " to get the description part
	parts := strings.Split(line, ": ")
	if len(parts) < 2 {
		return ""
	}

	desc := parts[1]

	// Remove PCI ID in brackets at the end
	pciIDRegex := regexp.MustCompile(`\s*\[[0-9a-fA-F]{4}:[0-9a-fA-F]{4}\].*$`)
	desc = pciIDRegex.ReplaceAllString(desc, "")

	// Remove vendor name from description
	vendorName := string(vendor)
	if vendor == preflight.GPUVendorAMD {
		// AMD cards often say "Advanced Micro Devices" or "ATI"
		desc = strings.ReplaceAll(desc, "Advanced Micro Devices, Inc.", "")
		desc = strings.ReplaceAll(desc, "ATI Technologies Inc", "")
	}
	desc = strings.ReplaceAll(desc, vendorName, "")
	desc = strings.ReplaceAll(desc, "Corporation", "")
	desc = strings.ReplaceAll(desc, "[", "")
	desc = strings.ReplaceAll(desc, "]", "")

	return strings.TrimSpace(desc)
}
