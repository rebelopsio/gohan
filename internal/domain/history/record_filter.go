package history

import "strings"

// RecordFilter is a value object for querying installation records
type RecordFilter struct {
	period      *InstallationPeriod
	outcome     *InstallationOutcome
	packageName string
}

// NewRecordFilter creates an empty filter that matches all records
func NewRecordFilter() RecordFilter {
	return RecordFilter{}
}

// WithPeriod adds a time period filter
func (f RecordFilter) WithPeriod(period InstallationPeriod) RecordFilter {
	f.period = &period
	return f
}

// WithOutcome adds an outcome filter
func (f RecordFilter) WithOutcome(outcome InstallationOutcome) RecordFilter {
	f.outcome = &outcome
	return f
}

// WithPackageName adds a package name filter
// Empty or whitespace-only names are ignored
func (f RecordFilter) WithPackageName(name string) RecordFilter {
	name = strings.TrimSpace(name)
	if name != "" {
		f.packageName = name
	}
	return f
}

// HasPeriodFilter returns true if a period filter is set
func (f RecordFilter) HasPeriodFilter() bool {
	return f.period != nil
}

// HasOutcomeFilter returns true if an outcome filter is set
func (f RecordFilter) HasOutcomeFilter() bool {
	return f.outcome != nil
}

// HasPackageFilter returns true if a package name filter is set
func (f RecordFilter) HasPackageFilter() bool {
	return f.packageName != ""
}

// IsEmpty returns true if no filters are set
func (f RecordFilter) IsEmpty() bool {
	return !f.HasPeriodFilter() && !f.HasOutcomeFilter() && !f.HasPackageFilter()
}

// Period returns the period filter if set
func (f RecordFilter) Period() InstallationPeriod {
	if f.period == nil {
		return InstallationPeriod{}
	}
	return *f.period
}

// Outcome returns the outcome filter if set
func (f RecordFilter) Outcome() InstallationOutcome {
	if f.outcome == nil {
		return InstallationOutcome("")
	}
	return *f.outcome
}

// PackageName returns the package name filter
func (f RecordFilter) PackageName() string {
	return f.packageName
}

// MatchesMetadata returns true if the metadata matches all set filters
func (f RecordFilter) MatchesMetadata(metadata InstallationMetadata) bool {
	// Empty filter matches everything
	if f.IsEmpty() {
		return true
	}

	// Check period filter
	if f.HasPeriodFilter() {
		if !f.period.Contains(metadata.InstalledAt()) {
			return false
		}
	}

	// Check package filter
	if f.HasPackageFilter() {
		if !metadata.HasPackage(f.packageName) {
			return false
		}
	}

	return true
}

// MatchesOutcome returns true if the outcome matches the filter
func (f RecordFilter) MatchesOutcome(outcome InstallationOutcome) bool {
	// Empty filter or no outcome filter matches everything
	if !f.HasOutcomeFilter() {
		return true
	}

	return f.outcome.Equals(outcome)
}
