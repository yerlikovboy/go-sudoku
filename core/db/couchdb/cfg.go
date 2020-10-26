package couchdb

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// used for basic auth.
type auth struct {
	useAuth bool
	admin   string
	pw      string
}

//TODO: does anyone outside this package need access to this?
type config struct {
	host  string
	port  string
	creds auth
}

func (c config) SetupRequest(req *http.Request) {
	// basic auth ...
	if c.creds.useAuth {
		req.SetBasicAuth(c.creds.admin, c.creds.pw)
	}

	// host + port
	req.URL.Host = fmt.Sprintf("%s:%s", c.host, c.port)

	// header
	req.Header.Add("Content-Type", "application/json")
}

func defaultConfig() config {
	return config{
		host: dbHostname(),
		port: dbPort(),
		creds: auth{
			useAuth: useAuth(),
			admin:   admin(),
			pw:      pw(),
		},
	}
}

func useAuth() bool {
	useAuth := os.Getenv("DB_NO_AUTH")
	if useAuth == "" {
		return true
	}
	return false
}

func dbPort() string {
	port := os.Getenv("DB_PORT")
	if port == "" {
		return fmt.Sprint(5984)
	}
	return port
}

func dbHostname() string {
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
	if len(pw) == 0 && useAuth() {
		log.Fatal("unable to retrieve db password (DB_PW not set)")
	}
	return pw
}
