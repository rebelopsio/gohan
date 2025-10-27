package history

import (
	"time"
)

// RetentionPolicy is a value object defining how long to keep records
type RetentionPolicy struct {
	retentionDays int
}

// NewRetentionPolicy creates a retention policy value object
func NewRetentionPolicy(days int) (RetentionPolicy, error) {
	if days < 1 {
		return RetentionPolicy{}, ErrInvalidRetentionPeriod
	}

	if days > MaxRetentionDays {
		return RetentionPolicy{}, ErrRetentionPeriodTooLong
	}

	return RetentionPolicy{
		retentionDays: days,
	}, nil
}

// DefaultRetentionPolicy returns the default 90-day policy
func DefaultRetentionPolicy() RetentionPolicy {
	return RetentionPolicy{retentionDays: DefaultRetentionDays}
}

// RetentionDays returns the retention period in days
func (r RetentionPolicy) RetentionDays() int {
	return r.retentionDays
}

// ShouldPurge checks if a record from the given time should be purged
func (r RetentionPolicy) ShouldPurge(recordTime time.Time) bool {
	cutoffDate := time.Now().AddDate(0, 0, -r.retentionDays)
	return recordTime.Before(cutoffDate)
}

// CutoffDate returns the earliest date to keep records
func (r RetentionPolicy) CutoffDate() time.Time {
	return time.Now().AddDate(0, 0, -r.retentionDays)
}
