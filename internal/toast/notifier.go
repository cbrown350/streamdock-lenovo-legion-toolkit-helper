package toast

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32                  = windows.NewLazySystemDLL("user32.dll")
	gdi32                   = windows.NewLazySystemDLL("gdi32.dll")
	procCreateWindowEx      = user32.NewProc("CreateWindowExW")
	procDefWindowProc       = user32.NewProc("DefWindowProcW")
	procDispatchMessage     = user32.NewProc("DispatchMessageW")
	procGetMessage          = user32.NewProc("GetMessageW")
	procRegisterClassEx     = user32.NewProc("RegisterClassExW")
	procPostQuitMessage     = user32.NewProc("PostQuitMessage")
	procShowWindow          = user32.NewProc("ShowWindow")
	procUpdateWindow        = user32.NewProc("UpdateWindow")
	procGetSystemMetrics    = user32.NewProc("GetSystemMetrics")
	procSetWindowPos        = user32.NewProc("SetWindowPos")
	procSetLayeredWindowAttributes = user32.NewProc("SetLayeredWindowAttributes")
	procGetDC               = user32.NewProc("GetDC")
	procReleaseDC           = user32.NewProc("ReleaseDC")
	procBeginPaint          = user32.NewProc("BeginPaint")
	procEndPaint            = user32.NewProc("EndPaint")
	procFillRect            = user32.NewProc("FillRect")
	procCreateSolidBrush    = gdi32.NewProc("CreateSolidBrush")
	procDeleteObject        = gdi32.NewProc("DeleteObject")
	procSetBkMode           = gdi32.NewProc("SetBkMode")
	procSetTextColor        = gdi32.NewProc("SetTextColor")
	procDrawText            = user32.NewProc("DrawTextW")
	procCreateFont          = gdi32.NewProc("CreateFontW")
	procSelectObject        = gdi32.NewProc("SelectObject")
	procSetTimer            = user32.NewProc("SetTimer")
	procKillTimer           = user32.NewProc("KillTimer")
	procDestroyWindow       = user32.NewProc("DestroyWindow")
	procTranslateMessage    = user32.NewProc("TranslateMessage")
)

const (
	WS_EX_LAYERED     = 0x00080000
	WS_EX_TOPMOST     = 0x00000008
	WS_EX_TOOLWINDOW  = 0x00000080
	WS_POPUP          = 0x80000000
	WS_VISIBLE        = 0x10000000
	SW_SHOW           = 5
	SWP_NOSIZE        = 0x0001
	SWP_NOMOVE        = 0x0002
	SWP_NOZORDER      = 0x0004
	SWP_SHOWWINDOW    = 0x0040
	HWND_TOPMOST      = ^uintptr(0)
	LWA_ALPHA         = 0x00000002
	SM_CXSCREEN       = 0
	SM_CYSCREEN       = 1
	WM_PAINT          = 0x000F
	WM_TIMER          = 0x0113
	WM_DESTROY        = 0x0002
	DT_CENTER         = 0x00000001
	DT_VCENTER        = 0x00000004
	DT_SINGLELINE     = 0x00000020
	TRANSPARENT       = 1
	FW_BOLD           = 700
	DEFAULT_CHARSET   = 1
)

type WNDCLASSEX struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   windows.Handle
	Icon       windows.Handle
	Cursor     windows.Handle
	Background windows.Handle
	MenuName   *uint16
	ClassName  *uint16
	IconSm     windows.Handle
}

type POINT struct {
	X, Y int32
}

