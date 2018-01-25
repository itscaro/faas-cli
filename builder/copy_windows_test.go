package builder

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func checkDestinationFiles(dir string, numberOfFiles, mode int) error {
	// Check each file inside the destination folder
	for i := 1; i <= numberOfFiles; i++ {
		f, err := os.Stat(fmt.Sprintf("%s/test-file-%d", dir, i))
		if os.IsNotExist(err) {
			return err
		}
		if f.IsDir() {
			return errors.New("expected a file not a directory")
		}
	}

	return nil
}
