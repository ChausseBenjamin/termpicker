package util

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea/v2"
)

// ClipboardManager holds pre-computed clipboard strategy information
type ClipboardManager struct {
	initialized           bool
	needsFallback         bool
	incompatibilityReason string
	terminalVersion       string
}

var clipboardManager ClipboardManager

// ClipboardResult represents the result of clipboard operations
type ClipboardResult struct {
	Success        bool
	Message        string
	UsedFallback   bool
	FallbackReason string
}

// ClipboardResultMsg is sent when clipboard operations complete
type ClipboardResultMsg ClipboardResult

// InitClipboard should be called early in the application to query terminal version
func InitClipboard() tea.Cmd {
	return func() tea.Msg {
		return tea.TerminalVersion()
	}
}

// HandleTerminalVersion processes terminal version and pre-computes clipboard strategy
func HandleTerminalVersion(version string) {
	clipboardManager.terminalVersion = version
	clipboardManager.needsFallback, clipboardManager.incompatibilityReason = computeClipboardStrategy(version)
	clipboardManager.initialized = true
}

// computeClipboardStrategy determines clipboard strategy and reason in one pass
func computeClipboardStrategy(version string) (needsFallback bool, reason string) {
	version = strings.ToLower(version)

	// Windows terminals
	if runtime.GOOS == "windows" {
		// Modern Windows terminals with good OSC52 support
		if strings.Contains(version, "windows terminal") ||
			strings.Contains(version, "wt") ||
			strings.Contains(version, "conemu") {
			return false, ""
		}
		// Legacy Windows terminals need fallback
		return true, fmt.Sprintf("Legacy Windows terminal: %s", version)
	}

	// Terminals with known OSC52 issues
	problematicTerminals := map[string]string{
		"screen":         "GNU Screen detected",
		"konsole":        "Older KDE Konsole version",
		"gnome-terminal": "Legacy GNOME Terminal",
		"xterm":          "Plain xterm without OSC52 patches",
		"vt100":          "Legacy VT100 terminal",
		"vt220":          "Legacy VT220 terminal",
		"vt320":          "Legacy VT320 terminal",
	}

	for term, description := range problematicTerminals {
		if strings.Contains(version, term) {
			// Special case: GNU Screen in tmux might work
			if term == "screen" && os.Getenv("TMUX") != "" {
				return false, ""
			}
			return true, fmt.Sprintf("%s (%s)", description, version)
		}
	}

	// Most modern terminals support OSC52
	return false, ""
}

// fallbackToEnvironmentDetection provides environment-based fallback when no terminal version
func fallbackToEnvironmentDetection() (needsFallback bool, reason string) {
	// Windows fallback
	if runtime.GOOS == "windows" {
		return true, "Windows system (environment detection)"
	}

	term := os.Getenv("TERM")

	// GNU Screen without tmux
	if strings.Contains(term, "screen") && os.Getenv("TMUX") == "" {
		return true, "GNU Screen without tmux"
	}

	// Legacy terminals
	if term == "xterm" || term == "vt100" || term == "vt220" {
		return true, "Legacy terminal detected"
	}

	// SSH without proper terminal program
	if os.Getenv("SSH_CONNECTION") != "" && os.Getenv("TERM_PROGRAM") == "" {
		return true, "SSH session with limited terminal support"
	}

	return false, ""
}

// SmartCopyToClipboard uses pre-computed strategy for efficient clipboard operations
func SmartCopyToClipboard(str string) tea.Cmd {
	// Use pre-computed strategy if available, otherwise fall back to environment detection
	var needsFallback bool
	var reason string

	if clipboardManager.initialized {
		needsFallback = clipboardManager.needsFallback
		reason = clipboardManager.incompatibilityReason
	} else {
		needsFallback, reason = fallbackToEnvironmentDetection()
	}

	// Build command list
	commands := []tea.Cmd{tea.SetClipboard(str)} // Always try OSC52

	// Add fallback only if needed
	if needsFallback {
		commands = append(commands, attemptFallbackClipboard(str))
	}

	// Add result message
	commands = append(commands, func() tea.Msg {
		return ClipboardResultMsg{
			Success:        true,
			Message:        fmt.Sprintf("Copied %s to clipboard!", str),
			UsedFallback:   needsFallback,
			FallbackReason: reason,
		}
	})

	return tea.Batch(commands...)
}

// attemptFallbackClipboard tries appropriate fallback methods based on the environment
func attemptFallbackClipboard(str string) tea.Cmd {
	return func() tea.Msg {
		// Try fallback method based on environment
		if err := tryFallbackClipboard(str); err != nil {
			return ClipboardResultMsg{
				Success:        false,
				Message:        fmt.Sprintf("OSC52 attempted, fallback failed: %v", err),
				UsedFallback:   true,
				FallbackReason: "Fallback method failed",
			}
		}
		return nil // Success case handled by main result
	}
}

// tryFallbackClipboard attempts the most appropriate fallback method
func tryFallbackClipboard(str string) error {
	// Try different methods based on environment
	switch runtime.GOOS {
	case "windows":
		return tryWindowsClipboard(str)
	case "darwin":
		return tryMacOSClipboard(str)
	case "linux":
		return tryLinuxClipboard(str)
	default:
		// Use atotto/clipboard as last resort
		return clipboard.WriteAll(str)
	}
}

// tryWindowsClipboard attempts Windows-specific clipboard methods
func tryWindowsClipboard(str string) error {
	// First try atotto/clipboard (uses Windows API)
	if err := clipboard.WriteAll(str); err == nil {
		return nil
	}

	// Fallback to PowerShell
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Set-Clipboard -Value '%s'", str))
	return cmd.Run()
}

// tryMacOSClipboard attempts macOS-specific clipboard methods
func tryMacOSClipboard(str string) error {
	// First try atotto/clipboard
	if err := clipboard.WriteAll(str); err == nil {
		return nil
	}

	// Fallback to pbcopy
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(str)
	return cmd.Run()
}

// tryLinuxClipboard attempts Linux-specific clipboard methods
func tryLinuxClipboard(str string) error {
	// First try atotto/clipboard
	if err := clipboard.WriteAll(str); err == nil {
		return nil
	}

	// Try xclip
	if cmd := exec.Command("xclip", "-selection", "clipboard"); cmd != nil {
		cmd.Stdin = strings.NewReader(str)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	// Try xsel
	if cmd := exec.Command("xsel", "--clipboard", "--input"); cmd != nil {
		cmd.Stdin = strings.NewReader(str)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	// Try wl-copy (Wayland)
	if cmd := exec.Command("wl-copy"); cmd != nil {
		cmd.Stdin = strings.NewReader(str)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("no working clipboard method found")
}

// CopyWarningMessage returns a warning message when fallback is used
func CopyWarningMessage(reason string) string {
	return fmt.Sprintf("Using fallback clipboard method (%s)", reason)
}
