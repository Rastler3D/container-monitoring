# Container Monitor

This project is a container monitoring system that pings Docker containers and displays their status on a web interface. It uses RabbitMQ for message queuing, Golang for backend, Nginx as a API Gateway, React.js for the frontend.

## Components

- Backend: Go service providing RESTful API and consuming messages from RabbitMQ
- Frontend: React.js application for displaying container statuses
- Pinger: Go service for pinging containers and sending data to RabbitMQ
- Database: PostgreSQL for storing container status data
- RabbitMQ: Message queue for communication between Pinger and Backend
- Nginx: API gateway for routing backend requests

## Prerequisites

- Docker
- Docker Compose

## Setup and Running

1. Clone the repository:
   ```
   git clone https://github.com/Rastler3D/container-monitoring.git
   cd container-monitoring
   ```

2. Build and run the application using Docker Compose:
   ```
   docker-compose up --build
   ```

3. The application will be available at `http://localhost:80`

## API Endpoints

### Get container statuses

```
GET /api/containers
```

### Add container statuses

```
GET /notes
Content-Type: application/json

[
    {
        "ip": "172.19.0.3",
        "ping_time": "0.32",
        "last_ping": "2025-02-09T19:53:05.467033Z"
    }
]
```

## Configuration

You can modify the configuration of each service by editing the respective configuration files in each service's directory.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
