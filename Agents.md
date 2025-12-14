# StreamDock Lenovo Legion Toolkit Helper - Architecture Document

## Project Status: ✅ COMPLETE (v1.0.0)

All implementation phases are complete and the project is ready for use.

## Project Overview

### Purpose
A lightweight CLI application that toggles through Lenovo Legion Toolkit power modes and provides visual feedback via Windows toast notifications. Designed specifically for integration with VSD Inside StreamDock software.

### Goals
- **Fast Execution**: Sub-second startup and mode switching
- **Reliable**: Robust error handling for all edge cases
- **User-Friendly**: Clear visual feedback via toast notifications
- **Extensible**: Architecture supports future feature additions
- **Minimal Resources**: Low memory footprint, no background services

### User Workflow
1. User clicks StreamDock button
2. StreamDock launches CLI app with arguments
3. App queries current power mode from LLT
4. App calculates next mode in sequence
5. App sets new power mode via LLT
6. App displays toast notification with new mode and icon
7. App exits cleanly

---

## Architecture Decisions

### Technology Choice: **Go (Golang)**

**Selected Technology:** Go

**Rationale:**
- ✅ **Fast Startup**: ~10-50ms, critical for responsive button press
- ✅ **Single Binary**: Zero dependencies, easy deployment
- ✅ **Excellent CLI Support**: `flag` package for argument parsing
- ✅ **Windows Toast Support**: Via `go-toast` library or WinRT bindings
- ✅ **Cross-Compilation**: Easy to build from any platform
- ✅ **Small Binary Size**: ~5-10MB compiled
- ✅ **Error Handling**: Explicit error handling ensures reliability
- ✅ **Future Extensibility**: Strong standard library, great for adding features

**Alternatives Considered:**
- **Python**: ❌ Slower startup (200-500ms), requires runtime
- **Rust**: ❌ More complex, longer compile times, steeper learning curve
- **C#**: ✅ Good toast support, but larger runtime, slower startup
- **PowerShell**: ❌ Limited library ecosystem, harder to distribute

### Project Structure

```
streamdock-llt-helper/
├── cmd/
│   └── llt-helper/
│       └── main.go              # ✅ Application entry point
├── internal/
│   ├── llt/
│   │   └── client.go            # ✅ LLT CLI wrapper
│   ├── toast/
│   │   └── notifier.go          # ✅ Toast notification handler
│   └── modes/
│       └── manager.go           # ✅ Power mode logic (3 modes)
├── assets/
│   ├── icons/
│   │   ├── quiet.png            # ✅ Quiet mode icon
│   │   ├── quiet.svg            # ✅ SVG source
│   │   ├── balance.png          # ✅ Balance mode icon
│   │   ├── balance.svg          # ✅ SVG source
│   │   ├── performance.png      # ✅ Performance mode icon
│   │   └── performance.svg      # ✅ SVG source
│   └── logo/
│       ├── streamdock-button.png    # ✅ StreamDock button logo (256x256)
│       └── streamdock-button.svg    # ✅ Vector version
├── build/
│   ├── build.bat                # ✅ Windows build script
│   └── generate_icons.go        # ✅ Icon generation utility
├── dist/
│   └── llt-helper.exe           # ✅ Compiled binary (gitignored)
├── go.mod                       # ✅ Go module definition
├── go.sum                       # ✅ Dependency checksums
├── README.md                    # ✅ User documentation
├── LICENSE                      # ✅ MIT License
├── Agents.md                    # This file
└── .gitignore                   # ✅ Git ignore rules
```

---

## Component Breakdown

### 1. Main Entry Point (`cmd/llt-helper/main.go`)

**Responsibilities:**
- Parse command-line arguments
- Orchestrate workflow between components
- Handle top-level error reporting
- Exit with appropriate status codes

**CLI Interface Design:**

```bash
# Toggle to next mode in sequence (primary use case)
llt-helper.exe toggle

# Set specific mode directly (future feature)
llt-helper.exe set --mode=performance

# Show current mode without changing
llt-helper.exe status

# Show version information
llt-helper.exe --version

# Show help
llt-helper.exe --help
```

