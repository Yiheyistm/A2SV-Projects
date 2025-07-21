# Task Manager Documentation

## Overview

Task Manager is a Go-based application for efficient task management, using MongoDB as its database. It supports creating, updating, deleting, and viewing tasks, suitable for both individuals and teams.

## Features

- **Task CRUD:** Create, read, update, and delete tasks.
- **Status Tracking:** Monitor task progress with status updates.
- **Due Dates:** Set and manage task deadlines.

## Installation

Make sure you have [MongoDB](https://www.mongodb.com/) installed and running.

```bash
git clone https://github.com/yourusername/task_manager.git
cd task_manager
go build
```

## Usage

Start the server (ensure MongoDB is running):

```bash
go run main.go
```

## API Documentation

For detailed API usage, request/response samples, and environment setup, see the official Postman documentation:

[Task Manager Postman Documentation](https://documenter.getpostman.com/view/37453586/2sB34ijej1)

This resource enables interactive API exploration and testing.

## API Reference

- `GET /tasks` — List all tasks
- `POST /tasks` — Create a new task
- `PUT /tasks/{id}` — Update a task
- `DELETE /tasks/{id}` — Delete a task

## Contributing

1. Fork the repository.
2. Create a feature branch.
3. Submit a pull request.
