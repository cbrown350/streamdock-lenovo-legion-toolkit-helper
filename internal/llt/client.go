package llt

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// Client wraps interactions with Lenovo Legion Toolkit CLI
type Client struct {
	lltPath string
}

// NewClient creates a new LLT client and auto-detects the LLT path
func NewClient() (*Client, error) {
	lltPath := os.Getenv("LOCALAPPDATA")
	if lltPath == "" {
		lltPath = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
	}
	lltPath = filepath.Join(lltPath, "Programs", "LenovoLegionToolkit", "llt.exe")

	if _, err := os.Stat(lltPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("LLT not found at %s", lltPath)
	}

	return &Client{lltPath: lltPath}, nil
}

// IsRunning checks if LLT is accessible
func (c *Client) IsRunning() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, c.lltPath, "f", "get", "power-mode")
	
	// Hide console window
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}

	err := cmd.Run()
	return err == nil
}

// GetCurrentMode retrieves the current power mode
func (c *Client) GetCurrentMode() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, c.lltPath, "f", "get", "power-mode")
	
	// Hide console window
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current mode: %w", err)
	}

	mode := strings.TrimSpace(string(output))
	return mode, nil
}

// SetMode sets the power mode to the specified value
func (c *Client) SetMode(mode string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, c.lltPath, "f", "set", "power-mode", mode)
	
	// Hide console window
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to set mode to %s: %w", mode, err)
	}

	return nil
}

// ListAvailableModes lists all available power modes
func (c *Client) ListAvailableModes() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, c.lltPath, "f", "set", "power-mode", "-l")
	
	// Hide console window
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list modes: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var modes []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			modes = append(modes, line)
		}
	}

	return modes, nil
}
