//go:build darwin

package keyboard

/*
#cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation -framework ApplicationServices

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

// Forward declaration of the Go callback.
extern void goEscapeTriggered();

static CFMachPortRef escTapRef = NULL;
static CFRunLoopSourceRef escRunLoopSource = NULL;

static CGEventRef escCallback(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *refcon) {
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
        kCGEventTapOptionListenOnly,
        CGEventMaskBit(kCGEventKeyDown),
        escCallback,
        NULL
    );
    if (escTapRef == NULL) return;
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
import "runtime"

// EscChan is the package-level ESC signal channel. Engine reads from this.
var EscChan = make(chan struct{}, 1)

//export goEscapeTriggered
func goEscapeTriggered() {
	select {
	case EscChan <- struct{}{}:
	default:
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

// StartEscListener starts the CGEventTap on a locked OS thread.
// Call this once at app startup. It runs until StopEscListener is called.
func StartEscListener() {
	go func() {
		runtime.LockOSThread()
		C.start_esc_tap()
	}()
}

// StopEscListener stops the CGEventTap.
func StopEscListener() {
	C.stop_esc_tap()
}
