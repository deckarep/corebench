package utility

import (
	"fmt"
	"log"
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

// PromptConfirmation asks the user for confirmation.
func PromptConfirmation(msg string) bool {
	fmt.Println(msg)

	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return PromptConfirmation(msg)
	}
}

func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}
