package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/lib/pq"
	csproto "github.com/tesujiro/smallfish/consumer/proto"
)

func connect(c *Config) (*sql.DB, error) {
	url := fmt.Sprintf("postgresql://%s@%s:%d/%s?sslmode=disable", c.db_user, c.db_host, c.db_port, c.db_name)
	db, err := sql.Open("postgres", url)
	if err != nil {
		//log.Fatal("error connecting to the database: ", err)
		fmt.Println("error connecting to the database: ", err)
	}
	return db, err
}

type Consumer struct {
	config Config
	conn   *sql.DB
}

func newConsumer(c *Config) *Consumer {
	//database connect
	db, err := connect(c)
	if err != nil {
		log.Printf("database connect failed!!")
		panic(err)
	}

	return &Consumer{config: *c, conn: db}
}

func (c *Consumer) finalize() error {
	return c.conn.Close()
}

func (c *Consumer) addConsumerGeo(db *sql.DB, geo *csproto.ConsumerGeo_Item) error {
	// Insert two rows into the "location" table.
	stmt, err := db.Prepare("INSERT INTO location (id, time, lat, lng) VALUES ($1,$2,$3,$4)")
	if err != nil {
		log.Printf("prepare statement faled!! %v", err)
		return err
	}
	defer stmt.Close() // danger!

	res, err := stmt.Exec(geo.ConsumerId, ptypes.TimestampString(geo.Timestamp), geo.Lat, geo.Lng)
	if err != nil {
		log.Printf("exec statement failed!!")
		return err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Printf("get rows affected failed!!")
		return err
	}
	log.Printf("affected = %d\n", rowCnt)

	return nil
}

func (c *Consumer) ConsumerGeoCollectionWriter(geos *csproto.ConsumerGeo) error {
	tx, err := c.conn.Begin()
	if err != nil {
		log.Printf("transaction begin failed!!")
		return err
	}
	defer tx.Rollback()

	for _, geo := range geos.ConsumerGeo {
		log.Printf("%v\n", geo)
		if err := c.addConsumerGeo(c.conn, geo); err != nil {
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

	c := newConsumer(config)
	defer c.finalize()

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
				log.Println("Received messages", string(msg.Key))
				geos := &csproto.ConsumerGeo{}
				if err := proto.Unmarshal(msg.Value, geos); err != nil {
					log.Printf("proto.Unmarshal failed!! %v", err)
					log.Printf("json:%s", string(msg.Value))
				}
				if err := c.ConsumerGeoCollectionWriter(geos); err != nil {
					log.Printf("ConsumerGeoCollectionWriter failed!! %v", err)
				}
			case <-signals:
				log.Println("SIGTERM is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed", msgCount, "messages")
}
