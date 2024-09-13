# Message Aggregator

The **Message Aggregator** project aims to consolidate messages from different platforms such as WhatsApp, Telegram, Email, Slack, and iMessage into a single dashboard. The system allows users to not only view messages but also reply directly from the dashboard. This project utilizes Matrix as the underlying open protocol for decentralized communication, enabling the aggregation of messages from various platforms into a unified channel.

## Features

- **Matrix Integration**: Messages from various platforms are aggregated using the Matrix protocol.
- **Slack Bridge**: Messages from Slack are received and can be replied to within the system.
- **Multi-platform Support**: The platform can be extended to include WhatsApp, Telegram, Email, and more.
- **Dockerized Setup**: The project is containerized using Docker and Docker Compose for easy deployment.

## Tech Stack

### Backend:
- **Language**: Go (Golang)
- **Web Framework**: Fiber
- **Database**: PostgreSQL with GORM
- **Message Queue**: RabbitMQ
- **Caching**: Redis

### Frontend:
- **Framework**: React with TypeScript
- **UI Library**: Material-UI
- **State Management**: React Query

## Prerequisites

- Docker and Docker Compose
- Go
- PostgreSQL
- Redis
- RabbitMQ

## Getting Started

### 1. Clone the repository:
```bash
git clone https://github.com/asmitsharp/message-aggregator.git
cd message-aggregator
```

### 2. Build and run the containers:
```bash
docker-compose up --build
```

## Matrix Integration

The project uses the gomatrix library for interaction with the Matrix homeserver. To set up the Matrix homeserver:
- **Configure Synapse**: Follow the official Matrix Synapse setup to deploy your homeserver.
- **Environment Variables**: Ensure all sensitive environment variables such as database credentials, Matrix homeserver details, and API keys are managed securely.

## Slack Bridge

The Slack bridge integration allows you to receive and reply to messages from Slack. To set up the Slack bridge:

- Set up your Slack App with the necessary permissions.
- Configure the bridge in your Matrix server to start receiving Slack messages in Matrix rooms.

## RoadMap

- Complete bridges for WhatsApp, Telegram, Email, and iMessage.
- Develop frontend for a more intuitive user experience.
