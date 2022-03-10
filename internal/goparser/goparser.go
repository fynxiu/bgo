package goparser

import (
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
	"strings"
)

func GetImports(serviePath string, modName string, wd string, imports map[string]struct{}) error {
	fset := token.NewFileSet()
	root := path.Join(wd, serviePath)
	fsys := os.DirFS(root)
	err := fs.WalkDir(fsys, ".", func(rpath string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			return nil
		}
		pkgs, err := parser.ParseDir(fset, path.Join(root, rpath), func(info fs.FileInfo) bool {
			return strings.HasSuffix(info.Name(), ".go") && !strings.HasSuffix(info.Name(), "_test.go")
		}, parser.ParseComments)
		if err != nil {
			return err
		}

		for _, v := range pkgs {
			for _, x := range v.Files {
				for _, y := range x.Imports {
					p := strings.TrimSuffix(strings.TrimPrefix(y.Path.Value, `"`), `"`)
					if strings.HasPrefix(p, path.Join(modName, serviePath)) {
						continue
					}
					if !strings.HasPrefix(p, modName) {
						continue
					}
					imports[p] = struct{}{}
					GetImports(strings.TrimPrefix(p, modName), modName, wd, imports)
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
