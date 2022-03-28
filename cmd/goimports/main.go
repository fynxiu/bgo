package main

import (
	"fmt"
	"os"

	"github.com/fynxiu/bgo/internal/gomod"
	"github.com/fynxiu/bgo/internal/goparser"
)

var (
	version   string
	buildTime string
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error
	fmt.Printf("version=%s, buildTime=%s", version, buildTime)

	modName, err := gomod.ModulePath("/home/fyn/Workspace/bntl/bgo")
	must(err)
	println("mod name:", modName)

	wd, err := os.Getwd()
	must(err)
	println("wd:", modName)

	modName = "jktcloud"
	wd = "/home/fyn/Workspace/bohai/jkt/jktcloud"
	serivcePath := "service/depa/rpc"
	imports := make(map[string]struct{})
	err = goparser.GetImports(serivcePath, modName, wd, imports)
	must(err)
	fmt.Printf("%v\n", imports)

	must(err)
}
