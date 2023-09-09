package odbc

import (
	"context"

	"github.com/turbot/go-kit/helpers"
    "database/sql"
	//"fmt"
    _ "github.com/alexbrainman/odbc"	
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"	

)

/*
func tableODBC(ctx context.Context, connection *plugin.Connection) (*plugin.Table, error) {
	plugin.Logger(ctx).Debug("tableODBC")
    // Hardcoding the columns "Title" and "Link"
    cols := []*plugin.Column{
        {Name: "Title", Type: proto.ColumnType_STRING, Description: "The title of the RSS item.", Transform: transform.FromField("Title")},
        {Name: "Link", Type: proto.ColumnType_STRING, Description: "The link of the RSS item.", Transform: transform.FromField("Link")},
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
*/

func tableODBC(ctx context.Context, connection *plugin.Connection) (*plugin.Table, error) {
	plugin.Logger(ctx).Debug("tableODBC")

	// Connect to the ODBC data source
	
	/*
	plugin.Logger(ctx).Debug("tableODBC open db")
    db, err := sql.Open("odbc", "DSN=CData RSS Source")
    if err != nil {
		plugin.Logger(ctx).Debug("tableODBC", "err", err)
        return nil, err
    }
	plugin.Logger(ctx).Debug("tableODBC defer close db")
    defer db.Close()

    // Fetch column names from the ODBC source
    rows, err := db.Query("SELECT * FROM rss WHERE 1=0")  // This will return an empty result set, but with column headers
    if err != nil {
		plugin.Logger(ctx).Debug("tableODBC", "err", err)
        return nil, err
    }

	plugin.Logger(ctx).Debug("tableODBC", "rows", rows)

    columns, err := rows.Columns()
    if err != nil {
		plugin.Logger(ctx).Debug("tableODBC", "err", err)
        return nil, err
    }
    rows.Close() // Close the rows immediately as we just want column names
	plugin.Logger(ctx).Debug("odbc", "columns", columns)


    // Construct the table columns based on the fetched column names
    cols := make([]*plugin.Column, len(columns))
    for i, columnName := range columns {
        cols[i] = &plugin.Column{
            Name:        columnName,
            Type:        proto.ColumnType_STRING, // Assuming string type for simplicity; a more robust approach would inspect column types
            Description: columnName,
            Transform:   transform.FromField(helpers.EscapePropertyName(columnName)),
        }
    }
	*/

	cols := []*plugin.Column{
        {Name: "Title", Type: proto.ColumnType_STRING, Description: "The title of the RSS item.", Transform: transform.FromField("Title")},
        {Name: "Link", Type: proto.ColumnType_STRING, Description: "The link of the RSS item.", Transform: transform.FromField("Link")},
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
    // Connect to the ODBC data source
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


