package main

import (
	"encoding/json"
	"fmt"
	"injector/services"
	"injector/types"
	"log"
	"os"
	"time"
)

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

	fmt.Print("ok")
	// Neo4j
	injector := services.NewInjectorService(uri, username, password)
	defer injector.Close()
	injector.TestConnection()

	// RabbitMQ
	rabbit, err := services.NewRabbitService(amqpURL, queueName)
	log.Printf("Escutando fila: %s no RabbitMQ: %s", queueName, amqpURL)

	if err != nil {
		log.Fatalf("Erro ao conectar RabbitMQ: %v", err)
	}
	defer rabbit.Close()

	ch := make(chan types.EdgeData, 100)

	go MessageWaiter(ch, rabbit)
	go EdgeListBuilder(ch, 5*time.Second, injector)

	select {}
}

func MessageWaiter(ch chan<- types.EdgeData, rabbit *services.RabbitService) {
	err := rabbit.Consume(func(body []byte) error {

		var e types.EdgeData
		if err := json.Unmarshal(body, &e); err != nil {
			return fmt.Errorf("erro ao decodificar mensagem: %w", err)
		}

		select {
		case ch <- e:
		default:
			log.Printf("Tipo: %s canal interno cheio, descartando mensagem de %s -> %f", e.DeviceType, e.DeviceId, e.Data.Flow)
		}

		return nil
	})

	if err != nil {
		log.Printf("erro ao iniciar consumo do RabbitMQ: %v", err)
	}

	log.Println("MessageWaiter iniciado e escutando mensagens...")
}

// Acumulador de Edges
func EdgeListBuilder(ch <-chan types.EdgeData, timeout time.Duration, injector *services.InjectorService) {
	log.Println("Mensagem consumida")
	buffer := make([]types.EdgeData, 0, 50)
	timer := time.NewTimer(timeout)

	for {
		select {
		case edge := <-ch:
			buffer = append(buffer, edge)

			if len(buffer) >= cap(buffer) {
				dataCopy := make([]types.EdgeData, len(buffer))
				copy(dataCopy, buffer)
				go injector.InjectBuffer(dataCopy)
				buffer = buffer[:0]
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(timeout)
			}

		case <-timer.C:
			if len(buffer) > 0 {
				dataCopy := make([]types.EdgeData, len(buffer))
				copy(dataCopy, buffer)
				go injector.InjectBuffer(dataCopy)
				buffer = buffer[:0]
			}
			timer.Reset(timeout)
		}
	}
}

// package main

// import (
// 	"log"
// 	"os"

// 	"github.com/streadway/amqp"
// )

// func main() {
// 	// Pega variáveis de ambiente
// 	amqpURL := os.Getenv("AMQP_URL")
// 	queueName := os.Getenv("QUEUE_NAME")

// 	if amqpURL == "" || queueName == "" {
// 		log.Fatal("Variáveis de ambiente AMQP_URL e QUEUE_NAME não configuradas")
// 	}

// 	log.Printf("Conectando ao RabbitMQ: %s", amqpURL)

// 	// Conexão
// 	conn, err := amqp.Dial(amqpURL)
// 	if err != nil {
// 		log.Fatalf("Erro ao conectar RabbitMQ: %v", err)
// 	}
// 	defer conn.Close()
// 	log.Println("✅ Conexão estabelecida")

// 	// Canal
// 	ch, err := conn.Channel()
// 	if err != nil {
// 		log.Fatalf("Erro ao abrir canal: %v", err)
// 	}
// 	defer ch.Close()

// 	// Garante que a fila existe
// 	q, err := ch.QueueDeclare(
// 		queueName, // nome da fila
// 		true,      // durável
// 		false,     // autodelete
// 		false,     // exclusiva
// 		false,     // no-wait
// 		nil,       // args
// 	)
// 	if err != nil {
// 		log.Fatalf("Erro ao declarar fila: %v", err)
// 	}

// 	log.Printf("🟢 Aguardando mensagens na fila: %s", q.Name)

// 	// Consumidor
// 	msgs, err := ch.Consume(
// 		q.Name, // fila
// 		"",     // consumer
// 		true,   // auto-ack
// 		false,  // exclusive
// 		false,  // no-local
// 		false,  // no-wait
// 		nil,    // args
// 	)
// 	if err != nil {
// 		log.Fatalf("Erro ao iniciar consumo: %v", err)
// 	}

// 	// Loop de consumo
// 	forever := make(chan bool)

// 	go func() {
// 		for d := range msgs {
// 			log.Printf("📩 Mensagem recebida: %s", d.Body)
// 		}
// 	}()

// 	log.Println("🚀 Aguardando mensagens... (CTRL+C para sair)")
// 	<-forever
// }
