package temp

import (
	"fmt"
	"os"
	"path"
	"time"
)

func TempDir(group string) string {
	return tempDir(group, time.Now)
}

func tempDir(group string, now func() time.Time) string {
	return path.Join(os.TempDir(), fmt.Sprintf("bgo-%v-%s", group, now().Format("20060102")))
}
