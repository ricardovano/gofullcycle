package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gofullcycle/internal/infra/database"
	"gofullcycle/internal/usecase"
	"gofullcycle/pkg/kafka"
	"gofullcycle/pkg/rabbitmq"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "usr_orders"
	password = "usr_orders"
	dbname   = "orders"
)

func getConnectionString(host string, port int, user string, password string, dbname string) string {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	return psqlInfo
}

func openDatabase() *sql.DB {
	db, err := sql.Open("postgres", getConnectionString(host, port, user, password, dbname))
	fmt.Println("Postgres opened")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func main() {

	repository := database.NewOrderRepository(openDatabase())
	usecase := usecase.CalculateFinalPrice{OrderRepository: repository}

	msgChanKafka := make(chan *ckafka.Message)

	topics := []string{"orders"}
	servers := "localhost:9092"
	fmt.Println("Kafka consumer has started with server!")
	go kafka.Consume(topics, servers, msgChanKafka)
	go kafkaWorker(msgChanKafka, usecase)

	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	msgRabbitmqChannel := make(chan amqp.Delivery)
	go rabbitmq.Consume(ch, msgRabbitmqChannel)
	rabbitmqWorker(msgRabbitmqChannel, usecase)
}

func kafkaWorker(msgChan chan *ckafka.Message, uc usecase.CalculateFinalPrice) {
	fmt.Println("Kafta worker has started!")
	for msg := range msgChan {
		var OrderInputDTO usecase.OrderInputDTO
		err := json.Unmarshal(msg.Value, &OrderInputDTO)
		if err != nil {
			panic(err)
		}
		outputDTO, err := uc.Execute(OrderInputDTO)
		if err != nil {
			panic(err)
		}
		fmt.Println("Kafka processed order %i\n", outputDTO.Id)
	}
}

func rabbitmqWorker(msgChan chan amqp.Delivery, uc usecase.CalculateFinalPrice) {
	fmt.Println("Rabbitmq worker has started")
	for msg := range msgChan {
		var OrderInputDTO usecase.OrderInputDTO
		err := json.Unmarshal(msg.Body, &OrderInputDTO)
		if err != nil {
			panic(err)
		}
		outputDTO, err := uc.Execute(OrderInputDTO)
		if err != nil {
			panic(err)
		}
		fmt.Println("Rabbitmq processed order %i\n", outputDTO.Id)
	}
}
