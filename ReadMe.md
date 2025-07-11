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

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/ElasticLogger.git
   cd ElasticLogger
   ```
2. Install dependencies:
```bash
    go mod tidy
```    

3. Edit ```config.yml`` file with your configs.

3. Start Elasticsearch and RabbitMQ:

## Usage

Running the Application :
```bash
    go run main.go
```    