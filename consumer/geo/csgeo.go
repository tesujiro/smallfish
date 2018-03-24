package csgeo

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func Sleeper(w http.ResponseWriter, r *http.Request) {
	log.Printf("Sleeper!!")
	q := r.URL.Query()
	i, _ := strconv.Atoi(q.Get("timer"))
	time.Sleep(time.Duration(i) * time.Millisecond)
	fmt.Fprintf(w, "Hello, World :slept %d msec\n", i)
}

type Consumer struct {
	config Config
}

func NewConsumer(c *Config) *Consumer {
	return &Consumer{config: *c}
}

type ConsumerGeoInfo struct {
	ConsumerId int       `json:"consumerId"`
	Timestamp  time.Time `json:"timestamp"`
	Lat        float64   `json:"latitude"`
	Lng        float64   `json:"longtitude"`
}

const topic = "new_topic"

func (c *Consumer) KafkaProduce(key string, value string) error {
	log.Printf("Start sending messages to Kafka\n")
	// Setup configuration
	config := sarama.NewConfig()
	// Return specifies what channels will be populated.
	// If they are set to true, you must read from
	// config.Producer.Return.Successes = true
	// The total number of times to retry sending a message (default 3).
	config.Producer.Retry.Max = 5
	// The level of acknowledgement reliability needed from the broker.
	config.Producer.RequiredAcks = sarama.WaitForAll
	//brokers := []string{"localhost:9092"}
	//brokers := []string{"my-kafka-kafka:9092"}
	brokers := []string{fmt.Sprintf("%s:%d", c.config.kafka_host, c.config.kafka_port)}
	log.Printf("brokers=%v\n", brokers)
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		// Should not reach here
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			// Should not reach here
			panic(err)
		}
	}()

	producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}
	log.Printf("key=%v value=%v\n", key, value)

	return nil
}

func (c *Consumer) GeoCollectionWriter(w http.ResponseWriter, r *http.Request) {
	log.Printf("ConsumerGeoCollectionWriter!!")

	if r.Header.Get("Content-Type") != "application/json" {
		log.Printf("bad Content-Type!!")
		log.Printf(r.Header.Get("Content-Type"))
	}

	//To allocate slice for request body
	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		log.Printf("Content-Length failed!!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Read body data to parse json
	body := make([]byte, length)
	length, err = r.Body.Read(body)
	if err != nil && err != io.EOF {
		log.Printf("read failed!!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("Content-Length:%v", length)

	//Kafka
	if err := c.KafkaProduce(r.RemoteAddr, string(body)); err != nil {
		log.Printf("Kafka produce failed!!")
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Consumer) Router() *mux.Router {
	r := mux.NewRouter()
	//r.HandleFunc("/consumer/@{latitude:[0-9]+.?[0-9]+},{longtitude:[0-9]+.?[0-9]+}", ConsumerHandler).Methods("GET")
	r.HandleFunc("/consumer/GeoCollection", c.GeoCollectionWriter).Methods("POST")
	r.HandleFunc("/", Sleeper)
	return r
}

func Run(ctx context.Context) {

	config, err := NewConfig()
	if err != nil {
		log.Printf("init config failed: %v", err)
	}

	consumer := NewConsumer(config)

	http.Handle("/", consumer.Router())

	log.Printf("Start Go HTTP Server")

	err = http.ListenAndServe(":"+strconv.Itoa(consumer.config.http_port), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
