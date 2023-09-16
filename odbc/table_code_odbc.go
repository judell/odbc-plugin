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

	columnNamesWithDSN := append([]string{"dsn"}, columnNames...)
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

	return &plugin.Table{
		Name:        tableName,
		Description: dsn,
		List: &plugin.ListConfig{
			Hydrate: listODBC,
			KeyColumns: plugin.SingleColumn("dsn"),
		},
		Columns: cols,
	}, nil
}

func listODBC(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listODBC start")

	dsn := d.EqualsQualString("dsn")
	plugin.Logger(ctx).Debug("listODBC", "dsn", dsn)

	db, err := sql.Open("odbc", "DSN="+dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Fetch all columns for demonstration; ideally, you'd limit columns or add conditions as necessary
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Iterate over the results and stream them
	for rows.Next() {
		// Create a slice of interface{}'s to represent each column, and a value to which the column's value will be scanned
		cols := make([]interface{}, len(columns))
		colPtrs := make([]interface{}, len(columns))
		for i := 0; i < len(columns); i++ {
			colPtrs[i] = &cols[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(colPtrs...); err != nil {
			return nil, err
		}

		// Create our map, and retrieve the value for each column from the pointers slice, then add it to our map
		m := make(map[string]interface{})
		m["dsn"] = dsn
		for i, colName := range columns {
			val := colPtrs[i].(*interface{})
			//m[colName] = *val
			m[helpers.EscapePropertyName(colName)] = *val
		}
		d.StreamListItem(ctx, m)
	}

	return nil, nil
}
