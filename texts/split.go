package texts

import (
	"regexp"
	"strings"
)

var rxEscape = regexp.MustCompile(`\\["'\\]`)
var rxDoubleQuoted = regexp.MustCompile(`"[^"]*"`)
var rxSingleQuoted = regexp.MustCompile(`'[^']*'`)
var rxNonSpace = regexp.MustCompile(`[^ ]+`)

func SplitLikeShell(line string) [][]int {
	line = rxEscape.ReplaceAllString(line, "\001\001")
	line = rxDoubleQuoted.ReplaceAllStringFunc(line, func(str string) string {
		str = strings.Replace(str, " ", "\001", -1)
		str = strings.Replace(str, `'`, "\001", -1)
		return str
	})
	line = rxSingleQuoted.ReplaceAllStringFunc(line, func(str string) string {
		return strings.Replace(str, " ", "\001", -1)
	})
	return rxNonSpace.FindAllStringIndex(line, -1)
}
