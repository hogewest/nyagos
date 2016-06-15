package commands

import (
	"io"
	"os"
	"os/exec"

	"github.com/mattn/go-colorable"

	"./ls"
)

func cmd_ls(cmd *exec.Cmd) (int, error) {
	var out io.Writer
	if cmd.Stdout == os.Stdout {
		out = colorable.NewColorableStdout()
	} else {
		out = cmd.Stdout
	}
	return 0, ls.Main(cmd.Args[1:], out, cmd.Stderr)
}
