package dmenu

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ActionType represents how the user exited the rofi menu.
type ActionType int

const (
	ActionAccept ActionType = iota // Enter (Exit 0)
	ActionCancel                   // Escape (Exit 1)
	ActionCustom1                  // Alt+Enter (Exit 10)
	ActionCustom2                  // Ctrl+Enter (Exit 11)
	ActionUnknown
)

// Result represents the selection and the action taken by the user.
type Result struct {
	Selection string
	Action    ActionType
}

// ShowMenu displays a rofi dmenu with the given prompt and items.
// It maps Alt+Enter to ActionCustom1 and Ctrl+Enter to ActionCustom2.
func ShowMenu(prompt string, items []string, message string) (Result, error) {
	cmd := exec.Command("rofi",
		"-dmenu",
		"-p", prompt,
		"-kb-accept-custom", "", // Clear Control+Return default binding to avoid conflicts
		"-kb-custom-1", "Alt+Return",
		"-kb-custom-2", "Control+Return",
		"-i", // Case-insensitive matching
	)

	if message != "" {
		cmd.Args = append(cmd.Args, "-mesg", message)
	}

	cmd.Stdin = strings.NewReader(strings.Join(items, "\n"))

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ProcessState.ExitCode()
		} else {
			return Result{}, fmt.Errorf("failed to run rofi: %w (stderr: %q)", err, stderr.String())
		}
	}

	selection := strings.TrimSpace(stdout.String())

	var action ActionType
	switch exitCode {
	case 0:
		action = ActionAccept
	case 1:
		action = ActionCancel
	case 10:
		action = ActionCustom1
	case 11:
		action = ActionCustom2
	default:
		return Result{}, fmt.Errorf("rofi exited with unexpected code %d (stderr: %q)", exitCode, stderr.String())
	}

	return Result{
		Selection: selection,
		Action:    action,
	}, nil
}

// ShowPasswordPrompt displays a rofi password-hidden input field.
func ShowPasswordPrompt(prompt string) (string, error) {
	cmd := exec.Command("rofi",
		"-dmenu",
		"-password",
		"-p", prompt,
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ProcessState.ExitCode() == 1 {
			return "", nil // User cancelled
		}
		return "", fmt.Errorf("failed to get password: %w (stderr: %q)", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

// ShowInputPrompt displays a generic rofi text input field.
func ShowInputPrompt(prompt string, defaultVal string) (string, error) {
	args := []string{"-dmenu", "-p", prompt}
	if defaultVal != "" {
		args = append(args, "-filter", defaultVal)
	}
	cmd := exec.Command("rofi", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ProcessState.ExitCode() == 1 {
			return "", nil // User cancelled
		}
		return "", fmt.Errorf("failed to get input: %w (stderr: %q)", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

// ShowMessage displays a pop-up dialog with the given message using rofi -e.
func ShowMessage(msg string) error {
	cmd := exec.Command("rofi", "-e", msg)
	return cmd.Run()
}

// NotifySend fires a desktop notification using notify-send.
func NotifySend(title, body string) {
	_ = exec.Command("notify-send", title, body).Run()
}
