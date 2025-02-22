package config

type CmdConfig struct {
	Config string `arg:"" name:"config" default:"config.yml" help:"config file path" type:"path"`
}
