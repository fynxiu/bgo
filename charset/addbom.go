package charset

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	flagExts []string
	flagDir  string
)

func addBomCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addbom",
		Short: "Add BOM to files",
		Long:  "Add BOM to files in the current directory, e.g. bgo charset addbom",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("addbom exts=%v", flagExts)
			filepath.Walk(flagDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				if !isQualifiedFile(path) {
					return nil
				}
				if err := addBomToUTF8(path); err != nil {
					return err
				}
				return nil
			})
		},
	}
	cmd.Flags().StringSliceVar(&flagExts, "exts", []string{"h", "hpp", "c", "cpp"}, "extensions to add BOM, default: h, hpp, c, cpp")
	cmd.Flags().StringVar(&flagDir, "dir", ".", "directory to add BOM, default: current directory")

	return cmd
}

var bom = []byte{0xef, 0xbb, 0xbf}

func addBomToUTF8(filepath string) error {
	f, err := os.OpenFile(filepath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	head := make([]byte, 3)
	n, err := f.Read(head)
	if err != nil {
		if err != io.EOF {
			return err
		} else {
			return nil
		}
	}
	if n == 3 && head[0] == bom[0] && head[1] == bom[1] && head[2] == bom[2] {
		return nil
	}
	old, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	if _, err := f.WriteAt(append(bom, old...), 0); err != nil {
		return err
	}
	return nil
}

func isQualifiedFile(filepath string) bool {
	for _, ext := range flagExts {
		if strings.HasSuffix(filepath, ext) {
			return true
		}
	}
	return false
}
