# StreamDock LLT Helper - Implementation Plan

## Overview
This document breaks down the architecture plan into specific, actionable tasks for implementation. Each phase builds on the previous one, ensuring a stable foundation before adding new features.

---

## Phase 1: Core Functionality (MVP)

### 1.1 Project Initialization
**Estimated Time:** 30 minutes

**Tasks:**
- [ ] Create project directory structure
- [ ] Initialize Git repository
- [ ] Create `.gitignore` file
- [ ] Initialize Go module: `go mod init github.com/yourusername/streamdock-llt-helper`
- [ ] Create basic directory structure (`cmd/`, `internal/`, `assets/`)
- [ ] Create placeholder `README.md`

**Deliverables:**
- Working Git repository
- Go module initialized
- Directory structure in place

**Validation:**
```bash
go mod verify
git status
```

---

### 1.2 LLT Client Implementation
**Estimated Time:** 2-3 hours

**Tasks:**
- [ ] Create `internal/llt/client.go`
- [ ] Implement `NewClient()` with LLT path auto-detection
- [ ] Implement `GetCurrentMode()` - execute `llt f get power-mode`
- [ ] Implement `SetMode()` - execute `llt f set power-mode <mode>`
- [ ] Implement `IsRunning()` - check if LLT process exists
- [ ] Add command timeout handling (5 seconds)
- [ ] Parse command output and errors
- [ ] Create `internal/llt/client_test.go` with unit tests

**Key Functions:**
```go
func NewClient() (*Client, error)
func (c *Client) GetCurrentMode() (string, error)
func (c *Client) SetMode(mode string) error
func (c *Client) IsRunning() bool
```

**Testing:**
- Mock LLT CLI responses
- Test error cases (not running, invalid mode)
- Test timeout scenarios

**Validation:**
```bash
go test ./internal/llt/
```

---

### 1.3 Power Mode Manager
**Estimated Time:** 1-2 hours

**Tasks:**
- [ ] Create `internal/modes/manager.go`
- [ ] Define `PowerMode` type and constants
- [ ] Implement mode sequence array
- [ ] Implement `GetNextMode()` - cycle logic
- [ ] Implement `IsValidMode()` - validation
- [ ] Define `ModeMetadata` struct
- [ ] Implement `GetModeMetadata()` - return mode info
- [ ] Create `internal/modes/manager_test.go`

**Key Data Structures:**
```go
type PowerMode string
const (
    Quiet       PowerMode = "quiet"
    Balance     PowerMode = "balance"
    Performance PowerMode = "performance"
    GodMode     PowerMode = "godmode"
)
```

**Testing:**
- Test full cycle: quiet → balance → performance → godmode → quiet
- Test edge cases (invalid modes)
- Test metadata retrieval

**Validation:**
```bash
go test ./internal/modes/
```

---

### 1.4 CLI Application Entry Point
**Estimated Time:** 2-3 hours

**Tasks:**
- [ ] Create `cmd/llt-helper/main.go`
- [ ] Implement command-line argument parsing
- [ ] Implement `toggle` command
- [ ] Implement error handling and exit codes
- [ ] Add basic logging to stderr
- [ ] Wire up LLT client and mode manager
- [ ] Test end-to-end toggle functionality

**CLI Structure:**
```go
func main() {
    // Parse args
    // Initialize components
    // Execute command
    // Handle errors
    // Exit with status code
}
```

**Exit Codes:**
- `0`: Success
- `1`: LLT not running
- `2`: Invalid arguments
- `3`: Unknown power mode
- `4`: Failed to set mode

**Validation:**
```bash
go build -o dist/llt-helper.exe cmd/llt-helper/main.go
.\dist\llt-helper.exe toggle
# Verify mode changes in LLT
```

---

### 1.5 Build and Test
**Estimated Time:** 1 hour

