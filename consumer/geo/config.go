package csgeo

import "flag"

const db_user = "maxroach"
const db_host_default = "localhost"
const db_port_default = 30257

//const db_host = "cockroachdb-public"

type Config struct {
	http_port int
	db_port   int
	db_host   string
	db_user   string
	db_name   string
}

func NewConfig() (*Config, error) {
	c := &Config{}
	port := flag.Int("port", 80, "port number")
	dbserver := flag.String("dbserver", db_host_default, "database server")
	dbport := flag.Int("dbport", db_port_default, "database port")
	flag.Parse()

	c.http_port = *port
	c.db_host = *dbserver
	c.db_port = *dbport
	c.db_user = db_user
	c.db_name = "consumer_geo"
	return c, nil
}

//func (c *Config) Set() error {
//}
