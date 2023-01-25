package chocolatemilk

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"
)

type App struct {
	Name          string
	DebugMode     bool
	Version       string
	ErrLog        *log.Logger
	InfoLogger    *log.Logger
	Addr          string
	templateCache map[string]*template.Template
}

func New() (*App, error) {
	app := &App{}

	// Setting up loggers
	app.ErrLog = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Lshortfile|log.Ltime)
	app.InfoLogger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)

	app.InfoLogger.Println("Setting up directory structure")
	var dirStructure = [12]string{
		"cmd",
		"cmd/web",
		"migrations",
		"ui",
		"ui/html",
		"ui/html/pages",
		"ui/html/partials",
		"ui/static",
		"ui/static/img",
		"ui/static/js",
		"ui/static/css",
		"logs",
	}

	for _, dirName := range dirStructure {
		err := os.Mkdir(dirName, 0644)
		if err != nil && !errors.Is(err, syscall.ERROR_ALREADY_EXISTS) {
			return nil, err
		}
	}

	app.InfoLogger.Println("Parsing environment variables from .env file")
	buf, err := os.ReadFile(".env")
	if err != nil {
		return nil, err
	}
	err = loadEnvironmentVariables(buf)

	app.Name = os.Getenv("APP_NAME")
	app.DebugMode = strings.ToLower(os.Getenv("DEBUG")) == "true"
	app.Version = os.Getenv("VERSION")
	app.Addr = os.Getenv("ADDR")

	app.templateCache, err = newTemplateCache()
	if err != nil {
		return nil, err
	}

	return app, err
}

func (app *App) Listen() error {
	s := http.Server{
		Addr:         app.Addr,
		ErrorLog:     app.ErrLog,
		Handler:      app.routes(),
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}
	app.InfoLogger.Printf("Listening at %s\n", app.Addr)
	return s.ListenAndServe()
}
