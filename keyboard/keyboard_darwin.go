//go:build darwin

package keyboard

/*
#cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation -framework ApplicationServices -framework Cocoa

#include <CoreGraphics/CGEvent.h>
#include <CoreGraphics/CGEventSource.h>
#include <CoreGraphics/CGRemoteOperation.h>
#include <ApplicationServices/ApplicationServices.h>
#include <stdint.h>

// Inject a single Unicode character.
static void inject_char(uint16_t uc) {
    UniChar c = (UniChar)uc;
    CGEventSourceRef src = CGEventSourceCreate(kCGEventSourceStateHIDSystemState);

    CGEventRef down = CGEventCreateKeyboardEvent(src, 0, true);
    CGEventKeyboardSetUnicodeString(down, 1, &c);
    CGEventPost(kCGHIDEventTap, down);
    CFRelease(down);

    CGEventRef up = CGEventCreateKeyboardEvent(src, 0, false);
    CGEventKeyboardSetUnicodeString(up, 1, &c);
    CGEventPost(kCGHIDEventTap, up);
    CFRelease(up);

    CFRelease(src);
}

// Inject the Return key (kVK_Return = 0x24).
static void inject_return() {
    CGEventSourceRef src = CGEventSourceCreate(kCGEventSourceStateHIDSystemState);

    CGEventRef down = CGEventCreateKeyboardEvent(src, 0x24, true);
    CGEventPost(kCGHIDEventTap, down);
    CFRelease(down);

    CGEventRef up = CGEventCreateKeyboardEvent(src, 0x24, false);
    CGEventPost(kCGHIDEventTap, up);
    CFRelease(up);

    CFRelease(src);
}

// Inject Backspace (kVK_Delete = 0x33).
static void inject_backspace() {
    CGEventSourceRef src = CGEventSourceCreate(kCGEventSourceStateHIDSystemState);

    CGEventRef down = CGEventCreateKeyboardEvent(src, 0x33, true);
    CGEventPost(kCGHIDEventTap, down);
    CFRelease(down);

    CGEventRef up = CGEventCreateKeyboardEvent(src, 0x33, false);
    CGEventPost(kCGHIDEventTap, up);
    CFRelease(up);

    CFRelease(src);
}

// Check Accessibility permission.
static int check_accessibility() {
    return AXIsProcessTrusted() ? 1 : 0;
}

// Poll whether the physical Escape key is currently held down.
static int is_esc_pressed() {
    return CGEventSourceKeyState(kCGEventSourceStateHIDSystemState, 53) ? 1 : 0;
}

// Forward declaration of the Go callback.
extern void goEscapeTriggered();

// Forward declaration of the NSEvent monitor bootstrap (implemented in keyboard_monitor_darwin.m).
extern void startNSEventMonitors();

static CFMachPortRef escTapRef = NULL;
static CFRunLoopSourceRef escRunLoopSource = NULL;

static CGEventRef escCallback(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *refcon) {
    // macOS disables event taps that are slow or unresponsive. Re-enable if asked.
    if (type == kCGEventTapDisabledByTimeout || type == kCGEventTapDisabledByUserInput) {
        if (escTapRef != NULL) {
            CGEventTapEnable(escTapRef, true);
        }
        return event;
    }
    if (type == kCGEventKeyDown) {
        CGKeyCode key = (CGKeyCode)CGEventGetIntegerValueField(event, kCGKeyboardEventKeycode);
        if (key == 53) { // kVK_Escape
            goEscapeTriggered();
        }
    }
    return event;
}

static void start_esc_tap() {
    escTapRef = CGEventTapCreate(
        kCGHIDEventTap,
        kCGHeadInsertEventTap,
        kCGEventTapOptionDefault,
        CGEventMaskBit(kCGEventKeyDown),
        escCallback,
        NULL
    );
    if (escTapRef == NULL) {
        // CGEventTapCreate failed (most commonly missing Accessibility permission).
        return;
    }
    escRunLoopSource = CFMachPortCreateRunLoopSource(kCFAllocatorDefault, escTapRef, 0);
    CFRunLoopAddSource(CFRunLoopGetCurrent(), escRunLoopSource, kCFRunLoopCommonModes);
    CGEventTapEnable(escTapRef, true);
    CFRunLoopRun();
}

static void stop_esc_tap() {
    if (escTapRef != NULL) {
        CGEventTapEnable(escTapRef, false);
        CFRunLoopRemoveSource(CFRunLoopGetCurrent(), escRunLoopSource, kCFRunLoopCommonModes);
        CFRelease(escRunLoopSource);
        CFRelease(escTapRef);
        escTapRef = NULL;
    }
}
*/
import "C"
import (
	"runtime"
	"sync/atomic"
	"time"

	"Peruzzi/logger"
)

// EscChan is the package-level ESC signal channel. Engine reads from this.
var EscChan = make(chan struct{}, 1)

//export goEscapeTriggered
func goEscapeTriggered() {
	logger.Log("ESC event triggered")
	select {
	case EscChan <- struct{}{}:
		logger.Log("ESC signal sent to EscChan")
	default:
		logger.Log("ESC signal dropped (EscChan full)")
	}
}

type darwinInjector struct{}

// NewInjector returns the Darwin Injector implementation.
func NewInjector() (Injector, error) {
	return &darwinInjector{}, nil
}

func (d *darwinInjector) InjectChar(r rune) error {
	C.inject_char(C.uint16_t(r))
	return nil
}

func (d *darwinInjector) InjectReturn() error {
	C.inject_return()
	return nil
}

func (d *darwinInjector) InjectBackspace() error {
	C.inject_backspace()
	return nil
}

// IsAccessibilityGranted returns true if the app has Accessibility permission.
func IsAccessibilityGranted() bool {
	return C.check_accessibility() == 1
}

// IsReady reports whether the keyboard subsystem can inject keystrokes.
func IsReady() bool {
	return IsAccessibilityGranted()
}

var pollStop int32

// StartEscListener starts the CGEventTap on a locked OS thread, registers
// NSEvent monitors on the calling thread, and starts a polling fallback.
// Call this once at app startup. It runs until StopEscListener is called.
func StartEscListener() {
	logger.Log("Starting ESC listener")

	// NSEvent monitors must be registered on the main thread; StartEscListener
	// is called from main() on the locked main thread.
	C.startNSEventMonitors()

	// Low-level CGEventTap on its own locked thread.
	go func() {
		runtime.LockOSThread()
		C.start_esc_tap()
	}()

	// Polling fallback: if the event tap and monitors miss the key (e.g. due
	// to focus races or the tap being disabled under load), this goroutine
	// checks the physical Escape key state every 10 ms.
	go func() {
		logger.Log("Starting ESC poll fallback")
		atomic.StoreInt32(&pollStop, 0)
		wasPressed := false
		for atomic.LoadInt32(&pollStop) == 0 {
			pressed := C.is_esc_pressed() == 1
			if pressed && !wasPressed {
				logger.Log("ESC poll fallback detected Escape")
				goEscapeTriggered()
			}
			wasPressed = pressed
			time.Sleep(10 * time.Millisecond)
		}
		logger.Log("ESC poll fallback stopped")
	}()
}

// StopEscListener stops the CGEventTap and the polling fallback.
func StopEscListener() {
	logger.Log("Stopping ESC listener")
	atomic.StoreInt32(&pollStop, 1)
	C.stop_esc_tap()
}
