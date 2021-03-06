package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/parnurzeal/gorequest"
	"golang.org/x/xerrors"
)

// MetasploitConf is metasploit config
type MetasploitConf struct {
	// DB type for metasploit dictionary (sqlite3, mysql, postgres or redis)
	Type string

	// http://metasploit-dictionary.com:1324 or DB connection string
	URL string `json:"-"`

	// /path/to/metasploit.sqlite3
	SQLite3Path string `json:"-"`
}

func (cnf *MetasploitConf) setDefault() {
	if cnf.Type == "" {
		cnf.Type = "sqlite3"
	}
	if cnf.URL == "" && cnf.SQLite3Path == "" {
		wd, _ := os.Getwd()
		cnf.SQLite3Path = filepath.Join(wd, "go-msfdb.sqlite3")
	}
}

const metasploitDBType = "METASPLOITDB_TYPE"
const metasploitDBURL = "METASPLOITDB_URL"
const metasploitDBPATH = "METASPLOITDB_SQLITE3_PATH"

// Init set options with the following priority.
// 1. Environment variable
// 2. config.toml
func (cnf *MetasploitConf) Init() {
	if os.Getenv(metasploitDBType) != "" {
		cnf.Type = os.Getenv(metasploitDBType)
	}
	if os.Getenv(metasploitDBURL) != "" {
		cnf.URL = os.Getenv(metasploitDBURL)
	}
	if os.Getenv(metasploitDBPATH) != "" {
		cnf.SQLite3Path = os.Getenv(metasploitDBPATH)
	}
	cnf.setDefault()
}

// IsFetchViaHTTP returns wether fetch via http
func (cnf *MetasploitConf) IsFetchViaHTTP() bool {
	return Conf.Metasploit.Type == "http"
}

// CheckHTTPHealth do health check
func (cnf *MetasploitConf) CheckHTTPHealth() error {
	if !cnf.IsFetchViaHTTP() {
		return nil
	}

	url := fmt.Sprintf("%s/health", cnf.URL)
	resp, _, errs := gorequest.New().Get(url).End()
	//  resp, _, errs = gorequest.New().SetDebug(config.Conf.Debug).Get(url).End()
	//  resp, _, errs = gorequest.New().Proxy(api.httpProxy).Get(url).End()
	if 0 < len(errs) || resp == nil || resp.StatusCode != 200 {
		return xerrors.Errorf("Failed to connect to metasploit server. url: %s, errs: %s", url, errs)
	}
	return nil
}
