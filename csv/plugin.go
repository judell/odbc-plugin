package csv

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
	plugin.Logger(ctx).Debug("odbc.PluginODBCTables starting")
	tables := map[string]*plugin.Table{}

	// You might want to list ODBC tables here. For simplicity, I'm assuming a single table named "rss"
	// You can make this dynamic based on configuration or introspection of the ODBC source.
	tableName := "rss"
	plugin.Logger(ctx).Debug("odbc.PluginODBCTables calling tableODBC")
	table, err := tableODBC(ctx, d.Connection)
	if err != nil {
		plugin.Logger(ctx).Debug("odbc.PluginODBCTables", "create_table_error", err, "table", tableName)
		return nil, err
	}
	tables[tableName] = table

	plugin.Logger(ctx).Debug("odbc.PluginODBCTables", "tables", tables)
	return tables, nil
}
