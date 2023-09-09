package csv

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type odbcConfig struct {
	// Placeholder for future ODBC related configurations
}

var ConfigSchema = map[string]*schema.Attribute{
	// Placeholder for future ODBC related configurations
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
