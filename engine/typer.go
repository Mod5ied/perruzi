package engine

import (
	"Peruzzi/keyboard"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var logFile *os.File
var logOnce sync.Once

func logEvent(format string, args ...interface{}) {
	logOnce.Do(func() {
		var err error
		logFile, err = os.OpenFile("/tmp/peruzzi.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return
		}
	})
	if logFile == nil {
		return
	}
	msg := fmt.Sprintf(format, args...)
	logFile.WriteString(time.Now().Format("15:04:05.000") + " " + msg + "\n")
	logFile.Sync()
}

// TypingEngine coordinates countdown and character-by-character injection.
type TypingEngine struct {
	text      string
	base      time.Duration
	humanise  bool
	injector  keyboard.Injector
	humaniser *Humaniser

	stopOnce sync.Once
	stopCh   chan struct{}

	OnTick        func(remaining int)
	OnTypingStart func()
	OnComplete    func()
	OnStopped     func()
}

// NewTypingEngine creates a new typing engine.
func NewTypingEngine(
	text string,
	base time.Duration,
	humanise bool,
	injector keyboard.Injector,
) *TypingEngine {
	return &TypingEngine{
		text:      text,
		base:      base,
		humanise:  humanise,
		injector:  injector,
		humaniser: NewHumaniser(base),
		stopCh:    make(chan struct{}),
	}
}

// Start begins the countdown and typing loop in a goroutine.
func (t *TypingEngine) Start() {
	go t.run()
}

// Stop signals the engine to stop.
func (t *TypingEngine) Stop() {
	t.stopOnce.Do(func() { close(t.stopCh) })
}

func (t *TypingEngine) run() {
	defer func() {
		if r := recover(); r != nil {
			logEvent("PANIC in typing goroutine: %v", r)
			if t.OnStopped != nil {
				t.OnStopped()
			}
		}
	}()

	logEvent("Engine started, text length=%d humanise=%v base=%v", len(t.text), t.humanise, t.base)

	// 1. Countdown from 5 to 1.
	for i := 5; i >= 1; i-- {
		if t.isStopped() {
			logEvent("Stopped during countdown")
			if t.OnStopped != nil {
				t.OnStopped()
			}
			return
		}
		if t.OnTick != nil {
			t.OnTick(i)
		}
		logEvent("Countdown %d", i)
		select {
		case <-t.stopCh:
			logEvent("Stopped during countdown sleep")
			if t.OnStopped != nil {
				t.OnStopped()
			}
			return
		case <-time.After(1 * time.Second):
		}
	}

	if t.OnTypingStart != nil {
		logEvent("Typing start")
		t.OnTypingStart()
	}

	// 2. Loop over lines and characters.
	lines := strings.Split(t.text, "\n")
	charCount := 0
	for lineIndex, line := range lines {
		if lineIndex > 0 {
			if t.isStopped() {
				logEvent("Stopped before newline return")
				if t.OnStopped != nil {
					t.OnStopped()
				}
				return
			}
			logEvent("Injecting Return")
			t.injector.InjectReturn()
			t.sleep(t.base)
		}

		for _, r := range line {
			if t.isStopped() {
				logEvent("Stopped before char %d", charCount)
				if t.OnStopped != nil {
					t.OnStopped()
				}
				return
			}

			if t.humanise {
				if mistype, wrong := t.humaniser.ShouldMistype(r); mistype {
					logEvent("Mistype char %d: injecting wrong %q then backspace", charCount, string(wrong))
					t.injector.InjectChar(wrong)
					t.sleep(80 * time.Millisecond)
					t.injector.InjectBackspace()
					t.sleep(40 * time.Millisecond)
				}
				logEvent("Injecting char %d: %q delay=%v", charCount, string(r), t.humaniser.Delay(r))
				t.injector.InjectChar(r)
				t.sleep(t.humaniser.Delay(r))
			} else {
				logEvent("Injecting char %d: %q delay=%v", charCount, string(r), t.base)
				t.injector.InjectChar(r)
				t.sleep(t.base)
			}
			charCount++
		}
	}

	logEvent("Typing complete, injected %d chars", charCount)
	if t.OnComplete != nil {
		t.OnComplete()
	}
}

func (t *TypingEngine) isStopped() bool {
	select {
	case <-t.stopCh:
		return true
	default:
		return false
	}
}

func (t *TypingEngine) sleep(d time.Duration) {
	select {
	case <-t.stopCh:
		logEvent("Sleep interrupted by stop")
	case <-time.After(d):
	}
}
