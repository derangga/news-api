package utils

import (
	"fmt"
	"regexp"

	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
)

func IsNoRowError(err error) bool {
	return err.Error() == "sql: no rows in result set"
}

func IsDuplicateKey(err error) (bool, string) {
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		log.Errorf("Postgre Error: %w", pqErr)

		re := regexp.MustCompile(`Key \((.+?)\)=\((.+?)\) already exists`)
		matches := re.FindStringSubmatch(pqErr.Detail)
		if len(matches) == 3 {
			field := matches[1]
			value := matches[2]
			return true, fmt.Sprintf("%s: %s already exists", field, value)
		}
		return true, "Duplicate key value"
	}

	return false, ""
}
