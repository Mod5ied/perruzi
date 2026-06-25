#import <Cocoa/Cocoa.h>

void setWindowSharingType(void *windowPtr) {
    NSWindow *window = (__bridge NSWindow *)windowPtr;
    [window setSharingType:NSWindowSharingNone];
}
