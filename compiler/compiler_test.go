package compiler

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/zegl/tre/cmd/tre/build"
)

func TestAllPrograms(t *testing.T) {
	files, _ := ioutil.ReadDir("testdata")
	if len(files) == 0 {
		t.Error("No test files found")
	}

	for _, file := range files {
		t.Run(file.Name(), func(t *testing.T) {
			if err := buildRunAndCheck(t, "testdata/"+file.Name()); err != nil {
				t.Error("failed: " + err.Error())
			}
		})
	}
}

func buildRunAndCheck(t *testing.T, path string) error {
	fp, err := os.Stat(path)
	if err != nil {
		return err
	}

	mainPath := path
	if fp.IsDir() {
		mainPath = path + "/main.go"
	}

	content, err := ioutil.ReadFile(mainPath)
	if err != nil {
		return err
	}

	expect := ""

	re, _ := regexp.Compile(`(?m)// (.*?)$`)
	for _, str := range re.FindAllString(string(content), -1) {
		expect += strings.Replace(str, "// ", "", -1) + "\n"
	}

	// Normalize newlines
	expect = strings.Replace(expect, "\r\n", "\n", -1)
	expect = strings.TrimSpace(expect)

	runProgram := true
	var output string

	outputBinaryPath := os.TempDir() + "/exec"

	// "GOROOT" (treroot?) detection
	_, testFilePath, _, _ := runtime.Caller(0)
	goroot := filepath.Clean(testFilePath + "/../../pkg/")

	err = build.Build(path,goroot, outputBinaryPath, false)
	if err != nil {
		output = strings.TrimSpace(err.Error())
		runProgram = false
	}

	// Run program output
	if runProgram {
		cmd := exec.Command(outputBinaryPath)
		stdout, err := cmd.CombinedOutput()
		if err != nil &&  err.Error() != "exit status 1" {
			output = output + strings.TrimSpace(err.Error())
		} else {
			output = output + strings.TrimSpace(string(stdout))
		}
	}

	if expect == output {
		return nil
	}

	t.Logf("Expected:\n---\n'%s'\n---\nResult:\n---\n'%s'\n---\n", expect, output)

	return errors.New("Unexpected result")
}
