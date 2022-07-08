package graphpg

import (
	"fmt"
	"strings"
)

func createEnumeration(column string, val string) string {
	var (
		condition string
		arr       []string
	)

	vals := strings.Split(val, "|")

	for _, v := range vals {
		arr = append(arr, fmt.Sprintf("%s = %s", column, v))
	}

	condition = "(" + strings.Join(arr, " OR ") + ")"

	return condition
}
