package fir

import (
	"log"

	"github.com/fynxiu/bgo/fir/internal/build"
	"github.com/fynxiu/bgo/fir/internal/config"
	"github.com/fynxiu/bgo/fir/internal/docker"
	"github.com/fynxiu/bgo/fir/internal/firlog"
	"github.com/fynxiu/bgo/fir/internal/git"
	"github.com/fynxiu/bgo/fir/internal/service"
	"github.com/spf13/cobra"
)

var (
	configFilename = ".fir.yml"
	logFilename    = ".firlog"
)

// CmdFir represents the upgrade command.
var CmdFir = &cobra.Command{
	Use:   "fir",
	Short: "build tools for mono-repo",
	Long:  "build tools for mono-repo, eg., go build, docker build, etc.",
	Run:   Run,
}

func init() {
	CmdFir.Flags().StringVarP(&configFilename, "config", "c", configFilename, "config file name")
	CmdFir.Flags().StringVarP(&logFilename, "log", "l", logFilename, "log file name")
}

// Run upgrade the bgo tools.
func Run(_ *cobra.Command, _ []string) {
	var (
		err             error
		changedFileList []string
		firLog          *firlog.FirLog
	)

	// read config
	c, err := config.FromFile(configFilename)
	if err != nil {
		log.Fatalf("failed to read config from %q, %v\n", configFilename, err)
	}
	services := c.Services

	// read .firlog
	firLog, err = firlog.FromFile(logFilename)
	if err != nil {
		log.Fatalf("failed to read %q, %v\n", logFilename, err)
		goto build
	}

	// build all
	if c.BuildAll {
		goto build
	}

	if firLog.NeedBuildAll() {
		goto build
	}
	if firLog.NoChange() {
		log.Println("no change")
		return
	}

	log.Printf("last commit is %s\n", firLog.Commit)
	// get all changed files
	changedFileList, err = git.ChangedFileList(firLog.Commit)
	log.Printf("changed files: \n%v\n", changedFileList)

	services, err = service.ServicesToBuild(services, changedFileList, c.Excluded)
	if err != nil {
		log.Fatalf("ServicesToBuild failed, %v\n", err)
	}
	log.Printf("services to build: \n%v\n", services)

build:
	log.Printf("build services, %v", services)
	if err = build.Build(services, c.BuildPath); err != nil {
		log.Fatalf("failed to build serives %v, %v\n", services, err)
	}

	exeChangedServices, err := firLog.ExeChangedServices(c)
	if err != nil {
		log.Fatalf("ExeChangedServices failed, %v", err)
	}
	log.Printf("ExeChanged services: %v", exeChangedServices)
	// overwrites .firlog
	firLog.Overwrite(logFilename)

	if c.Dockerizing {
		ss := service.ServicesToDockerize(services, changedFileList, exeChangedServices)
		log.Printf("services to dockerize, %v\n", ss)

		// rebuild services again to embed version info
		if err := build.BuildWithVersion(ss, c.BuildPath); err != nil {
			log.Fatalf("BuildWithVersion failed, %v", err)
		}

		// find related service according config and dependency graph, do build (given qualified services)
		// if dockerizing enabled, do docker build
		// should ignore unchanged artifacts
		// should tag images properly
		if err := docker.Build(ss, c); err != nil {
			log.Fatalf("docker process failed, %v", err)
		}
	}

	log.Println("fir done.")
}