**Tasks:**
- [ ] Create `build/build.bat` script
- [ ] Test build on Windows
- [ ] Verify binary size (<10MB)
- [ ] Test execution time (<100ms)
- [ ] Test all error scenarios
- [ ] Document build process

**Build Script:**
```batch
@echo off
echo Building llt-helper.exe...
go build -ldflags="-s -w" -o dist/llt-helper.exe cmd/llt-helper/main.go
echo Build complete: dist/llt-helper.exe
```

**Validation:**
- Binary builds successfully
- Execution time measured with `Measure-Command`
- All error cases handled gracefully

---

## Phase 2: Visual Feedback (Toast Notifications)

### 2.1 Toast Notification Library Integration
**Estimated Time:** 2-3 hours

**Tasks:**
- [ ] Add `go-toast` dependency: `go get github.com/go-toast/toast`
- [ ] Create `internal/toast/notifier.go`
- [ ] Implement `NewNotifier()` constructor
- [ ] Implement `ShowModeChange()` - display mode change toast
- [ ] Implement `ShowError()` - display error toast
- [ ] Add graceful degradation if toast fails
- [ ] Create `internal/toast/notifier_test.go`

**Key Functions:**
```go
func NewNotifier() *Notifier
func (n *Notifier) ShowModeChange(mode string, iconPath string) error
func (n *Notifier) ShowError(message string) error
```

**Toast Configuration:**
```go
toast.Notification{
    AppID:    "LenovoLegionToolkit.Helper",
    Title:    "Power Mode Changed",
    Message:  "Switched to Performance Mode",
    Icon:     "assets/icons/performance.ico",
    Duration: toast.Short, // 3-5 seconds
}
```

**Testing:**
- Test toast display on Windows 10/11
- Test with valid and invalid icon paths
- Test error handling

---

### 2.2 Power Mode Icon Creation
**Estimated Time:** 3-4 hours

**Tasks:**
- [ ] Design minimalist flat icons for each mode
- [ ] Create SVG source files for each icon
- [ ] Convert SVGs to multi-resolution ICO files
  - Quiet mode: Feather/leaf (soft blue)
  - Balance mode: Scale/balance (green/yellow)
  - Performance mode: Lightning bolt (orange)
  - God mode: Flame/rocket (red)
- [ ] Test icons at 16x16, 32x32, 48x48, 256x256
- [ ] Place icons in `assets/icons/` directory

**Tools:**
- Design: Figma, Illustrator, or Inkscape
- Conversion: ImageMagick or online ICO converter

**Icon Files:**
```
assets/icons/quiet.ico
assets/icons/balance.ico
assets/icons/performance.ico
assets/icons/godmode.ico
```

**Validation:**
- Icons display correctly in Windows Explorer
- Icons look good at all sizes
- Consistent visual style across all modes

---

### 2.3 Integrate Toast with CLI
**Estimated Time:** 1-2 hours

