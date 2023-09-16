package odbc

import (
	"strings"

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

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) odbcConfig {
	if connection == nil || connection.Config == nil {
		return odbcConfig{}
	}
	config, _ := connection.Config.(odbcConfig)
	return config
}

// Additional helper function to split data source and table from a string
func splitDataSourceAndTable(s string) (string, string) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		// Handle this case, maybe return an error or default values
		return "", ""
	}
	return parts[0], parts[1]
}
