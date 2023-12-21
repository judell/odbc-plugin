package odbc

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-odbc",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromGo().NullIfZero(),
		SchemaMode:       plugin.SchemaModeDynamic,
		TableMapFunc:     PluginODBCTables,
	}
	return p
}

func PluginODBCTables(ctx context.Context, d *plugin.TableMapData) (map[string]*plugin.Table, error) {
	plugin.Logger(ctx).Debug("PluginODBCTables")
	tables := map[string]*plugin.Table{}

	config := GetConfig(d.Connection)
	for _, dataSource := range config.DataSources {
		plugin.Logger(ctx).Debug("PluginODBCTables", "dataSource", dataSource)
		dsn, tablename := splitDataSourceAndTable(dataSource)
		// Combine the DSN and tablename to form the new table name
		newTableName := strings.ToLower(dsn) + "_" + tablename
		tableCtx := context.WithValue(ctx, "dsn", dsn)
		tableCtx = context.WithValue(tableCtx, "tablename", tablename)
		table, err := tableODBC(tableCtx, d.Connection)
		if err != nil {
			return nil, err
		}
		tables[newTableName] = table
		plugin.Logger(ctx).Debug("PluginODBCTables", "adding", newTableName)
	}

	return tables, nil
}
