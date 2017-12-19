package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/mitchellh/cli"
)

func main() {
	runtime.GOMAXPROCS(4)
	os.Exit(realMain())
}

func realMain() int {

	// setup application logging
	log.SetFlags(log.LstdFlags)

	os.Mkdir("logs", 0666)
	logFile, err := os.OpenFile("logs/application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file: ", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	return wrappedMain()
}

func wrappedMain() int {

	// Get the command line args. We shortcut "--version" and "-v" to
	// just show the version.
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			newArgs := make([]string, len(args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], args)
			args = newArgs
			break
		}
	}

	cli := &cli.CLI{
		Args:     args,
		Commands: Commands,
		HelpFunc: cli.BasicHelpFunc("datica-users"),
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}
