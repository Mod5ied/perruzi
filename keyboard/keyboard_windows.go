//go:build windows

package keyboard

import (
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	INPUT_KEYBOARD    = 1
	KEYEVENTF_UNICODE = 0x0004
	KEYEVENTF_KEYUP   = 0x0002
	VK_RETURN         = 0x0D
	VK_BACK           = 0x08
)

type KEYBDINPUT struct {
	WVk         uint16
	WScan       uint16
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

type INPUT struct {
	Type uint32
	Ki   KEYBDINPUT
	_    [8]byte // padding to match C struct size on 64-bit
}

var sendInput = windows.NewLazySystemDLL("user32.dll").NewProc("SendInput")

func sendUnicode(r rune, keyUp bool) {
	flags := uint32(KEYEVENTF_UNICODE)
	if keyUp {
		flags |= KEYEVENTF_KEYUP
	}
	inp := INPUT{
		Type: INPUT_KEYBOARD,
		Ki: KEYBDINPUT{
			WScan:   uint16(r),
			DwFlags: flags,
		},
	}
	sendInput.Call(1, uintptr(unsafe.Pointer(&inp)), unsafe.Sizeof(inp))
}

func sendVKey(vk uint16, keyUp bool) {
	flags := uint32(0)
	if keyUp {
		flags = KEYEVENTF_KEYUP
	}
	inp := INPUT{
		Type: INPUT_KEYBOARD,
		Ki: KEYBDINPUT{
			WVk:     vk,
			DwFlags: flags,
		},
	}
	sendInput.Call(1, uintptr(unsafe.Pointer(&inp)), unsafe.Sizeof(inp))
}

type windowsInjector struct{}

// NewInjector returns the Windows Injector implementation.
func NewInjector() (Injector, error) {
	return &windowsInjector{}, nil
}

func (w *windowsInjector) InjectChar(r rune) error {
	sendUnicode(r, false)
	sendUnicode(r, true)
	return nil
}

func (w *windowsInjector) InjectReturn() error {
	sendVKey(VK_RETURN, false)
	sendVKey(VK_RETURN, true)
	return nil
}

func (w *windowsInjector) InjectBackspace() error {
	sendVKey(VK_BACK, false)
	sendVKey(VK_BACK, true)
	return nil
}

// IsReady reports whether the keyboard subsystem can inject keystrokes.
func IsReady() bool {
	return true
}

// IsAccessibilityGranted is always true on Windows; no equivalent permission is required.
func IsAccessibilityGranted() bool {
	return true
}

// EscChan is the package-level ESC signal channel. Engine reads from this.
var EscChan = make(chan struct{}, 1)

var escRunning int32

var getAsyncKeyState = windows.NewLazySystemDLL("user32.dll").NewProc("GetAsyncKeyState")

const VK_ESCAPE = 0x1B

// StartEscListener starts polling for the ESC key.
func StartEscListener() {
	atomic.StoreInt32(&escRunning, 1)
	go func() {
		for atomic.LoadInt32(&escRunning) == 1 {
			ret, _, _ := getAsyncKeyState.Call(VK_ESCAPE)
			if ret&0x8000 != 0 {
				select {
				case EscChan <- struct{}{}:
				default:
				}
				time.Sleep(300 * time.Millisecond) // debounce
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
}

// StopEscListener stops polling for the ESC key.
func StopEscListener() {
	atomic.StoreInt32(&escRunning, 0)
}
