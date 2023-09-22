package odbc

import (
	"context"

	//"encoding/json"
	"fmt"
	//"os"
	"strings"

	"database/sql"

	"github.com/turbot/go-kit/helpers"

	//"fmt"
	_ "github.com/alexbrainman/odbc"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func getSchemas(ctx context.Context, dataSource string, tablename string) ([]*plugin.Column, error) {
    plugin.Logger(ctx).Debug("odbc.getSchema")

    db, err := sql.Open("odbc", "DSN="+dataSource)
    if err != nil {
        return nil, err
    }
    defer db.Close()

    rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s WHERE 1=0", tablename))
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    columnNames, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    plugin.Logger(ctx).Debug("odbc.getSchema", "columnNames", columnNames)

    cols := make([]*plugin.Column, len(columnNames))
    for i, columnName := range columnNames {
        cols[i] = &plugin.Column{
            Name:        columnName,
            Type:        proto.ColumnType_STRING,
            Description: dataSource + " " + columnName,
            Transform:   transform.FromField(helpers.EscapePropertyName(columnName)),
        }
    }

    plugin.Logger(ctx).Debug("odbc.getSchema", "cols", cols)
    return cols, nil
}


func tableODBC(ctx context.Context, connection *plugin.Connection) (*plugin.Table, error) {
    dsn := ctx.Value("dsn").(string)
    tablename := ctx.Value("tablename").(string)

    cols, err := getSchemas(ctx, dsn, tablename)
    if err != nil {
        return nil, err
    }

    return &plugin.Table{
        Name:        strings.ToLower(dsn) + "_" + tablename,
        Description: dsn,
        List: &plugin.ListConfig{
            Hydrate: listODBC,
        },
        Columns: cols,
    }, nil
}

func listODBC(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
    plugin.Logger(ctx).Debug("listODBC start")

    // Split the table name to get dsn and tablename
    parts := strings.Split(d.Table.Name, "_")
    if len(parts) != 2 {
        return nil, fmt.Errorf("invalid table name format")
    }
    dsn, tablename := parts[0], parts[1]

    plugin.Logger(ctx).Debug("listODBC", "dsn", dsn, "tablename", tablename)

    // Fetch data from the database
    results, err := fetchFromDatabase(ctx, dsn, tablename)
    if err != nil {
        return nil, err
    }

    // Stream the results
    for _, result := range results {
        d.StreamListItem(ctx, result)
    }

    return nil, nil
}


func fetchFromDatabase(ctx context.Context, dsn string, tablename string) ([]map[string]interface{}, error) {
	plugin.Logger(ctx).Debug("odbc: fetchFromDatabase", "dsn", dsn, "tablename", tablename)
	db, err := sql.Open("odbc", "DSN="+dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Fetch all columns for demonstration
	rows, err := db.Query("SELECT * FROM " + tablename)
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
		m["tablename"] = tablename
		for i, colName := range columns {
			val := colPtrs[i].(*interface{})
			m[helpers.EscapePropertyName(colName)] = *val
		}

		results = append(results, m)
	}

	return results, nil
}
