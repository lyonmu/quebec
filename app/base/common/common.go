package common

import baseconfig "github.com/lyonmu/quebec/app/base/config"

var (
	System baseconfig.SystemConfig

	SkipApiPath = []string{
		"/quebec/base/metrics",
		"/quebec/base/swagger",
		"/quebec/base/swagger/index.html",
		"/quebec/base/swagger/swagger-ui.css",
		"/quebec/base/swagger/swagger-ui-bundle.js",
		"/quebec/base/swagger/swagger-ui-standalone-preset.js",
		"/quebec/base/swagger/doc.json",
		"/quebec/base/swagger/favicon-32x32.png",
	}
)
