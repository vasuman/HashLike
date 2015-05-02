package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	"github.com/vasuman/HashLike/db"
	"github.com/vasuman/HashLike/handlers"
)

var logger *log.Logger

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

// Command-line arguments
var (
	logPath string
	dbPath  string
	port    int
)

func init() {
	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.StringVar(&logPath, "logfile", "", "Path to log file. Defaults to stdout if not specified")
	flag.StringVar(&dbPath, "db", "test.db", "Path to bolt database")
	flag.Parse()
}

const dbFileMode = 0660

func main() {
	var err error
	logger, err = setupLogger(logPath)
	if err != nil {
		fmt.Printf("failed to setup logger - %v\n", err)
		return
	}
	dbInst, err := bolt.Open(dbPath, dbFileMode, bolt.DefaultOptions)
	if err != nil {
		fmt.Printf("failed to initialize database - %v\n", err)
		return
	}
	defer dbInst.Close()
	err = db.Init(dbInst)
	if err != nil {
		fmt.Printf("failed to setup database - %v\n", err)
		return
	}
	logger.Printf("setup database")
	addr := fmt.Sprintf(":%d", port)
	logger.Println("listening on address, ", addr)
	logger.Println("starting server...")
	err = http.ListenAndServe(addr, handlers.GetRootHandler(logger))
	logger.Fatal(err)
}
