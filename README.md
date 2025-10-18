Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass

# run.ps1
$env:AMQP_URL="amqp://user:pass@localhost:5672/"
$env:QUEUE_NAME="test-queue"
$env:GRAPH_DB_URI="bolt://localhost:7687"
$env:GRAPH_DB_USERNAME="neo4j"
$env:GRAPH_DB_PASSWORD="neo4j"

go run main.go

