package main

import (
	"fmt"
	"gopassrofi/internal/dmenu"
	"gopassrofi/internal/pass"
)

func main() {
	for {
		passwords, err := pass.GetPasswords()
		if err != nil {
			_ = dmenu.ShowMessage(fmt.Sprintf("Error loading passwords: %v", err))
			return
		}

		// Prepend a create option to the list
		items := append([]string{"+ Create New Password"}, passwords...)

		res, err := dmenu.ShowMenu("pass", items, "")
		if err != nil {
			_ = dmenu.ShowMessage(fmt.Sprintf("Error: %v", err))
			return
		}

		if res.Action == dmenu.ActionCancel {
			break
		}

		// Handle selection of "+ Create New Password"
		if res.Selection == "+ Create New Password" {
			if handleCreatePassword("") {
				return
			}
			continue
		}

		// Check if selected item is an existing password
		isExisting := false
		for _, p := range passwords {
			if p == res.Selection {
				isExisting = true
				break
			}
		}

		// If user typed a new password name that doesn't exist, prompt to create it
		if !isExisting {
			if res.Selection != "" {
				confirm, err := dmenu.ShowMenu("Create Password?", []string{"Yes, create", "Cancel"}, fmt.Sprintf("Create new password entry '%s'?", res.Selection))
				if err == nil && confirm.Selection == "Yes, create" {
					if handleCreatePassword(res.Selection) {
						return
					}
				}
			}
			continue
		}

		// Existing password selected
		switch res.Action {
		case dmenu.ActionAccept: // Enter -> Copy Password
			err := pass.CopyPassword(res.Selection, "gopassrofi", fmt.Sprintf("Password for '%s' copied to clipboard!", res.Selection))
			if err != nil {
				_ = dmenu.ShowMessage(fmt.Sprintf("Error copying password: %v", err))
			} else {
				return
			}

		case dmenu.ActionCustom1: // Alt+Enter -> Copy OTP
			err := pass.CopyOTP(res.Selection, "gopassrofi", fmt.Sprintf("OTP for '%s' copied to clipboard!", res.Selection))
			if err != nil {
				_ = dmenu.ShowMessage(fmt.Sprintf("Error copying OTP: %v", err))
			} else {
				return
			}

		case dmenu.ActionCustom2: // Ctrl+Enter -> Other Options
			if showOtherOptions(res.Selection) {
				return
			}
		}
	}
}

// handleCreatePassword guides the user through creating a new password.
// Returns true if a password was successfully created and copied (so the application should exit),
// or false if the user cancelled or creation failed (so it should return to the main menu).
func handleCreatePassword(defaultName string) bool {
	name, err := dmenu.ShowInputPrompt("Enter path/name for new password:", defaultName)
	if err != nil {
		_ = dmenu.ShowMessage(fmt.Sprintf("Error: %v", err))
		return true // Exit on critical error
	}
	if name == "" {
		return false // Cancelled
	}

	createOpts := []string{
		"Generate password (24 chars)",
		"Generate password (16 chars)",
		"Generate password (32 chars)",
		"Enter password manually",
	}

	choice, err := dmenu.ShowMenu("Create Options", createOpts, "Configure password for "+name)
	if err != nil {
		_ = dmenu.ShowMessage(fmt.Sprintf("Error: %v", err))
		return true // Exit on critical error
	}
	if choice.Action == dmenu.ActionCancel {
		return false // Cancelled
	}

	switch choice.Selection {
	case "Generate password (24 chars)":
		_, err = pass.GeneratePassword(name, 24)
	case "Generate password (16 chars)":
		_, err = pass.GeneratePassword(name, 16)
	case "Generate password (32 chars)":
		_, err = pass.GeneratePassword(name, 32)
	case "Enter password manually":
		var pwd string
		pwd, err = dmenu.ShowPasswordPrompt("Enter password for " + name)
		if err != nil {
			_ = dmenu.ShowMessage(fmt.Sprintf("Error: %v", err))
			return true // Exit on critical error
		}
		if pwd == "" {
			return false // Cancelled
		}
		err = pass.CreatePassword(name, pwd)
	default:
		return false
	}

	if err != nil {
		_ = dmenu.ShowMessage(fmt.Sprintf("Error creating password: %v", err))
		return false
	}

	// Copy the newly created password to clipboard
	err = pass.CopyPassword(name, "gopassrofi", fmt.Sprintf("Password '%s' created and copied to clipboard!", name))
	if err != nil {
		_ = dmenu.ShowMessage(fmt.Sprintf("Password created, but copying failed: %v", err))
		return false
	}

	return true
}