**Tasks:**
- [ ] Update `main.go` to initialize notifier
- [ ] Show toast after successful mode change
- [ ] Pass correct icon path based on mode
- [ ] Handle toast errors gracefully (don't fail mode switch)
- [ ] Add `--no-toast` flag to suppress notifications
- [ ] Test complete workflow: toggle → LLT → toast

**Integration:**
```go
// After mode change succeeds
metadata := modeManager.GetModeMetadata(newMode)
notifier.ShowModeChange(string(newMode), metadata.IconPath)
```

**Validation:**
```bash
.\dist\llt-helper.exe toggle
# Verify toast appears with correct icon
```

---

## Phase 3: StreamDock Integration

### 3.1 StreamDock Button Logo Creation
**Estimated Time:** 2-3 hours

**Tasks:**
- [ ] Design StreamDock button logo (256x256 PNG)
- [ ] Create SVG source file
- [ ] Export PNG with transparency
- [ ] Create smaller variants (128x128, 64x64) if needed
- [ ] Place in `assets/logo/` directory

**Design Requirements:**
- Minimalist, recognizable shape
- High contrast for visibility
- Lenovo red (#E2231A) or neutral colors
- Transparent background

**Validation:**
- Logo looks good on StreamDock button
- Clear at small sizes

---

### 3.2 Binary Optimization
**Estimated Time:** 1-2 hours

**Tasks:**
- [ ] Optimize build flags for size
- [ ] Strip debug symbols
- [ ] Test with UPX compression (optional)
- [ ] Verify binary size <10MB
- [ ] Test execution speed <200ms total

**Optimized Build:**
```bash
go build -ldflags="-s -w" -o dist/llt-helper.exe cmd/llt-helper/main.go
```

**Validation:**
```powershell
# Measure execution time
Measure-Command { .\dist\llt-helper.exe toggle }
# Check file size
(Get-Item .\dist\llt-helper.exe).Length / 1MB
```

---

### 3.3 StreamDock Testing
**Estimated Time:** 1-2 hours

**Tasks:**
- [ ] Configure StreamDock button with CLI app
- [ ] Test button responsiveness
- [ ] Test rapid button presses
- [ ] Verify logo display
- [ ] Document StreamDock configuration
- [ ] Create screenshots for documentation

**StreamDock Configuration:**
```
Button Name: LLT Power Mode
Icon: streamdock-button.png
Command: C:\Path\To\llt-helper.exe
Arguments: toggle
```

**Validation:**
- Button press feels instant
- No lag or delays
- Toast displays correctly

---

## Phase 4: Polish & Documentation

### 4.1 Additional CLI Commands
**Estimated Time:** 2-3 hours

**Tasks:**
- [ ] Implement `status` command - show current mode
- [ ] Implement `set --mode=<mode>` command
- [ ] Add `--version` flag
- [ ] Add `--help` flag with usage info
- [ ] Improve error messages
- [ ] Add optional debug logging

**Commands:**
```bash
llt-helper.exe status
llt-helper.exe set --mode=performance
llt-helper.exe --version
llt-helper.exe --help
```

**Validation:**
- All commands work correctly
- Help text is clear and comprehensive

---

### 4.2 Documentation
**Estimated Time:** 3-4 hours

**Tasks:**
- [ ] Write comprehensive `README.md`
  - Project description
  - Features
  - Requirements
  - Installation instructions
  - StreamDock setup guide
  - Usage examples
  - Troubleshooting
  - Screenshots
- [ ] Create `CHANGELOG.md`
- [ ] Update `Agents.md` with final architecture notes
- [ ] Create developer documentation

**README Structure:**
```markdown
# StreamDock Lenovo Legion Toolkit Helper

## Features
## Requirements
## Installation
## StreamDock Setup
## Usage
## Troubleshooting
## Development
## License
```

**Validation:**
- Documentation is clear and complete
- Screenshots illustrate key features
- Installation steps verified

---

### 4.3 Release Preparation
**Estimated Time:** 2 hours

**Tasks:**
- [ ] Run full test suite
- [ ] Create GitHub release (v1.0.0)
- [ ] Upload compiled binary
- [ ] Upload assets (icons, logo)
- [ ] Tag release in Git
- [ ] Write release notes

**Release Checklist:**
- [ ] All tests passing
- [ ] Binary size <10MB
- [ ] Execution time <200ms
- [ ] Toast notifications working
- [ ] Icons displaying correctly
- [ ] Documentation complete
- [ ] No antivirus false positives

---

## Phase 5: Future Enhancements (Post-MVP)

### Future Features (Priority Order)

1. **Custom Mode Sequences** (2-3 hours)
   - Add config file support
   - Allow user-defined mode order
   - Validate custom sequences

2. **System Tray Icon** (4-6 hours)
   - Add optional background service
   - Right-click menu for mode selection
   - Status indicator

3. **Keyboard Shortcuts** (3-4 hours)
   - Global hotkey registration
   - Configurable key bindings
   - Handle conflicts

4. **Configuration UI** (8-10 hours)
   - Simple settings GUI
   - Mode sequence editor
   - Icon customization

5. **Mode Profiles** (4-5 hours)
   - Save/load mode profiles
   - Quick switch between profiles
   - Profile management UI

---

## Testing Matrix

### Manual Testing Checklist

**Environment:**
- [ ] Windows 10 (21H2 or later)
- [ ] Windows 11
- [ ] LLT installed and running
- [ ] LLT CLI enabled in settings

**Functionality:**
- [ ] Toggle command works
- [ ] Full cycle: quiet → balance → performance → godmode → quiet
- [ ] Status command shows correct mode
- [ ] Set command works for all modes
- [ ] Error handling: LLT not running
- [ ] Error handling: CLI disabled
- [ ] Error handling: Invalid mode

**Performance:**
- [ ] Execution time <200ms
- [ ] Binary size <10MB
- [ ] Memory usage <20MB
- [ ] No memory leaks

**Visual:**
- [ ] Toast notifications appear
- [ ] Icons display correctly
- [ ] Toast duration 3-5 seconds
- [ ] No UI blocking

**StreamDock:**
- [ ] Button press responsive
- [ ] Logo displays correctly
- [ ] Rapid presses handled
- [ ] No crashes or hangs

**Edge Cases:**
- [ ] Antivirus compatibility
- [ ] UAC compatibility
- [ ] Multiple rapid executions
- [ ] Network disconnected
- [ ] Low disk space

---

## Time Estimates

### Phase Breakdown

| Phase | Tasks | Estimated Time |
|-------|-------|----------------|
| Phase 1: Core Functionality | 5 tasks | 8-12 hours |
| Phase 2: Visual Feedback | 3 tasks | 6-9 hours |
| Phase 3: StreamDock Integration | 3 tasks | 4-7 hours |
| Phase 4: Polish & Documentation | 3 tasks | 7-9 hours |
| **Total MVP** | **14 tasks** | **25-37 hours** |

### Recommended Schedule

**Week 1:**
- Phase 1 complete (core functionality)
- Basic testing

**Week 2:**
- Phase 2 complete (toast notifications)
- Icon creation
- Integration testing

**Week 3:**
- Phase 3 complete (StreamDock integration)
- Phase 4 started (documentation)

**Week 4:**
- Phase 4 complete (polish & docs)
- Release v1.0.0

---

## Success Criteria

### MVP Launch Requirements

**Functionality:**
- [x] Toggle between all 4 power modes
- [x] Display toast notifications with icons
- [x] Execute in <200ms
- [x] Handle all error cases gracefully

**Quality:**
- [x] Zero crashes in 1000 test cycles
- [x] Binary size <10MB
- [x] Professional documentation
- [x] StreamDock integration tested

**User Experience:**
- [x] Button press feels instant
- [x] Toast notifications clear and informative
- [x] Easy installation (<5 minutes)
- [x] Reliable operation (99%+ success rate)

---

## Risk Mitigation

### High Priority Risks

1. **LLT CLI Changes**
   - *Mitigation:* Implement version detection early
   - *Fallback:* Graceful degradation to basic functionality

2. **Toast Notification Failures**
   - *Mitigation:* Test on multiple Windows versions
   - *Fallback:* Alternative notification library ready

3. **Performance Issues**
   - *Mitigation:* Benchmark early and often
   - *Fallback:* Optimize critical path, consider Rust

4. **StreamDock Compatibility**
   - *Mitigation:* Test with actual hardware early
   - *Fallback:* Provide alternative launch methods

---

## Next Steps

1. **Review this implementation plan**
2. **Approve architecture and approach**
3. **Switch to Code mode to begin Phase 1**
4. **Implement iteratively with testing at each step**
5. **Gather feedback and adjust as needed**

This plan provides a clear roadmap from initial setup to production-ready release, with flexibility for adjustments based on testing and feedback.