package usecases_test

import (
	"context"
	"testing"

	"github.com/rebelopsio/gohan/internal/application/installation/dto"
	"github.com/rebelopsio/gohan/internal/application/installation/usecases"
	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartInstallationUseCase_Execute(t *testing.T) {
	t.Run("successfully starts installation with valid request", func(t *testing.T) {
		sessionRepo := repository.NewMemorySessionRepository()
		useCase := usecases.NewStartInstallationUseCase(sessionRepo)
		ctx := context.Background()

		request := dto.InstallationRequest{
			Components: []dto.ComponentRequest{
				{
					Name:    "hyprland",
					Version: "0.35.0",
				},
			},
			AvailableSpace: 100 * uint64(installation.GB),
			RequiredSpace:  10 * uint64(installation.GB),
		}

		response, err := useCase.Execute(ctx, request)

		require.NoError(t, err)
		assert.NotEmpty(t, response.SessionID)
		assert.Equal(t, "pending", response.Status)
		assert.Equal(t, 1, response.ComponentCount)
		assert.NotEmpty(t, response.StartedAt)
	})

	t.Run("returns error for insufficient disk space", func(t *testing.T) {
		sessionRepo := repository.NewMemorySessionRepository()
		useCase := usecases.NewStartInstallationUseCase(sessionRepo)
		ctx := context.Background()

		request := dto.InstallationRequest{
			Components: []dto.ComponentRequest{
				{
					Name:    "hyprland",
					Version: "0.35.0",
				},
			},
			AvailableSpace: 5 * uint64(installation.GB),
			RequiredSpace:  10 * uint64(installation.GB),
		}

		_, err := useCase.Execute(ctx, request)

		assert.Error(t, err)
		assert.ErrorIs(t, err, installation.ErrInsufficientDiskSpace)
	})

	t.Run("returns error when no components specified", func(t *testing.T) {
		sessionRepo := repository.NewMemorySessionRepository()
		useCase := usecases.NewStartInstallationUseCase(sessionRepo)
		ctx := context.Background()

		request := dto.InstallationRequest{
			Components:     []dto.ComponentRequest{},
			AvailableSpace: 100 * uint64(installation.GB),
			RequiredSpace:  10 * uint64(installation.GB),
		}

		_, err := useCase.Execute(ctx, request)

		assert.Error(t, err)
	})

	t.Run("returns error when core component missing", func(t *testing.T) {
		sessionRepo := repository.NewMemorySessionRepository()
		useCase := usecases.NewStartInstallationUseCase(sessionRepo)
		ctx := context.Background()

		request := dto.InstallationRequest{
			Components: []dto.ComponentRequest{
				{
					Name:    "waybar", // Not core component
					Version: "0.9.20",
				},
			},
			AvailableSpace: 100 * uint64(installation.GB),
			RequiredSpace:  10 * uint64(installation.GB),
		}

		_, err := useCase.Execute(ctx, request)

		assert.Error(t, err)
		assert.ErrorIs(t, err, installation.ErrInvalidConfiguration)
	})

	t.Run("includes GPU support when provided", func(t *testing.T) {
		sessionRepo := repository.NewMemorySessionRepository()
		useCase := usecases.NewStartInstallationUseCase(sessionRepo)
		ctx := context.Background()

		request := dto.InstallationRequest{
			Components: []dto.ComponentRequest{
				{
					Name:    "hyprland",
					Version: "0.35.0",
				},
			},
			GPU: &dto.GPURequest{
				Vendor:         "amd",
				RequiresDriver: true,
				DriverName:     "amd_driver",
			},
			AvailableSpace: 100 * uint64(installation.GB),
			RequiredSpace:  10 * uint64(installation.GB),
		}

		response, err := useCase.Execute(ctx, request)

		require.NoError(t, err)
		assert.NotEmpty(t, response.SessionID)
	})

	t.Run("validates component names", func(t *testing.T) {
		sessionRepo := repository.NewMemorySessionRepository()
		useCase := usecases.NewStartInstallationUseCase(sessionRepo)
		ctx := context.Background()

		request := dto.InstallationRequest{
			Components: []dto.ComponentRequest{
				{
					Name:    "", // Empty name
					Version: "0.35.0",
				},
			},
			AvailableSpace: 100 * uint64(installation.GB),
			RequiredSpace:  10 * uint64(installation.GB),
		}

		_, err := useCase.Execute(ctx, request)

		assert.Error(t, err)
	})

	t.Run("validates component versions", func(t *testing.T) {
		sessionRepo := repository.NewMemorySessionRepository()
		useCase := usecases.NewStartInstallationUseCase(sessionRepo)
		ctx := context.Background()

		request := dto.InstallationRequest{
			Components: []dto.ComponentRequest{
				{
					Name:    "hyprland",
					Version: "", // Empty version
				},
			},
			AvailableSpace: 100 * uint64(installation.GB),
			RequiredSpace:  10 * uint64(installation.GB),
		}

		_, err := useCase.Execute(ctx, request)

		assert.Error(t, err)
	})
}

func TestStartInstallationUseCase_ConvertComponentName(t *testing.T) {
	t.Run("converts known component names", func(t *testing.T) {
		tests := []struct {
			input    string
			expected installation.ComponentName
		}{
			{"hyprland", installation.ComponentHyprland},
			{"waybar", installation.ComponentWaybar},
			{"rofi", installation.ComponentFuzzel},
			{"kitty", installation.ComponentKitty},
			{"amd_driver", installation.ComponentAMDDriver},
			{"nvidia_driver", installation.ComponentNVIDIADriver},
			{"intel_driver", installation.ComponentIntelDriver},
		}

		for _, tt := range tests {
			t.Run(tt.input, func(t *testing.T) {
				result := usecases.ConvertComponentName(tt.input)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("returns original string for unknown components", func(t *testing.T) {
		result := usecases.ConvertComponentName("unknown_component")
		assert.Equal(t, installation.ComponentName("unknown_component"), result)
	})
}
