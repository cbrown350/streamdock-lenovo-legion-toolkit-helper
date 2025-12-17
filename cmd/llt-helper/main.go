package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"unsafe"

	"github.com/cbrown350/streamdock-lenovo-legion-toolkit-helper/internal/llt"
	"github.com/cbrown350/streamdock-lenovo-legion-toolkit-helper/internal/modes"
	"github.com/cbrown350/streamdock-lenovo-legion-toolkit-helper/internal/toast"
	"golang.org/x/sys/windows"
)

const version = "1.0.0"

var consoleHandle uintptr

func attachConsole() {
	const ATTACH_PARENT_PROCESS = ^uint32(0) // (DWORD)-1
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")

	attachConsoleProc := kernel32.NewProc("AttachConsole")
	ret, _, _ := attachConsoleProc.Call(uintptr(ATTACH_PARENT_PROCESS))

	if ret == 0 {
		// Couldn't attach to parent, not running from console
		return
	}

	// Get stderr handle for output
	const STD_ERROR_HANDLE = ^uintptr(11) + 1 // -12
	getStdHandleProc := kernel32.NewProc("GetStdHandle")
	handle, _, _ := getStdHandleProc.Call(STD_ERROR_HANDLE)

	if handle != 0 && handle != uintptr(windows.InvalidHandle) {
		consoleHandle = handle
	}
}

// writeToConsole writes directly to the console using Windows API
func writeToConsole(message string) {
	if consoleHandle == 0 {
		return
	}

	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	writeFileProc := kernel32.NewProc("WriteFile")

	data := []byte(message)
	var written uint32
	writeFileProc.Call(
		consoleHandle,
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(unsafe.Pointer(&written)),
		0,
	)
}

func main() {
	// Attempt to attach to parent console for CLI output
	attachConsole()

	// Check for global flags first
	if len(os.Args) > 1 {
		if os.Args[1] == "--version" || os.Args[1] == "-version" {
			versionMsg := fmt.Sprintf("llt-helper version %s\n", version)
			writeToConsole(versionMsg)
			fmt.Print(versionMsg)
			os.Exit(0)
		}
		if os.Args[1] == "--help" || os.Args[1] == "-help" || os.Args[1] == "-h" {
			printUsage()
			os.Exit(0) // Exit code 0 for help
		}
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1) // Exit code 1 for missing arguments
	}

	command := os.Args[1]

	// Parse command-specific flags
	var modeFlag string
	var noToast bool
	var modesFlag string
	var helpFlag bool

	fs := flag.NewFlagSet(command, flag.ExitOnError)
	fs.StringVar(&modeFlag, "mode", "", "Target mode for set command (quiet|balance|performance)")
	fs.BoolVar(&noToast, "no-toast", false, "Suppress toast notification")
	fs.StringVar(&modesFlag, "modes", "", "Comma-separated list of modes to cycle through for toggle command (e.g., quiet,performance)")
	fs.BoolVar(&helpFlag, "help", false, "Show help message")
	fs.BoolVar(&helpFlag, "h", false, "Show help message (shorthand)")

	fs.Usage = func() {
		printUsage()
	}

	// Parse flags after the command
	if len(os.Args) > 2 {
		if err := fs.Parse(os.Args[2:]); err != nil {
			// flag.ExitOnError handles exit usually, but if we catch it:
			os.Exit(2)
		}
	}

	if helpFlag {
		printUsage()
		os.Exit(0)
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
			printUsage() // Helpful to show usage on error
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
	usage := fmt.Sprintf(`Usage: %s [command] [flags]

Commands:
  toggle              Cycle to next power mode in sequence
  set --mode=MODE     Set specific power mode
  status              Show current power mode

Global Flags:
  --version           Show version information
  --help, -h          Show this help message

Command Flags:
  --mode string       Target mode (quiet|balance|performance)
  --modes string      Comma-separated modes for toggle (e.g., quiet,performance)
  --no-toast          Suppress toast notification

Examples:
  %s toggle
  %s set --mode=balance
  %s toggle --no-toast
  %s toggle --modes=quiet,performance
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])

	writeToConsole(usage)
	// Also write to stderr for non-console contexts
	fmt.Fprint(os.Stderr, usage)
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
	statusMsg := fmt.Sprintf("Current Mode: %s (%s)\n", meta.Name, current)
	writeToConsole(statusMsg)
	fmt.Print(statusMsg)
	return nil
}
