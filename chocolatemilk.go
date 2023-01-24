package chocolatemilk

import (
	"errors"
	"os"
	"syscall"
)

type App struct {
	Name      string
	DebugMode bool
	Version   string
}

func New(name string, debug bool, version string) (*App, error) {
	app := &App{}

	var dirStructure = [5]string{
		"cmd",
		"cmd/web",
		"migrations",
		"ui",
		"logs",
	}

	for _, dirName := range dirStructure {
		err := os.Mkdir(dirName, 0644)
		if err != nil && !errors.Is(err, syscall.ERROR_ALREADY_EXISTS) {
			return nil, err
		}
	}

	buf, err := os.ReadFile(".env")
	if err != nil {
		return nil, err
	}
	err = loadEnvironmentVariables(buf)

	return app, err
}
