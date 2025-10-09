package services

import (
	"context"
	"fmt"
	"injector/types"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type InjectorService struct {
	uri      string
	username string
	password string
	driver   neo4j.DriverWithContext
}

func NewInjectorService(uri, username, password string) *InjectorService {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Fatalf("Erro ao criar driver Neo4j: %v", err)
	}
	return &InjectorService{
		uri:      uri,
		username: username,
		password: password,
		driver:   driver,
	}
}

func (s *InjectorService) Close() {
	if s.driver != nil {
		s.driver.Close(context.Background())
	}
}

func (s *InjectorService) TestConnection() {
	ctx := context.Background()
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.Run(ctx, "RETURN 'Conexão OK!' AS message", nil)
	if err != nil {
		log.Fatalf("Erro ao executar query: %v", err)
	}

	if result.Next(ctx) {
		fmt.Println(result.Record().Values[0])
	} else if err = result.Err(); err != nil {
		log.Fatalf("Erro ao ler resultado: %v", err)
	}
}

func (s *InjectorService) InjectBuffer(data []types.EdgeData) {

}
