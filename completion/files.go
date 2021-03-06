package completion

import (
	"os"
	"regexp"
	"strings"

	"../dos"
	"../interpreter"
)

const (
	STD_SLASH = "\\"
	OPT_SLASH = "/"
)

var rxEnvPattern = regexp.MustCompile("%[^%]+%")
var rxTilde = regexp.MustCompile("^~[/\\\\]")

func listUpFiles(str string) ([]string, error) {
	orgSlash := STD_SLASH[0]
	if pos := strings.IndexAny(str, STD_SLASH+OPT_SLASH); pos >= 0 {
		orgSlash = str[pos]
	}
	str = rxEnvPattern.ReplaceAllStringFunc(str, func(p string) string {
		if len(p) == 2 {
			return "%"
		} else if val := os.Getenv(p[1 : len(p)-1]); val != "" {
			return val
		} else if f, ok := interpreter.PercentFunc[p[1:len(p)-1]]; ok {
			return f()
		} else {
			return p
		}
	})

	str = rxTilde.ReplaceAllStringFunc(str, func(p string) string {
		if home := dos.GetHome(); home != "" {
			return home + "\\"
		} else {
			return p
		}
	})
	str = strings.Replace(strings.Replace(str, OPT_SLASH, STD_SLASH, -1), "\"", "", -1)
	directory := dos.DirName(str)
	wildcard := dos.Join(directory, "*")

	// Drive letter
	cutprefix := 0
	if strings.HasPrefix(directory, STD_SLASH) {
		wd, _ := os.Getwd()
		directory = wd[0:2] + directory
		cutprefix = 2
	}
	commons := make([]string, 0)
	STR := strings.ToUpper(str)
	fdErr := dos.ForFiles(wildcard, func(fd *dos.FileInfo) bool {
		if fd.Name() == "." || fd.Name() == ".." || fd.IsHidden() {
			return true
		}
		name := dos.Join(directory, fd.Name())
		if fd.IsDir() {
			name += STD_SLASH
		}
		if cutprefix > 0 {
			name = name[2:]
		}
		nameUpr := strings.ToUpper(name)
		if strings.HasPrefix(nameUpr, STR) {
			if orgSlash != STD_SLASH[0] {
				name = strings.Replace(name, STD_SLASH, OPT_SLASH, -1)
			}
			commons = append(commons, name)
		}
		return true
	})
	return commons, fdErr
}
