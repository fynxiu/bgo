package service

import (
	"log"
	"os"
	"path"
	"strings"

	set "github.com/deckarep/golang-set"
	"github.com/fynxiu/bgo/fir/internal/config"
	"github.com/fynxiu/bgo/internal/gomod"
	"github.com/fynxiu/bgo/internal/goparser"
)

func ServicesToBuild(services []config.Service, changedFileList []string, excluded []string) ([]config.Service, error) {

	// filter service exclusive changed files
	_changedFileList := filterExcluded(excluded, filterNoSourceFiles(changedFileList))

	log.Printf("remain: %v", _changedFileList)
	// TODO: analyses go.mod
	if shouldBuildAll(_changedFileList) {
		return services, nil
	}

	var ret []config.Service
	for _, x := range services {
		for _, y := range changedFileList {
			if x.IsEmbeded(y) {
				ret = append(ret, x)
				break
			}
		}
	}
	if len(ret) == len(services) {
		return ret, nil
	}

	pathSet := set.NewSet()
	servicePaths := set.NewSet()
	for _, x := range services {
		servicePaths.Add(x.Path)
	}
	for _, x := range _changedFileList {
		p := path.Dir(x)
		if servicePaths.Contains(p) {
			pathSet.Add(p)
		}
	}

	var servicesNeedAnalysis []config.Service

	for _, x := range services {
		if !(pathSet.Contains(x.Path)) {
			servicesNeedAnalysis = append(servicesNeedAnalysis, x)
		} else {
			ret = append(ret, x)
		}
	}
	if len(ret) == len(services) {
		return ret, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	modName, err := gomod.ModulePath(wd)
	if err != nil {
		return nil, err
	}

	// remain changed files
	// get service's import
	changed := func(imports map[string]struct{}) bool {
		// TODO: _changedFileList should remove service-specific entities
		for _, x := range _changedFileList {
			if _, ok := imports[path.Dir(x)]; ok {
				return true
			}
		}
		return false
	}
	for _, x := range servicesNeedAnalysis {
		imports := make(map[string]struct{})
		if err := goparser.GetImports(x.Path, modName, wd, imports); err != nil {
			return nil, err
		}
		if changed(imports) {
			ret = append(ret, x)
		}
	}

	return ret, nil
}

func ServicesToDockerize(services []config.Service, changedFileList []string, exeChangedServiceNames []string) []config.Service {
	var ss []config.Service
	for _, x := range services {
		for _, y := range exeChangedServiceNames {
			if x.Name == y {
				ss = append(ss, x)
				break
			}
		}
	}

	for _, x := range changedFileList {
		for _, y := range services {
			if y.IsRelative(x) {
				ss = append(ss, y)
				// do not break
			}
		}
	}

	return ss
}

// filterNoSourceFiles only source code remains
func filterNoSourceFiles(changedFileList []string) (ret []string) {
	for _, x := range changedFileList {
		if strings.HasSuffix(x, ".go") && !strings.HasSuffix(x, "_test.go") {
			ret = append(ret, x)
		}
	}
	return
}

func filterExcluded(excluded []string, changedFileList []string) []string {
	isExclude := func(changedFile string) bool {
		for _, x := range excluded {
			if strings.HasPrefix(changedFile, x) {
				return true
			}
		}
		return false
	}
	var ret []string
	for _, x := range changedFileList {
		if !isExclude(x) {
			ret = append(ret, x)
		}
	}

	return ret
}

func shouldBuildAll(changedFileList []string) bool {
	for _, x := range changedFileList {
		if x == "go.mod" || x == "go.sum" {
			return true
		}
	}
	return false
}
