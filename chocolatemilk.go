package chocolatemilk

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const version = "1.0.0"

type config struct {
	port     string
	renderer string
}

type ChocolateMilk struct {
	AppName    string
	Debug      bool
	Version    string
	ErrLogger  *log.Logger
	InfoLogger *log.Logger
	RootPath   string
	Routes     http.Handler
	config     config
}

func (c *ChocolateMilk) Init() error {
	var err error

	c.RootPath, err = os.Getwd()
	if err != nil {
		return err
	}

	err = c.CreateDirectoryStructure()
	if err != nil {
		return err
	}

	err = c.CreateEnvironmentFile()
	if err != nil {
		return err
	}

	err = godotenv.Load(filepath.FromSlash(c.RootPath + "/.env"))
	if err != nil {
		return err
	}

	c.Version = version
	c.AppName = os.Getenv("APP_NAME")
	c.config.port = os.Getenv("SERVER_PORT")
	c.config.renderer = os.Getenv("RENDERER")
	c.Routes = c.routes()
	c.Debug, err = strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		return err
	}

	c.ErrLogger = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	c.InfoLogger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)

	return nil
}

func (c *ChocolateMilk) ListenAndServe() error {
	s := &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%s", c.config.port),
		Handler:      c.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 300 * time.Second,
	}
	c.InfoLogger.Printf("Starting the http server at %s", s.Addr)
	return s.ListenAndServe()
}

func (c *ChocolateMilk) CreateDirectoryStructure() error {
	initialDirectories := []string{"handlers", "migrations", "views", "data", "static", "tmp", "logs", "middleware"}
	for _, dirName := range initialDirectories {
		path := fmt.Sprintf("%s/%s", c.RootPath, dirName)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err := os.Mkdir(path, 0755)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *ChocolateMilk) CreateEnvironmentFile() error {
	envFilePath := fmt.Sprintf("%s/.env", c.RootPath)

	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		envFile, err := os.Create(envFilePath)
		if err != nil {
			return err
		}
		defer func(envFile *os.File) {
			err := envFile.Close()
			if err != nil {
				panic(err)
			}
		}(envFile)
	}

	return nil
}
