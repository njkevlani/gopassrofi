package pass

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// runPassCmd runs the pass command with the given arguments and optional stdin.
// It injects PASSWORD_STORE_GPG_OPTS=--trust-model always to avoid interactive GPG trust prompts.
func runPassCmd(args []string, stdinInput string) (string, error) {
	cmd := exec.Command("pass", args...)
	cmd.Env = append(os.Environ(), "PASSWORD_STORE_GPG_OPTS=--trust-model always")
	if stdinInput != "" {
		cmd.Stdin = strings.NewReader(stdinInput)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command 'pass %s' failed: %w (stderr: %q)", strings.Join(args, " "), err, stderrBuf.String())
	}

	return stdoutBuf.String(), nil
}

// GetPasswords runs the pass command to list all passwords and parses the tree output.
func GetPasswords() ([]string, error) {
	output, err := runPassCmd(nil, "")
	if err != nil {
		return nil, err
	}
	return parseTree(output), nil
}

// GetPassword retrieves the decrypted password (first line) for the given name.
func GetPassword(name string) (string, error) {
	output, err := runPassCmd([]string{"show", name}, "")
	if err != nil {
		return "", err
	}

	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("empty output from pass show")
	}

	return strings.TrimRight(lines[0], "\r\n"), nil
}

// GetPasswordEntry retrieves the full decrypted text of the password file.
func GetPasswordEntry(name string) (string, error) {
	return runPassCmd([]string{"show", name}, "")
}

// GetOTP retrieves the current OTP code for the given name using pass otp.
func GetOTP(name string) (string, error) {
	output, err := runPassCmd([]string{"otp", name}, "")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// CreatePassword saves a password manually to the given name.
func CreatePassword(name string, password string) error {
	_, err := runPassCmd([]string{"insert", "-f", "-e", name}, password+"\n")
	return err
}

// GeneratePassword generates a password for the given name and length, returning the generated password.
func GeneratePassword(name string, length int) (string, error) {
	args := []string{"generate", "-f", name}
	if length > 0 {
		args = append(args, fmt.Sprintf("%d", length))
	}

	_, err := runPassCmd(args, "")
	if err != nil {
		return "", err
	}

	return GetPassword(name)
}

// CopyPassword copies the password to the clipboard using pass show -c asynchronously.
// If notifyTitle is not empty, it will display a desktop notification using notify-send after the copy finishes.
func CopyPassword(name string, notifyTitle, notifyBody string) error {
	var script string
	if notifyTitle != "" {
		script = fmt.Sprintf("pass show -c %q && notify-send %q %q", name, notifyTitle, notifyBody)
	} else {
		script = fmt.Sprintf("pass show -c %q", name)
	}
	cmd := exec.Command("sh", "-c", script)
	cmd.Env = append(os.Environ(), "PASSWORD_STORE_GPG_OPTS=--trust-model always")
	return cmd.Start()
}

// CopyOTP copies the OTP to the clipboard using pass otp -c asynchronously.
// If notifyTitle is not empty, it will display a desktop notification using notify-send after the copy finishes.
func CopyOTP(name string, notifyTitle, notifyBody string) error {
	var script string
	if notifyTitle != "" {
		script = fmt.Sprintf("pass otp -c %q && notify-send %q %q", name, notifyTitle, notifyBody)
	} else {
		script = fmt.Sprintf("pass otp -c %q", name)
	}
	cmd := exec.Command("sh", "-c", script)
	cmd.Env = append(os.Environ(), "PASSWORD_STORE_GPG_OPTS=--trust-model always")
	return cmd.Start()
}

// DeletePassword deletes the password entry.
func DeletePassword(name string) error {
	_, err := runPassCmd([]string{"rm", "-f", name}, "")
	return err
}


