package fs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

const dirPerm = 0o700

func BgoHomeWithDir(dir string) string {
	home := path.Join(bgoHome(), dir)
	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := os.MkdirAll(home, dirPerm); err != nil {
			log.Fatal(err)
		}
	}
	return home
}

func Tree(path string, dir string) {
	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fmt.Printf("%s %s (%v bytes)\n", color.GreenString("CREATED"), strings.Replace(path, dir+"/", "", -1), info.Size())
		}
		return nil
	})
}

func bgoHome() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	home := path.Join(dir, ".bgo")
	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := os.MkdirAll(home, dirPerm); err != nil {
			log.Fatal(err)
		}
	}
	return home
}

func CopyDir(src, dst string, replaces, ignores []string) error {
	var err error
	var fds []os.FileInfo
	var srcInfo os.FileInfo

	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		if hasSets(fd.Name(), ignores) {
			continue
		}
		srcPath := path.Join(src, fd.Name())
		dstPath := path.Join(dst, fd.Name())
		if fd.IsDir() {
			if err = CopyDir(srcPath, dstPath, replaces, ignores); err != nil {
				return err
			}
		} else {
			if err = copyFile(srcPath, dstPath, replaces); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string, replaces []string) error {
	var err error
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	var old string
	for i, next := range replaces {
		if i%2 == 0 {
			old = next
			continue
		}
		buf = bytes.ReplaceAll(buf, []byte(old), []byte(next))
	}
	return ioutil.WriteFile(dst, buf, srcInfo.Mode())
}

func hasSets(name string, sets []string) bool {
	for _, ig := range sets {
		if ig == name {
			return true
		}
	}
	return false
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
