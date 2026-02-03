package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/lmagdanello/tasks/internal/store"
)

func parsePriority(s string) store.Priority {
	switch strings.ToLower(s) {
	case "alta", "high", "h":
		return store.PriorityHigh
	case "baixa", "low", "l":
		return store.PriorityLow
	default:
		return store.PriorityMedium
	}
}

func parseDue(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "none" || s == "-" {
		return time.Time{}, fmt.Errorf("empty date")
	}
	return time.Parse("2006-01-02", s)
}
