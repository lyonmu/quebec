package tools

import (
	"testing"

	"log"
)

func TestUserAgent(t *testing.T) {
	raw := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0"
	ua := ParseUserAgent(raw)
	if ua == nil {
		log.Fatal("Expected non-nil UserAgent, got nil")
	}
	log.Printf("ua : %+v", ua)
}
