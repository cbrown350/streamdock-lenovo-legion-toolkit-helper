package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cbrown350/streamdock-lenovo-legion-toolkit-helper/internal/llt"
	"github.com/cbrown350/streamdock-lenovo-legion-toolkit-helper/internal/modes"
	"github.com/cbrown350/streamdock-lenovo-legion-toolkit-helper/internal/toast"
)

const version = "1.0.0"

func main() {
	// Check for global flags first
	if len(os.Args) > 1 {
		if os.Args[1] == "--version" || os.Args[1] == "-version" {
			fmt.Println("llt-helper version", version)
			os.Exit(0)
		}
		if os.Args[1] == "--help" || os.Args[1] == "-help" || os.Args[1] == "-h" {
			printUsage()
			os.Exit(0)
		}
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	command := os.Args[1]

	// Parse command-specific flags
	var modeFlag string
	var noToast bool
	var modesFlag string

	fs := flag.NewFlagSet(command, flag.ExitOnError)
	fs.StringVar(&modeFlag, "mode", "", "Target mode for set command (quiet|balance|performance)")
	fs.BoolVar(&noToast, "no-toast", false, "Suppress toast notification")
	fs.StringVar(&modesFlag, "modes", "", "Comma-separated list of modes to cycle through for toggle command (e.g., quiet,performance)")

	fs.Usage = func() {
		printUsage()
	}

	// Parse flags after the command
	if len(os.Args) > 2 {
		fs.Parse(os.Args[2:])
	}

	// Initialize components
	lltClient, err := llt.NewClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if !lltClient.IsRunning() {
		fmt.Fprintf(os.Stderr, "Error: LLT not running or CLI disabled\n")
		os.Exit(1)
	}

	modeManager := modes.NewManager()
	var notifier *toast.Notifier
	if !noToast {
		notifier = toast.NewNotifier()
	}

	switch command {
	case "toggle":
		err = handleToggle(lltClient, modeManager, notifier, modesFlag)
	case "set":
		if modeFlag == "" {
			fmt.Fprintf(os.Stderr, "Error: --mode flag required for set command\n")
			os.Exit(2)
		}
		err = handleSet(lltClient, modeManager, modeFlag, notifier)
	case "status":
		err = handleStatus(lltClient, modeManager)
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n\n", command)
		printUsage()
		os.Exit(2)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(4)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [command] [flags]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Commands:\n")
	fmt.Fprintf(os.Stderr, "  toggle              Cycle to next power mode in sequence\n")
	fmt.Fprintf(os.Stderr, "  set --mode=MODE     Set specific power mode\n")
	fmt.Fprintf(os.Stderr, "  status              Show current power mode\n")
	fmt.Fprintf(os.Stderr, "\nGlobal Flags:\n")
	fmt.Fprintf(os.Stderr, "  --version           Show version information\n")
	fmt.Fprintf(os.Stderr, "  --help              Show this help message\n")
	fmt.Fprintf(os.Stderr, "\nCommand Flags:\n")
	fmt.Fprintf(os.Stderr, "  --mode string       Target mode (quiet|balance|performance)\n")
	fmt.Fprintf(os.Stderr, "  --modes string      Comma-separated modes for toggle (e.g., quiet,performance)\n")
	fmt.Fprintf(os.Stderr, "  --no-toast          Suppress toast notification\n")
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  %s toggle\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s set --mode=balance\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s toggle --no-toast\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s toggle --modes=quiet,performance\n", os.Args[0])
}

func handleToggle(client *llt.Client, manager *modes.Manager, notifier *toast.Notifier, modesFlag string) error {
	current, err := client.GetCurrentMode()
	if err != nil {
		return err
	}

	var allowedModes []modes.PowerMode
	if modesFlag != "" {
		// Parse comma-separated modes
		parts := strings.Split(modesFlag, ",")
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed == "" {
				continue
			}
			if !manager.IsValidMode(trimmed) {
				return fmt.Errorf("invalid mode '%s' in --modes flag", trimmed)
			}
			allowedModes = append(allowedModes, modes.PowerMode(trimmed))
		}
		if len(allowedModes) == 0 {
			return fmt.Errorf("no valid modes specified in --modes flag")
		}
	}

	next := manager.GetNextModeFromList(modes.PowerMode(current), allowedModes)
	err = client.SetMode(string(next))
	if err != nil {
		return err
	}

	if notifier != nil {
		meta := manager.GetModeMetadata(next)
		if err := notifier.ShowModeChange(meta.Name, meta.IconPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: toast notification failed: %v\n", err)
			// Don't exit, as mode was set successfully
		}
	}

	return nil
}

func handleSet(client *llt.Client, manager *modes.Manager, mode string, notifier *toast.Notifier) error {
	if !manager.IsValidMode(mode) {
		return fmt.Errorf("unknown power mode: %s", mode)
	}

	err := client.SetMode(mode)
	if err != nil {
		return err
	}

	if notifier != nil {
		meta := manager.GetModeMetadata(modes.PowerMode(mode))
		if err := notifier.ShowModeChange(meta.Name, meta.IconPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: toast notification failed: %v\n", err)
		}
	}

	return nil
}

func handleStatus(client *llt.Client, manager *modes.Manager) error {
	current, err := client.GetCurrentMode()
	if err != nil {
		return err
	}

	meta := manager.GetModeMetadata(modes.PowerMode(current))
	fmt.Printf("Current Mode: %s (%s)\n", meta.Name, current)
	return nil
}
