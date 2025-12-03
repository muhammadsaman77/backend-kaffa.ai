package pkg

import (
	"regexp"
	"strings"
)

func CamelToSnake(s string) string {
	// Case: "RoleID" → "Role_Id"
	snake := regexp.MustCompile("(.)([A-Z][a-z]+)").ReplaceAllString(s, "${1}_${2}")
	// Case: "Role_Id" → "Role_ID"
	snake = regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
