# Peruzzi

Peruzzi is a free, open-source keyboard automation app for macOS and Windows.

Most tools that inject keystrokes at the OS level hide behind subscriptions, activation keys, or cloud accounts. Peruzzi does not. It was built to give anyone a simple, local, forever-free way to automate typing into any window on their machine.

Paste your text, set the speed, toggle Humanise Mode if you want natural variation, and press **Start**. A short countdown gives you time to switch to the target input; Peruzzi handles the rest. Press **ESC** at any moment and it stops instantly.

## Why it is different

- **Native OS injection** — CGEvent on macOS, SendInput on Windows. No brittle clipboard pasting or app-specific scripting.
- **Keyboard-layout blind** — it outputs the right characters regardless of QWERTY, AZERTY, Dvorak, or any other layout.
- **Humanise Mode** — adds realistic timing variation, hesitations, and occasional typo-then-correction loops so the output looks hand-typed.
- **No accounts, no licence checks, no telemetry** — download, open, and use it. That is the whole process.
- **Open source** — MIT licence. Read the code, fork it, or build on top of it.

## Downloads

- [Download for macOS](https://github.com/Mod5ied/perruzi/releases)
- [Download for Windows](https://github.com/Mod5ied/perruzi/releases)

## Build from source

```bash
git clone https://github.com/Mod5ied/perruzi.git
cd perruzi
go build -o Peruzzi .
```

On Windows you will need MinGW-w64 installed for CGO.

## Licence

MIT
