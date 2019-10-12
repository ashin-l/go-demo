package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ashin-l/go-demo/pkg/util"
)

const (
	defLogLevel = "error"
	defTLS      = "false"
	defDBHost   = "localhost"
	defDBPort   = "5432"
	defNum      = "37"

	envLogLevel = "MY_LOG_LEVEL"
	envTLS      = "MY_TLS"
	envDBHost   = "MY_DB_HOST"
	envDBPort   = "MY_DB_PORT"
	envNum      = "MY_NUM"
)

type config struct {
	logLevel string
	dbtls    bool
	dbhost   string
	dbport   string
	num      int64
}

func loadConfig() config {
	tls, err := strconv.ParseBool(util.Env(envTLS, defTLS))
	if err != nil {
		log.Fatalf("Invalid value passed for %s\n", envTLS)
	}

	tnum, err := strconv.ParseInt(util.Env(envNum, defNum), 10, 64)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", envNum, err.Error())
	}

	return config{
		logLevel: util.Env(envLogLevel, defLogLevel),
		dbtls:    tls,
		dbhost:   util.Env(envDBHost, defDBHost),
		dbport:   util.Env(envDBPort, defDBPort),
		num:      tnum,
	}
}

func main() {
	cfg := loadConfig()
	fmt.Printf("%+v\n", cfg)
}
