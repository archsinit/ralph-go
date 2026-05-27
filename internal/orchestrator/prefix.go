package orchestrator

import (
	"strings"
)

// ParsePrefix parses an optional leading 'name:' prefix from input.
// If input starts with "<agentname>:" where agentname is one of the configured agents (case-insensitive),
// returns that target and the trimmed remainder.
// Otherwise returns empty target and input unchanged.
func ParsePrefix(input string, agents []string) (target, body string) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", input
	}

	// Find the first colon
	colonIdx := strings.Index(trimmed, ":")
	if colonIdx <= 0 {
		return "", input
	}

	// Extract the potential agent name
	potentialAgent := trimmed[:colonIdx]

	// Check if it matches any agent (case-insensitive)
	for _, agent := range agents {
		if strings.EqualFold(potentialAgent, agent) {
			// Found a match, return the target and the trimmed body
			body := strings.TrimSpace(trimmed[colonIdx+1:])
			return agent, body
		}
	}

	// No match found, return empty target and unchanged input
	return "", input
}
