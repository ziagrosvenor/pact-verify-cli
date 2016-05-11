package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/ziagrosvenor/pact-verify-cli/src/pact-verify/run"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func PrintHelp() {
	cmd := exec.Command("pact-verify", "help")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func main() {
	var pactFilePath string
	var providerUrl string
	var stateServerUrl string

	app := cli.NewApp()
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "pact",
			Usage:       "Read a Pact file from `PATH` and process it",
			Destination: &pactFilePath,
		},

		cli.StringFlag{
			Name:        "provider, prov",
			Usage:       "The `URL` of the provider service to verify the pact with",
			Destination: &providerUrl,
		},

		cli.StringFlag{
			Name:        "setup, s",
			Usage:       "The `URL` of the provider state server - This is used to process provider states",
			Destination: &stateServerUrl,
		},
	}

	app.Name = "pact-verify"
	app.Usage = "Command line interface for Pact verification"
	app.Action = func(c *cli.Context) error {
		if pactFilePath == "" {
			return cli.NewExitError("\nEXITED \nA Pact file path is required i.e /tmp/pacts/pact.json \n", 86)
		}

		if providerUrl == "" {
			return cli.NewExitError("\nEXITED \nA provider service URL is required\n", 86)
		}

		if stateServerUrl == "" {
			return cli.NewExitError("\nEXITED \nA provider states setup service URL is required\n", 86)
		}

		var ROOT_DIR = GetSrcDir()
		var PWD = GetPwd()

		var pactHelperStr = run.BuildPactHelperFromPactJson(PWD, pactFilePath)
		run.WritePactHelperFile(ROOT_DIR, pactHelperStr)
		run.PactVerify(ROOT_DIR, PWD, pactFilePath, providerUrl, stateServerUrl)

		return nil
	}

	app.Run(os.Args)
}

var errorExitMessage = `
EXITING due to an error.

Please ensure all CLI options are correctly configured.
Try pact-verify help for more information.

ERROR:
`

func check(e error) {
	if e != nil {
		fmt.Println(errorExitMessage)
		log.Fatal(e)
	}
}

func GetPwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	check(err)

	return dir
}

func GetSrcDir() string {
	var env = os.Getenv("CLI_SRC_DIR")

	if env != "" {
		return env
	}

	return os.Getenv("GOPATH")
}
