package main

import (
	"github.com/turbot/steampipe-plugin-odbc/odbc"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: odbc.Plugin})
}
