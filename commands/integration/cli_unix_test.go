// +build !windows

package integration

import "os"

func getTmpBinaryPath() string {
	return os.TempDir() + "/faas-cli"
}
