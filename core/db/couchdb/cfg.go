package couchdb

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// used for basic auth.
type auth struct {
	use_auth bool
	admin    string
	pw       string
}

//TODO: does anyone outside this package need access to this?
type Config struct {
	host  string
	port  string
	creds auth
}

func (c Config) SetupRequest(req *http.Request) {
	// basic auth ...
	if c.creds.use_auth {
		req.SetBasicAuth(c.creds.admin, c.creds.pw)
	}

	// host + port
	req.URL.Host = fmt.Sprintf("%s:%s", c.host, c.port)

	// header
	req.Header.Add("Content-Type", "application/json")
}

func DefaultConfig() Config {
	return Config{
		host: db_hostname(),
		port: db_port(),
		creds: auth{
			use_auth: use_auth(),
			admin:    admin(),
			pw:       pw(),
		},
	}
}

func use_auth() bool {
	use_auth := os.Getenv("DB_NO_AUTH")
	if use_auth == "" {
		return true
	}
	return false
}

func db_port() string {
	port := os.Getenv("DB_PORT")
	if port == "" {
		return fmt.Sprint(5984)
	}
	return port
}

func db_hostname() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return "localhost"
	}
	return host
}

func admin() string {
	user := os.Getenv("DB_USER")
	if user == "" {
		return "admin"
	}
	return user
}

func pw() string {
	pw := os.Getenv("DB_PW")
	if len(pw) == 0 && use_auth() {
		log.Fatal("unable to retrieve db password (DB_PW not set)")
	}
	return pw
}
