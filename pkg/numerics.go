package pkg

import (
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

func Float64ToNumeric(f float64) (pgtype.Numeric, error) {
	s := strconv.FormatFloat(f, 'f', -1, 64)
	var num pgtype.Numeric

	err := num.Scan(s)
	if err != nil {
		return pgtype.Numeric{}, err
	}

	return num, nil
}
