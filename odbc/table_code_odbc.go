package odbc

import (
	"context"

	"encoding/json"
	"fmt"
	"os"

	"database/sql"
	"github.com/turbot/go-kit/helpers"
	//"fmt"
	_ "github.com/alexbrainman/odbc"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func getSchema(ctx context.Context, dataSource string, tableName string) ([]*plugin.Column, error) {
	plugin.Logger(ctx).Debug("odbc.getSchema")

	// Check if schema file exists
	if _, err := os.Stat("/tmp/schema_cache"); err == nil {
		// Read the schema from the file
		fileData, err := os.ReadFile("/tmp/schema_cache")
		if err != nil {
			return nil, err
		}
		var schemas map[string]map[string][]string
		err = json.Unmarshal(fileData, &schemas)
		if err != nil {
			return nil, err
		}

		// If the schema for the current data source and table is present in the cache
		if tables, ok := schemas[dataSource]; ok {
			if columnNames, ok := tables[tableName]; ok {
				// Reconstruct the columns from the cached data
				cols := make([]*plugin.Column, len(columnNames))
				for i, columnName := range columnNames {
					cols[i] = &plugin.Column{
						Name:        columnName,
						Type:        proto.ColumnType_STRING,
						Description: dataSource + " " + columnName,
						Transform:   transform.FromField(helpers.EscapePropertyName(columnName)),
					}
				}
				return cols, nil
			}
		}
	}

	// If the schema for the current data source isn't present in the cache (or if the cache file doesn't exist), fetch the schema from the database
	db, err := sql.Open("odbc", "DSN="+dataSource)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s WHERE 1=0", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	plugin.Logger(ctx).Debug("odbc.getSchema", "columnNames", columnNames)

	// Create a new map (or use the previously read map) and add the fetched schema for the current data source
	var schemas map[string]map[string][]string
	if _, err := os.Stat("/tmp/schema_cache"); err == nil {
		fileData, _ := os.ReadFile("/tmp/schema_cache")
		json.Unmarshal(fileData, &schemas)
	}

	if schemas == nil {
		schemas = make(map[string]map[string][]string)
	}
	if schemas[dataSource] == nil {
		schemas[dataSource] = make(map[string][]string)
	}

	columnNamesWithDSN := append([]string{"dsn", "tableName"}, columnNames...)
	schemas[dataSource][tableName] = columnNamesWithDSN

	// Serialize and save to /tmp/schema_cache
	jsonData, err := json.Marshal(schemas)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile("/tmp/schema_cache", jsonData, 0644)
	if err != nil {
		return nil, err
	}

	// Reconstruct the columns from the cached data
	cols := make([]*plugin.Column, len(columnNames))
	for i, columnName := range columnNames {
		if columnName == "dsn" {
			cols[i] = &plugin.Column{
				Name:        "dsn",
				Type:        proto.ColumnType_STRING,
				Description: "Data Source Name for the ODBC connection",
				Transform:   transform.FromQual("dsn"),
			}
		} else if columnName == "tableName" {
			cols[i] = &plugin.Column{
				Name:        "tableName",
				Type:        proto.ColumnType_STRING,
				Description: "Table name",
				Transform:   transform.FromQual("tableName"),
			}
		} else {
			cols[i] = &plugin.Column{
				Name:        columnName,
				Type:        proto.ColumnType_STRING,
				Description: dataSource + " " + columnName,
				Transform:   transform.FromField(helpers.EscapePropertyName(columnName)),
			}
		}
	}

	plugin.Logger(ctx).Debug("odbc.getSchema", "cols", cols)
	return cols, nil
}

func tableODBC(ctx context.Context, connection *plugin.Connection) (*plugin.Table, error) {
	dsn := ctx.Value("dsn").(string)
	tableName := ctx.Value("tableName").(string)

	cols, err := getSchema(ctx, dsn, tableName)

	plugin.Logger(ctx).Debug("tableODBC", "cols", cols)

	if err != nil {
		return nil, err
	}

	dsnCol := &plugin.KeyColumn{Name: "dsn"}
	tableCol := &plugin.KeyColumn{Name: "tableName"}

	return &plugin.Table{
		Name:        tableName,
		Description: dsn,
		List: &plugin.ListConfig{
			Hydrate: listODBC,
			KeyColumns: []*plugin.KeyColumn{
				dsnCol,
				tableCol,
			},
		},
		Columns: cols,
	}, nil
}

func listODBC(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listODBC start")

	dsn := d.EqualsQualString("dsn")
	tableName := d.EqualsQualString("tableName")

	plugin.Logger(ctx).Debug("listODBC", "dsn", dsn, "tableName", tableName)

	// Fetch data from the database
	results, err := fetchFromDatabase(ctx, dsn, tableName)
	if err != nil {
		return nil, err
	}

	// Stream the results
	for _, result := range results {
		d.StreamListItem(ctx, result)
	}

	return nil, nil
}

func fetchFromDatabase(ctx context.Context, dsn string, tableName string) ([]map[string]interface{}, error) {
	plugin.Logger(ctx).Debug("odbc: fetchFromDatabase", "dsn", dsn, "tableName", tableName)
	db, err := sql.Open("odbc", "DSN="+dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Fetch all columns for demonstration
	rows, err := db.Query("SELECT * FROM " + tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		cols := make([]interface{}, len(columns))
		colPtrs := make([]interface{}, len(columns))
		for i := 0; i < len(columns); i++ {
			colPtrs[i] = &cols[i]
		}

		if err := rows.Scan(colPtrs...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		m["dsn"] = dsn
		m["tableName"] = tableName
		for i, colName := range columns {
			val := colPtrs[i].(*interface{})
			m[helpers.EscapePropertyName(colName)] = *val
		}

		results = append(results, m)
	}

	return results, nil
}
