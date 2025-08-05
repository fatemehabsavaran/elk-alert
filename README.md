# ELK Alert

`elk-alert` is a lightweight service for sending alerts based on APM data from **Elasticsearch**. It supports alerting through **Email**, **Slack**, **SMS**, and **Telegram**, making it easy to monitor and react to application performance issues.

---

## Project Structure

```bash

elk-alert/
├── cmd/ # Main application entry point
├── config/ # Configuration files 
├── internal/ # Internal packages (services, repositories)
├── pkg/ # Shared libraries (Elasticsearch, Redis, Telegram, etc.)
├── Dockerfile # Containerization setup
├── go.mod / go.sum # Go module dependencies
└── README.md

```


---

## Features

- Pulls APM metrics from Elasticsearch
- Supports multiple alert channels:
  - 📧 Email
  - 💬 Slack
  - 📱 SMS
  - 📲 Telegram
- Redis for caching or state management

---

## Configuration

Modify the `config/config.yaml` file to set your connection details for:

- Elasticsearch
- Redis
- Email/Slack/SMS/Telegram providers
- Alert thresholds and rules


## Docker

To build and run the project using Docker:
```bash
docker build -t elk-alert .
docker run -p 8080:8080 elk-alert
```