package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vasuman/HashLike/models"
	"gopkg.in/yaml.v2"
)

var (
	logger *log.Logger
	cfg    *config
)

type config struct {
	Port int
	Db   struct {
		Driver string
		Source string
	}
	AllowedSites []struct {
		Domain string
		Paths  []string
	} `yaml:"allowed_sites"`
}

func (cfg *config) parseSites() error {
	var err error
	for _, allowedSite := range cfg.AllowedSites {
		s := new(models.Site)
		s.Domain, err = regexp.Compile(allowedSite.Domain)
		if err != nil {
			return fmt.Errorf("domain regex (%s) - %v", allowedSite.Domain, err)
		}
		for _, path := range allowedSite.Paths {
			pExp, err := regexp.Compile(path)
			if err != nil {
				return fmt.Errorf("path regex (%s) - %v", path, err)
			}
			s.Paths = append(s.Paths, pExp)
		}
		models.AddSite(s)
	}
	return nil
}

func loadConfig(configPath string) (*config, error) {
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	cfg := new(config)
	return cfg, yaml.Unmarshal(b, cfg)
}

func setupLogger(logPath string) (*log.Logger, error) {
	var logFile *os.File
	if logPath != "" {
		var err error
		lfs := os.O_CREATE | os.O_APPEND | os.O_WRONLY
		logFile, err = os.OpenFile(logPath, lfs, 0666)
		if err != nil {
			return nil, err
		}
	} else {
		logFile = os.Stdout
	}
	lFlags := log.Lshortfile
	return log.New(logFile, "", lFlags), nil
}

func main() {
	// Command-line arguments
	var (
		staticDir  string
		configPath string
		logPath    string
	)
	flag.StringVar(&configPath, "config", "config.yml", "Path to config file")
	flag.StringVar(&staticDir, "static", "./static/", "Path to static directory")
	flag.StringVar(&logPath, "logfile", "", "Path to log file. Defaults to stdout if not specified")
	flag.Parse()

	var err error
	cfg, err = loadConfig(configPath)
	if err != nil {
		fmt.Printf("failed to load config - %v\n", err)
		return
	}
	logger, err = setupLogger(logPath)
	if err != nil {
		fmt.Printf("failed to setup logger - %v\n", err)
		return
	}
	logger.Printf("using config,\n%+v\n", cfg)
	var db *sql.DB
	db, err = models.InitDb(cfg.Db.Driver, cfg.Db.Source)
	if err != nil {
		fmt.Printf("failed to initialize database - %v\n", err)
		return
	}
	defer db.Close()
	logger.Printf("initialized database")
	err = cfg.parseSites()
	if err != nil {
		fmt.Printf("error parsing 'allowed_sites' regexes\n%v\n", err)
		return
	}
	handler := getRootHandler(staticDir)
	addr := fmt.Sprintf(":%d", cfg.Port)
	logger.Println("listening on address, ", addr)
	logger.Println("starting server...")
	err = http.ListenAndServe(addr, handler)
	logger.Fatal(err)
}
