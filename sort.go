package graphpg

import (
	"strings"

	"github.com/go-pg/pg/orm"
)

// BuildPagination adds to the query a limit offset and the total number of records.
func BuildPagination(query *orm.Query, sort string, limit, offset int) (*orm.Query, int, error) {
	var (
		tmp    string
		result []string
		count  int
	)

	sortItem := strings.Split(sort, ",")

	for _, item := range sortItem {
		if strings.HasPrefix(item, "-") {
			columnName := item
			tmp = columnName[1:] + " DESC"
		} else {
			tmp = item
		}

		result = append(result, tmp)
	}

	for _, srt := range result {
		if len(srt) == 0 {
			continue
		}

		query.Order(srt)
	}

	count, err := query.Count()
	if err != nil {
		return query, count, err
	}

	if limit > 0 {
		query.Limit(limit)
	}

	if offset > 0 {
		query.Offset(offset)
	}

	return query, count, nil
}
