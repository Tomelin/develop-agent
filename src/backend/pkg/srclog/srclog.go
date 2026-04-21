package srclog

import (
	"fmt"
	"runtime"
	"strings"
)

// GetComponent returns a formatted string representing the component
// that called the function where GetComponent is invoked.
//
// Rules:
// - Matches /pkg/ -> "third party"
// - Matches /internal/infrastructure/ -> "infra - <pkg_name>"
// - Matches /internal/ -> "business - <pkg_name>"
// - Otherwise -> "unknown"
func GetComponent() string {
	// Call depth is 1: we want to know who called GetComponent.
	// If a wrapper function is used, this depth would need to be adjusted.
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}

	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown"
	}

	return ParseCallerName(f.Name())
}

// ParseCallerName parses a fully qualified function name and applies the component rules.
// Exposed for testing purposes.
func ParseCallerName(fullName string) string {
	lastSlash := strings.LastIndex(fullName, "/")
	if lastSlash == -1 {
		return "unknown"
	}

	afterSlash := fullName[lastSlash+1:]
	firstDot := strings.Index(afterSlash, ".")

	var pkgName string
	if firstDot != -1 {
		pkgName = afterSlash[:firstDot]
	} else {
		pkgName = afterSlash
	}

	// We check paths in order of specificity
	if strings.Contains(fullName, "/pkg/") {
		return "third party"
	}

	if strings.Contains(fullName, "/internal/infrastructure/") {
		return fmt.Sprintf("infra - %s", pkgName)
	}

	if strings.Contains(fullName, "/internal/") {
		return fmt.Sprintf("business - %s", pkgName)
	}

	return "unknown"
}
