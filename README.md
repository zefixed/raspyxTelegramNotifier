# Raspyx Telegram Notifier

A Telegram notifier service written in Go, designed to receive and forward kafka notifications directly to Telegram chats. It was developed as a microservice for Raspyx API.
## 📖 Table of Contents

- [📄 Description](#-description)
- [✨ Features](#-features)
- [🛠️ Technologies](#%EF%B8%8F-technologies)
- [📥 Installation and Setup](#-installation-and-setup)
  - [🖥️ Local Setup](#%EF%B8%8F-local-setup)
  - [📦 Docker Setup](#-docker-setup)
- [✅ Testing](#-testing)
- [📜 License](#-license)

## 📄 Description

This service acts as a bridge between Raspyx API and Telegram, forwarding messages to bot’s chat. Useful as monitoring tools.

## ✨ Features

- Database connection with PostgreSQL
- Database migration with goose
- Message streaming with Kafka
- Request and error logging
- CI/CD integration with Jenkins

## 🛠️ Technologies

- Go (1.24)
- PostgreSQL
- Goose
- Kafka
- Docker
- Jenkins

## 📥 Installation and Setup

1. Clone the repository:

```bash
git clone https://github.com/zefixed/raspyxTelegramNotifier.git
cd raspyxTelegramNotifier
```

> ❗ Before using rename `.env.example` to `.env` and set up your parameters

### 🖥️ Local Setup

To run the application locally, follow these steps:

1. Install go

   Arch

   ```bash
   yay -Sy go
   ```

   Debian

   ```bash
   sudo apt install golang
   ```
   
2. Run app

   ```bash
   make all
   ```

### 📦 Docker Setup

To run the application with Docker, follow these steps:

1. Run the docker-compose:

   ```bash
   docker compose up --build -d
   ```


## ✅ Testing

To test the API, you can use [Postman](https://www.postman.com/) or [cURL](https://curl.se/). You can also set up unit tests in the project using:

```bash
go test -v ./...
```

## 📜 License

This project is licensed under the GNU License v3 - see the [LICENSE](LICENSE) file for details.
