package util

import (
	"os/exec"
	"strings"
)

func GhUsername() (string, error) {
	ghCmd := exec.Command("gh", "api", "user", "--jq", ".login")
	var ghCmdOutBuf strings.Builder
	ghCmd.Stdout = &ghCmdOutBuf
	if err := ghCmd.Run(); err != nil {
		return "", err
	} else {
		return strings.TrimSpace(ghCmdOutBuf.String()), nil
	}
}

func GitDisplayname() (string, error) {
	gitCmd := exec.Command("git", "config", "user.name")
	var ghCmdOutBuf strings.Builder
	gitCmd.Stdout = &ghCmdOutBuf
	if err := gitCmd.Run(); err != nil {
		return "", err
	} else {
		return strings.TrimSpace(ghCmdOutBuf.String()), nil
	}
}
