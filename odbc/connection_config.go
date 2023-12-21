package odbc

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type odbcConfig struct {
	DataSources []string `cty:"data_sources"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"data_sources": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
	},
}

func ConfigInstance() interface{} {
	return &odbcConfig{}
}

func GetConfig(connection *plugin.Connection) odbcConfig {
	if connection == nil || connection.Config == nil {
		return odbcConfig{}
	}
	config, _ := connection.Config.(odbcConfig)
	return config
}
