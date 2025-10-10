package main

import (
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

    go MessageWaiter(ch, rabbit)
    go EdgeListBuilder(ch, 5*time.Second)

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
            log.Printf("canal interno cheio, descartando mensagem de %s -> %d", e.DeviceID, e.Location.To)
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

