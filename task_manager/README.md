# Task Manager Documentation

## Overview

Task Manager is a Go-based application designed to help users manage tasks efficiently. It provides features for creating, updating, deleting, and viewing tasks, supporting both individual and team workflows.

## Features

- **Task CRUD:** Create, read, update, and delete tasks.
- **User Management:** Assign tasks to users and manage user roles.
- **Status Tracking:** Track task progress with status updates.
- **Due Dates:** Set and manage deadlines for tasks.
- **Filtering & Search:** Find tasks by status, assignee, or due date.

## Installation

```bash
git clone https://github.com/yourusername/task_manager.git
cd task_manager
go build
```

## Usage

Start the server:

```bash
go run main.go
```

## Postman API Documentation

For detailed API usage, request/response examples, and environment setup, refer to the official Postman documentation:

[Task Manager Postman Documentation](https://documenter.getpostman.com/view/37453586/2sB34ijej1)

This resource provides interactive API exploration and testing capabilities.

## API Reference

- `GET /tasks` - List all tasks
- `POST /tasks` - Create a new task
- `PUT /tasks/{id}` - Update a task
- `DELETE /tasks/{id}` - Delete a task

## Contributing

1. Fork the repository.
2. Create a feature branch.
3. Submit a pull request.
