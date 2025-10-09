package main

import (
	"injector/services"
	"injector/types"
	"log"
	"os"
	"time"
)

// type EdgeData struct {
// 	DeviceID string `json:"deviceId"`
// 	Location struct {
// 		To   int `json:"to"`
// 		From int `json:"from"`
// 	} `json:"location"`
// 	Data struct {
// 		Confiability float32 `json:"confiability"`
// 		Flow         float32 `json:"flow"`
// 	} `json:"data"`
// }

func main() {
	amqpURL := os.Getenv("AMQP_URL")
	queueName := os.Getenv("QUEUE_NAME")
	if amqpURL == "" || queueName == "" {
		log.Fatal("Variáveis de ambiente do broker não configuradas")
	}

	uri := os.Getenv("GRAPH_DB_URI")
	username := os.Getenv("GRAPH_DB_USERNAME")
	password := os.Getenv("GRAPH_DB_PASSWORD")
	if uri == "" || username == "" || password == "" {
		log.Fatal("Variáveis de ambiente do banco de grafos não configuradas")
	}

	// Neo4j
	injector := services.NewInjectorService(uri, username, password)
	defer injector.Close()
	injector.TestConnection()

	// RabbitMQ
	rabbit, err := services.NewRabbitService(amqpURL, queueName)
	if err != nil {
		log.Fatalf("Erro ao conectar RabbitMQ: %v", err)
	}
	defer rabbit.Close()

	ch := make(chan types.EdgeData, 100)

	go MessageWaiter(ch)
	go EdgeListBuilder(ch, 5*time.Second)

	select {}
}

func MessageWaiter(ch chan<- types.EdgeData) types.EdgeData {
	var e types.EdgeData
	e.DeviceID = "sensor-01"
	e.Location.From = 1
	e.Location.To = 2
	e.Data.Confiability = 0.93
	e.Data.Flow = 10.5
	return e
}

// Acumulador de Edges
func EdgeListBuilder(ch <-chan types.EdgeData, timeout time.Duration) {
	buffer := make([]types.EdgeData, 0, 50)
	timer := time.NewTimer(timeout)

	for {
		select {
		case edge := <-ch:
			buffer = append(buffer, edge)

			if len(buffer) >= cap(buffer) {
				Injector(buffer)
				buffer = buffer[:0]
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(timeout)
			}

		case <-timer.C:
			if len(buffer) > 0 {
				Injector(buffer)
				buffer = buffer[:0]
			}
			timer.Reset(timeout)
		}
	}
}

func Injector(data []EdgeData) {

}
