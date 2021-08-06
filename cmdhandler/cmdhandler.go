package cmdhandler

import (
	"cmdtools/core/dragon/conf"
	"log"
	"os/exec"
)

// HandleArgs ./bin/dragon gen domain [table name]
func HandleArgs(args []string) {
	if len(args) < 4 {
		log.Fatalln("cmd args err, create domain use cmd:\n ./bin/dragon gen domain [table_name] [table_name]")
	}
	switch {
	case args[1] == "gen" && args[2] == "domain":
		// generate domain files(entity, repository, service)
		GenDomain(args)
	}

	// format code

	path := conf.ExecDir+ "../domain"
	path = conf.FmtSlash(path)
	exec.Command("gofmt", "-w", "-l", path).Run()
}
