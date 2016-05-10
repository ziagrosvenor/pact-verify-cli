package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"os"
	"regexp"
	"testing"
)

func TestPrintHelp(t *testing.T) {
	runStdoutTest(t, PrintHelp, helpOutput)
}

func TestGetSrcDir(t *testing.T) {
	var testSrcDir = "test/src/dir"
	var oldSrcDir = os.Getenv("CLI_SRC_DIR")

	os.Setenv("CLI_SRC_DIR", testSrcDir)
	var srcDir = GetSrcDir()

	assert.Equal(t, testSrcDir, srcDir)

	os.Setenv("CLI_SRC_DIR", oldSrcDir)
}

func TestGetPwd(t *testing.T) {
	var pwd = GetPwd()
	var expectedPath = "pact-verify/_test"

	matched, err := regexp.MatchString(expectedPath, pwd)

	if err != nil {
		log.Fatal(err)
	}

	if !matched {
		t.Error("Expected:\n"+expectedPath+"\n\n", pwd)
	}
}

func runStdoutTest(t *testing.T, command func(), expectedOutput string) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	command()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	assert.Equal(t, expectedOutput, out)
}

var helpOutput = "NAME:\n   verify - Command line interface for Pact verification\n\nUSAGE:\n   pact-verify [global options] command [command options] [arguments...]\n   \nVERSION:\n   0.0.0\n   \nCOMMANDS:\nGLOBAL OPTIONS:\n   --pact PATH\t\t\tRead a Pact file from PATH and process it\n   --provider URL, --prov URL\tThe URL of the provider service to verify the pact with\n   --setup URL, -s URL\t\tThe URL of the provider state server - This is used to process provider states\n   --help, -h\t\t\tshow help\n   --version, -v\t\tprint the version\n   \n"
