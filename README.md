# Go API Gateway

This project implements an API Gateway in Go, utilizing gRPC for communication between services. The gateway is designed to handle incoming requests, route them to the appropriate microservices, and aggregate responses.

## Features

- **gRPC Integration**: The API Gateway uses gRPC for efficient and type-safe communication with microservices.
- **Environment Variables**: Environment variables are used to configure the application, allowing for easy deployment across different environments.

## Prerequisites

- Go
- Docker
- Make

## Installation

1. Clone the repository:
```bash
git clone git@github.com:SeiFlow-3P2/api_gateway.git
cd api_gateway
```

2. Install dependencies:
```bash
go mod download
```

3. Configure the application:
  - Update the `config/config.yaml` file with your service configurations.
  - Set environment variables in the `.env` file.

## Usage

To run the API Gateway, execute the following command:

```bash
make env
make api-gateway-up
```

Or running while developing:

```bash
go run cmd/main.go
```

## Configuration

The `config.yaml` file contains settings for the API Gateway, including:

- Service name, host, shutdown timeout
- Protected routes

## Environment Variables

The `.env` file is used to set environment-specific variables, such as:

- `APP_MODE`: The current environment (e.g., development, production).
- `PORT`: The port on which the API Gateway will listen.
- `BOARD_SERVICE_ADDR`: The address of the board service.
- `PAYMENT_SERVICE_ADDR`: The address of the payment service.
- `CALENDAR_SERVICE_ADDR`: The address of the calendar service.
- `AUTH_SERVICE_ADDR`: The address of the auth service.
- `OTEL_ADDR`: The address of the Open Telemetry Collector.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
