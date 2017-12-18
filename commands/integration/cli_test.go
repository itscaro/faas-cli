package integration

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

var tmpBinaryPath string

func TestMain(m *testing.M) {
	setup()

	exitCode := m.Run()

	shutdown()

	os.Exit(exitCode)
}

func setup() {
	tmpBinaryPath := getTmpBinaryPath()

	fmt.Printf(tmpBinaryPath)

	build := exec.Command("go")
	build.Args = append(build.Args, "build", "-o", tmpBinaryPath, "../..")
	if err := build.Run(); err != nil {
		fmt.Printf("could not make binary for %s: %v", tmpBinaryPath, err)
		os.Exit(1)
	}
}

func shutdown() {
	os.RemoveAll(tmpBinaryPath)
}