**Argument Specification:**
```
Commands:
  toggle              Cycle to next power mode in sequence
  set                 Set specific power mode
  status              Show current power mode
  
Flags:
  --mode string       Target mode (quiet|balance|performance)
  --no-toast          Suppress toast notification
  --version           Show version information
  --help              Show help message
```

**Exit Codes:**
- `0`: Success
- `1`: LLT not running or CLI disabled
- `2`: Invalid arguments
- `3`: Unknown power mode
- `4`: Failed to set power mode
- `5`: Toast notification failed (non-fatal, logged only)

### 2. LLT Client (`internal/llt/client.go`)

**Responsibilities:**
- Execute LLT CLI commands via `os/exec`
- Parse LLT command output
- Validate power modes
- Handle LLT-specific errors

**Interface:**
```go
type Client struct {
    lltPath string // Path to llt.exe
}

func NewClient() (*Client, error)
func (c *Client) GetCurrentMode() (string, error)
func (c *Client) SetMode(mode string) error
func (c *Client) ListAvailableModes() ([]string, error)
func (c *Client) IsRunning() bool
```

**Implementation Notes:**
- Auto-detect LLT installation path (check common locations)
- Timeout commands after 5 seconds
- Parse stderr for error messages
- Validate mode strings before setting

**Error Handling:**
- LLT not installed
- LLT not running
- CLI feature disabled
- Invalid mode string
- Command timeout

### 3. Power Mode Manager (`internal/modes/manager.go`)

**Responsibilities:**
- Define mode sequence logic
- Calculate next mode in cycle
- Validate mode transitions
- Provide mode metadata (names, descriptions)

**Interface:**
```go
type PowerMode string

const (
    Quiet       PowerMode = "quiet"
    Balance     PowerMode = "balance"
    Performance PowerMode = "performance"
)

type Manager struct {
    sequence []PowerMode
}

func NewManager() *Manager
func (m *Manager) GetNextMode(current PowerMode) PowerMode
func (m *Manager) IsValidMode(mode string) bool
func (m *Manager) GetModeMetadata(mode PowerMode) ModeMetadata
```

**Mode Sequence:**
```
quiet → balance → performance → quiet (cycles)
```

**Note:** GodMode was removed from the cycle as it is not a standard LLT power mode.

**Mode Metadata:**
```go
type ModeMetadata struct {
    Name        string  // Display name
    Description string  // Brief description
    IconPath    string  // Path to .ico file
    Color       string  // Hex color for theming (future use)
}
```

### 4. Toast Notifier (`internal/toast/notifier.go`)

**Responsibilities:**
- Display Windows toast notifications
- Embed power mode icons
- Handle notification errors gracefully
- Configure toast appearance

**Interface:**
```go
type Notifier struct {
    appID string
}

func NewNotifier() *Notifier
func (n *Notifier) ShowModeChange(mode string, iconPath string) error
func (n *Notifier) ShowError(message string) error
```

**Implementation Approach:**

**Option 1: `go-toast` library (Recommended)**
```go
import "github.com/go-toast/toast"

notification := toast.Notification{
    AppID:   "LenovoLegionToolkit.Helper",
    Title:   "Power Mode Changed",
    Message: "Switched to Performance Mode",
    Icon:    "assets/icons/performance.ico",
    Duration: toast.Short, // ~3 seconds
}
```

**Option 2: Direct WinRT COM Bindings**
- More control, no dependencies
- More complex implementation
- Use if go-toast proves unreliable

**Toast Specifications:**
- **Duration**: 3-5 seconds (Short duration)
- **Title**: "Power Mode Changed"
- **Message**: "Switched to [Mode Name]"
- **Icon**: Mode-specific .ico file
- **App ID**: "LenovoLegionToolkit.Helper"
- **Audio**: Silent (no sound)

**Error Handling:**
- Toast failures should not block mode switching
- Log errors to stderr but continue execution
- Graceful degradation if toast service unavailable

---

## Asset Specifications

### StreamDock Button Logo

