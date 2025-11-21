package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/lyonmu/quebec/cmd/gateway/internal/global"
	"github.com/lyonmu/quebec/pkg/constant"
	"github.com/prometheus/common/version"
)

func main() {
	kong.Parse(&global.Cfg,
		kong.Name(string(constant.ModuleNameGateway)),
		kong.Description(string(constant.ModuleNameGateway)),
		kong.UsageOnError(),
		kong.HelpOptions{
			Compact: true,
			Summary: true,
		},
	)
	if global.Cfg.Version {
		fmt.Println(version.Print(string(constant.ModuleNameGateway)))
		os.Exit(0)
	}

}
