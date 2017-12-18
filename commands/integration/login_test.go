package integration

import (
	"os/exec"
	"strings"
	"testing"
)

func Test_Login_PasswordStdinAndPassword(t *testing.T) {
	tmpBinaryPath := getTmpBinaryPath()
	proc := exec.Command(tmpBinaryPath, "login", "--username=username_test", "--password=password_test", "--password-stdin")
	proc.Stdin = strings.NewReader("password_std")
	output, _ := proc.Output()

	if !strings.Contains(string(output), "--password and --password-stdin are mutually exclusive") {
		t.Fatalf("Output is not what expected:\n%s", output)
	}
}
