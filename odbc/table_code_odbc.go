package odbc

import (
	"context"

	//"github.com/turbot/go-kit/helpers"
    "database/sql"
    _ "github.com/alexbrainman/odbc"	
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"	

)

func tableODBC(ctx context.Context, connection *plugin.Connection) (*plugin.Table, error) {
    // Hardcoding the columns "Title" and "Link"
    cols := []*plugin.Column{
        {Name: "Title", Type: proto.ColumnType_STRING, Description: "The title of the RSS item.", Transform: transform.FromField("Title")},
        {Name: "Link", Type: proto.ColumnType_STRING, Description: "The link of the RSS item.", Transform: transform.FromField("Link")},
    }

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

    // Query for the desired columns "Title" and "Link"
    rows, err := db.Query("SELECT Title, Link FROM rss")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Iterate over the results and stream them
    for rows.Next() {
        var title, link string
        if err := rows.Scan(&title, &link); err != nil {
            return nil, err
        }
        d.StreamListItem(ctx, map[string]string{
            "Title": title,
            "Link":  link,
        })
    }

    return nil, nil
}


/*
func tableODBC(ctx context.Context, connection *plugin.Connection) (*plugin.Table, error) {
	// For demonstration purposes, hardcoding column names.
	// In a real-world scenario, this would involve introspecting the ODBC source to get the column names/types.
	cols := []*plugin.Column{
		{Name: "example_column1", Type: proto.ColumnType_STRING, Transform: transform.FromField("example_column1"), Description: "Example column 1."},
		{Name: "example_column2", Type: proto.ColumnType_STRING, Transform: transform.FromField("example_column2"), Description: "Example column 2."},
	}
	

	return &plugin.Table{
		Name:        "odbc_table",
		Description: "ODBC Table",
		List: &plugin.ListConfig{
			Hydrate: dummyList, // A dummy list function for testing purposes
		},
		Columns: cols,
	}, nil
}

func dummyList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Return a dummy row for testing
	row := map[string]string{
		"example_column1": "one",
		"example_column2": "two",
	}
	plugin.Logger(ctx).Error("odbc", "dummyList row", row)
	d.StreamListItem(ctx, row)

	return nil, nil
}
*/