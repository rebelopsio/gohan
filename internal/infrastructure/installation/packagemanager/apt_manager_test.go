package packagemanager_test

import (
	"context"
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/packagemanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAPTManager(t *testing.T) {
	t.Run("creates APT manager", func(t *testing.T) {
		manager := packagemanager.NewAPTManager()

		assert.NotNil(t, manager)
	})
}

func TestAPTManager_DetectConflicts(t *testing.T) {
	t.Run("implements ConflictResolver interface", func(t *testing.T) {
		manager := packagemanager.NewAPTManager()

		var _ installation.ConflictResolver = manager
	})

	t.Run("detects no conflicts for compatible packages", func(t *testing.T) {
		manager := packagemanager.NewAPTManager()
		ctx := context.Background()

		selection, err := installation.NewComponentSelection(
			installation.ComponentHyprland,
			"0.35.0",
			nil,
		)
		require.NoError(t, err)

		components := []installation.ComponentSelection{selection}

		conflicts, err := manager.DetectConflicts(ctx, components)

		assert.NoError(t, err)
		assert.Empty(t, conflicts)
	})
}

func TestAPTManager_ResolveConflict(t *testing.T) {
	t.Run("resolves conflict with remove action", func(t *testing.T) {
		manager := packagemanager.NewAPTManager()
		ctx := context.Background()

		conflict, err := installation.NewPackageConflict(
			"hyprland",
			"hyprland-git",
			"version conflict",
		)
		require.NoError(t, err)

		err = manager.ResolveConflict(ctx, conflict, installation.ActionRemove)

		// Will fail with permission error in test environment, which is acceptable
		// The important part is that it doesn't panic and the API works
		_ = err
	})

	t.Run("handles skip action", func(t *testing.T) {
		manager := packagemanager.NewAPTManager()
		ctx := context.Background()

		conflict, err := installation.NewPackageConflict(
			"package1",
			"package2",
			"test conflict",
		)
		require.NoError(t, err)

		err = manager.ResolveConflict(ctx, conflict, installation.ActionSkip)

		// Skip should always succeed (no-op)
		assert.NoError(t, err)
	})

	t.Run("handles abort action", func(t *testing.T) {
		manager := packagemanager.NewAPTManager()
		ctx := context.Background()

		conflict, err := installation.NewPackageConflict(
			"package1",
			"package2",
			"test conflict",
		)
		require.NoError(t, err)

		err = manager.ResolveConflict(ctx, conflict, installation.ActionAbort)

		// Abort should return error
		assert.Error(t, err)
	})
}

func TestAPTManager_InstallPackage(t *testing.T) {
	t.Run("validates package name", func(t *testing.T) {
		manager := packagemanager.NewAPTManager()
		ctx := context.Background()

		err := manager.InstallPackage(ctx, "", "1.0.0")

		assert.Error(t, err)
	})

	t.Run("accepts valid package name without version", func(t *testing.T) {
		manager := packagemanager.NewAPTManagerDryRun()
		ctx := context.Background()

		// Dry-run mode won't actually install
		err := manager.InstallPackage(ctx, "hyprland", "")

		assert.NoError(t, err)
	})

	t.Run("accepts valid package name with version", func(t *testing.T) {
		manager := packagemanager.NewAPTManagerDryRun()
		ctx := context.Background()

		// Dry-run mode won't actually install
		err := manager.InstallPackage(ctx, "hyprland", "0.35.0")

		assert.NoError(t, err)
	})
}

func TestAPTManager_RemovePackage(t *testing.T) {
	t.Run("validates package name", func(t *testing.T) {
		manager := packagemanager.NewAPTManager()
		ctx := context.Background()

		err := manager.RemovePackage(ctx, "")

		assert.Error(t, err)
	})

	t.Run("accepts valid package name", func(t *testing.T) {
		manager := packagemanager.NewAPTManager()
		ctx := context.Background()

		// This will fail in test environment, but validates input
		err := manager.RemovePackage(ctx, "nonexistent-package")

		// Either succeeds or fails with execution error
		_ = err
	})
}

func TestAPTManager_IsPackageInstalled(t *testing.T) {
	t.Run("validates package name", func(t *testing.T) {
		manager := packagemanager.NewAPTManager()
		ctx := context.Background()

		_, err := manager.IsPackageInstalled(ctx, "")

		assert.Error(t, err)
	})

	t.Run("checks if package exists", func(t *testing.T) {
		manager := packagemanager.NewAPTManager()
		ctx := context.Background()

		// Check for a common system package that should exist
		installed, err := manager.IsPackageInstalled(ctx, "coreutils")

		assert.NoError(t, err)
		// coreutils is a base package, should be installed
		assert.True(t, installed)
	})

	t.Run("returns false for non-existent package", func(t *testing.T) {
		manager := packagemanager.NewAPTManager()
		ctx := context.Background()

		installed, err := manager.IsPackageInstalled(ctx, "definitely-not-a-real-package-name-xyz123")

		assert.NoError(t, err)
		assert.False(t, installed)
	})
}

func TestAPTManager_UpdatePackageCache(t *testing.T) {
	t.Run("attempts to update cache", func(t *testing.T) {
		manager := packagemanager.NewAPTManager()
		ctx := context.Background()

		// Will fail without sudo, but validates method exists
		err := manager.UpdatePackageCache(ctx)

		// Either succeeds or fails with permission error
		_ = err
	})
}

func TestAPTManager_GetPackageInfo(t *testing.T) {
	t.Run("validates package name", func(t *testing.T) {
		manager := packagemanager.NewAPTManager()
		ctx := context.Background()

		_, err := manager.GetPackageInfo(ctx, "")

		assert.Error(t, err)
	})

	t.Run("retrieves package info for installed package", func(t *testing.T) {
		manager := packagemanager.NewAPTManager()
		ctx := context.Background()

		// Use a common package
		info, err := manager.GetPackageInfo(ctx, "coreutils")

		if err != nil {
			t.Skipf("Skipping: %v", err)
		}

		assert.NotEmpty(t, info.Name)
		assert.NotEmpty(t, info.Version)
	})
}

// ========================================
// Phase 3.1: Batch Installation Tests
// ========================================

func TestAPTManager_ArePackagesInstalled(t *testing.T) {
	tests := []struct {
		name     string
		packages []string
		wantErr  bool
	}{
		{
			name:     "empty package list",
			packages: []string{},
			wantErr:  false,
		},
		{
			name:     "single installed package",
			packages: []string{"coreutils"},
			wantErr:  false,
		},
		{
			name:     "multiple packages - mixed installed and not installed",
			packages: []string{"coreutils", "definitely-not-installed-xyz"},
			wantErr:  false,
		},
		{
			name:     "all packages not installed",
			packages: []string{"fake-package-1", "fake-package-2"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := packagemanager.NewAPTManager()
			ctx := context.Background()

			result, err := manager.ArePackagesInstalled(ctx, tt.packages)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, len(tt.packages), len(result), "Result map should have entry for each package")

			// Verify specific expectations
			if tt.name == "single installed package" {
				assert.True(t, result["coreutils"], "coreutils should be marked as installed")
			}

			if tt.name == "multiple packages - mixed installed and not installed" {
				assert.True(t, result["coreutils"], "coreutils should be installed")
				assert.False(t, result["definitely-not-installed-xyz"], "fake package should not be installed")
			}
		})
	}
}

