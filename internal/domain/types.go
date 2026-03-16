package domain

import (
	"time"
)

// generateID genera un ID único simple.
func generateID() string {
	return time.Now().Format("20060102150405") + "-stub"
}
