package tools

import (
	"github.com/mssola/useragent"
)

func ParseUserAgent(ua string) *useragent.UserAgent {
	return useragent.New(ua)
}
