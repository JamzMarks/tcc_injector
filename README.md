Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass

# run.ps1
$env:AMQP_URL="amqp://user:pass@host.docker.internal:5672/"
$env:QUEUE_NAME="test-queue"
$env:GRAPH_DB_URI="bolt://host.docker.internal:7687"
$env:GRAPH_DB_USERNAME="neo4j"
$env:GRAPH_DB_PASSWORD="neo4j"

go run main.go

module github.com/JamzMarks/tcc_injector.git

go 1.25.2

require (
	github.com/neo4j/neo4j-go-driver/v5 v5.28.4
	github.com/rabbitmq/amqp091-go v1.10.0
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/streadway/amqp v1.1.0
)
github.com/joho/godotenv v1.5.1 h1:7eLL/+HRGLY0ldzfGMeQkb7vMd0as4CfYvUVzLqw0N0=
github.com/joho/godotenv v1.5.1/go.mod h1:f4LDr5Voq0i2e/R5DDNOoa2zzDfwtkZa6DnEwAbqwq4=
github.com/neo4j/neo4j-go-driver/v5 v5.28.4 h1:7toxehVcYkZbyxV4W3Ib9VcnyRBQPucF+VwNNmtSXi4=
github.com/neo4j/neo4j-go-driver/v5 v5.28.4/go.mod h1:Vff8OwT7QpLm7L2yYr85XNWe9Rbqlbeb9asNXJTHO4k=
github.com/rabbitmq/amqp091-go v1.10.0 h1:STpn5XsHlHGcecLmMFCtg7mqq0RnD+zFr4uzukfVhBw=
github.com/rabbitmq/amqp091-go v1.10.0/go.mod h1:Hy4jKW5kQART1u+JkDTF9YYOQUHXqMuhrgxOEeS7G4o=
github.com/streadway/amqp v1.1.0 h1:py12iX8XSyI7aN/3dUT8DFIDJazNJsVJdxNVEpnQTZM=
github.com/streadway/amqp v1.1.0/go.mod h1:WYSrTEYHOXHd0nwFeUXAe2G2hRnQT+deZJJf88uS9Bg=