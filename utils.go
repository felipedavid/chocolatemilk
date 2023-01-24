package chocolatemilk

import (
	"bytes"
	"errors"
	"os"
	"runtime"
)

// loadEnvironmentVariables takes the contents of a file, parses the environment variables definitions setting them.
// OBS: This implementation is dump and hopefully temporary
func loadEnvironmentVariables(buf []byte) error {
	sep := []byte{'\n'}
	if runtime.GOOS == "windows" {
		sep = []byte{'\r', '\n'}
	}

	pairs := bytes.Split(buf, sep)
	for _, pair := range pairs {
		if pair[0] == '#' {
			continue
		}

		keyValue := bytes.Split(pair, []byte{'='})
		if len(keyValue) != 2 {
			return errors.New("malformed .env file")
		}

		err := os.Setenv(string(keyValue[0]), string(keyValue[1]))
		if err != nil {
			return err
		}
	}
	return nil
}
