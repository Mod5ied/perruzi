package ui

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"Peruzzi/engine"
	"Peruzzi/keyboard"
	"Peruzzi/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/systray"
)

// MainWindow wraps the Fyne application window.
type MainWindow struct {
	Window fyne.Window
}

// NewMainWindow builds and wires the main application window.
func NewMainWindow(app fyne.App) *MainWindow {
	w := app.NewWindow("Peruzzi")
	w.Resize(fyne.NewSize(380, 300))
	w.SetFixedSize(true)
	w.SetCloseIntercept(func() { w.Hide() })

	textInput := widget.NewMultiLineEntry()
	textInput.SetPlaceHolder("// paste text here...")
	textInput.TextStyle = fyne.TextStyle{Monospace: true}
	textInput.Wrapping = fyne.TextWrapWord
	textInput.SetMinRowsVisible(3)

	speedSlider := widget.NewSlider(5, 200)
	speedSlider.Value = 90
	speedSlider.Step = 1

	delayLabel := widget.NewLabel("0.090s")
	delayLabel.Alignment = fyne.TextAlignTrailing
	delayLabel.TextStyle = fyne.TextStyle{Monospace: true}

	humaniseCheck := NewHumaniseCheck(true, nil)

	startBtn := NewNeuButton("Start", nil, neuTeal, neuTeal, nil)

	statusLabel := widget.NewLabel("Ready")
	statusLabel.TextStyle = fyne.TextStyle{Monospace: true}
	statusLabel.Alignment = fyne.TextAlignCenter

	var (
		activeEngine *engine.TypingEngine
		engineMu     sync.Mutex
	)

	// Single ESC listener for the lifetime of the window.
	go func() {
		for {
			<-keyboard.EscChan
			logger.Log("UI received ESC signal")
			engineMu.Lock()
			eng := activeEngine
			engineMu.Unlock()
			if eng != nil {
				logger.Log("UI stopping active engine")
				eng.Stop()
			} else {
				logger.Log("UI received ESC but no active engine")
			}
		}
	}()

	speedSlider.OnChanged = func(v float64) {
		seconds := v / 1000.0
		delayLabel.SetText(fmt.Sprintf("%.3fs", seconds))
	}

	startBtn.OnTapped = func() {
		text := strings.TrimSpace(textInput.Text)
		if text == "" {
			return
		}

		if !keyboard.IsReady() {
			statusLabel.SetText("Permission needed")
			if runtime.GOOS == "darwin" {
				dialog.ShowInformation(
					"Permission Required",
					"Peruzzi needs Accessibility access to inject keystrokes.\n\nGo to:\nSystem Settings > Privacy & Security > Accessibility\n\nEnable Peruzzi, then try again.",
					w,
				)
				exec.Command("open", "x-apple.systempreferences:com.apple.preference.security?Privacy_Accessibility").Start()
			}
			return
		}

		// Drain any stale ESC signal before starting a new session.
		select {
		case <-keyboard.EscChan:
		default:
		}

		startBtn.Disable()
		humaniseCheck.Disable()

		injector, err := keyboard.NewInjector()
		if err != nil {
			statusLabel.SetText("Keyboard error")
			startBtn.Enable()
			humaniseCheck.Enable()
			return
		}

		base := time.Duration(speedSlider.Value) * time.Millisecond
		currentEngine := engine.NewTypingEngine(text, base, humaniseCheck.IsActive(), injector)

		engineMu.Lock()
		activeEngine = currentEngine
		engineMu.Unlock()

		currentEngine.OnTick = func(remaining int) {
			statusLabel.SetText(fmt.Sprintf("%d...", remaining))
			systray.SetTitle(fmt.Sprintf("%d", remaining))
		}
		currentEngine.OnTypingStart = func() {
			w.Hide()
		}
		currentEngine.OnComplete = func() {
			statusLabel.SetText("Done!")
			startBtn.Enable()
			humaniseCheck.Enable()
			systray.SetTitle("")
			w.Show()
			engineMu.Lock()
			if activeEngine == currentEngine {
				activeEngine = nil
			}
			engineMu.Unlock()
		}
		currentEngine.OnStopped = func() {
			statusLabel.SetText("Stopped by ESC or error")
			startBtn.Enable()
			humaniseCheck.Enable()
			systray.SetTitle("")
			w.Show()
			engineMu.Lock()
			if activeEngine == currentEngine {
				activeEngine = nil
			}
			engineMu.Unlock()
		}

		// Hide Peruzzi so it cannot steal focus during the countdown/typing.
		w.Hide()
		currentEngine.Start()
		statusLabel.SetText("Switch to target input...")
	}

	// Input area wrapped in a neubrutalist box with internal padding.
	inputBox := NewNeuBox(container.NewPadded(textInput), neuWhite, neuBlack)

	// Speed label + value row.
	speedHeader := container.NewBorder(nil, nil, nil, delayLabel,
		widget.NewLabelWithStyle("TYPING SPEED:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	// Humanise toggle row — checkbox is wrapped to keep it a perfect square,
	// label is vertically centred with the checkbox.
	humaniseSquare := container.NewGridWrap(fyne.NewSize(16, 16), humaniseCheck)
	humaniseLabel := canvas.NewText("Humanise Mode (Natural variation)", neuBlack)
	humaniseLabel.TextStyle = fyne.TextStyle{Bold: true}
	humaniseLabelBox := container.NewVBox(layout.NewSpacer(), humaniseLabel, layout.NewSpacer())
	humaniseRow := container.NewHBox(
		humaniseSquare,
		humaniseLabelBox,
		layout.NewSpacer(),
	)

	// Status sits just below the controls with no extra filler.
	controls := container.NewVBox(
		inputBox,
		speedHeader,
		speedSlider,
		humaniseRow,
		startBtn,
		statusLabel,
	)

	card := NewNeuBox(controls, neuWhite, neuTeal)

	w.SetContent(container.NewPadded(card))

	// Hide the window from screen capture / screen sharing.
	SetContentProtection(w)

	return &MainWindow{Window: w}
}
