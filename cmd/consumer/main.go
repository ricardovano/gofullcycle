package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gofullcycle/internal/infra/database"
	"gofullcycle/internal/usecase"
	"gofullcycle/pkg/kafka"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
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
	fmt.Println(servers)
	fmt.Println("Kafka consumer has started with server!")
	go kafka.Consume(topics, servers, msgChanKafka)
	kafkaWorker(msgChanKafka, usecase)
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
