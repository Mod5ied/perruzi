package engine

import (
	"math/rand"
	"time"
)

var fastBigrams = map[string]bool{
	"th": true, "he": true, "in": true, "er": true, "an": true,
	"on": true, "re": true, "nd": true, "at": true, "es": true,
	"en": true, "of": true, "to": true, "it": true, "is": true,
	"or": true, "te": true, "ti": true, "ar": true, "st": true,
}

var slowBigrams = map[string]bool{
	"zx": true, "xz": true, "qj": true, "jq": true,
	"vb": true, "bv": true, "km": true, "mk": true,
}

var adjacentKeys = map[rune]string{
	'a': "sq", 'b': "vn", 'c': "xv", 'd': "sf", 'e': "wr",
	'f': "dg", 'g': "fh", 'h': "gj", 'i': "uo", 'j': "hk",
	'k': "jl", 'l': "k", 'm': "n", 'n': "mb", 'o': "ip",
	'p': "o", 'q': "w", 'r': "et", 's': "ad", 't': "ry",
	'u': "yi", 'v': "cb", 'w': "qe", 'x': "zc", 'y': "tu", 'z': "x",
}

type Humaniser struct {
	baseMs          float64
	prevChar        rune
	velocityCounter int
	velocityTarget  int
	inBurst         bool
	burstCount      int
}

func NewHumaniser(base time.Duration) *Humaniser {
	return &Humaniser{
		baseMs:         float64(base) / float64(time.Millisecond),
		velocityTarget: 8 + rand.Intn(8),
	}
}

// Delay returns the delay to wait before injecting curr.
func (h *Humaniser) Delay(curr rune) time.Duration {
	// Step 1: Base jitter ±20%.
	jitter := h.baseMs * (rand.Float64()*0.4 - 0.2)
	delay := h.baseMs + jitter

	// Step 2: Bigram multiplier.
	if h.prevChar != 0 {
		pair := string(h.prevChar) + string(curr)
		if fastBigrams[pair] {
			delay *= 0.6 + rand.Float64()*0.25
		} else if slowBigrams[pair] {
			delay *= 1.2 + rand.Float64()*0.4
		}
	}
	h.prevChar = curr

	// Step 3: Word boundary pause.
	if curr == ' ' {
		delay += 20 + rand.Float64()*60
	}

	// Step 4: Punctuation hesitation.
	switch curr {
	case '.', ',', '!', '?', ';', ':':
		delay += 50 + rand.Float64()*100
	}

	// Step 5: Flow velocity.
	h.velocityCounter++
	if h.inBurst {
		delay *= 0.7 + rand.Float64()*0.1
		h.burstCount++
		if h.burstCount >= 3+rand.Intn(3) { // 3 to 5 chars in burst
			h.inBurst = false
			h.burstCount = 0
		}
	}
	if h.velocityCounter >= h.velocityTarget {
		h.inBurst = true
		h.burstCount = 0
		h.velocityCounter = 0
		h.velocityTarget = 8 + rand.Intn(8)
	}

	if delay < 5 {
		delay = 5
	}
	return time.Duration(delay) * time.Millisecond
}

// ShouldMistype returns whether to simulate a typo and which wrong character to inject.
func (h *Humaniser) ShouldMistype(curr rune) (bool, rune) {
	if h.baseMs <= 40 {
		return false, 0
	}
	if rand.Float64() > 0.02 {
		return false, 0
	}
	if curr < 'a' || curr > 'z' {
		return false, 0
	}
	adjacent, ok := adjacentKeys[curr]
	if !ok || adjacent == "" {
		return false, 0
	}
	return true, rune(adjacent[rand.Intn(len(adjacent))])
}
