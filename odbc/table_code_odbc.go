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

