#import <Cocoa/Cocoa.h>
#import <CoreGraphics/CoreGraphics.h>

static void runOnMain(void (^block)(void)) {
    if ([NSThread isMainThread]) {
        block();
    } else {
        dispatch_sync(dispatch_get_main_queue(), block);
    }
}

// Apply the full invisibility configuration:
// - sharingType = none (legacy CGWindowList capture)
// - level = assistiveTechHighWindow (modern ScreenCaptureKit / display capture)
// - collectionBehavior = CanJoinAllSpaces | Stationary | IgnoresCycle
static void applyInvisibility(void *windowPtr) {
    runOnMain(^ {
        NSWindow *window = (__bridge NSWindow *)windowPtr;
        if (window == nil) return;

        [window setSharingType:NSWindowSharingNone];
        [window setLevel:CGWindowLevelForKey(kCGAssistiveTechHighWindowLevelKey)];
        [window setCollectionBehavior:
            NSWindowCollectionBehaviorCanJoinAllSpaces |
            NSWindowCollectionBehaviorStationary |
            NSWindowCollectionBehaviorIgnoresCycle];
    });
}

void setWindowSharingType(void *windowPtr) {
    applyInvisibility(windowPtr);
}

// Aggressively hides the window during typing.
void hideWindowFromCapture(void *windowPtr) {
    runOnMain(^ {
        NSWindow *window = (__bridge NSWindow *)windowPtr;
        if (window == nil) return;

        [window setSharingType:NSWindowSharingNone];
        [window setAlphaValue:0.0];
        [window orderOut:nil];

        NSRect frame = [window frame];
        frame.origin.x = -20000;
        frame.origin.y = -20000;
        [window setFrame:frame display:NO];
    });
}

// Restores the window after typing is done.
void showWindowForCapture(void *windowPtr) {
    runOnMain(^ {
        NSWindow *window = (__bridge NSWindow *)windowPtr;
        if (window == nil) return;

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
        applyInvisibility(windowPtr);
    });
}
