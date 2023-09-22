package odbc

import (
	"strings"
)

func splitDataSourceAndTable(s string) (string, string) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}
