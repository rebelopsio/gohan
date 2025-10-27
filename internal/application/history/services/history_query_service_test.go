package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/application/history/services"
	"github.com/rebelopsio/gohan/internal/domain/history"
	"github.com/rebelopsio/gohan/internal/infrastructure/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHistoryQueryService_ListRecords(t *testing.T) {
	t.Run("lists all records with no filter", func(t *testing.T) {
		repo := memory.NewHistoryRepository()
		service := services.NewHistoryQueryService(repo)
		ctx := context.Background()

		// Create and save test records
		record1 := createSuccessRecord(t, "hyprland", time.Now().Add(-2*time.Hour))
		record2 := createSuccessRecord(t, "waybar", time.Now().Add(-1*time.Hour))
		require.NoError(t, repo.Save(ctx, record1))
		require.NoError(t, repo.Save(ctx, record2))

		// Query all records
		records, err := service.ListRecords(ctx, history.NewRecordFilter())

		require.NoError(t, err)
		assert.Len(t, records, 2)
	})

	t.Run("filters by outcome", func(t *testing.T) {
		repo := memory.NewHistoryRepository()
		service := services.NewHistoryQueryService(repo)
		ctx := context.Background()

		// Create mixed records
		success := createSuccessRecord(t, "hyprland", time.Now())
		failed := createFailedRecord(t, "waybar", time.Now())
		require.NoError(t, repo.Save(ctx, success))
		require.NoError(t, repo.Save(ctx, failed))

		// Filter for successful only
		outcome, _ := history.NewInstallationOutcome("success")
		filter := history.NewRecordFilter().WithOutcome(outcome)
		records, err := service.ListRecords(ctx, filter)

		require.NoError(t, err)
		assert.Len(t, records, 1)
		assert.True(t, records[0].WasSuccessful())
	})

	t.Run("returns most recent records first", func(t *testing.T) {
		repo := memory.NewHistoryRepository()
		service := services.NewHistoryQueryService(repo)
		ctx := context.Background()

		// Create records with different times
		old := createSuccessRecordAtTime(t, "hyprland", time.Now().Add(-3*time.Hour))
		middle := createSuccessRecordAtTime(t, "waybar", time.Now().Add(-2*time.Hour))
		recent := createSuccessRecordAtTime(t, "kitty", time.Now().Add(-1*time.Hour))

		require.NoError(t, repo.Save(ctx, old))
		require.NoError(t, repo.Save(ctx, middle))
		require.NoError(t, repo.Save(ctx, recent))

		records, err := service.ListRecent(ctx, 10)

		require.NoError(t, err)
		assert.Len(t, records, 3)
		// Should be sorted newest first
		assert.Equal(t, "kitty", records[0].PackageName())
		assert.Equal(t, "waybar", records[1].PackageName())
		assert.Equal(t, "hyprland", records[2].PackageName())
	})
}

func TestHistoryQueryService_GetRecordByID(t *testing.T) {
	t.Run("retrieves record by ID", func(t *testing.T) {
		repo := memory.NewHistoryRepository()
		service := services.NewHistoryQueryService(repo)
		ctx := context.Background()

		record := createSuccessRecord(t, "hyprland", time.Now())
		require.NoError(t, repo.Save(ctx, record))

		retrieved, err := service.GetRecordByID(ctx, record.ID())

		require.NoError(t, err)
		assert.Equal(t, record.ID(), retrieved.ID())
		assert.Equal(t, record.PackageName(), retrieved.PackageName())
	})

	t.Run("returns error for non-existent ID", func(t *testing.T) {
		repo := memory.NewHistoryRepository()
		service := services.NewHistoryQueryService(repo)
		ctx := context.Background()

		nonExistentID, _ := history.NewRecordID()
		_, err := service.GetRecordByID(ctx, nonExistentID)

		assert.ErrorIs(t, err, history.ErrRecordNotFound)
	})
}