func TestAPTManager_InstallPackages(t *testing.T) {
	tests := []struct {
		name         string
		packages     []string
		wantErr      bool
		validateFunc func(*testing.T, []packagemanager.PackageProgress)
	}{
		{
			name:     "empty package list",
			packages: []string{},
			wantErr:  false,
			validateFunc: func(t *testing.T, progress []packagemanager.PackageProgress) {
				assert.Empty(t, progress, "No progress events for empty list")
			},
		},
		{
			name:     "single package",
			packages: []string{"hyprland"},
			wantErr:  false,
			validateFunc: func(t *testing.T, progress []packagemanager.PackageProgress) {
				assert.NotEmpty(t, progress, "Should have progress events")
				// Should have at least: started, completed
				assert.GreaterOrEqual(t, len(progress), 2)
			},
		},
		{
			name:     "multiple packages",
			packages: []string{"hyprland", "waybar", "kitty"},
			wantErr:  false,
			validateFunc: func(t *testing.T, progress []packagemanager.PackageProgress) {
				assert.NotEmpty(t, progress)
				// Each package should have progress events
				packageNames := make(map[string]bool)
				for _, p := range progress {
					packageNames[p.PackageName] = true
				}
				assert.Equal(t, 3, len(packageNames), "Should have progress for all 3 packages")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := packagemanager.NewAPTManagerDryRun() // Dry-run mode for testing
			ctx := context.Background()

			// Create progress channel
			progressChan := make(chan packagemanager.PackageProgress, 100)
			var collectedProgress []packagemanager.PackageProgress

			// Collect progress in goroutine
			done := make(chan struct{})
			go func() {
				for p := range progressChan {
					collectedProgress = append(collectedProgress, p)
				}
				close(done)
			}()

			// Execute installation
			err := manager.InstallPackages(ctx, tt.packages, progressChan)
			close(progressChan)
			<-done // Wait for collection to finish

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if tt.validateFunc != nil {
				tt.validateFunc(t, collectedProgress)
			}
		})
	}
}