// showOtherOptions displays the secondary options menu for an existing password.
// Returns true if the user triggered an action that copies the password/OTP and should exit the script.
func showOtherOptions(entryName string) bool {
	options := []string{
		"Copy Password",
		"Copy OTP",
		"Show Password",
		"Show Full Entry",
		"Generate New Password (overwrite)",
		"Edit Password Manually",
		"Delete Password",
	}

	res, err := dmenu.ShowMenu(entryName, options, "Select action for "+entryName)
	if err != nil {
		_ = dmenu.ShowMessage(fmt.Sprintf("Error: %v", err))
		return true // Exit on critical error
	}

	if res.Action == dmenu.ActionCancel {
		return false
	}

	switch res.Selection {
	case "Copy Password":
		err := pass.CopyPassword(entryName, "gopassrofi", fmt.Sprintf("Password for '%s' copied to clipboard!", entryName))
		if err != nil {
			_ = dmenu.ShowMessage(fmt.Sprintf("Error copying password: %v", err))
		} else {
			return true
		}

	case "Copy OTP":
		err := pass.CopyOTP(entryName, "gopassrofi", fmt.Sprintf("OTP for '%s' copied to clipboard!", entryName))
		if err != nil {
			_ = dmenu.ShowMessage(fmt.Sprintf("Error copying OTP: %v", err))
		} else {
			return true
		}

	case "Show Password":
		pwd, err := pass.GetPassword(entryName)
		if err != nil {
			_ = dmenu.ShowMessage(fmt.Sprintf("Error getting password: %v", err))
		} else {
			_ = dmenu.ShowMessage(fmt.Sprintf("Password for %s:\n\n%s", entryName, pwd))
		}

	case "Show Full Entry":
		entry, err := pass.GetPasswordEntry(entryName)
		if err != nil {
			_ = dmenu.ShowMessage(fmt.Sprintf("Error getting password entry: %v", err))
		} else {
			_ = dmenu.ShowMessage(fmt.Sprintf("Entry for %s:\n\n%s", entryName, entry))
		}

	case "Generate New Password (overwrite)":
		confirm, err := dmenu.ShowMenu("Confirm Overwrite", []string{"Cancel", "Yes, overwrite"}, "Are you sure you want to overwrite password for "+entryName+"?")
		if err != nil {
			_ = dmenu.ShowMessage(fmt.Sprintf("Error: %v", err))
			return true // Exit on critical error
		}
		if confirm.Selection != "Yes, overwrite" {
			return false
		}

		lenOpts := []string{"24 chars", "16 chars", "32 chars"}
		lenRes, err := dmenu.ShowMenu("Password Length", lenOpts, "Select length for generated password")
		if err != nil {
			_ = dmenu.ShowMessage(fmt.Sprintf("Error: %v", err))
			return true // Exit on critical error
		}
		if lenRes.Action == dmenu.ActionCancel {
			return false
		}

		length := 24
		switch lenRes.Selection {
		case "16 chars":
			length = 16
		case "32 chars":
			length = 32
		}

		_, err = pass.GeneratePassword(entryName, length)
		if err != nil {
			_ = dmenu.ShowMessage(fmt.Sprintf("Error generating password: %v", err))
		} else {
			_ = pass.CopyPassword(entryName, "gopassrofi", fmt.Sprintf("New password generated and copied for '%s'!", entryName))
			return true
		}

	case "Edit Password Manually":
		pwd, err := dmenu.ShowPasswordPrompt("Enter new password for " + entryName)
		if err != nil {
			_ = dmenu.ShowMessage(fmt.Sprintf("Error: %v", err))
			return true // Exit on critical error
		}
		if pwd == "" {
			return false
		}

		err = pass.CreatePassword(entryName, pwd)
		if err != nil {
			_ = dmenu.ShowMessage(fmt.Sprintf("Error updating password: %v", err))
		} else {
			_ = pass.CopyPassword(entryName, "gopassrofi", fmt.Sprintf("Password updated and copied for '%s'!", entryName))
			return true
		}

	case "Delete Password":
		confirm, err := dmenu.ShowMenu("Confirm Delete", []string{"Cancel", "Yes, delete"}, "Are you sure you want to delete password for "+entryName+"?")
		if err != nil {
			_ = dmenu.ShowMessage(fmt.Sprintf("Error: %v", err))
			return true // Exit on critical error
		}
		if confirm.Selection != "Yes, delete" {
			return false
		}

		err = pass.DeletePassword(entryName)
		if err != nil {
			_ = dmenu.ShowMessage(fmt.Sprintf("Error deleting password: %v", err))
		} else {
			dmenu.NotifySend("gopassrofi", fmt.Sprintf("Password for '%s' deleted!", entryName))
		}
	}

	return false
}
