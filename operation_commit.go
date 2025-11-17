package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"strings"
)

var errCommitCancelled = errors.New("commit cancelled")

func init() {
	registerMiddleware("commit", protectProtectedBranches)
}

func protectProtectedBranches(ctx *CommandContext) error {
	branchBytes, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return fmt.Errorf("unable to determine current branch: %w", err)
	}

	branch := strings.TrimSpace(string(branchBytes))
	if !isProtectedBranch(branch) {
		return nil
	}

	if ctx.Flags[FlagForceMaster] {
		color.Yellow("Forced commit on protected branch: %s", branch)
		return nil
	}

	color.Yellow("You are on %s branch. Continue? (y/n): ", branch)
	answer := strings.ToLower(readInput())
	if answer != "y" {
		return errCommitCancelled
	}

	return nil
}

func isProtectedBranch(branch string) bool {
	switch branch {
	case "main", "master":
		return true
	default:
		return false
	}
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	return strings.TrimSpace(answer)
}
