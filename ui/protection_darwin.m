#import <Cocoa/Cocoa.h>

// Hides the window from legacy screen-capture APIs.
void setWindowSharingType(void *windowPtr) {
    NSWindow *window = (__bridge NSWindow *)windowPtr;
    if (window == nil) return;
    [window setSharingType:NSWindowSharingNone];
}

// Aggressively hides the window from screen-capture pickers (Chrome/ScreenCaptureKit).
// We order it out, make it fully transparent, and move it off-screen so the
// window server stops including it in capture enumerations.
void hideWindowFromCapture(void *windowPtr) {
    NSWindow *window = (__bridge NSWindow *)windowPtr;
    if (window == nil) return;

    [window setSharingType:NSWindowSharingNone];
    [window setAlphaValue:0.0];
    [window orderOut:nil];

    NSRect frame = [window frame];
    // Move far off-screen so picker thumbnails don't show it.
    frame.origin.x = -20000;
    frame.origin.y = -20000;
    [window setFrame:frame display:NO];
}

// Restores the window after typing is done.
void showWindowForCapture(void *windowPtr) {
    NSWindow *window = (__bridge NSWindow *)windowPtr;
    if (window == nil) return;

    // Center on screen if it ended up off-screen.
    NSRect frame = [window frame];
    if (frame.origin.x <= -10000 || frame.origin.y <= -10000) {
        NSScreen *screen = [window screen] ?: [NSScreen mainScreen];
        NSRect screenFrame = [screen visibleFrame];
        frame.origin.x = screenFrame.origin.x + (screenFrame.size.width - frame.size.width) / 2.0;
        frame.origin.y = screenFrame.origin.y + (screenFrame.size.height - frame.size.height) / 2.0;
        [window setFrame:frame display:NO];
    }

    [window setAlphaValue:1.0];
    [window orderFront:nil];
    [window setSharingType:NSWindowSharingNone];
}
