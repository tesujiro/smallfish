package csgeo

import "flag"

const kafka_host_default = "localhost"
const kafka_port_default = 31090
const kafka_topic_default = "new_topic"

type Config struct {
	http_port   int
	https_port  int
	kafka_host  string
	kafka_port  int
	kafka_topic string
}

func NewConfig() (*Config, error) {
	c := &Config{}
	http_port := flag.Int("http_port", 80, "http port number")
	https_port := flag.Int("https_port", 443, "https port number")
	kafka_server := flag.String("kafka_server", kafka_host_default, "kafka server")
	kafka_port := flag.Int("kafka_port", kafka_port_default, "kafka port")
	kafka_topic := flag.String("kafka_topic", kafka_topic_default, "kafka topic")
	flag.Parse()

	c.http_port = *http_port
	c.https_port = *https_port
	c.kafka_host = *kafka_server
	c.kafka_port = *kafka_port
	c.kafka_topic = *kafka_topic
	return c, nil
}

//func (c *Config) Set() error {
//}
