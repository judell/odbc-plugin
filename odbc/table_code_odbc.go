package odbc

import (
	"context"

	"encoding/json"
	"os"

	"github.com/turbot/go-kit/helpers"
    "database/sql"
	//"fmt"
    _ "github.com/alexbrainman/odbc"	
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"	

)

func getSchema(ctx context.Context, dataSource string) ([]*plugin.Column, error) {
	plugin.Logger(ctx).Debug("odbc.getSchema")

	// Check if schema file exists
	if _, err := os.Stat("/tmp/schema_cache"); err == nil {
		// Read the schema from the file
		fileData, err := os.ReadFile("/tmp/schema_cache")
		if err != nil {
			return nil, err
		}
		var schemas map[string][]string
		err = json.Unmarshal(fileData, &schemas)
		if err != nil {
			return nil, err
		}

		// If the schema for the current data source is present in the cache
		if columnNames, ok := schemas[dataSource]; ok {
			// Reconstruct the columns from the cached data
			cols := make([]*plugin.Column, len(columnNames))
			for i, columnName := range columnNames {
				cols[i] = &plugin.Column{
					Name:        columnName,
					Type:        proto.ColumnType_STRING,
					Description: columnName,
					Transform:   transform.FromField(helpers.EscapePropertyName(columnName)),
				}
			}
			return cols, nil
		}
	}

	// If the schema for the current data source isn't present in the cache (or if the cache file doesn't exist), fetch the schema from the database
	db, err := sql.Open("odbc", "DSN=" + dataSource)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM rss WHERE 1=0") // This will return an empty result set, but with column headers
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
	schemas := make(map[string][]string)
	if _, err := os.Stat("/tmp/schema_cache"); err == nil {
		fileData, _ := os.ReadFile("/tmp/schema_cache")
		json.Unmarshal(fileData, &schemas)
	}
	schemas[dataSource] = columnNames

	// Serialize and save to /tmp/schema_cache
	jsonData, err := json.Marshal(schemas)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile("/tmp/schema_cache", jsonData, 0644)
	if err != nil {
		return nil, err
	}

	cols := make([]*plugin.Column, len(columnNames))
	for i, columnName := range columnNames {
		cols[i] = &plugin.Column{
			Name:        columnName,
			Type:        proto.ColumnType_STRING,
			Description: columnName,
			Transform:   transform.FromField(helpers.EscapePropertyName(columnName)),
		}
	}
	return cols, nil
}

func tableODBC(ctx context.Context, connection *plugin.Connection) (*plugin.Table, error) {
	plugin.Logger(ctx).Debug("tableODBC")

	config := GetConfig(connection)
	plugin.Logger(ctx).Debug("tableODBC",  "config", config)


	cols, err := getSchema(ctx, "CData RSS Source")
	if err != nil {
		return nil, err
	}

	plugin.Logger(ctx).Debug("odbc", "cols", cols)

	return &plugin.Table{
		Name:        "rss",
		Description: "RSS ODBC Table",
		List: &plugin.ListConfig{
			Hydrate: listODBC,
		},
		Columns: cols,
	}, nil
}

func listODBC(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listODBC start")

    db, err := sql.Open("odbc", "DSN=CData RSS Source")
    if err != nil {
        return nil, err
    }
    defer db.Close()

    // Fetch all columns for demonstration; ideally, you'd limit columns or add conditions as necessary
    rows, err := db.Query("SELECT * FROM rss")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Get column names
    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }

	plugin.Logger(ctx).Debug("listODBC", "columns", columns)
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
        for i, colName := range columns {
            val := colPtrs[i].(*interface{})
            //m[colName] = *val
			m[helpers.EscapePropertyName(colName)] = *val
        }
        d.StreamListItem(ctx, m)
    }

    return nil, nil
}


