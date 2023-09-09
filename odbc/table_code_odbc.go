package odbc

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableODBC(ctx context.Context, connection *plugin.Connection) (*plugin.Table, error) {
	// For demonstration purposes, hardcoding column names.
	// In a real-world scenario, this would involve introspecting the ODBC source to get the column names/types.
	cols := []*plugin.Column{
		{Name: "example_column1", Type: proto.ColumnType_STRING, Description: "Example column 1."},
		{Name: "example_column2", Type: proto.ColumnType_STRING, Description: "Example column 2."},
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
		"example_column1": "Value 1",
		"example_column2": "Value 2",
	}
	d.StreamListItem(ctx, row)

	return nil, nil
}
