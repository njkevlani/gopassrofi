package pass

import (
	"regexp"
	"strings"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripANSI(str string) string {
	return ansiRegex.ReplaceAllString(str, "")
}

// parseTree parses the tree-like output of 'pass' and returns a slice of selectable password paths.
func parseTree(output string) []string {
	lines := strings.Split(output, "\n")
	type lineInfo struct {
		level int
		name  string
	}

	var parsedLines []lineInfo

	for _, line := range lines {
		line = stripANSI(line)
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			continue
		}

		runes := []rune(line)
		// Count prefix runes
		prefixRunesCount := 0
		for _, r := range runes {
			if isTreeRune(r) {
				prefixRunesCount++
			} else {
				break
			}
		}

		if prefixRunesCount == 0 {
			// Skip the root line (e.g. "Password Store" or subdirectory name)
			continue
		}

		// Each level is 4 runes
		level := prefixRunesCount / 4
		if level == 0 {
			continue
		}

		// The name starts after the prefix runes
		name := string(runes[level*4:])
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		parsedLines = append(parsedLines, lineInfo{
			level: level,
			name:  name,
		})
	}

	var passwords []string
	var pathParts []string

	for i, pl := range parsedLines {
		// Update pathParts based on level
		if pl.level > len(pathParts) {
			pathParts = append(pathParts, pl.name)
		} else {
			pathParts = append(pathParts[:pl.level-1], pl.name)
		}

		// Check if it's a leaf node.
		// A node is a leaf if:
		// 1. It is the last line.
		// 2. The next line's level is less than or equal to the current line's level.
		isLeaf := false
		if i == len(parsedLines)-1 {
			isLeaf = true
		} else if parsedLines[i+1].level <= pl.level {
			isLeaf = true
		}

		if isLeaf {
			passwords = append(passwords, strings.Join(pathParts, "/"))
		}
	}

	return passwords
}

func isTreeRune(r rune) bool {
	switch r {
	case '│', '├', '─', '└', ' ', '\u00a0':
		return true
	default:
		return false
	}
}
