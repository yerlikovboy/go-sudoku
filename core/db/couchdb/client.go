package couchdb

import (
	"fmt"
	"log"
	"net/http"
)

type couchDBClient struct {
	clnt *http.Client
	cfg  config
}

func NewClient() couchDBClient {
	return couchDBClient{
		clnt: &http.Client{},
		cfg:  defaultConfig(),
	}

}
func (c *couchDBClient) Do(req *http.Request) (*http.Response, error) {
	if c.cfg.creds.useAuth {
		req.SetBasicAuth(c.cfg.creds.admin, c.cfg.creds.pw)
	}

	// host + port
	req.URL.Host = fmt.Sprintf("%s:%s", c.cfg.host, c.cfg.port)

	// header
	req.Header.Add("Content-Type", "application/json")

	log.Printf("request %v", req.URL)

	return c.clnt.Do(req)

}
