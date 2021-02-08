package dtaservice

import (
	"fmt"
	"strings"
)

const (
	// RepoName is the name of this repository
	RepoName string = "github.com/theovassiliou/doctrans-framework"

	// Version contains the actuall version number. Might be replaces using the LDFLAGS.
	Version = "1.1.0-src"
)

// FormatFullVersion formats for a cmdName the version number based on version, branch and commit
func FormatFullVersion(cmdName, version, branch, commit string) string {
	var parts = []string{cmdName}

	if version != "" {
		parts = append(parts, version)
	} else {
		parts = append(parts, "unknown")
	}

	if branch != "" || commit != "" {
		if branch == "" {
			branch = "unknown"
		}
		if commit == "" {
			commit = "unknown"
		}
		git := fmt.Sprintf("(git: %s %s)", branch, commit)
		parts = append(parts, git)
	}

	return strings.Join(parts, " ")
}
