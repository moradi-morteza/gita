package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
)

func main() {

	if len(os.Args) < 2 {
		color.Red("usage : gita <git arguments>")
		fmt.Println(`example : gita commit -m "my commit message"`)
		os.Exit(1)
	}

	args := os.Args[1:]
	command := args[0]

	normalizedArgs, flags := cleanArgs(args)

	ctx := &CommandContext{
		Command:   command,
		Args:      args,
		CleanArgs: normalizedArgs,
		Flags:     flags,
	}

	if err := runMiddlewares(ctx); err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	cmd := exec.Command("git", ctx.CleanArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		color.Red(err.Error())
	}
}

func cleanArgs(args []string) ([]string, map[string]bool) {
	cleanArgs := []string{}

	flags := map[string]bool{
		FlagForceMaster: false,
		FlagAICommit:    false,
		FlagDryRun:      false,
		FlagDebug:       false,
	}

	for _, arg := range args {
		if _, exists := flags[arg]; exists {
			flags[arg] = true
		} else {
			cleanArgs = append(cleanArgs, arg)
		}
	}
	return cleanArgs, flags
}
