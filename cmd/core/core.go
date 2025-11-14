package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/common"
	"github.com/prometheus/common/version"
)

func main() {
	kong.Parse(&global.Cfg,
		kong.Name(string(common.ModuleNameCore)),
		kong.Description(string(common.ModuleNameCore)),
		kong.UsageOnError(),
		kong.HelpOptions{
			Compact: true,
			Summary: true,
		},
	)
	if global.Cfg.Version {
		fmt.Println(version.Print(string(common.ModuleNameCore)))
		os.Exit(0)
	}

}
