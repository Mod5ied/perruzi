#import <Cocoa/Cocoa.h>

// Go callback declared in keyboard_darwin.go.
extern void goEscapeTriggered();

static id globalMonitor = nil;
static id localMonitor  = nil;

// startNSEventMonitors registers Cocoa key-event monitors on the main thread.
// Global monitors fire when Peruzzi is NOT the active app; local monitors fire
// when it IS the active app. Together they cover the cases where the low-level
// CGEventTap may miss the Escape key.
void startNSEventMonitors() {
    if (![NSThread isMainThread]) {
        dispatch_sync(dispatch_get_main_queue(), ^{
            startNSEventMonitors();
        });
        return;
    }

    NSEventMask mask = NSEventMaskKeyDown;

    globalMonitor = [NSEvent addGlobalMonitorForEventsMatchingMask:mask handler:^(NSEvent *event) {
        if ([event keyCode] == 53) { // kVK_Escape
            goEscapeTriggered();
        }
    }];

    localMonitor = [NSEvent addLocalMonitorForEventsMatchingMask:mask handler:^NSEvent *(NSEvent *event) {
        if ([event keyCode] == 53) { // kVK_Escape
            goEscapeTriggered();
        }
        return event;
    }];
}
