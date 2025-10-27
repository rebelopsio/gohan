package history

import (
	"time"
)

// InstallationPeriod is a value object representing a time range for queries
type InstallationPeriod struct {
	start time.Time
	end   time.Time
}

// NewInstallationPeriod creates a period value object
func NewInstallationPeriod(start, end time.Time) (InstallationPeriod, error) {
	if start.IsZero() || end.IsZero() {
		return InstallationPeriod{}, ErrInvalidPeriod
	}

	if end.Before(start) {
		return InstallationPeriod{}, ErrInvalidTimeRange
	}

	return InstallationPeriod{
		start: start,
		end:   end,
	}, nil
}

// NewPeriodFromDaysAgo creates a period from N days ago until now
func NewPeriodFromDaysAgo(days int) (InstallationPeriod, error) {
	if days < 0 {
		return InstallationPeriod{}, ErrInvalidPeriod
	}

	end := time.Now()
	start := end.AddDate(0, 0, -days)

	return InstallationPeriod{
		start: start,
		end:   end,
	}, nil
}

// Start returns the period start time
func (p InstallationPeriod) Start() time.Time {
	return p.start
}

// End returns the period end time
func (p InstallationPeriod) End() time.Time {
	return p.end
}

// Contains checks if a time falls within this period (inclusive)
func (p InstallationPeriod) Contains(t time.Time) bool {
	return (t.Equal(p.start) || t.After(p.start)) &&
		(t.Equal(p.end) || t.Before(p.end))
}

// Duration returns the length of this period
func (p InstallationPeriod) Duration() time.Duration {
	return p.end.Sub(p.start)
}
