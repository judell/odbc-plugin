package odbc

import (
	"context"

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
	tables := map[string]*plugin.Table{}

	config := GetConfig(d.Connection)
	for _, dataSource := range config.DataSources {
		dsn, tableName := splitDataSourceAndTable(dataSource)
		tableCtx := context.WithValue(ctx, "dsn", dsn)
		tableCtx = context.WithValue(tableCtx, "tableName", tableName)
		table, err := tableODBC(tableCtx, d.Connection)
		if err != nil {
			return nil, err
		}
		tables[tableName] = table
	}

	return tables, nil
}