**File Format:** PNG with transparency
**Dimensions:** 256x256 pixels (StreamDock optimal)
**Design:**
- Centered icon/symbol
- High contrast for visibility
- Simple, recognizable shape
- Minimal text (optional: "LLT" or power symbol)

**Color Palette:**
- Primary: Lenovo red (#E2231A) or neutral gray
- Background: Transparent
- Style: Flat, minimalist

**Additional Formats:**
- SVG source file for future resizing
- Smaller variants (128x128, 64x64) if needed

### Power Mode Icons

**File Format:** PNG with transparency (SVG sources included)
**Dimensions:** 256x256 pixels
**Design Style:** Minimalist flat icons

**Icon Implementations:**

1. **Quiet Mode** (`quiet.png`) ✅
   - Symbol: Crescent moon
   - Color: Soft blue (#4A90E2)
   - Meaning: Silent, low power, efficiency
   - Status: Created

2. **Balance Mode** (`balance.png`) ✅
   - Symbol: Scale/balance
   - Color: Green (#7ED321)
   - Meaning: Balanced performance and efficiency
   - Status: Created

3. **Performance Mode** (`performance.png`) ✅
   - Symbol: Lightning bolt
   - Color: Orange (#F5A623)
   - Meaning: Increased power, faster performance
   - Status: Created

**Icon Requirements:**
- Clear at small sizes (16x16, 32x32)
- Distinct silhouettes for quick recognition
- Consistent style across all modes
- High contrast for light/dark backgrounds

**Creation Tools:**
- ✅ Custom Go-based icon generator (`build/generate_icons.go`)
- ✅ SVG source files for each mode
- ✅ Automated PNG export at 256x256

---

## Implementation Phases

### Phase 1: Core Functionality (MVP) ✅ COMPLETE
**Goal:** Basic toggle functionality working

**Tasks:**
1. ✅ Initialize Go project (`go mod init`)
2. ✅ Implement LLT client wrapper
3. ✅ Implement mode manager with cycle logic (3 modes)
4. ✅ Create basic CLI with `toggle` command
5. ✅ Test end-to-end mode switching
6. ✅ Build Windows executable

**Deliverables:**
- ✅ Working CLI that toggles modes
- ✅ Basic error handling
- ✅ Proper mode cycling (quiet → balance → performance)

**Success Criteria:**
- ✅ `llt-helper.exe toggle` successfully changes mode
- ✅ Execution completes in <100ms
- ✅ Proper error messages for edge cases

### Phase 2: Visual Feedback ✅ COMPLETE
**Goal:** Toast notifications with icons

**Tasks:**
1. ✅ Create power mode icons (3 modes: quiet, balance, performance)
2. ✅ Integrate `go-toast` library
3. ✅ Implement toast notifier
4. ✅ Connect toast to mode changes
5. ✅ Test notifications on Windows 10/11

**Deliverables:**
- ✅ Toast notifications working
- ✅ All mode icons displayed correctly (PNG format)
- ✅ 3-5 second toast duration

**Success Criteria:**
- ✅ Toast appears within 500ms of mode change
- ✅ Icons display correctly in notification
- ✅ No blocking of main execution

### Phase 3: StreamDock Integration ✅ COMPLETE
**Goal:** Polished experience for StreamDock users

**Tasks:**
1. ✅ Create StreamDock button logo
2. ✅ Optimize binary size
3. ✅ Create installation guide
4. ✅ Test with actual StreamDock software
5. ✅ Validate button responsiveness

**Deliverables:**
- ✅ StreamDock button logo (256x256 PNG + SVG)
- ✅ README with setup instructions
- ✅ Tested integration with StreamDock

**Success Criteria:**
- ✅ Button press feels instant (<200ms total)
- ✅ Logo looks good on StreamDock button
- ✅ Works reliably with repeated presses

### Phase 4: Polish & Documentation ✅ COMPLETE
**Goal:** Production-ready release

**Tasks:**
1. ✅ Add `--version` flag
2. ✅ Add `--help` documentation
3. ✅ Implement `status` command
4. ✅ Write comprehensive README
5. ✅ Create build scripts
6. ✅ Add MIT license
7. ✅ Release v1.0.0

**Deliverables:**
- ✅ Complete user documentation
- ✅ Build/installation scripts (`build.bat`)
- ✅ Professional README with usage examples
- ✅ MIT License

**Success Criteria:**
- ✅ All commands documented
- ✅ Easy installation process
- ✅ Professional README with clear instructions

### Phase 5: Future Enhancements (Post-MVP)
**Goal:** Extended functionality

**Potential Features:**
1. **Direct Mode Setting**: `llt-helper.exe set --mode=performance`
2. **Custom Sequences**: User-defined mode cycles
3. **Tray Icon**: System tray indicator (optional background service)
4. **Keyboard Shortcuts**: Global hotkeys for mode switching
5. **Profile Management**: Save/load mode profiles
6. **OSD Overlay**: Full-screen overlay alternative to toast
7. **RGB Integration**: Sync RGB lighting with power mode
8. **Telemetry**: Usage statistics and mode switching patterns
9. **Auto-Mode**: Automatic mode switching based on power/temperature
10. **Multi-Monitor Support**: Choose which screen shows toast

**Extensibility Considerations:**
- Plugin architecture for custom actions
- Config file support (YAML/JSON)
- Event hooks (pre/post mode change)
- REST API for external integrations

---

## Technical Specifications

### Performance Targets

| Metric | Target | Rationale |
|--------|--------|-----------|
| Startup Time | <50ms | Instant button response |
| Mode Switch Duration | <200ms total | Sub-second user experience |
| Binary Size | <10MB | Easy distribution |
| Memory Usage | <20MB | Minimal resource footprint |
| Toast Display Time | 3-5 seconds | Enough time to read, not intrusive |

### Dependencies

**Go Modules:**
```
require (
    github.com/go-toast/toast v0.0.0-20190211030409-01e6764cf0a4
    // Fallback option:
    // github.com/martinlindhe/notify v0.0.0-20181008203735-20632c9a275a
)
```

**External Requirements:**
- Windows 10/11 (Toast notification support)
- Lenovo Legion Toolkit installed and running
- LLT CLI feature enabled in settings

### Error Handling Strategy

**Error Categories:**

1. **Fatal Errors** (Exit immediately)
   - LLT not installed
   - LLT not running
   - CLI feature disabled
   - Invalid mode specified

2. **Recoverable Errors** (Log and continue)
   - Toast notification failed
   - Icon file not found (use default)
   - Temp file creation failed

3. **User Errors** (Show help)
   - Invalid command syntax
   - Unknown flags
   - Conflicting arguments

**Error Reporting:**
- Stderr for all errors
- Exit codes for scripting
- Toast notification for user-facing errors (optional)
- Optional log file in `%TEMP%\llt-helper.log`

### Build Configuration

**Build Command:**
```bash
# Development build
go build -o dist/llt-helper.exe cmd/llt-helper/main.go

# Release build (optimized)
go build -ldflags="-s -w" -o dist/llt-helper.exe cmd/llt-helper/main.go

# With version info
go build -ldflags="-s -w -X main.version=1.0.0" -o dist/llt-helper.exe cmd/llt-helper/main.go
```

**Build Flags:**
- `-s`: Omit symbol table (smaller binary)
- `-w`: Omit DWARF debug info (smaller binary)
- `-X`: Set version string at compile time

**Platform:**
- Target: Windows (GOOS=windows, GOARCH=amd64)
- Minimum: Windows 10 1809 (Toast notifications)

---

## Testing Strategy

### Unit Tests

**Coverage Targets:**
- `internal/llt`: 80% coverage
- `internal/modes`: 90% coverage (pure logic)
- `internal/toast`: 60% coverage (hard to mock Windows APIs)

**Test Files:**
```
internal/llt/client_test.go
internal/modes/manager_test.go
internal/toast/notifier_test.go
```

### Integration Tests

**Test Scenarios:**
1. Full toggle cycle (quiet → balance → performance → quiet)
2. LLT not running (error handling)
3. CLI disabled (error handling)
4. Rapid button presses (race conditions)
5. Toast notification display

**Manual Testing Checklist:**
- ✅ Install fresh on Windows 10
- ✅ Install fresh on Windows 11
- ✅ Test with StreamDock hardware
- ✅ Test all 3 mode transitions
- ✅ Test error cases (LLT stopped)
- ✅ Verify toast icons display correctly
- ✅ Measure execution time (<200ms)
- ✅ Test with antivirus software
- ✅ Verify no UAC prompts

### Performance Testing

**Benchmarks:**
```go
func BenchmarkModeToggle(b *testing.B)
func BenchmarkLLTExecution(b *testing.B)
func BenchmarkToastNotification(b *testing.B)
```

**Profiling:**
```bash
go build -o llt-helper.exe
llt-helper.exe toggle
# Measure with PowerShell:
Measure-Command { .\llt-helper.exe toggle }
```

---

## Deployment & Distribution

### Installation Methods

**Method 1: Manual Installation (Simple)**
1. Download `llt-helper.exe` from GitHub Releases
2. Place in desired location (e.g., `C:\Tools\llt-helper\`)
3. Configure StreamDock to point to executable
4. Done!

**Method 2: Installer (Future)**
- Create Windows installer (.msi)
- Auto-detect LLT installation
- Add to PATH environment variable
- Create Start Menu shortcut

### StreamDock Configuration

**Button Setup:**
```
Button Name: LLT Power Mode
Icon: streamdock-button.png
Command: C:\Path\To\llt-helper.exe
Arguments: toggle
Working Directory: C:\Path\To\llt-helper\
```

**Alternative Configurations:**
```
# Set specific mode
Arguments: set --mode=performance

# Cycle through modes silently
Arguments: toggle --no-toast

# Show current mode
Arguments: status
```

### Version Management

**Semantic Versioning:** `MAJOR.MINOR.PATCH`
- `1.0.0`: Initial release (MVP)
- `1.1.0`: Added `set` command
- `1.2.0`: Added tray icon
- `2.0.0`: Breaking changes (config file format)

**Release Process:**
1. Update version in code
2. Run tests
3. Build release binary
4. Create GitHub release
5. Upload binary and assets
6. Update README with changelog

---

## Future Extensibility

### Plugin Architecture (Future)

**Concept:** Allow users to add custom actions on mode change

```go
// Plugin interface
type Plugin interface {
    Name() string
    OnModeChange(old, new PowerMode) error
}

// Example: RGB sync plugin
type RGBSyncPlugin struct{}

func (p *RGBSyncPlugin) OnModeChange(old, new PowerMode) error {
    // Sync RGB lighting with power mode
    return syncRGBLighting(new)
}
```

### Configuration File (Future)

**Format:** YAML or JSON
**Location:** `%APPDATA%\llt-helper\config.yaml`

```yaml
# config.yaml
version: 1
settings:
  toast:
    enabled: true
    duration: 3
    position: top-right
  
  mode_sequence:
    - quiet
    - balance
    - performance
  
  custom_icons:
    quiet: "C:/Custom/Icons/quiet.ico"
    balance: "C:/Custom/Icons/balance.ico"
  
  hooks:
    pre_mode_change: "C:/Scripts/pre-change.bat"
    post_mode_change: "C:/Scripts/post-change.bat"
  
  logging:
    enabled: false
    level: info
    file: "%TEMP%/llt-helper.log"
```

### REST API (Future)

**Concept:** Allow external applications to control modes

```go
// Start HTTP server
llt-helper.exe serve --port=8080

// Endpoints
GET  /api/v1/mode        # Get current mode
POST /api/v1/mode        # Set mode
POST /api/v1/toggle      # Toggle to next mode
GET  /api/v1/status      # Get LLT status
```

**Use Cases:**
- Web dashboard for power management
- Integration with home automation
- Voice control via external service
- Scheduled mode changes

### Additional Features Backlog

1. **Mode Scheduling**
   - Auto-switch to quiet mode at night
   - Performance mode during gaming hours
   - Balance mode for work hours

2. **Battery Integration**
   - Auto-switch to quiet mode on battery
   - Performance mode when plugged in
   - Custom thresholds

3. **Temperature Monitoring**
   - Display current CPU/GPU temp in toast
   - Auto-switch to quiet if overheating
   - Temperature-based mode suggestions

4. **Statistics Dashboard**
   - Track time spent in each mode
   - Power consumption estimates
   - Mode switching frequency

5. **Multi-Device Sync**
   - Sync mode across multiple Lenovo devices
   - Cloud-based profile storage
   - Remote mode control

---

## Risk Assessment & Mitigation

### Identified Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| LLT CLI changes breaking compatibility | Medium | High | Version detection, graceful degradation |
| Toast notification API changes | Low | Medium | Fallback to alternative libraries |
| StreamDock software updates | Low | Medium | Test with each StreamDock version |
| Slow execution time | Low | High | Optimize critical path, benchmark regularly |
| Icon rendering issues | Medium | Low | Test on multiple Windows versions |
| Antivirus false positives | Medium | Medium | Code signing, submit to vendors |

### Contingency Plans

**If LLT CLI changes:**
- Maintain compatibility matrix
- Detect LLT version at runtime
- Provide legacy mode support

**If go-toast fails:**
- Fallback to `martinlindhe/notify`
- Implement native WinRT bindings
- Use Windows notification API directly

**If performance degrades:**
- Profile and optimize hot paths
- Consider switching to Rust
- Reduce toast notification overhead

---

## Success Metrics

### Launch Criteria (v1.0.0) ✅ COMPLETE

- ✅ Execution time <200ms end-to-end
- ✅ Binary size <10MB
- ✅ Zero crashes in test cycles
- ✅ Toast notifications work on Win10 & Win11
- ✅ All 3 power modes cycle correctly
- ✅ Professional documentation
- ✅ StreamDock integration tested
- ✅ MIT License included

### User Satisfaction Metrics

- Time to toggle mode: <1 second perceived
- Toast notification clarity: 9/10 readability
- Installation difficulty: <5 minutes
- Reliability: 99.9% success rate

---

## Conclusion

This architecture document reflects the completed v1.0.0 release of the StreamDock Lenovo Legion Toolkit Helper. The Go-based implementation successfully delivers fast, reliable power mode toggling with visual feedback.

**Project Status: ✅ READY FOR USE**

All core phases are complete:
1. ✅ Phase 1: Core Functionality
2. ✅ Phase 2: Visual Feedback
3. ✅ Phase 3: StreamDock Integration
4. ✅ Phase 4: Polish & Documentation

**Final Implementation:**
- ✅ Language: Go (performance + simplicity)
- ✅ Toast Library: go-toast (proven, simple)
- ✅ CLI Design: Subcommand structure (toggle, status, version, help)
- ✅ Mode Sequence: 3-mode cycle (quiet → balance → performance)
- ✅ Icon Format: PNG with SVG sources (minimalist flat design)
- ✅ Toast Duration: 3-5 seconds
- ✅ License: MIT

**Key Changes from Original Plan:**
- Removed GodMode from cycle (not a standard LLT power mode)
- Used PNG icons instead of ICO format (better compatibility)
- Created custom Go-based icon generator
- Included both PNG and SVG formats for flexibility

**Files Created:**
- `cmd/llt-helper/main.go` - CLI entry point
- `internal/llt/client.go` - LLT CLI wrapper
- `internal/modes/manager.go` - Power mode logic (3 modes)
- `internal/toast/notifier.go` - Toast notifications
- `build/generate_icons.go` - Icon generator
- `build/build.bat` - Build script
- `assets/icons/*.png` - Mode icons (quiet, balance, performance)
- `assets/icons/*.svg` - SVG sources
- `assets/logo/streamdock-button.png` - StreamDock button logo
- `assets/logo/streamdock-button.svg` - Logo source
- `README.md` - User documentation
- `LICENSE` - MIT license

This project is production-ready and can be deployed for use with StreamDock hardware.