# ElasticLogger

ElasticLogger is a Go-based logging service that integrates with Elasticsearch for storing logs and RabbitMQ for consuming log messages. It provides a simple interface for logging structured log entries and supports event-driven logging via RabbitMQ.

## Features

- **Elasticsearch Integration**: Logs are stored in Elasticsearch for efficient querying and analysis.
- **RabbitMQ Consumer**: Consumes log messages from RabbitMQ and logs them to Elasticsearch.
- **REST API**: Provides an HTTP endpoint for logging entries via POST requests.

## Requirements

- Go 1.18 or later
- Elasticsearch 8.x
- RabbitMQ

## Installation(Manual Setup)

1. Clone the repository:
```bash
   git clone https://github.com/yourusername/ElasticLogger.git
   cd ElasticLogger
```
2. Install dependencies:
```bash
    go mod tidy
```    

3. Customize the ```config.yml``` file as needed (optional).

3. Start Elasticsearch and RabbitMQ:

4. Run the application:
```bash
   go run main.go
```

##  🐳 Running with Docker Compose (Recommended)

### Prerequisites
Docker and Docker Compose installed and running.

### Running the Services :
1. Clone the repository:
```bash
   git clone https://github.com/yourusername/ElasticLogger.git
   cd ElasticLogger
```
   
3. Customize the ```config.yml``` file as needed (optional).
4. Build and run the application and services:

```bash
   docker-compose up --build
```

### Stopping the Services
```bash
docker-compose down
```

