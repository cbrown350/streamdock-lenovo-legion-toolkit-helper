# StreamDock Lenovo Legion Toolkit Helper

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Windows](https://img.shields.io/badge/Windows-10%2F11-0078D6?style=flat&logo=windows)](https://www.microsoft.com/windows)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A fast, lightweight CLI tool that toggles through Lenovo Legion Toolkit power modes with visual toast notifications. Designed specifically for seamless integration with **VSD Inside StreamDock** software.

![StreamDock Button Demo](assets/logo/streamdock-button.png)

---

## ✨ Features

- **🔄 Power Mode Cycling** - Toggle through quiet → balance → performance → quiet
- **🔔 Toast Notifications** - Visual feedback with mode-specific icons
- **⚡ Lightning Fast** - Sub-200ms execution time
- **📦 Single Binary** - No dependencies, no installation required
- **🎮 StreamDock Ready** - Perfect for one-button power mode switching
- **🛠️ Flexible CLI** - Set specific modes or check current status

---

## 📋 Requirements

Before using this tool, ensure you have:

| Requirement | Details |
|-------------|---------|
| **Operating System** | Windows 10 (1809+) or Windows 11 |
| **Lenovo Legion Toolkit** | [Download here](https://github.com/BartoszCichworski/LenovoLegionToolkit) - Must be installed and running |
| **LLT CLI Feature** | Must be enabled in LLT settings (see [Setup Guide](#enabling-llt-cli-feature)) |
| **StreamDock** (optional) | [VSD Inside StreamDock](https://www.vsd-inside.com/) for button integration |

---

## 🚀 Quick Start

### 1. Download

Download the latest `llt-helper.exe` from the [Releases](https://github.com/cbrown350/streamdock-llt-helper/releases) page.

### 2. Place the Executable

Put `llt-helper.exe` in a convenient location, for example:
```
C:\Tools\llt-helper\llt-helper.exe
```

### 3. Enable LLT CLI Feature

1. Open **Lenovo Legion Toolkit**
2. Go to **Settings** (gear icon)
3. Find **"CLI"** or **"Command Line Interface"** option
4. **Enable** the CLI feature
5. Keep LLT running in the background

### 4. Test It

Open Command Prompt or PowerShell and run:
```bash
C:\Tools\llt-helper\llt-helper.exe toggle
```

You should see a toast notification showing the new power mode! 🎉

---

## 💻 Usage

### Basic Commands

```bash
# Toggle to next power mode (primary use case)
llt-helper.exe toggle

# Set a specific power mode
llt-helper.exe set --mode=quiet
llt-helper.exe set --mode=balance
llt-helper.exe set --mode=performance

# Check current power mode
llt-helper.exe status

# Show version information
llt-helper.exe --version

# Show help
llt-helper.exe --help
```

### Advanced Options

```bash
# Toggle without showing toast notification
llt-helper.exe toggle --no-toast

# Set mode silently
llt-helper.exe set --mode=performance --no-toast
```

### Power Mode Cycle

The `toggle` command cycles through modes in this order:

```
┌─────────┐     ┌─────────┐     ┌─────────────┐
│  Quiet  │ ──► │ Balance │ ──► │ Performance │
└─────────┘     └─────────┘     └─────────────┘
     ▲                                  │
     └──────────────────────────────────┘
```

---

## 🎮 StreamDock Setup

Setting up a StreamDock button for one-touch power mode switching:

### Step 1: Add a New Button

1. Open the **StreamDock** software
2. Select an empty button slot
3. Choose **"Open"** or **"Run Program"** action type

### Step 2: Configure the Button

| Setting | Value |
|---------|-------|
| **Program/Path** | `C:\Tools\llt-helper\llt-helper.exe` |
| **Arguments** | `toggle` |
| **Start In** | `C:\Tools\llt-helper\` (optional) |

### Step 3: Set the Button Icon

1. Click on the button icon area
2. Navigate to: `C:\Tools\llt-helper\assets\logo\`
3. Select `streamdock-button.png`

### Step 4: Test

Press the StreamDock button - you should see:
1. The power mode changes
2. A toast notification appears showing the new mode

<!-- 
### Screenshots

![StreamDock Configuration](docs/images/streamdock-config.png)
*StreamDock button configuration*

![Toast Notification](docs/images/toast-example.png)
*Toast notification showing mode change*
-->

---

## 🔧 Building from Source

### Prerequisites

- [Go 1.21+](https://golang.org/dl/) installed
- Git (optional, for cloning)

### Build Steps

```bash
# Clone the repository
git clone https://github.com/cbrown350/streamdock-llt-helper.git
cd streamdock-llt-helper

# Download dependencies
go mod download

# Build the executable
go build -ldflags="-s -w" -o dist/llt-helper.exe ./cmd/llt-helper

# (Optional) Generate icon assets
go run build/generate_icons.go
```

### Build Flags Explained

| Flag | Purpose |
|------|---------|
| `-s` | Omit symbol table (smaller binary) |
| `-w` | Omit DWARF debug info (smaller binary) |

---

## 📁 Project Structure

```
streamdock-llt-helper/
├── cmd/
│   └── llt-helper/
│       └── main.go           # CLI entry point
├── internal/
│   ├── llt/
│   │   └── client.go         # LLT CLI wrapper
│   ├── modes/
│   │   └── manager.go        # Power mode logic
│   └── toast/
│       └── notifier.go       # Toast notifications
├── assets/
│   ├── icons/                # Mode icons (PNG/SVG)
│   │   ├── quiet.png
│   │   ├── balance.png
│   │   └── performance.png
│   └── logo/
│       └── streamdock-button.png  # StreamDock button icon
├── build/
│   └── generate_icons.go     # Icon generation script
├── dist/
│   └── llt-helper.exe        # Compiled binary
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
└── README.md                 # This file
```

---

## ❌ Exit Codes

For scripting and automation, the tool returns these exit codes:

| Code | Meaning |
|------|---------|
| `0` | Success - operation completed |
| `1` | LLT not running or CLI feature disabled |
| `2` | Invalid command-line arguments |
| `3` | Unknown power mode specified |
| `4` | Failed to set power mode |

---

## 🔍 Troubleshooting

### "LLT not running" Error

**Problem:** The tool reports that Lenovo Legion Toolkit is not running.

**Solutions:**
1. Make sure LLT is running (check system tray)
2. Try restarting LLT
3. Run LLT as administrator if issues persist

### "CLI feature disabled" Error

**Problem:** The CLI feature is not enabled in LLT.

**Solution:**
1. Open Lenovo Legion Toolkit
2. Go to Settings → CLI
3. Enable the CLI feature
4. Restart LLT

### No Toast Notification Appears

**Problem:** Mode changes but no notification shows.

**Solutions:**
1. Check Windows notification settings
2. Ensure "Lenovo Legion Toolkit Helper" notifications are allowed
3. Try running with administrator privileges
4. Check if Focus Assist/Do Not Disturb is enabled

### Mode Doesn't Change

**Problem:** Toast shows but power mode doesn't actually change.

**Solutions:**
1. Verify LLT is running and responsive
2. Try changing mode manually in LLT first
3. Check if your laptop supports all power modes
4. Some modes may require AC power

### StreamDock Button Not Working

**Problem:** Pressing the StreamDock button does nothing.

**Solutions:**
1. Verify the path to `llt-helper.exe` is correct
2. Check that the file exists at the specified location
3. Try running the command manually in Command Prompt
4. Ensure StreamDock software is running

### Slow Execution

**Problem:** There's a noticeable delay when pressing the button.

**Solutions:**
1. Ensure LLT is already running (first launch is slower)
2. Check for antivirus interference
3. Try placing the executable on an SSD

---

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

### Reporting Issues

1. Check existing [Issues](https://github.com/cbrown350/streamdock-llt-helper/issues) first
2. Include your Windows version and LLT version
3. Provide steps to reproduce the problem
4. Include any error messages

### Submitting Changes

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Test thoroughly on Windows
5. Commit with clear messages: `git commit -m "Add amazing feature"`
6. Push to your fork: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Development Guidelines

- Follow Go best practices and conventions
- Add tests for new functionality
- Update documentation as needed
- Keep commits focused and atomic

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

```
MIT License

Copyright (c) 2024 cbrown350

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

## 🙏 Acknowledgments

- [Lenovo Legion Toolkit](https://github.com/BartoszCichworski/LenovoLegionToolkit) - The amazing toolkit this helper integrates with
- [VSD Inside StreamDock](https://www.vsd-inside.com/) - StreamDock hardware and software
- [go-toast](https://github.com/go-toast/toast) - Windows toast notification library

---

## 📞 Support

- **Issues:** [GitHub Issues](https://github.com/cbrown350/streamdock-llt-helper/issues)
- **Discussions:** [GitHub Discussions](https://github.com/cbrown350/streamdock-llt-helper/discussions)

---

<p align="center">
  Made with ❤️ for the Lenovo Legion community
</p>
