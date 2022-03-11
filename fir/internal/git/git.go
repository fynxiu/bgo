package git

import (
	"os"
	"os/exec"
	"strings"
)

const (
	headCommit = "HEAD"
)

func GetVersion() (string, error) {
	return gitRun("describe", "--tags")
}

func GetHeadHash() (string, error) {
	return gitRun("rev-parse", headCommit)
}

func ChangedFileList(commit string) ([]string, error) {

	return changedFilefList(commit, headCommit)
}

func changedFilefList(prevCommit, currCommit string) ([]string, error) {
	prevCommit = strings.TrimSpace(prevCommit)
	currCommit = strings.TrimSpace(currCommit)
	output, err := gitRun("diff", "--name-only", prevCommit, currCommit)
	if err != nil {
		return nil, err
	}
	return strings.Split(output, "\n"), nil
}

func gitRun(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
