package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	_ "github.com/lib/pq"
)

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

func (c *Consumer) ConsumerGeoCollectionWriter(key, value []byte) error {
	var geos []ConsumerGeoInfo
	err := json.Unmarshal(value, &geos)
	if err != nil {
		log.Printf("json.Unmarshal failed!! %v", err)
		log.Printf("json:%s", string(value))
		return err
	}

	//database
	db, err := c.connect()
	if err != nil {
		log.Printf("database connect failed!!")
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("transaction begin failed!!")
		return err
	}
	defer tx.Rollback()

	for _, geo := range geos {
		log.Printf("%v\n", geo)
		if err := c.addConsumerGeo(db, geo); err != nil {
			log.Fatal(err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("commit faled!!")
		return err
	}

	return nil
}

func main() {

	config, err := NewConfig()
	if err != nil {
		log.Printf("init config failed: %v", err)
	}

	c := NewConsumer(config)

	srmConf := sarama.NewConfig()
	srmConf.Consumer.Return.Errors = true

	// Specify brokers address. This is default one
	brokers := []string{fmt.Sprintf("%s:%d", config.kafka_host, config.kafka_port)}
	log.Println("brokers=" + brokers[0])

	// Create new consumer
	master, err := sarama.NewConsumer(brokers, srmConf)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	// How to decide partition, is it fixed value...?
	consumer, err := master.ConsumePartition(config.kafka_topic, 0, sarama.OffsetNewest)
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
				c.ConsumerGeoCollectionWriter(msg.Key, msg.Value)
			case <-signals:
				log.Println("SIGTERM is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed", msgCount, "messages")
}
