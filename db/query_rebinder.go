package db

import (
	"strings"
	"fmt"
	"errors"
)

func RebindForSQLDialect(query, dialect string) (string, error) {
	switch dialect {
	case "mysql":
		return query, nil
	case "sqlserver":
		strParts := strings.Split(query, "?")
		for i := 1; i < len(strParts); i++ {
			strParts[i-1] = fmt.Sprintf("%s@p%d", strParts[i-1], i)
		}
		return strings.Join(strParts, ""), nil
	case "postgres":
		strParts := strings.Split(query, "?")
		for i := 1; i < len(strParts); i++ {
			strParts[i-1] = fmt.Sprintf("%s$%d", strParts[i-1], i)
		}
		return strings.Join(strParts, ""), nil
	default:
		return "", errors.New(fmt.Sprintf("Unrecognized DB dialect '%s'", dialect))
	}
}
