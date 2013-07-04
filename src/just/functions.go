package just

import (
	"strings"
)

func Trim(s string) string {
	return strings.Trim(s, " \t\n\r")
}
