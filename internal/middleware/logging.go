package middleware

import (
	"log"
	"time"
)

func logAccessDenied(reason string, details map[string]interface{}) {
	log.Printf("[ACCESS DENIED] %s | %v | %s\n", reason, details, time.Now().Format(time.RFC3339))
}
