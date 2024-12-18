package main

import (
	"fmt"
	"os"

	"github.com/mikolajgs/broccli"
)

func main() {
	cli := broccli.NewCLI("website-generator", "Generates static HTML", "Mikolaj Gasior")
	cmdGen := cli.AddCmd("generate", "Generates HTML from a specified directory", generateHandler)
	cmdGen.AddFlag("source", "s", "", "Path to source directory", broccli.TypePathFile, broccli.IsExistent|broccli.IsDirectory|broccli.IsRequired)
	cmdGen.AddFlag("destination", "d", "", "Path to target directory", broccli.TypePathFile, broccli.IsExistent|broccli.IsRequired)
	_ = cli.AddCmd("version", "Prints version", versionHandler)
	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"App", "version"}
	}
	os.Exit(cli.Run())
}

func versionHandler(c *broccli.CLI) int {
	fmt.Fprintf(os.Stdout, VERSION+"\n")
	return 0
}

func generateHandler(c *broccli.CLI) int {
	website := Website{
		SourcePath: c.Flag("source"),
	}

	if err := website.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "!!!! Error with website initialization: %s\n", err.Error())
		return 1
	}

	gen := Generator{
		DestinationPath: c.Flag("destination"),
	}

	if err := gen.Generate(&website); err != nil {
		fmt.Fprintf(os.Stderr, "!!!! Error with generation: %s\n", err.Error())
		return 1
	}

	return 0
}
