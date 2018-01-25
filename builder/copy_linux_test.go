package builder

import (
	"errors"
	"fmt"
	"os"
)

func checkDestinationFiles(dir string, numberOfFiles, mode int) error {
	// Check each file inside the destination folder
	for i := 1; i <= numberOfFiles; i++ {
		f, err := os.Stat(fmt.Sprintf("%s/test-file-%d", dir, i))
		if os.IsNotExist(err) {
			return err
		}
		if f.Mode() != os.FileMode(mode) {
			return errors.New("expected mode did not match")
		}
		if f.IsDir() {
			return errors.New("expected a file not a directory")
		}
	}

	return nil
}
