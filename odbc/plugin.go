package odbc

import (
	"context"
	"os"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
    os.RemoveAll("/tmp/schema_cache")	
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
		dsn, tablename := splitDataSourceAndTable(dataSource)
		tableCtx := context.WithValue(ctx, "dsn", dsn)
		tableCtx = context.WithValue(tableCtx, "tablename", tablename)
		table, err := tableODBC(tableCtx, d.Connection)
		if err != nil {
			return nil, err
		}
		tables[tablename] = table
	}

	return tables, nil
}
