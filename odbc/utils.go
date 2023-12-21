package odbc

import (
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
)

func splitDataSourceAndTable(s string) (string, string) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func protoToODBCValue(val *proto.QualValue) string {
	switch val.GetValue().(type) {
	case *proto.QualValue_BoolValue:
		if val.GetBoolValue() {
			return "TRUE"
		}
		return "FALSE"
	case *proto.QualValue_DoubleValue:
		return fmt.Sprintf("%f", val.GetDoubleValue())
	case *proto.QualValue_Int64Value:
		return fmt.Sprintf("%d", val.GetInt64Value())
	case *proto.QualValue_StringValue:
		return fmt.Sprintf("'%s'", val.GetStringValue())
	case *proto.QualValue_TimestampValue:
		t := val.GetTimestampValue().AsTime()
		return fmt.Sprintf("'%s'", t.Format("2006-01-02 15:04:05"))
	// Add cases for other data types as needed.
	// Note: For types like CIDR or Timestamp which may not be universally supported
	// across all databases behind ODBC, you might need more nuanced handling.
	default:
		return "<INVALID>" // this will probably cause an error on the query, which might be acceptable
	}
}
