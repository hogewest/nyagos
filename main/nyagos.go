package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/shiena/ansicolor"

	"../alias"
	"../commands"
	"../completion"
	"../conio"
	"../dos"
	"../history"
	"../interpreter"
	"../lua"
)

var rxAnsiEscCode = regexp.MustCompile("\x1b[^a-zA-Z]*[a-zA-Z]")

var stamp string
var commit string
var version string
var ansiOut io.Writer

func nyagosPrompt(L lua.Lua) int {
	template, err := L.ToString(-1)
	if err != nil {
		template = "[" + err.Error() + "]"
	}
	text := Format2Prompt(template)
	fmt.Fprint(ansiOut, text)
	text = rxAnsiEscCode.ReplaceAllString(text, "")
	lfPos := strings.LastIndex(text, "\n")
	if lfPos >= 0 {
		text = text[lfPos+1:]
	}
	L.PushInteger(lua.Integer(conio.GetStringWidth(text)))
	return 1
}

var prompt_hook lua.Pushable = lua.TGoFunction{nyagosPrompt}

func printPrompt(this *conio.LineEditor) (int, error) {
	L := NewNyagosLua()
	defer L.Close()
	L.Push(prompt_hook)

	if !L.IsFunction(-1) {
		L.Pop(1)
		return 0, nil
	}
	L.PushString(os.Getenv("PROMPT"))
	if err := L.Call(1, 1); err != nil {
		return 0, err
	}
	length, lengthErr := L.ToInteger(-1)
	L.Pop(1)
	if lengthErr == nil {
		return length, nil
	} else {
		return 0, fmt.Errorf("nyagos.prompt: return-value(length) is invalid: %s", lengthErr.Error())
	}
}

func when_panic() {
	err := recover()
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, "************ Panic Occured. ***********")
	fmt.Fprintln(os.Stderr, err)
	debug.PrintStack()
	fmt.Fprintln(os.Stderr, "*** Please copy these error message ***")
	fmt.Fprintln(os.Stderr, "*** And hit ENTER key to quit.      ***")
	var dummy [1]byte
	os.Stdin.Read(dummy[:])
}

var luaFilter lua.Pushable = lua.TNil{}

func main() {
	defer when_panic()

	dos.CoInitializeEx(0, dos.COINIT_MULTITHREADED)
	defer dos.CoUninitialize()

	conio.DisableCtrlC()

	completion := conio.KeyGoFuncT{F: completion.KeyFuncCompletion}

	if err := conio.BindKeySymbolFunc(conio.K_CTRL_I, "COMPLETE", &completion); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	// ANSI Escape Sequence Support
	ansiOut = ansicolor.NewAnsiColorWriter(os.Stdout)

	commands.Init()
	alias.Init()

	// Lua extension
	L := NewNyagosLua()
	defer L.Close()

	if !optionParse(L) {
		return
	}

	appData := filepath.Join(os.Getenv("APPDATA"), "NYAOS_ORG")
	os.Mkdir(appData, 0777)
	histPath := filepath.Join(appData, "nyagos.history")
	history.Load(histPath)
	history.Save(histPath) // cut over max-line

	exeName, exeNameErr := dos.GetModuleFileName()
	if exeNameErr != nil {
		fmt.Fprintln(os.Stderr, exeNameErr)
	}
	exeFolder := filepath.Dir(exeName)
	nyagos_lua := filepath.Join(exeFolder, "nyagos.lua")
	if _, err := os.Stat(nyagos_lua); err == nil {
		err := L.Source(nyagos_lua)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	conio.DefaultEditor.Prompt = printPrompt

	for {
		it := interpreter.New()
		it.Tag = L
		it.CloneHook = func(this *interpreter.Interpreter) error {
			this.Tag = NewNyagosLua()
			it.CloseHook = func(this *interpreter.Interpreter) {
				if this.Tag != nil {
					if L, ok := this.Tag.(lua.Lua); ok {
						L.Close()
					}
					this.Tag = nil
				}
			}
			return nil
		}
		conio.DefaultEditor.Tag = it

		wd, wdErr := os.Getwd()
		if wdErr == nil {
			conio.SetTitle("NYAGOS - " + wd)
		} else {
			conio.SetTitle("NYAGOS - " + wdErr.Error())
		}
		history_count := conio.DefaultEditor.HistoryLen()
		line, err := conio.DefaultEditor.ReadLine()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			break
		}

		var isReplaced bool
		line, isReplaced = history.Replace(line)
		if isReplaced {
			fmt.Fprintln(os.Stdout, line)
		}
		if line == "" {
			continue
		}
		if conio.DefaultEditor.HistoryLen() > history_count {
			fd, err := os.OpenFile(histPath, os.O_APPEND, 0600)
			if err == nil {
				fmt.Fprintln(fd, line)
				fd.Close()
			} else {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		} else {
			conio.DefaultEditor.HistoryResetPointer()
		}

		stackPos := L.GetTop()
		L.Push(luaFilter)
		if L.IsFunction(-1) {
			L.PushString(line)
			err := L.Call(1, 1)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				if L.IsString(-1) {
					line, err = L.ToString(-1)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
					}
				}
			}
		}
		L.SetTop(stackPos)

		whatToDo, err := it.Interpret(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if whatToDo == interpreter.SHUTDOWN {
			break
		}
	}
}
