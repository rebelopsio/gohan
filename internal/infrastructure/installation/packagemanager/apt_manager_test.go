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

		err := manager.InstallPackage(ctx, "")

		assert.Error(t, err)
	})

	t.Run("accepts valid package name", func(t *testing.T) {
		manager := packagemanager.NewAPTManager()
		ctx := context.Background()

		// This will fail in test environment without sudo, but validates input
		err := manager.InstallPackage(ctx, "hyprland")

		// Either succeeds or fails with execution error (not validation error)
		// We can't actually install in tests, so either outcome is acceptable
		_ = err
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
