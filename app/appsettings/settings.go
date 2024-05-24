package appsettings

import (
	"os"
	"time"
)

const MaxIdleConns = 2
const MaxOpenConns = 5
const MaxIdleTime = time.Minute * 5
const ConnMaxLifetime = time.Hour
const SelectRowsLimit = 50

const DateLayout string = "20060102"
const WebDir = "./web"
const HostName = "localhost"

var EnvPort = ""
var EnvDbfile = ""
var EnvPassword = ""

func GetEnvVariables() {
	EnvPort = os.Getenv("TODO_PORT")
	EnvDbfile = os.Getenv("TODO_DBFILE")
	EnvPassword = os.Getenv("TODO_PASSWORD")
}
