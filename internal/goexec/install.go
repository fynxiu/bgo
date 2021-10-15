package goexec

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Install go get path.
func Install(path ...string) error {
	for _, p := range path {
		log.Printf("go install %s@latest\n", p)
		cmd := exec.Command("go", "install", fmt.Sprintf("%s@latest", p))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
