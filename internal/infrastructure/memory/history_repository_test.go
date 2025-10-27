package memory_test

import (
	"context"
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/history"
	"github.com/rebelopsio/gohan/internal/infrastructure/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHistoryRepository(t *testing.T) {
	repo := memory.NewHistoryRepository()
	assert.NotNil(t, repo)
}

func TestHistoryRepository_Save(t *testing.T) {
	repo := memory.NewHistoryRepository()
	ctx := context.Background()

	record := createTestRecord(t, "success", 1)

	err := repo.Save(ctx, record)
	assert.NoError(t, err)

	// Verify record was saved
	found, err := repo.FindByID(ctx, record.ID())
	assert.NoError(t, err)
	assert.Equal(t, record.ID(), found.ID())
}

func TestHistoryRepository_Save_Overwrite(t *testing.T) {
	repo := memory.NewHistoryRepository()
	ctx := context.Background()

	record := createTestRecord(t, "success", 1)

	// Save once
	err := repo.Save(ctx, record)
	require.NoError(t, err)

	// Save again (should overwrite)
	err = repo.Save(ctx, record)
	assert.NoError(t, err)

	// Verify only one record exists
	count, err := repo.Count(ctx, history.NewRecordFilter())
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestHistoryRepository_FindByID(t *testing.T) {
	repo := memory.NewHistoryRepository()
	ctx := context.Background()

	record := createTestRecord(t, "success", 1)
	err := repo.Save(ctx, record)
	require.NoError(t, err)

	t.Run("existing record", func(t *testing.T) {
		found, err := repo.FindByID(ctx, record.ID())
		assert.NoError(t, err)
		assert.Equal(t, record.ID(), found.ID())
		assert.Equal(t, record.SessionID(), found.SessionID())
	})

	t.Run("non-existent record", func(t *testing.T) {
		nonExistentID, _ := history.NewRecordID()
		_, err := repo.FindByID(ctx, nonExistentID)
		assert.ErrorIs(t, err, history.ErrRecordNotFound)
	})
}

func TestHistoryRepository_FindAll(t *testing.T) {
	repo := memory.NewHistoryRepository()
	ctx := context.Background()

	// Save multiple records
	record1 := createTestRecordWithTime(t, "success", 1, time.Date(2025, 10, 26, 14, 0, 0, 0, time.UTC))
	record2 := createTestRecordWithTime(t, "failed", 0, time.Date(2025, 10, 25, 14, 0, 0, 0, time.UTC))
	record3 := createTestRecordWithTime(t, "success", 2, time.Date(2025, 10, 24, 14, 0, 0, 0, time.UTC))

	require.NoError(t, repo.Save(ctx, record1))
	require.NoError(t, repo.Save(ctx, record2))
	require.NoError(t, repo.Save(ctx, record3))

	t.Run("empty filter returns all", func(t *testing.T) {
		records, err := repo.FindAll(ctx, history.NewRecordFilter())
		assert.NoError(t, err)
		assert.Len(t, records, 3)
	})

	t.Run("filter by outcome", func(t *testing.T) {
		outcome, _ := history.NewInstallationOutcome("success")
		filter := history.NewRecordFilter().WithOutcome(outcome)

		records, err := repo.FindAll(ctx, filter)
		assert.NoError(t, err)
		assert.Len(t, records, 2)
	})

	t.Run("filter by period", func(t *testing.T) {
		period, _ := history.NewInstallationPeriod(
			time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
			time.Date(2025, 10, 27, 0, 0, 0, 0, time.UTC),
		)
		filter := history.NewRecordFilter().WithPeriod(period)

		records, err := repo.FindAll(ctx, filter)
		assert.NoError(t, err)
		assert.Len(t, records, 2) // record1 and record2
	})
}

func TestHistoryRepository_FindRecent(t *testing.T) {
	repo := memory.NewHistoryRepository()
	ctx := context.Background()

	// Save records with different recorded times
	now := time.Now()
	record1 := createTestRecordAtRecordedTime(t, "success", 1, now.Add(-3*time.Hour))
	record2 := createTestRecordAtRecordedTime(t, "success", 1, now.Add(-2*time.Hour))
	record3 := createTestRecordAtRecordedTime(t, "success", 1, now.Add(-1*time.Hour))

	require.NoError(t, repo.Save(ctx, record1))
	require.NoError(t, repo.Save(ctx, record2))
	require.NoError(t, repo.Save(ctx, record3))

	t.Run("limit 2 returns most recent", func(t *testing.T) {
		records, err := repo.FindRecent(ctx, 2)
		assert.NoError(t, err)
		assert.Len(t, records, 2)

		// Should be ordered newest first
		assert.True(t, records[0].RecordedAt().After(records[1].RecordedAt()))
		assert.Equal(t, record3.ID(), records[0].ID())
		assert.Equal(t, record2.ID(), records[1].ID())
	})

	t.Run("limit exceeds total returns all", func(t *testing.T) {
		records, err := repo.FindRecent(ctx, 10)
		assert.NoError(t, err)
		assert.Len(t, records, 3)
	})

	t.Run("zero limit returns empty", func(t *testing.T) {
		records, err := repo.FindRecent(ctx, 0)
		assert.NoError(t, err)
		assert.Len(t, records, 0)
	})
}

func TestHistoryRepository_Count(t *testing.T) {
	repo := memory.NewHistoryRepository()
	ctx := context.Background()

	record1 := createTestRecord(t, "success", 1)
	record2 := createTestRecord(t, "failed", 0)

	require.NoError(t, repo.Save(ctx, record1))
	require.NoError(t, repo.Save(ctx, record2))

	t.Run("empty filter counts all", func(t *testing.T) {
		count, err := repo.Count(ctx, history.NewRecordFilter())
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("filter counts matching", func(t *testing.T) {
		outcome, _ := history.NewInstallationOutcome("success")
		filter := history.NewRecordFilter().WithOutcome(outcome)

		count, err := repo.Count(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}

func TestHistoryRepository_Delete(t *testing.T) {
	repo := memory.NewHistoryRepository()
	ctx := context.Background()

	record := createTestRecord(t, "success", 1)
	require.NoError(t, repo.Save(ctx, record))

	t.Run("delete existing record", func(t *testing.T) {
		err := repo.Delete(ctx, record.ID())
		assert.NoError(t, err)

		// Verify record is gone
		_, err = repo.FindByID(ctx, record.ID())
		assert.ErrorIs(t, err, history.ErrRecordNotFound)
	})

	t.Run("delete non-existent record", func(t *testing.T) {
		nonExistentID, _ := history.NewRecordID()
		err := repo.Delete(ctx, nonExistentID)
		assert.ErrorIs(t, err, history.ErrRecordNotFound)
	})
}

func TestHistoryRepository_PurgeOlderThan(t *testing.T) {
	repo := memory.NewHistoryRepository()
	ctx := context.Background()

	now := time.Now()

	// Save records with different recorded times
	oldRecord1 := createTestRecordAtRecordedTime(t, "success", 1, now.AddDate(0, 0, -100))
	oldRecord2 := createTestRecordAtRecordedTime(t, "success", 1, now.AddDate(0, 0, -95))
	recentRecord := createTestRecordAtRecordedTime(t, "success", 1, now.AddDate(0, 0, -10))

	require.NoError(t, repo.Save(ctx, oldRecord1))
	require.NoError(t, repo.Save(ctx, oldRecord2))
	require.NoError(t, repo.Save(ctx, recentRecord))

	// Purge records older than 90 days
	cutoffDate := now.AddDate(0, 0, -90)
	deleted, err := repo.PurgeOlderThan(ctx, cutoffDate)
	assert.NoError(t, err)
	assert.Equal(t, 2, deleted)

	// Verify only recent record remains
	count, err := repo.Count(ctx, history.NewRecordFilter())
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

// Helper functions

func createTestRecord(t *testing.T, outcomeStr string, packageCount int) history.InstallationRecord {
	return createTestRecordWithTime(t, outcomeStr, packageCount, time.Now())
}

func createTestRecordWithTime(t *testing.T, outcomeStr string, packageCount int, installedAt time.Time) history.InstallationRecord {
	completedAt := installedAt.Add(time.Minute)

	var packages []history.InstalledPackage
	for i := 0; i < packageCount; i++ {
		pkg, _ := history.NewInstalledPackage("test-package", "1.0.0", 1024)
		packages = append(packages, pkg)
	}

	metadata, _ := history.NewInstallationMetadata(
		"test-package",
		"1.0.0",
		installedAt,
		completedAt,
		packages,
	)

	systemCtx, _ := history.NewSystemContext("OS", "", "", "")
	outcome, _ := history.NewInstallationOutcome(outcomeStr)

	var failureDetails *history.FailureDetails
	if outcomeStr == "failed" {
		fd, _ := history.NewFailureDetails("test failure", completedAt, "test", "ERR")
		failureDetails = &fd
	}

	record, err := history.NewInstallationRecord(
		"session-123",
		outcome,
		metadata,
		systemCtx,
		failureDetails,
		time.Now(),
	)
	require.NoError(t, err)

	return record
}

func createTestRecordAtRecordedTime(t *testing.T, outcomeStr string, packageCount int, recordedAt time.Time) history.InstallationRecord {
	installedAt := recordedAt.Add(-time.Minute)
	completedAt := recordedAt

	var packages []history.InstalledPackage
	for i := 0; i < packageCount; i++ {
		pkg, _ := history.NewInstalledPackage("test-package", "1.0.0", 1024)
		packages = append(packages, pkg)
	}

	metadata, _ := history.NewInstallationMetadata(
		"test-package",
		"1.0.0",
		installedAt,
		completedAt,
		packages,
	)

	systemCtx, _ := history.NewSystemContext("OS", "", "", "")
	outcome, _ := history.NewInstallationOutcome(outcomeStr)

	var failureDetails *history.FailureDetails
	if outcomeStr == "failed" {
		fd, _ := history.NewFailureDetails("test failure", completedAt, "test", "ERR")
		failureDetails = &fd
	}

	record, err := history.NewInstallationRecord(
		"session-123",
		outcome,
		metadata,
		systemCtx,
		failureDetails,
		recordedAt,
	)
	require.NoError(t, err)

	return record
}
