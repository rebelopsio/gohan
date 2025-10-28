package theme

// ThemeInfo represents theme information for presentation
type ThemeInfo struct {
	Name        string
	DisplayName string
	Author      string
	Description string
	Variant     string // "dark" or "light"
	IsActive    bool
	PreviewURL  string
	ColorScheme map[string]string
}

// ThemeProgress reports theme operation progress
type ThemeProgress struct {
	Component       string
	Status          string // "started", "processing", "completed", "failed"
	PercentComplete float64
	Error           error
}

// ThemePreview provides a preview of a theme without applying it
type ThemePreview struct {
	Name         string
	DisplayName  string
	Author       string
	Description  string
	Variant      string
	ColorScheme  map[string]string
	PreviewText  string // Visual representation
}
