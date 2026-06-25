package keyboard

// Injector is the platform-specific keyboard injection implementation.
// Both keyboard_darwin.go and keyboard_windows.go implement this.
type Injector interface {
	// InjectChar injects a single Unicode character into the active window.
	InjectChar(r rune) error

	// InjectReturn injects a Return/Enter keypress.
	InjectReturn() error

	// InjectBackspace injects a single Backspace keypress.
	InjectBackspace() error
}
