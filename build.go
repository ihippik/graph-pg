package graphpg

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-pg/pg/orm"
)

type (
	Schema struct {
		tableName  struct{} `sql:"information_schema.columns"`
		ColumnName string   `json:"column_name"`
		ColumnType string   `json:"udt_name" sql:"udt_name"`
	}
)

const (
	int4Type      = "int4"
	varcharType   = "varchar"
	timestampType = "timestamp"
)

var rangeRegexp = regexp.MustCompile("\\d+;;\\d+")

// BuildQuery build a query from user query parameters.
func BuildQuery(db orm.DB, query *orm.Query, tableName, q string) (*orm.Query, error) {
	var (
		columns   []Schema
		condition string
	)

	if err := db.Model(&columns).Where("table_name=?", tableName).Select(); err != nil {
		return nil, fmt.Errorf("get table columns: %w", err)
	}

	fieldsType := make(map[string]string, len(columns))
	for _, column := range columns {
		fieldsType[column.ColumnName] = column.ColumnType
	}

	queryItem := strings.Split(q, ",")
	for _, item := range queryItem {
		queryVal := strings.Split(item, ":")
		if len(queryVal) != 2 {
			return query, nil
		}

		columns, val := queryVal[0], queryVal[1]
		columnName := strings.Split(columns, "|")

		valueType, ok := fieldsType[columnName[0]]
		if !ok {
			return query, errUnknownColumn
		}

		for i := range columnName {
			columnName[i] = fmt.Sprintf(`"%s".%s`, tableName, columnName[i])
		}

		switch valueType {
		case int4Type:
			condition = conditionForInt(val, columnName)
		case varcharType:
			condition = conditionForVarchar(val, columnName)
		case timestampType:
			condition = conditionForTimestamp(val, columnName)
		default:
			condition = fmt.Sprintf("%s = '%s'", columnName[0], val)
		}

		query.Where(condition)
	}

	return query, nil
}

func conditionForTimestamp(val string, columnName []string) string {
	var condition string

	switch {
	case strings.HasPrefix(val, "||"):
		condition = fmt.Sprintf("extract(epoch from %s) < %s", columnName[0], val[2:])
	case strings.HasSuffix(val, "||"):
		condition = fmt.Sprintf("extract(epoch from %s) > %s", columnName[0], val[:len(val)-2])
	case rangeRegexp.Match([]byte(val)):
		values := strings.Split(val, ";;")
		condition = fmt.Sprintf(
			"extract(epoch from %s) > %s AND extract(epoch from %s) <= %s",
			columnName[0],
			values[0],
			columnName[0],
			values[1],
		)
	default:
		condition = fmt.Sprintf("%s = %s", columnName[0], val)
	}

	return condition
}

func conditionForVarchar(val string, columnName []string) string {
	var condition string

	switch {
	case !strings.Contains(val, "||") && strings.Contains(val, "|"):
		condition = createEnumeration(columnName[0], val)
	case strings.HasPrefix(val, "~"):
		condition = fmt.Sprintf("%s ILIKE '%%%s%%'", columnName[0], val[1:])
	default:
		condition = fmt.Sprintf("%s = '%s'", columnName[0], val)
	}

	return condition
}

func conditionForInt(val string, columnName []string) string {
	var condition string

	switch {
	// TODO: not processed more than two
	case len(columnName) > 1:
		condition = fmt.Sprintf("(%s = %s OR %s = %s)", columnName[0], val, columnName[1], val)
	case strings.HasPrefix(val, "||"):
		condition = fmt.Sprintf("%s <= %s", columnName[0], val[2:])
	case strings.HasSuffix(val, "||"):
		condition = fmt.Sprintf("%s >= %s", columnName[0], val[:len(val)-2])
	case rangeRegexp.Match([]byte(val)):
		values := strings.Split(val, ";;")
		condition = fmt.Sprintf(
			"%s > %s AND %s <= %s",
			columnName[0],
			values[0],
			columnName[0],
			values[1],
		)
	case len(strings.Split(val, "~")) > 1 && !strings.Contains(val, "~~"):
		condition = fmt.Sprintf(
			"%s IN(%s)",
			columnName[0],
			strings.Replace(val, "~", ",", -1),
		)
	default:
		condition = fmt.Sprintf("%s = %s", columnName[0], val)
	}

	return condition
}
