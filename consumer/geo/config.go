package csgeo

import "flag"

const db_host_default = "localhost"
const db_port_default = 30257
const db_user = "maxroach"
const kafka_host_default = "localhost"
const kafka_port_default = 31090

//const db_host = "cockroachdb-public"

type Config struct {
	http_port  int
	db_host    string
	db_port    int
	db_user    string
	db_name    string
	kafka_host string
	kafka_port int
}

func NewConfig() (*Config, error) {
	c := &Config{}
	port := flag.Int("port", 80, "port number")
	dbserver := flag.String("dbserver", db_host_default, "database server")
	dbport := flag.Int("dbport", db_port_default, "database port")
	kafka_server := flag.String("kafka_server", kafka_host_default, "kafka server")
	kafka_port := flag.Int("kafka_port", kafka_port_default, "kafka port")
	flag.Parse()

	c.http_port = *port
	c.db_host = *dbserver
	c.db_port = *dbport
	c.db_user = db_user
	c.db_name = "consumer_geo"
	c.kafka_host = *kafka_server
	c.kafka_port = *kafka_port
	return c, nil
}

//func (c *Config) Set() error {
//}
