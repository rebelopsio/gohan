package theme

// InitializeStandardThemes creates and registers all standard themes
func InitializeStandardThemes(registry ThemeRegistry) error {
	themes := []*Theme{
		createMochaTheme(),
		createLatteTheme(),
		createFrappeTheme(),
		createMacchiatoTheme(),
		createGohanTheme(),
	}

	for _, theme := range themes {
		if err := registry.Register(theme); err != nil {
			return err
		}
	}

	return nil
}

// createMochaTheme creates the Catppuccin Mocha theme (warm dark)
func createMochaTheme() *Theme {
	metadata := ThemeMetadata{
		displayName: "Catppuccin Mocha",
		author:      "Catppuccin",
		description: "Soothing pastel theme for the high-spirited!",
		variant:     ThemeVariantDark,
		previewURL:  "https://github.com/catppuccin/catppuccin",
	}

	colorScheme := ColorScheme{
		// Base colors
		base:    Color("#1e1e2e"),
		surface: Color("#313244"),
		overlay: Color("#45475a"),
		text:    Color("#cdd6f4"),
		subtext: Color("#bac2de"),

		// Accent colors
		rosewater: Color("#f5e0dc"),
		flamingo:  Color("#f2cdcd"),
		pink:      Color("#f5c2e7"),
		mauve:     Color("#cba6f7"),
		red:       Color("#f38ba8"),
		maroon:    Color("#eba0ac"),
		peach:     Color("#fab387"),
		yellow:    Color("#f9e2af"),
		green:     Color("#a6e3a1"),
		teal:      Color("#94e2d5"),
		sky:       Color("#89dceb"),
		sapphire:  Color("#74c7ec"),
		blue:      Color("#89b4fa"),
		lavender:  Color("#b4befe"),
	}

	theme, _ := NewTheme(ThemeMocha, metadata, colorScheme)
	return theme
}

// createLatteTheme creates the Catppuccin Latte theme (light)
func createLatteTheme() *Theme {
	metadata := ThemeMetadata{
		displayName: "Catppuccin Latte",
		author:      "Catppuccin",
		description: "Soothing pastel theme for the high-spirited!",
		variant:     ThemeVariantLight,
		previewURL:  "https://github.com/catppuccin/catppuccin",
	}

	colorScheme := ColorScheme{
		// Base colors
		base:    Color("#eff1f5"),
		surface: Color("#e6e9ef"),
		overlay: Color("#dce0e8"),
		text:    Color("#4c4f69"),
		subtext: Color("#5c5f77"),

		// Accent colors
		rosewater: Color("#dc8a78"),
		flamingo:  Color("#dd7878"),
		pink:      Color("#ea76cb"),
		mauve:     Color("#8839ef"),
		red:       Color("#d20f39"),
		maroon:    Color("#e64553"),
		peach:     Color("#fe640b"),
		yellow:    Color("#df8e1d"),
		green:     Color("#40a02b"),
		teal:      Color("#179299"),
		sky:       Color("#04a5e5"),
		sapphire:  Color("#209fb5"),
		blue:      Color("#1e66f5"),
		lavender:  Color("#7287fd"),
	}

	theme, _ := NewTheme(ThemeLatte, metadata, colorScheme)
	return theme
}

// createFrappeTheme creates the Catppuccin Frappe theme (muted dark)
func createFrappeTheme() *Theme {
	metadata := ThemeMetadata{
		displayName: "Catppuccin Frappe",
		author:      "Catppuccin",
		description: "Soothing pastel theme for the high-spirited!",
		variant:     ThemeVariantDark,
		previewURL:  "https://github.com/catppuccin/catppuccin",
	}

	colorScheme := ColorScheme{
		// Base colors
		base:    Color("#303446"),
		surface: Color("#414559"),
		overlay: Color("#51576d"),
		text:    Color("#c6d0f5"),
		subtext: Color("#b5bfe2"),

		// Accent colors
		rosewater: Color("#f2d5cf"),
		flamingo:  Color("#eebebe"),
		pink:      Color("#f4b8e4"),
		mauve:     Color("#ca9ee6"),
		red:       Color("#e78284"),
		maroon:    Color("#ea999c"),
		peach:     Color("#ef9f76"),
		yellow:    Color("#e5c890"),
		green:     Color("#a6d189"),
		teal:      Color("#81c8be"),
		sky:       Color("#99d1db"),
		sapphire:  Color("#85c1dc"),
		blue:      Color("#8caaee"),
		lavender:  Color("#babbf1"),
	}

	theme, _ := NewTheme(ThemeFrappe, metadata, colorScheme)
	return theme
}

// createMacchiatoTheme creates the Catppuccin Macchiato theme (warmer dark)
func createMacchiatoTheme() *Theme {
	metadata := ThemeMetadata{
		displayName: "Catppuccin Macchiato",
		author:      "Catppuccin",
		description: "Soothing pastel theme for the high-spirited!",
		variant:     ThemeVariantDark,
		previewURL:  "https://github.com/catppuccin/catppuccin",
	}

	colorScheme := ColorScheme{
		// Base colors
		base:    Color("#24273a"),
		surface: Color("#363a4f"),
		overlay: Color("#494d64"),
		text:    Color("#cad3f5"),
		subtext: Color("#b8c0e0"),

		// Accent colors
		rosewater: Color("#f4dbd6"),
		flamingo:  Color("#f0c6c6"),
		pink:      Color("#f5bde6"),
		mauve:     Color("#c6a0f6"),
		red:       Color("#ed8796"),
		maroon:    Color("#ee99a0"),
		peach:     Color("#f5a97f"),
		yellow:    Color("#eed49f"),
		green:     Color("#a6da95"),
		teal:      Color("#8bd5ca"),
		sky:       Color("#91d7e3"),
		sapphire:  Color("#7dc4e4"),
		blue:      Color("#8aadf4"),
		lavender:  Color("#b7bdf8"),
	}

	theme, _ := NewTheme(ThemeMacchiato, metadata, colorScheme)
	return theme
}

// createGohanTheme creates the custom Gohan theme (dark brand theme)
func createGohanTheme() *Theme {
	metadata := ThemeMetadata{
		displayName: "Gohan",
		author:      "Gohan Team",
		description: "Custom brand theme for Gohan",
		variant:     ThemeVariantDark,
		previewURL:  "",
	}

	// Gohan uses Mocha as base but could be customized
	colorScheme := ColorScheme{
		// Base colors
		base:    Color("#1e1e2e"),
		surface: Color("#313244"),
		overlay: Color("#45475a"),
		text:    Color("#cdd6f4"),
		subtext: Color("#bac2de"),

		// Accent colors with custom Gohan branding
		rosewater: Color("#f5e0dc"),
		flamingo:  Color("#f2cdcd"),
		pink:      Color("#f5c2e7"),
		mauve:     Color("#cba6f7"),
		red:       Color("#f38ba8"),
		maroon:    Color("#eba0ac"),
		peach:     Color("#fab387"),
		yellow:    Color("#f9e2af"),
		green:     Color("#a6e3a1"),
		teal:      Color("#94e2d5"),
		sky:       Color("#89dceb"),
		sapphire:  Color("#74c7ec"),
		blue:      Color("#89b4fa"),
		lavender:  Color("#b4befe"),
	}

	theme, _ := NewTheme(ThemeGohan, metadata, colorScheme)
	return theme
}
