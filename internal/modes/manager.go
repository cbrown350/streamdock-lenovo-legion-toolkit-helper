package modes

import (
	"os"
	"path/filepath"
)

// PowerMode represents a Lenovo Legion Toolkit power mode
type PowerMode string

const (
	Quiet       PowerMode = "quiet"
	Balance     PowerMode = "balance"
	Performance PowerMode = "performance"
	GodMode     PowerMode = "godmode"
)

// ModeMetadata contains display information for a power mode
type ModeMetadata struct {
	Name        string
	Description string
	IconPath    string
	Color       string // Future use
}

// Manager handles power mode operations
type Manager struct {
	sequence []PowerMode
}

// NewManager creates a new power mode manager
func NewManager() *Manager {
	return &Manager{
		sequence: []PowerMode{Quiet, Balance, Performance},
	}
}

// GetNextMode returns the next power mode in the sequence
func (m *Manager) GetNextMode(current PowerMode) PowerMode {
	currentIndex := -1
	for i, mode := range m.sequence {
		if mode == current {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		// Invalid current mode, default to first
		return m.sequence[0]
	}

	nextIndex := (currentIndex + 1) % len(m.sequence)
	return m.sequence[nextIndex]
}

// GetNextModeFromList returns the next power mode from the provided list
func (m *Manager) GetNextModeFromList(current PowerMode, allowedModes []PowerMode) PowerMode {
	if len(allowedModes) == 0 {
		return m.GetNextMode(current)
	}

	currentIndex := -1
	for i, mode := range allowedModes {
		if mode == current {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		// Current mode not in list, start from first
		return allowedModes[0]
	}

	nextIndex := (currentIndex + 1) % len(allowedModes)
	return allowedModes[nextIndex]
}

// IsValidMode checks if the given mode string is valid
func (m *Manager) IsValidMode(mode string) bool {
	for _, pm := range m.sequence {
		if string(pm) == mode {
			return true
		}
	}
	return false
}

// GetModeMetadata returns metadata for the given power mode
func (m *Manager) GetModeMetadata(mode PowerMode) ModeMetadata {
	// Find the assets directory relative to the executable
	baseDir := findAssetsDir()

	metadata := map[PowerMode]ModeMetadata{
		Quiet: {
			Name:        "Quiet",
			Description: "Silent operation with minimal power consumption",
			IconPath:    filepath.Join(baseDir, "assets", "icons", "quiet.png"),
			Color:       "#4A90E2",
		},
		Balance: {
			Name:        "Balance",
			Description: "Balanced performance and efficiency",
			IconPath:    filepath.Join(baseDir, "assets", "icons", "balance.png"),
			Color:       "#7ED321",
		},
		Performance: {
			Name:        "Performance",
			Description: "Increased power for better performance",
			IconPath:    filepath.Join(baseDir, "assets", "icons", "performance.png"),
			Color:       "#F5A623",
		},
	}

	if meta, exists := metadata[mode]; exists {
		return meta
	}

	// Default metadata for unknown modes
	return ModeMetadata{
		Name:        string(mode),
		Description: "Unknown power mode",
		IconPath:    "",
		Color:       "#000000",
	}
}

// findAssetsDir locates the assets directory relative to the executable
func findAssetsDir() string {
	// Try to get executable path
	exePath, err := os.Executable()
	if err != nil {
		// Fallback to current working directory
		cwd, _ := os.Getwd()
		return cwd
	}

	exeDir := filepath.Dir(exePath)

	// Check if assets exists in executable directory
	assetsPath := filepath.Join(exeDir, "assets")
	if _, err := os.Stat(assetsPath); err == nil {
		return exeDir
	}

	// Check if assets exists in parent directory (for dist/ subdirectory)
	parentDir := filepath.Dir(exeDir)
	assetsPath = filepath.Join(parentDir, "assets")
	if _, err := os.Stat(assetsPath); err == nil {
		return parentDir
	}

	// Last resort: try current working directory
	cwd, _ := os.Getwd()
	assetsPath = filepath.Join(cwd, "assets")
	if _, err := os.Stat(assetsPath); err == nil {
		return cwd
	}

	// Default to executable directory even if assets not found
	return exeDir
}
