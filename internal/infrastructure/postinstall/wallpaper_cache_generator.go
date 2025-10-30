package postinstall

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rebelopsio/gohan/internal/domain/postinstall"
)

// WallpaperCacheGenerator generates wallpaper cache
type WallpaperCacheGenerator struct {
	wallpaperDir string
	cacheDir     string
}

// NewWallpaperCacheGenerator creates a new wallpaper cache generator
func NewWallpaperCacheGenerator(wallpaperDir string) *WallpaperCacheGenerator {
	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".cache/gohan/wallpapers")

	return &WallpaperCacheGenerator{
		wallpaperDir: wallpaperDir,
		cacheDir:     cacheDir,
	}
}

// Name returns the installer name
func (g *WallpaperCacheGenerator) Name() string {
	return "Wallpaper Cache"
}

// Component returns the component type
func (g *WallpaperCacheGenerator) Component() postinstall.ComponentType {
	return postinstall.ComponentWallpaper
}

// Install performs the cache generation
func (g *WallpaperCacheGenerator) Install(ctx context.Context) (postinstall.ComponentResult, error) {
	result := postinstall.NewComponentResult(
		postinstall.ComponentWallpaper,
		postinstall.StatusInProgress,
		"Generating wallpaper cache",
	)

	details := []string{}

	// Ensure wallpaper directory exists
	if _, err := os.Stat(g.wallpaperDir); os.IsNotExist(err) {
		return result.
			WithDetails("Wallpaper directory does not exist: " + g.wallpaperDir).
			Complete(postinstall.StatusSkipped), nil
	}

	// Create cache directory
	if err := os.MkdirAll(g.cacheDir, 0755); err != nil {
		return postinstall.NewComponentResultWithError(
			postinstall.ComponentWallpaper,
			"Failed to create cache directory",
			err,
		), err
	}
	details = append(details, "Cache directory created: "+g.cacheDir)

	// Scan wallpaper directory
	files, err := g.scanWallpapers()
	if err != nil {
		return postinstall.NewComponentResultWithError(
			postinstall.ComponentWallpaper,
			"Failed to scan wallpaper directory",
			err,
		), err
	}

	details = append(details, fmt.Sprintf("Found %d wallpapers", len(files)))

	// Generate cache file with wallpaper list
	cacheFile := filepath.Join(g.cacheDir, "wallpaper-list.txt")
	if err := g.writeCacheFile(cacheFile, files); err != nil {
		return postinstall.NewComponentResultWithError(
			postinstall.ComponentWallpaper,
			"Failed to write cache file",
			err,
		), err
	}
	details = append(details, "Cache file created: "+cacheFile)

	// Set default wallpaper if available
	if len(files) > 0 {
		defaultWallpaper := files[0]
		details = append(details, "Default wallpaper: "+defaultWallpaper)
	}

	return result.
		WithDetails(details...).
		Complete(postinstall.StatusCompleted), nil
}

func (g *WallpaperCacheGenerator) scanWallpapers() ([]string, error) {
	var wallpapers []string

	// Walk the wallpaper directory
	err := filepath.Walk(g.wallpaperDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if file is an image
		ext := filepath.Ext(path)
		switch ext {
		case ".jpg", ".jpeg", ".png", ".webp":
			wallpapers = append(wallpapers, path)
		}

		return nil
	})

	return wallpapers, err
}

func (g *WallpaperCacheGenerator) writeCacheFile(path string, wallpapers []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, wp := range wallpapers {
		if _, err := f.WriteString(wp + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// Verify checks if the cache is generated
func (g *WallpaperCacheGenerator) Verify(ctx context.Context) (bool, error) {
	cacheFile := filepath.Join(g.cacheDir, "wallpaper-list.txt")
	_, err := os.Stat(cacheFile)
	return err == nil, nil
}

// Rollback reverts the cache generation
func (g *WallpaperCacheGenerator) Rollback(ctx context.Context) error {
	return os.RemoveAll(g.cacheDir)
}
