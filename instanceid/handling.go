package instanceid

import (
	"strings"
)

// Returns a miid as string in the form appName/version/branch-commit% consisting of the
// appName - micro service name
// version - the version number of the servce
// branch - a branch name of the repository
// commit - a commit id
func BuildVBC(appName, version, branch, commit string) string {
	branchB := strings.Builder{}
	if branch != "" || commit != "" {
		branchB.WriteString("/")
	}
	if branch != "" {
		branchB.WriteString(branch)
	}

	if commit != "" {
		branchB.WriteString("-")
		branchB.WriteString(commit)
	}

	return appName + "/" + version + branchB.String() + "%"
}
