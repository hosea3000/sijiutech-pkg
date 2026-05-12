package helper

import (
	"database/sql"
	"strings"
)

func ToSqlNullString(s string) sql.NullString {
	s = strings.TrimSpace(s)
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func ToSqlNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}
