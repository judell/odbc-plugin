package odbc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"database/sql"
	_ "github.com/alexbrainman/odbc"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func getSchemas(ctx context.Context, dataSource string, tablename string) ([]*plugin.Column, error) {
	plugin.Logger(ctx).Debug("odbc.getSchema", "dataSource", dataSource, "tablename", tablename)

	db, err := sql.Open("odbc", "DSN="+dataSource)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Fetch a single row to probe types
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s LIMIT 1", tablename))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		// Handle the case where there are no rows
		return nil, fmt.Errorf("unable to retrieve column data for table %s", tablename)
	}

	values := make([]interface{}, len(columnNames))
	valuePtrs := make([]interface{}, len(columnNames))
	for i := 0; i < len(columnNames); i++ {
		valuePtrs[i] = &values[i]
	}

	// Scan values from the row into the slice
	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	cols := make([]*plugin.Column, len(columnNames))
	for i, colName := range columnNames {
		var columnType proto.ColumnType
		switch values[i].(type) {
		case int, int32, int64:
			columnType = proto.ColumnType_INT
		case float32, float64:
			columnType = proto.ColumnType_DOUBLE
		case string:
			columnType = proto.ColumnType_STRING
		case time.Time:
			columnType = proto.ColumnType_TIMESTAMP
		default:
			columnType = proto.ColumnType_STRING
		}
		cols[i] = &plugin.Column{
			Name:        colName,
			Type:        columnType,
			Description: dataSource + " " + colName,
			Transform:   transform.FromField(helpers.EscapePropertyName(colName)),
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

	// Add the _metadata column to the list of columns
	cols = append(cols, &plugin.Column{
		Name:        "_metadata",
		Type:        proto.ColumnType_JSON,
		Transform:   transform.FromField("_metadata"),
		Description: "Metadata related to the data source",
	})

	return &plugin.Table{
		Name:        strings.ToLower(dsn) + "_" + tablename,
		Description: dsn + ":" + tablename,
		List: &plugin.ListConfig{
			Hydrate:    listODBC,
			KeyColumns: getKeyColumns(cols),
		},
		Columns: cols,
	}, nil
}

func getKeyColumns(columns []*plugin.Column) plugin.KeyColumnSlice {
	var keyCols plugin.KeyColumnSlice

	for _, col := range columns {
		keyCols = append(keyCols, &plugin.KeyColumn{
			Name:    col.Name,
			Require: plugin.Optional, // This means the column can be used as a qualifier, but it's optional
		})
	}

	return keyCols
}

func listODBC(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listODBC start")

	// Split the table name to get dsn and tablename
	parts := strings.Split(d.Table.Name, "_")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid table name format")
	}
	dsn, tablename := parts[0], parts[1]

	plugin.Logger(ctx).Debug("listODBC", "dsn", dsn, "tablename", tablename, "quals", d.Quals)

	// Fetch data from the database using the qualifiers
	results, err := fetchFromDatabase(ctx, dsn, tablename, d.Quals)
	if err != nil {
		return nil, err
	}

	// Stream the results
	for _, result := range results {
		d.StreamListItem(ctx, result)
	}

	return nil, nil
}

func fetchFromDatabase(ctx context.Context, dsn string, tablename string, quals plugin.KeyColumnQualMap) ([]map[string]interface{}, error) {
	plugin.Logger(ctx).Debug("odbc: fetchFromDatabase", "dsn", dsn, "tablename", tablename, "quals", quals)
	db, err := sql.Open("odbc", "DSN="+dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT * FROM " + tablename

	// Construct the WHERE clause using protoToODBCValue
	conds := make([]string, 0)
	for _, qualsForCol := range quals {
		for _, qual := range qualsForCol.Quals {
			if qual.Value.Value == nil {
				conds = append(conds, fmt.Sprintf("%s %s", qual.Column, qual.Operator))
			} else {
				valueStr := protoToODBCValue(qual.Value)
				conds = append(conds, fmt.Sprintf("%s %s %s", qual.Column, qual.Operator, valueStr))
			}
		}
	}
	if len(conds) > 0 {
		query = query + " WHERE " + strings.Join(conds, " AND ")
	}

	plugin.Logger(ctx).Debug("listODBC", "adjust query", query)

	rows, err := db.Query(query)
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

		m["_metadata"] = map[string]interface{}{
			"connection_name": "odbc",
			"dsn":             dsn,
		}

		for i, colName := range columns {
			val := colPtrs[i].(*interface{})
			m[helpers.EscapePropertyName(colName)] = *val
		}

		results = append(results, m)
	}

	return results, nil
}
