package csgeo

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	csproto "github.com/tesujiro/smallfish/consumer/proto"
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

func (c *Consumer) ConsumerAutoPost(w http.ResponseWriter, r *http.Request) {
	log.Printf("Consumer Auto Post Page!!")
	tpl := template.Must(template.ParseFiles("template/AutoPost.html"))
	w.Header().Set("Content-Type", "text/html")

	err := tpl.Execute(w, map[string]string{"APIKEY": os.Getenv("APIKEY")})
	if err != nil {
		panic(err)
	}
}

func (c *Consumer) ConsumerManualTester(w http.ResponseWriter, r *http.Request) {
	log.Printf("Consumer Manual Tester Page!!")
	tpl := template.Must(template.ParseFiles("template/ManualTester.html"))
	w.Header().Set("Content-Type", "text/html")

	err := tpl.Execute(w, map[string]string{"APIKEY": os.Getenv("APIKEY")})
	if err != nil {
		panic(err)
	}
}

//func (c *Consumer) KafkaProduce(key string, value string) error {
func (c *Consumer) KafkaProduce(key string, value []byte) error {
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

	// JSON to ProtoBuf
	var geos []ConsumerGeoInfo
	if err := json.Unmarshal(value, &geos); err != nil {
		log.Printf("json.Unmarshal failed!! %v\n", err)
		log.Printf("json:%s\n", string(value))
		return err
	}

	geos_proto := &csproto.ConsumerGeo{}
	for _, geo := range geos {
		ts, err := ptypes.TimestampProto(geo.Timestamp)
		if err != nil {
			log.Printf("ptypes.TimestampProto failed!! %v\n", err)
			return err
		}
		geos_proto.ConsumerGeo = append(geos_proto.ConsumerGeo, &csproto.ConsumerGeo_Item{
			ConsumerId: int64(geo.ConsumerId),
			Timestamp:  ts,
			Lat:        geo.Lat,
			Lng:        geo.Lng,
		})
	}

	msg, err := proto.Marshal(geos_proto)
	if err != nil {
		log.Printf("proto.Marshal failed!! %v\n", err)
		log.Printf("json:%s\n", geos)
		return err
	}

	// send a message to kafka
	producer.Input() <- &sarama.ProducerMessage{
		Topic: c.config.kafka_topic,
		Key:   sarama.StringEncoder(key),
		//Value: sarama.StringEncoder(value),
		Value: sarama.ByteEncoder(msg),
	}
	log.Printf("key=%v value=%v\n", key, geos)

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
	//if err := c.KafkaProduce(r.RemoteAddr, string(body)); err != nil {
	if err := c.KafkaProduce(r.RemoteAddr, body); err != nil {
		log.Printf("Kafka produce failed!!")
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Consumer) Router() *mux.Router {
	r := mux.NewRouter()
	//r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	//r.HandleFunc("/consumer/@{latitude:[0-9]+.?[0-9]+},{longtitude:[0-9]+.?[0-9]+}", ConsumerHandler).Methods("GET")
	r.HandleFunc("/consumer/GeoCollection", c.GeoCollectionWriter).Methods("POST")
	r.HandleFunc("/consumer/manualTester", c.ConsumerManualTester)
	r.HandleFunc("/consumer/auto", c.ConsumerAutoPost)
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
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Printf("Start Go HTTP Server")

	go func() {
		err = http.ListenAndServe(":"+strconv.Itoa(consumer.config.http_port), nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	err = http.ListenAndServeTLS(":"+strconv.Itoa(consumer.config.https_port), "ssl/server.crt", "ssl/server.key", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
