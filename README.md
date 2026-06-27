# gopassrofi

A Golang-based tool that integrates the standard Unix Password Store (`pass`) with `rofi` (in `dmenu` mode) to provide a rich, interactive graphical password manager.

## Features

- **Select Passwords**: Fuzzy search and select your passwords using `rofi`.
- **Keyboard Shortcuts**:
  - `Enter`: Copies the password to the clipboard (and clears it after 45 seconds).
  - `Alt + Enter`: Copies the 2FA OTP code to the clipboard (using `pass otp`).
  - `Ctrl + Enter`: Displays the **Other Options** submenu.
- **Other Options Submenu**:
  - Copy Password
  - Copy OTP
  - Show Password (displays in a `rofi` pop-up dialog)
  - Show Full Entry (displays the password file metadata/additional lines)
  - Generate New Password (prompts for length and overwrites the entry)
  - Edit Password Manually (prompts securely for a custom password)
  - Delete Password (safely deletes the password from the store)
- **Create New Password**:
  - Selection of `+ Create New Password` (or typing a name that does not exist and pressing Enter) prompts to create a new password.
  - Choose between generating a random password (16, 24, or 32 characters) or entering one manually.
- **Desktop Notifications**: Uses `notify-send` to confirm when passwords/OTPs are copied or created.

## Installation

### Prerequisites

1. **Golang**: Ensure Go is installed (version 1.16+ recommended).
2. **Password Store**: Make sure `pass` is installed and initialized.
3. **Rofi**: Ensure `rofi` is installed.
4. **Notify Send**: (Optional) For desktop notifications.

### Building

Initialize the module and compile the binary:

```bash
go build -o gopassrofi main.go
```

This will produce a `gopassrofi` executable in the root directory.

## Usage

Run the compiled binary:

```bash
./gopassrofi
```

### Integration with Window Managers

To assign `gopassrofi` to a global hotkey, add a configuration to your window manager:

#### i3wm / Sway
```i3
bindsym $mod+p exec --no-startup-id /path/to/gopassrofi
```

#### sxhkd
```sxhkdrc
super + p
    /path/to/gopassrofi
```

## Project Structure

- `internal/pass/`: Handles interacting with the `pass` CLI, parsing the output tree, and managing passwords/OTPs.
  - [pass.go](file:///home/njkevlani/git/gopassrofi/internal/pass/pass.go): Exported library APIs.
  - [parser.go](file:///home/njkevlani/git/gopassrofi/internal/pass/parser.go): Tree-parsing functions.
  - [parser_test.go](file:///home/njkevlani/git/gopassrofi/internal/pass/parser_test.go): Tests for parser correctness.
- `internal/dmenu/`: Wraps `rofi` execution and inputs.
  - [dmenu.go](file:///home/njkevlani/git/gopassrofi/internal/dmenu/dmenu.go): Dmenu display, inputs, messages, and notifications.
- [main.go](file:///home/njkevlani/git/gopassrofi/main.go): Core orchestrator loop.
