package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/lyonmu/quebec/cmd/gateway/internal/global"
	"github.com/lyonmu/quebec/pkg/common"
	"github.com/prometheus/common/version"
)

func main() {
	kong.Parse(&global.Cfg,
		kong.Name(string(common.ModuleNameGateway)),
		kong.Description(string(common.ModuleNameGateway)),
		kong.UsageOnError(),
		kong.HelpOptions{
			Compact: true,
			Summary: true,
		},
	)
	if global.Cfg.Version {
		fmt.Println(version.Print(string(common.ModuleNameGateway)))
		os.Exit(0)
	}

}
