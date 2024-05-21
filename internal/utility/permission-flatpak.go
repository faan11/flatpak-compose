package utility

import (
	"fmt"
	"strings"
)

// MapPermissionsToFlatpakOverrideFlags maps permissions to Flatpak override flags
func MapPermissionsToFlatpakOverrideFlags(permissionContext string) []string {
	return ParseFlatpakPermissions(permissionContext)
}

// ParseFlatpakPermissions parses the given permissions and returns Flatpak override flags
func ParseFlatpakPermissions(permissionContext string) []string {
	lines := strings.Split(permissionContext, "\n")
	flags := []string{}
	parsingSection := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "[") {
			parsingSection = line
			continue
		}

		switch parsingSection {
		case "[Context]":
			// Logic for processing [Context] section key-value pairs
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := parts[0]
			values := strings.Split(parts[1], ";")
			for _, value := range values {
				flag := ""
				if value == "" {
					continue
				}
				if strings.HasPrefix(value, "!") {
					flag = getNegativeContextFlag(key, strings.TrimPrefix(value, "!"))
				} else {
					flag = getContextFlag(key, value)
				}
				flags = append(flags, flag)
			}

		case "[Session Bus Policy]", "[System Bus Policy]":
			// Logic for processing [Session Bus Policy] and [System Bus Policy] sections
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				flag := fmt.Sprintf("--%s-name=%s", getContextPolicyFlag(parts[1]), parts[0])
				flags = append(flags, flag)
			}
		}
	}

	return flags
}

// Helper function to get context flags for [Context] section
func getContextFlag(key, value string) string {
	switch key {
	case "shared":
		return fmt.Sprintf("--share=%s", value)
	case "sockets":
		return fmt.Sprintf("--socket=%s", value)
	case "devices":
		return fmt.Sprintf("--device=%s", value)
	case "features":
		return fmt.Sprintf("--allow=%s", value)
	case "filesystems":
		return fmt.Sprintf("--filesystem=%s", value)
	case "persistent":
		return fmt.Sprintf("--persist=%s", value)
	default:
		return ""
	}
}

// Helper function to get negative context flags for [Context] section
func getNegativeContextFlag(key, value string) string {
	switch key {
	case "shared":
		return fmt.Sprintf("--unshare=%s", value)
	case "sockets":
		return fmt.Sprintf("--nosocket=%s", value)
	case "devices":
		return fmt.Sprintf("--nodevice=%s", value)
	case "features":
		return fmt.Sprintf("--disallow=%s", value)
	case "filesystems":
		return fmt.Sprintf("--nofilesystem=%s", value)
	default:
		return ""
	}
}

// Helper function to get context policy flag for [Session Bus Policy] and [System Bus Policy] sections
func getContextPolicyFlag(value string) string {
	switch value {
	case "own":
		return "own"
	case "talk":
		return "talk"
	default:
		return fmt.Sprintf("no-talk-%s", value)
	}
}
