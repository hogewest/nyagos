package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"../alias"
	"../dos"
	. "../interpreter"
)

const (
	WHICH_NOT_FOUND ErrorLevel = 1
)

func envToList(first1, env string) []string {
	list1 := strings.Split(os.Getenv(env), ";")
	result := make([]string, 1, len(list1))
	result[0] = first1
	result = append(result, list1...)
	return result
}

func cmd_which(cmd *Interpreter) (ErrorLevel, error) {
	all := false
	var pathList []string
	var extList []string
	for _, name := range cmd.Args[1:] {
		if name == "-a" {
			all = true
			pathList = envToList(".", "PATH")
			extList = envToList("", "PATHEXT")
			continue
		}
		if a, ok := alias.Table[strings.ToLower(name)]; ok {
			fmt.Fprintf(cmd.Stdout, "%s: aliased to %s\n", name, a.String())
			if !all {
				continue
			}
		}
		if _, ok := BuildInCommand[name]; ok {
			fmt.Fprintf(cmd.Stdout, "%s: built-in command\n", name)
			if !all {
				continue
			}
		}
		if all {
			for _, dir1 := range pathList {
				for _, ext1 := range extList {
					fullpath1 := dos.Join(dir1, name)
					fullpath1 = fullpath1 + ext1
					if _, err1 := os.Stat(fullpath1); err1 == nil {
						fmt.Fprintln(cmd.Stdout, fullpath1)
					}
				}
			}

		} else {
			path, err := exec.LookPath(name)
			if err != nil {
				return WHICH_NOT_FOUND, err
			}
			fmt.Fprintln(cmd.Stdout, dos.YenYen2Yen(path))
		}
	}
	return NOERROR, nil
}
