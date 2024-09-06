package main

import (
	"github.com/aswinkevin/steampipe-plugin-quark/quark"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: quark.Plugin})
}
