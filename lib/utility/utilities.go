package utility

import (
	"strings"
)

// GitPathLast returns the last component of a git path.
func GitPathLast(repo string) string {
	if strings.Contains(repo, "/") {
		parts := strings.Split(repo, "/")
		return parts[len(parts)-1]
	}
	return ""
}
