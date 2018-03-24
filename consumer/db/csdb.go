package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	_ "github.com/lib/pq"
)

/*
type Consumer struct {
	//config Config
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

func (c *Consumer) connect() (*sql.DB, error) {
	url := fmt.Sprintf("postgresql://%s@%s:%d/%s?sslmode=disable", c.config.db_user, c.config.db_host, c.config.db_port, c.config.db_name)
	db, err := sql.Open("postgres", url)
	if err != nil {
		//log.Fatal("error connecting to the database: ", err)
		fmt.Println("error connecting to the database: ", err)
	}
	return db, err
}

func (c *Consumer) addConsumerGeo(db *sql.DB, geo ConsumerGeoInfo) error {
	// Insert two rows into the "location" table.
	stmt, err := db.Prepare("INSERT INTO location (id, time, lat, lng) VALUES ($1,$2,$3,$4)")
	if err != nil {
		log.Printf("prepare statement faled!! %v", err)
		return err
	}
	defer stmt.Close() // danger!

	res, err := stmt.Exec(geo.ConsumerId, geo.Timestamp, geo.Lat, geo.Lng)
	if err != nil {
		log.Printf("exec statement faled!!")
		return err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Printf("get rows affected faled!!")
		return err
	}
	log.Printf("affected = %d\n", rowCnt)

	return nil
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
*/

const topic = "new_topic"
const kafka_host_default = "localhost"
const kafka_port_default = 31090
const zookeeper_host_default = "my-kafka-zookeeper"
const zookeeper_port_default = 2181

func main() {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Specify brokers address. This is default one
	//brokers := []string{"localhost:9092"}
	//brokers := []string{fmt.Sprintf("%s:%d", zookeeper_host_default, zookeeper_port_default)}
	brokers := []string{fmt.Sprintf("%s:%d", kafka_host_default, kafka_port_default)}

	// Create new consumer
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	// How to decide partition, is it fixed value...?
	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)

	// Count how many message processed
	msgCount := 0

	// Get signnal for finish
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Println(err)
			case msg := <-consumer.Messages():
				msgCount++
				log.Println("Received messages", string(msg.Key), string(msg.Value))
			case <-signals:
				log.Println("SIGTERM is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed", msgCount, "messages")
}
