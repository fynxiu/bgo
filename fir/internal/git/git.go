package git

import (
	"os"
	"os/exec"
	"strings"
)

const (
	headCommit = "HEAD"
)

func GetHeadHash() (string, error) {
	cmd := exec.Command("git", "rev-parse", headCommit)
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func ChangedFileList(commit string) ([]string, error) {

	return changedFilefList(commit, headCommit)
}

func changedFilefList(prevCommit, currCommit string) ([]string, error) {
	prevCommit = strings.TrimSpace(prevCommit)
	currCommit = strings.TrimSpace(currCommit)
	cmd := exec.Command("git", "diff", "--name-only", prevCommit, currCommit)
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
}
