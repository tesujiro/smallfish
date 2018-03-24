package main

import "flag"

const db_host_default = "localhost"
const db_port_default = 30257
const db_user = "maxroach"
const kafka_host_default = "localhost"
const kafka_port_default = 31090
const kafka_topic_default = "new_topic"

//const db_host = "cockroachdb-public"

type Config struct {
	db_host     string
	db_port     int
	db_user     string
	db_name     string
	kafka_host  string
	kafka_port  int
	kafka_topic string
}

func NewConfig() (*Config, error) {
	c := &Config{}
	dbserver := flag.String("dbserver", db_host_default, "database server")
	dbport := flag.Int("dbport", db_port_default, "database port")
	kafka_server := flag.String("kafka_server", kafka_host_default, "kafka server")
	kafka_port := flag.Int("kafka_port", kafka_port_default, "kafka port")
	kafka_topic := flag.String("kafka_topic", kafka_topic_default, "kafka topic")
	flag.Parse()

	c.db_host = *dbserver
	c.db_port = *dbport
	c.db_user = db_user
	c.db_name = "consumer_geo"
	c.kafka_host = *kafka_server
	c.kafka_port = *kafka_port
	c.kafka_topic = *kafka_topic
	return c, nil
}

//func (c *Config) Set() error {
//}
