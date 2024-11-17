# Tender Management Backend

## Overview

This project is a backend system for managing tenders and bids. It supports features like user authentication, tender creation, bid submission, and bid evaluation. Additionally, real-time notifications are provided via WebSockets. The backend uses Golang with Gin, PostgreSQL for database management, and Redis for caching and pub/sub functionality.

## Requirements

- Docker
- Docker Compose

## Setup

Follow the steps below to set up and run the project locally.

### 1. Enter project

```bash
cd tender-management-backend
```
### 2. Init postgresql database and redis client
```bash
make run_db
```

### 3. Run application

```bash
make run
```

### 4. Enter swagger for seeing routes
```bash
htttp://localhost:8888/swagger/index.html
```
