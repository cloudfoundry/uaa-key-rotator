package db

import (
	"errors"
	"fmt"
	"strings"
)

func RebindForSQLDialect(query, dialect string) (string, error) {
	switch dialect {
	case "mysql":
		return query, nil
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
