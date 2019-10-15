package graphpg

import (
	"errors"
	"fmt"
	"github.com/go-pg/pg"
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

var (
	rangeRegexp = regexp.MustCompile("\\d+;;\\d+")
)

// BuildQuery build a query from user query parameters.
func BuildQuery(db *pg.DB, query *orm.Query, tableName, q string) (*orm.Query, error) {
	var (
		columns   []Schema
		condition string
	)

	if err := db.Model(&columns).Where("table_name=?", tableName).Select(); err != nil {
		return nil, err
	}

	fieldsType := make(map[string]string)
	for _, column := range columns {
		fieldsType[column.ColumnName] = column.ColumnType
	}

	queryItem := strings.Split(q, ",")
	for _, item := range queryItem {
		queryVal := strings.Split(item, ":")
		if len(queryVal) != 2 {
			return query, nil
		}
		columnName := queryVal[0]
		val := queryVal[1]
		valueType, ok := fieldsType[columnName]
		if !ok {
			return query, errors.New("unknown column")
		}
		columnName = fmt.Sprintf(`"%s".%s`, tableName, columnName)

		switch valueType {
		case int4Type:
			if strings.HasPrefix(val, "||") {
				condition = fmt.Sprintf("%s <= %s", columnName, val[2:])
			} else if strings.HasSuffix(val, "||") {
				condition = fmt.Sprintf("%s >= %s", columnName, val[:len(val)-2])
			} else if rangeRegexp.Match([]byte(val)) {
				values := strings.Split(val, ";;")
				condition = fmt.Sprintf("%s > %s AND %s <= %s", columnName, values[0], columnName, values[1])
			} else if len(strings.Split(val, "~")) > 1 && !strings.Contains(val, "~~") {
				condition = fmt.Sprintf("%s IN(%s)", columnName, strings.Replace(val, "~", ",", -1))
			} else {
				condition = fmt.Sprintf("%s = %s", columnName, val)
			}
		case varcharType:
			if !strings.Contains(val, "||") && strings.Contains(val, "|") {
				condition = createEnumeration(columnName, val)
			} else if !strings.HasPrefix(val, "~") {
				condition = fmt.Sprintf("%s ILIKE '%%%s%%'", columnName, val[1:])
			} else {
				condition = fmt.Sprintf("%s = %s", columnName, val)
			}
		case timestampType:
			if strings.HasPrefix(val, "||") {
				condition = fmt.Sprintf("extract(epoch from %s) < %s", columnName, val[2:])
			} else if strings.HasSuffix(val, "||") {
				condition = fmt.Sprintf("extract(epoch from %s) > %s", columnName, val[:len(val)-2])
			} else if rangeRegexp.Match([]byte(val)) {
				values := strings.Split(val, ";;")
				condition = fmt.Sprintf("extract(epoch from %s) > %s AND extract(epoch from %s) <= %s", columnName, values[0], columnName, values[1])
			} else {
				condition = fmt.Sprintf("%s = %s", columnName, val)
			}
		default:
			condition = fmt.Sprintf("%s = %s", columnName, val)
		}
		query.Where(condition)
	}

	return query, nil
}
