package instanceid

import (
	"strings"
)

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