func TestHistoryQueryService_ListRecent(t *testing.T) {
	t.Run("limits results to specified count", func(t *testing.T) {
		repo := memory.NewHistoryRepository()
		service := services.NewHistoryQueryService(repo)
		ctx := context.Background()

		// Create 10 records
		for i := 0; i < 10; i++ {
			record := createSuccessRecord(t, "package", time.Now().Add(time.Duration(-i)*time.Hour))
			require.NoError(t, repo.Save(ctx, record))
		}

		// Request only 5
		records, err := service.ListRecent(ctx, 5)

		require.NoError(t, err)
		assert.Len(t, records, 5)
	})

	t.Run("returns empty list when limit is zero", func(t *testing.T) {
		repo := memory.NewHistoryRepository()
		service := services.NewHistoryQueryService(repo)
		ctx := context.Background()

		records, err := service.ListRecent(ctx, 0)

		require.NoError(t, err)
		assert.Empty(t, records)
	})
}

func TestHistoryQueryService_CountRecords(t *testing.T) {
	t.Run("counts all records", func(t *testing.T) {
		repo := memory.NewHistoryRepository()
		service := services.NewHistoryQueryService(repo)
		ctx := context.Background()

		record1 := createSuccessRecord(t, "hyprland", time.Now())
		record2 := createSuccessRecord(t, "waybar", time.Now())
		require.NoError(t, repo.Save(ctx, record1))
		require.NoError(t, repo.Save(ctx, record2))

		count, err := service.CountRecords(ctx, history.NewRecordFilter())

		require.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("counts with filter", func(t *testing.T) {
		repo := memory.NewHistoryRepository()
		service := services.NewHistoryQueryService(repo)
		ctx := context.Background()

		success1 := createSuccessRecord(t, "hyprland", time.Now())
		success2 := createSuccessRecord(t, "waybar", time.Now())
		failed := createFailedRecord(t, "kitty", time.Now())
		require.NoError(t, repo.Save(ctx, success1))
		require.NoError(t, repo.Save(ctx, success2))
		require.NoError(t, repo.Save(ctx, failed))

		outcome, _ := history.NewInstallationOutcome("success")
		filter := history.NewRecordFilter().WithOutcome(outcome)
		count, err := service.CountRecords(ctx, filter)

		require.NoError(t, err)
		assert.Equal(t, 2, count)
	})
}

// Helper functions

func createSuccessRecord(t *testing.T, packageName string, installedAt time.Time) history.InstallationRecord {
	return createSuccessRecordAtTime(t, packageName, installedAt)
}

func createSuccessRecordAtTime(t *testing.T, packageName string, installedAt time.Time) history.InstallationRecord {
	completedAt := installedAt.Add(time.Minute)

	pkg, _ := history.NewInstalledPackage(packageName, "1.0.0", 1024)
	packages := []history.InstalledPackage{pkg}

	metadata, _ := history.NewInstallationMetadata(
		packageName,
		"1.0.0",
		installedAt,
		completedAt,
		packages,
	)

	sysCtx, _ := history.NewSystemContext("Debian GNU/Linux", "6.1.0", "1.0.0", "testhost")
	outcome, _ := history.NewInstallationOutcome("success")

	record, err := history.NewInstallationRecord(
		"session-123",
		outcome,
		metadata,
		sysCtx,
		nil,
		completedAt,
	)
	require.NoError(t, err)

	return record
}

func createFailedRecord(t *testing.T, packageName string, installedAt time.Time) history.InstallationRecord {
	completedAt := installedAt.Add(time.Minute)

	pkg, _ := history.NewInstalledPackage(packageName, "1.0.0", 1024)
	packages := []history.InstalledPackage{pkg}

	metadata, _ := history.NewInstallationMetadata(
		packageName,
		"1.0.0",
		installedAt,
		completedAt,
		packages,
	)

	sysCtx, _ := history.NewSystemContext("Debian GNU/Linux", "6.1.0", "1.0.0", "testhost")
	outcome, _ := history.NewInstallationOutcome("failed")

	failureDetails, _ := history.NewFailureDetails(
		"Installation failed",
		completedAt,
		"installing",
		"ERR001",
	)

	record, err := history.NewInstallationRecord(
		"session-456",
		outcome,
		metadata,
		sysCtx,
		&failureDetails,
		completedAt,
	)
	require.NoError(t, err)

	return record
}
