package chocolatemilk

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
)

type App struct {
	Name           string
	DebugMode      bool
	Version        string
	ErrLog         *log.Logger
	InfoLogger     *log.Logger
	Addr           string
	TemplateEngine string
}

func New() (*App, error) {
	app := &App{}

	// Setting up loggers
	app.ErrLog = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Lshortfile|log.Ltime)
	app.InfoLogger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)

	app.InfoLogger.Println("Setting up directory structure")
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
	app.TemplateEngine = os.Getenv("TEMPLATE_ENGINE")

	return app, err
}

func (a *App) Listen() error {
	s := http.Server{
		Addr:     a.Addr,
		ErrorLog: a.ErrLog,
		Handler:  http.HandlerFunc(a.WelcomePage),
	}
	a.InfoLogger.Printf("Listening at %s", a.Addr)
	return s.ListenAndServe()
}

func (a *App) WelcomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Chocolate Milk! Your sweetest Go web framework :)")
}
