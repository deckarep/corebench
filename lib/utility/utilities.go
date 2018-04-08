package utility

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

// GitPathLast returns the last component of a git path.
func GitPathLast(repo string) string {
	if strings.Contains(repo, "/") {
		parts := strings.Split(repo, "/")
		return parts[len(parts)-1]
	}
	return ""
}

// NewInstanceID returns the last component of a V4 Guid, to keep the names reasonably small.
func NewInstanceID() string {
	g := uuid.Must(uuid.NewV4())
	p := strings.Split(g.String(), "-")
	return p[len(p)-1]
}
