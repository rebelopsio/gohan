package theme

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

var (
	// ErrInvalidTheme indicates theme validation failure
	ErrInvalidTheme = errors.New("invalid theme")
	// ErrInvalidColor indicates color validation failure
	ErrInvalidColor = errors.New("invalid color format")
)

// ThemeName represents a unique theme identifier
type ThemeName string

// Standard theme names
const (
	ThemeMocha      ThemeName = "mocha"
	ThemeLatte      ThemeName = "latte"
	ThemeFrappe     ThemeName = "frappe"
	ThemeMacchiato  ThemeName = "macchiato"
	ThemeGohan      ThemeName = "gohan"
)

// ThemeVariant indicates if a theme is suitable for day or night use
type ThemeVariant string

const (
	ThemeVariantDark  ThemeVariant = "dark"
	ThemeVariantLight ThemeVariant = "light"
)

// Color represents a hex color code
type Color string

// hexColorRegex matches valid hex color codes (#RRGGBB)
var hexColorRegex = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

// Validate checks if the color is a valid hex color
func (c Color) Validate() error {
	if c == "" {
		return fmt.Errorf("%w: color cannot be empty", ErrInvalidColor)
	}
	if !hexColorRegex.MatchString(string(c)) {
		return fmt.Errorf("%w: must be in format #RRGGBB", ErrInvalidColor)
	}
	return nil
}

// String returns the color as a string
func (c Color) String() string {
	return string(c)
}

// ThemeMetadata contains theme information and properties
type ThemeMetadata struct {
	displayName string
	author      string
	description string
	variant     ThemeVariant
	previewURL  string
}

// DisplayName returns the human-readable theme name
func (m ThemeMetadata) DisplayName() string {
	return m.displayName
}

// Author returns the theme creator
func (m ThemeMetadata) Author() string {
	return m.author
}

// Description returns the theme description
func (m ThemeMetadata) Description() string {
	return m.description
}

// Variant returns the theme variant (dark/light)
func (m ThemeMetadata) Variant() ThemeVariant {
	return m.variant
}

// PreviewURL returns the URL to a theme preview image
func (m ThemeMetadata) PreviewURL() string {
	return m.previewURL
}

// IsDark returns true if this is a dark theme
func (m ThemeMetadata) IsDark() bool {
	return m.variant == ThemeVariantDark
}

// IsLight returns true if this is a light theme
func (m ThemeMetadata) IsLight() bool {
	return m.variant == ThemeVariantLight
}

// ColorScheme defines all colors in a theme
type ColorScheme struct {
	// Base colors
	base    Color
	surface Color
	overlay Color
	text    Color
	subtext Color

	// Accent colors (Catppuccin palette)
	rosewater Color
	flamingo  Color
	pink      Color
	mauve     Color
	red       Color
	maroon    Color
	peach     Color
	yellow    Color
	green     Color
	teal      Color
	sky       Color
	sapphire  Color
	blue      Color
	lavender  Color
}

// Base returns the base background color
func (cs ColorScheme) Base() Color {
	return cs.base
}

// Surface returns the surface color
func (cs ColorScheme) Surface() Color {
	return cs.surface
}

// Overlay returns the overlay color
func (cs ColorScheme) Overlay() Color {
	return cs.overlay
}

// Text returns the primary text color
func (cs ColorScheme) Text() Color {
	return cs.text
}

// Subtext returns the secondary text color
func (cs ColorScheme) Subtext() Color {
	return cs.subtext
}

// Rosewater returns the rosewater accent color
func (cs ColorScheme) Rosewater() Color {
	return cs.rosewater
}

// Flamingo returns the flamingo accent color
func (cs ColorScheme) Flamingo() Color {
	return cs.flamingo
}

// Pink returns the pink accent color
func (cs ColorScheme) Pink() Color {
	return cs.pink
}

// Mauve returns the mauve accent color
func (cs ColorScheme) Mauve() Color {
	return cs.mauve
}

// Red returns the red accent color
func (cs ColorScheme) Red() Color {
	return cs.red
}

// Maroon returns the maroon accent color
func (cs ColorScheme) Maroon() Color {
	return cs.maroon
}

// Peach returns the peach accent color
func (cs ColorScheme) Peach() Color {
	return cs.peach
}

// Yellow returns the yellow accent color
func (cs ColorScheme) Yellow() Color {
	return cs.yellow
}

// Green returns the green accent color
func (cs ColorScheme) Green() Color {
	return cs.green
}

// Teal returns the teal accent color
func (cs ColorScheme) Teal() Color {
	return cs.teal
}

// Sky returns the sky accent color
func (cs ColorScheme) Sky() Color {
	return cs.sky
}

// Sapphire returns the sapphire accent color
func (cs ColorScheme) Sapphire() Color {
	return cs.sapphire
}

// Blue returns the blue accent color
func (cs ColorScheme) Blue() Color {
	return cs.blue
}

// Lavender returns the lavender accent color
func (cs ColorScheme) Lavender() Color {
	return cs.lavender
}

// Theme represents a complete visual theme
type Theme struct {
	name        ThemeName
	metadata    ThemeMetadata
	colorScheme ColorScheme
	createdAt   time.Time
}

// NewTheme creates a new theme with validation
func NewTheme(name ThemeName, metadata ThemeMetadata, colorScheme ColorScheme) (*Theme, error) {
	// Validate theme name
	if name == "" {
		return nil, fmt.Errorf("%w: theme name cannot be empty", ErrInvalidTheme)
	}

	// Validate metadata
	if metadata.displayName == "" {
		return nil, fmt.Errorf("%w: display name cannot be empty", ErrInvalidTheme)
	}
	if metadata.author == "" {
		return nil, fmt.Errorf("%w: author cannot be empty", ErrInvalidTheme)
	}
	if metadata.variant != ThemeVariantDark && metadata.variant != ThemeVariantLight {
		return nil, fmt.Errorf("%w: variant must be 'dark' or 'light'", ErrInvalidTheme)
	}

	return &Theme{
		name:        name,
		metadata:    metadata,
		colorScheme: colorScheme,
		createdAt:   time.Now(),
	}, nil
}

// Name returns the theme name
func (t *Theme) Name() ThemeName {
	return t.name
}

// DisplayName returns the human-readable theme name
func (t *Theme) DisplayName() string {
	return t.metadata.displayName
}

// Author returns the theme creator
func (t *Theme) Author() string {
	return t.metadata.author
}

// Description returns the theme description
func (t *Theme) Description() string {
	return t.metadata.description
}

// Variant returns the theme variant
func (t *Theme) Variant() ThemeVariant {
	return t.metadata.variant
}

// PreviewURL returns the preview URL
func (t *Theme) PreviewURL() string {
	return t.metadata.previewURL
}

// ColorScheme returns the theme's color scheme
func (t *Theme) ColorScheme() ColorScheme {
	return t.colorScheme
}

// CreatedAt returns when the theme was created
func (t *Theme) CreatedAt() time.Time {
	return t.createdAt
}

// IsDark returns true if this is a dark theme
func (t *Theme) IsDark() bool {
	return t.metadata.IsDark()
}

// IsLight returns true if this is a light theme
func (t *Theme) IsLight() bool {
	return t.metadata.IsLight()
}

// Equals checks if two themes are equal
func (t *Theme) Equals(other *Theme) bool {
	if other == nil {
		return false
	}
	return t.name == other.name &&
		t.metadata.displayName == other.metadata.displayName &&
		t.metadata.author == other.metadata.author &&
		t.metadata.variant == other.metadata.variant &&
		t.createdAt.Equal(other.createdAt)
}
