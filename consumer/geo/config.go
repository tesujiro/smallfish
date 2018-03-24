package csgeo

import "flag"

const kafka_host_default = "localhost"
const kafka_port_default = 31090

type Config struct {
	http_port  int
	kafka_host string
	kafka_port int
}

func NewConfig() (*Config, error) {
	c := &Config{}
	port := flag.Int("port", 80, "port number")
	kafka_server := flag.String("kafka_server", kafka_host_default, "kafka server")
	kafka_port := flag.Int("kafka_port", kafka_port_default, "kafka port")
	flag.Parse()

	c.http_port = *port
	c.kafka_host = *kafka_server
	c.kafka_port = *kafka_port
	return c, nil
}

//func (c *Config) Set() error {
//}