type MSG struct {
	Hwnd    windows.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

type RECT struct {
	Left, Top, Right, Bottom int32
}

type PAINTSTRUCT struct {
	Hdc         windows.Handle
	FErase      int32
	RcPaint     RECT
	FRestore    int32
	FIncUpdate  int32
	RgbReserved [32]byte
}

// Notifier handles OSD-style overlay notifications
type Notifier struct {
	appID string
}

// NewNotifier creates a new OSD notifier
func NewNotifier() *Notifier {
	return &Notifier{
		appID: "LenovoLegionToolkit.Helper",
	}
}

var globalMessage string
var globalTitle string

// ShowModeChange displays an OSD overlay notification for power mode change
func (n *Notifier) ShowModeChange(modeName, iconPath string) error {
	globalTitle = "Power Mode Changed"
	globalMessage = fmt.Sprintf("Switched to %s Mode", modeName)
	
	// Show OSD (blocks for duration, but that's OK - we want the notification to stay)
	if err := showOSD(globalTitle, globalMessage, 3*time.Second); err != nil {
		return fmt.Errorf("OSD notification error: %w", err)
	}
	
	return nil
}

// ShowError displays an error OSD notification
func (n *Notifier) ShowError(message string) error {
	globalTitle = "Power Mode Error"
	globalMessage = message
	
	if err := showOSD(globalTitle, globalMessage, 3*time.Second); err != nil {
		return fmt.Errorf("OSD notification error: %w", err)
	}
	
	return nil
}

func showOSD(title, message string, duration time.Duration) error {
	className, _ := syscall.UTF16PtrFromString("LLTHelperOSD")
	
	instance := windows.Handle(0)
	modhandle, err := syscall.LoadLibrary("kernel32.dll")
	if err == nil {
		proc, _ := syscall.GetProcAddress(modhandle, "GetModuleHandleW")
		if proc != 0 {
			instance = windows.Handle(proc)
		}
	}

	wndProc := syscall.NewCallback(wndProcCallback)
	
	wc := WNDCLASSEX{
		Size:      uint32(unsafe.Sizeof(WNDCLASSEX{})),
		WndProc:   wndProc,
		Instance:  instance,
		ClassName: className,
	}

	ret, _, _ := procRegisterClassEx.Call(uintptr(unsafe.Pointer(&wc)))
	if ret == 0 {
		// Class might already be registered, continue anyway
	}

	// Get screen dimensions
	screenWidth, _, _ := procGetSystemMetrics.Call(SM_CXSCREEN)
	screenHeight, _, _ := procGetSystemMetrics.Call(SM_CYSCREEN)

	// OSD dimensions and position
	osdWidth := 400
	osdHeight := 100
	osdX := int((int(screenWidth) - osdWidth) / 2)
	osdY := int(screenHeight) - int(float64(screenHeight)*0.15) // 15% from bottom

	windowName, _ := syscall.UTF16PtrFromString("LLT Helper OSD")
	
	hwnd, _, _ := procCreateWindowEx.Call(
		WS_EX_LAYERED|WS_EX_TOPMOST|WS_EX_TOOLWINDOW,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		WS_POPUP,
		uintptr(osdX),
		uintptr(osdY),
		uintptr(osdWidth),
		uintptr(osdHeight),
		0,
		0,
		uintptr(instance),
		0,
	)

	if hwnd == 0 {
		return fmt.Errorf("CreateWindowEx failed")
	}

	// Set window transparency (220 = ~86% opacity)
	procSetLayeredWindowAttributes.Call(hwnd, 0, 220, LWA_ALPHA)

	// Show window
	procShowWindow.Call(hwnd, SW_SHOW)
	procUpdateWindow.Call(hwnd)

	// Set timer to close window after duration
	timerID := uintptr(1)
	procSetTimer.Call(hwnd, timerID, uintptr(duration.Milliseconds()), 0)

	// Message loop
	var msg MSG
	for {
		ret, _, _ := procGetMessage.Call(
			uintptr(unsafe.Pointer(&msg)),
			0,
			0,
			0,
		)
		if ret == 0 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}

	return nil
}

func wndProcCallback(hwnd windows.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_PAINT:
		var ps PAINTSTRUCT
		hdc, _, _ := procBeginPaint.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&ps)))
		
		// Create dark background
		bgBrush, _, _ := procCreateSolidBrush.Call(0x00202020) // Dark gray
		var rect RECT
		rect.Left = 0
		rect.Top = 0
		rect.Right = 400
		rect.Bottom = 100
		procFillRect.Call(hdc, uintptr(unsafe.Pointer(&rect)), bgBrush)
		procDeleteObject.Call(bgBrush)

		// Set text properties
		procSetBkMode.Call(hdc, TRANSPARENT)
		procSetTextColor.Call(hdc, 0x00FFFFFF) // White text

		// Create fonts
		titleFont, _, _ := procCreateFont.Call(
			24, 0, 0, 0,
			FW_BOLD,
			0, 0, 0,
			DEFAULT_CHARSET,
			0, 0, 0, 0,
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Segoe UI"))),
		)
		messageFont, _, _ := procCreateFont.Call(
			18, 0, 0, 0,
			0,
			0, 0, 0,
			DEFAULT_CHARSET,
			0, 0, 0, 0,
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Segoe UI"))),
		)

		// Draw title
		oldFont, _, _ := procSelectObject.Call(hdc, titleFont)
		titleRect := RECT{Left: 10, Top: 15, Right: 390, Bottom: 45}
		titleText, _ := syscall.UTF16PtrFromString(globalTitle)
		procDrawText.Call(
			hdc,
			uintptr(unsafe.Pointer(titleText)),
			uintptr(^uint(0)), // -1 as uintptr
			uintptr(unsafe.Pointer(&titleRect)),
			DT_CENTER|DT_VCENTER|DT_SINGLELINE,
		)

		// Draw message
		procSelectObject.Call(hdc, messageFont)
		messageRect := RECT{Left: 10, Top: 50, Right: 390, Bottom: 85}
		messageText, _ := syscall.UTF16PtrFromString(globalMessage)
		procDrawText.Call(
			hdc,
			uintptr(unsafe.Pointer(messageText)),
			uintptr(^uint(0)), // -1 as uintptr
			uintptr(unsafe.Pointer(&messageRect)),
			DT_CENTER|DT_VCENTER|DT_SINGLELINE,
		)

		procSelectObject.Call(hdc, oldFont)
		procDeleteObject.Call(titleFont)
		procDeleteObject.Call(messageFont)
		
		procEndPaint.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&ps)))
		return 0

	case WM_TIMER:
		procDestroyWindow.Call(uintptr(hwnd))
		return 0

	case WM_DESTROY:
		procPostQuitMessage.Call(0)
		return 0
	}

	ret, _, _ := procDefWindowProc.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
	return ret
}