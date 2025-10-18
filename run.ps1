# run.ps1
$env:AMQP_URL="amqp://user:pass@localhost:5672/"
$env:QUEUE_NAME="injector_queue"
$env:GRAPH_DB_URI="neo4j://127.0.0.1:7687"
$env:GRAPH_DB_USERNAME="neo4j"
$env:GRAPH_DB_PASSWORD="senha123"

go run main.go