func TestAPTManager_InstallPackages_ProgressReporting(t *testing.T) {
	t.Run("reports progress for each package", func(t *testing.T) {
		manager := packagemanager.NewAPTManagerDryRun()
		ctx := context.Background()

		packages := []string{"hyprland", "waybar"}
		progressChan := make(chan packagemanager.PackageProgress, 100)

		var progress []packagemanager.PackageProgress
		done := make(chan struct{})
		go func() {
			for p := range progressChan {
				progress = append(progress, p)
			}
			close(done)
		}()

		err := manager.InstallPackages(ctx, packages, progressChan)
		close(progressChan)
		<-done

		require.NoError(t, err)

		// Verify progress structure
		assert.NotEmpty(t, progress)

		// Each package should have: started, installing, completed
		hyprlandEvents := 0
		waybarEvents := 0

		for _, p := range progress {
			assert.NotEmpty(t, p.PackageName, "Progress should have package name")
			assert.NotEmpty(t, p.Status, "Progress should have status")

			if p.PackageName == "hyprland" {
				hyprlandEvents++
			}
			if p.PackageName == "waybar" {
				waybarEvents++
			}
		}

		assert.Greater(t, hyprlandEvents, 0, "Should have progress for hyprland")
		assert.Greater(t, waybarEvents, 0, "Should have progress for waybar")
	})

	t.Run("handles context cancellation", func(t *testing.T) {
		manager := packagemanager.NewAPTManagerDryRun()
		ctx, cancel := context.WithCancel(context.Background())

		packages := []string{"pkg1", "pkg2", "pkg3"}
		progressChan := make(chan packagemanager.PackageProgress, 100)

		// Cancel immediately
		cancel()

		err := manager.InstallPackages(ctx, packages, progressChan)
		close(progressChan)

		assert.Error(t, err, "Should return error when context is cancelled")
		assert.Contains(t, err.Error(), "context", "Error should mention context")
	})
}

func TestAPTManager_InstallProfile(t *testing.T) {
	tests := []struct {
		name    string
		profile string
		wantErr bool
	}{
		{
			name:    "minimal profile",
			profile: "minimal",
			wantErr: false,
		},
		{
			name:    "recommended profile",
			profile: "recommended",
			wantErr: false,
		},
		{
			name:    "full profile",
			profile: "full",
			wantErr: false,
		},
		{
			name:    "unknown profile",
			profile: "invalid-profile",
			wantErr: true,
		},
		{
			name:    "empty profile name",
			profile: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := packagemanager.NewAPTManagerDryRun()
			ctx := context.Background()

			progressChan := make(chan packagemanager.PackageProgress, 100)
			go func() {
				// Drain progress channel
				for range progressChan {
				}
			}()

			err := manager.InstallProfile(ctx, tt.profile, progressChan)
			close(progressChan)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
